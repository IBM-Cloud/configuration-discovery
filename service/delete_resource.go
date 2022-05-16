package service

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"time"

	"github.com/IBM-Cloud/configuration-discovery/base"
	"github.com/IBM-Cloud/configuration-discovery/tfplugin"
	"github.com/IBM-Cloud/configuration-discovery/utils"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	tfjson "github.com/hashicorp/terraform-json"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

// Remove de-provisioned resources from local .tf/state files
func DeleteDriftResources(ctx context.Context, services []string, confDir, randomID string, planTimeOut time.Duration) error {
	logger := utils.GetLogger(ctx)
	logger.Say("INFO: Remove the drift resources from .tf/state file from local repo.")

	//Get drift resources
	driftResourceList, err := getDriftResources(ctx, services, confDir, randomID)
	if err != nil {
		logger.Failed("ERROR: %v.", err)
		return err
	}

	if len(driftResourceList) > 0 {
		//Check any resource depends on drift resources
		err := getDependentResource(ctx, driftResourceList)
		if err != nil {
			logger.Failed("ERROR:  %v", err)
			return err
		}

		//Create backup folder
		backupFolder, err := utils.CreateBackupFolder(base.DiscoveryBackupFolder)
		if err != nil {
			logger.Failed("ERROR:  %v", err)
			return err
		}

		//Take backup of terraform.tfstate file
		srcFile := confDir + "/terraform.tfstate"
		destFile := confDir + "/" + backupFolder + "/terraform.tfstate_bkp"
		err = utils.Copy(srcFile, destFile)
		if err != nil {
			logger.Failed("ERROR:  %v", err)
			return err
		}

		//Remove resources from .tf/state file
		r, _ := regexp.Compile("^[A-Za-z-_]+\\.tf$")
		files, err := ioutil.ReadDir(".")
		if err != nil {
			logger.Failed("ERROR:  %v", err)
		}

		var toBeWritten *hclwrite.File
		for _, file := range files {
			if r.MatchString(file.Name()) {
				info, err := os.Lstat(file.Name())
				b, err := ioutil.ReadFile(file.Name())
				if err != nil {
					logger.Failed("ERROR:  %v", err)
					return err
				}

				backupCount := 0
				toBeWritten, _ = hclwrite.ParseConfig(b, file.Name(), hcl.InitialPos)
				for _, block := range toBeWritten.Body().Blocks() {
					resources := block.Labels()
					if block.Type() == "resource" {
						resTypeAndName := resources[0] + "." + resources[1]
						_, found := Find(driftResourceList, resTypeAndName)
						if found {
							backupCount++
							//Take Backup of .TF file.
							if backupCount == 1 {
								srcFile = confDir + "/" + file.Name()
								destFile = confDir + "/" + backupFolder + "/" + file.Name() + "_bkp"
								err = utils.Copy(srcFile, destFile)
								if err != nil {
									logger.Failed("ERROR:  %v", err)
									return err
								}
							}

							//Remove resource block from .tf file
							toBeWritten.Body().RemoveBlock(block)
							//Remove resource from terraform.tfstate file
							err = tfplugin.TerraformStateRemove(confDir, resTypeAndName, randomID, planTimeOut)
							if err != nil {
								logger.Say("INFO: Failed to remove resource '%v' from .tf file.", resTypeAndName)
								return err
							}
							logger.Say("INFO: Removed resource '%v' from .tf/state file sucessfully.", resTypeAndName)
							continue
						}
					}
				}
				// Write back the update content to .tf file
				err = ioutil.WriteFile(file.Name(), toBeWritten.Bytes(), info.Mode())
				if err != nil {
					logger.Failed("ERROR:  %v", err)
					return nil
				}

			}
		}
		logger.Say("INFO:  Discovery deleted resources '%v' successfully.", driftResourceList)
	} else {
		logger.Say("INFO: No drift detected in local state file!!")
	}

	return nil
}

// getDriftResources ..
func getDriftResources(ctx context.Context, services []string, confDir, randomID string) ([]string, error) {
	logger := utils.GetLogger(ctx)
	var driftResourceList []string

	//Generate terraform plan
	planText := "-out=" + base.PlanTextFile
	err := tfplugin.TerraformPlan(confDir, planText, planTimeOut, randomID)
	if err != nil {
		logger.Failed("ERROR:  Failed to generate terraform plan : %v", err)
		return nil, err
	}

	//Convert plan file to json format
	planJSON := base.PlanTextFile + " > " + base.PlanJSONFile
	err = tfplugin.GenerateTerraformPlanJson(confDir, planJSON, planTimeOut, randomID)
	if err != nil {
		logger.Failed("ERROR:  Failed to convert terraform plan file to json : %v", err)
		return nil, err
	}

	//Unmarshal plan json file
	var planIntf tfjson.Plan
	planData, err := ReadFile(ctx, base.PlanJSONFile)
	planErr := planIntf.UnmarshalJSON(planData)
	if planErr != nil {
		logger.Failed("ERROR:  Failed to unmarshal the plan file", zap.Error(planErr))
		return nil, err
	}

	//Get drift resources
	for _, resource := range planIntf.ResourceChanges {
		for _, action := range resource.Change.Actions {
			if action == tfjson.ActionCreate {
				if ValidateResourceType(services, resource.Type) {
					driftResourceList = append(driftResourceList, resource.Address)
				} else {
					logger.Warn("WARN: : Drift resource '%v' is not part of services '%v'", resource.Address, services)
				}
			}
		}
	}

	return driftResourceList, nil
}

// getDependentResource ..
func getDependentResource(ctx context.Context, driftResources []string) error {
	logger := utils.GetLogger(ctx)
	var dependentResourceList []string
	r, _ := regexp.Compile("^[A-Za-z-_]+\\.tf$")

	// Check is there any other resource which has got reference of the removing resource
	// If yes, Prepare the list and ask user to fix the reference issue before delete operation
	var toBeWritten *hclwrite.File
	files, err := ioutil.ReadDir(".")
	if err != nil {
		logger.Failed("ERROR:  %v", err)
		return err
	}

	//Prepare resource reference list
	for _, file := range files {
		if r.MatchString(file.Name()) {
			b, err := ioutil.ReadFile(file.Name())
			if err != nil {
				logger.Failed("ERROR:  %v", err)
				return err
			}

			toBeWritten, _ = hclwrite.ParseConfig(b, file.Name(), hcl.InitialPos)
			for _, block := range toBeWritten.Body().Blocks() {
				var resTypeAndName string
				resources := block.Labels()
				if block.Type() == "data" {
					resTypeAndName = "data." + resources[0] + "." + resources[1]
				} else if block.Type() == "resource" {
					resTypeAndName = resources[0] + "." + resources[1]
				} else if block.Type() == "output" {
					resTypeAndName = "output." + resources[0]
				} else {
					resTypeAndName = block.Type()
				}

				//Validate resource is not part of drift resource
				if !isValueInList(driftResources, resTypeAndName, false) {
					//Validate non-drift resource is dependent on drift resources
					if checkMapContainsValue(driftResources, block.Body().Attributes()) {
						dependentResourceList = append(dependentResourceList, resTypeAndName)
					}
				}
			}
		}
	}

	if len(dependentResourceList) > 0 {
		logger.Warn("WARN: Discovery tool going to delete drift resources '%v' from .tf/state file.", driftResources)
		logger.Warn("WARN: But there is a references of drift resources in following resources '%v', Please remove the drift resource reference from .tf file and run the command again..", dependentResourceList)
		return fmt.Errorf("Failed to remove the drift resource '%v' reference from .tf file.", driftResources)
	}

	return nil
}

// ValidateResourceType ..
func ValidateResourceType(services []string, resourceType string) bool {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	resourceFile, err := ioutil.ReadFile(basepath + "/resources.yml")
	if err != nil {
		panic(err)
	}

	resourceMap := make(map[string][]string)
	err = yaml.Unmarshal(resourceFile, resourceMap)
	if err != nil {
		panic(err)
	}

	for _, service := range services {
		if resourceMap[service] != nil {
			for _, rType := range resourceMap[service] {
				if resourceType == rType {
					return true
				}
			}
		}
	}

	return false
}
