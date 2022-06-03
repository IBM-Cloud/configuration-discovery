package service

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"regexp"

	"time"

	"github.com/IBM-Cloud/configuration-discovery/base"
	"github.com/IBM-Cloud/configuration-discovery/tfplugin"
	"github.com/IBM-Cloud/configuration-discovery/utils"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	tfjson "github.com/hashicorp/terraform-json"
	"go.uber.org/zap"
)

// Remove de-provisioned resources from local .tf/state files
func DeleteDriftResources(ctx context.Context, services []string, localTFDir, randomID string, planTimeOut time.Duration) error {
	logger := utils.GetLogger(ctx)
	logger.Say("INFO: Remove the drift resources from .tfstate file from local repo.")

	// get drift resources
	driftResourceList, err := getDriftResources(ctx, services, localTFDir, randomID)
	if err != nil {
		logger.Failed("ERROR: %v.", err)
		return err
	}

	if len(driftResourceList) > 0 {
		// check any resource depends on drift resources
		err := getDependentResource(ctx, driftResourceList)
		if err != nil {
			logger.Failed("ERROR:  %v", err)
			return err
		}

		// create backup folder
		backupFolder, err := utils.CreateBackupFolder(base.DiscoveryBackupFolder)
		if err != nil {
			logger.Failed("ERROR:  %v", err)
			return err
		}

		// take backup of terraform.tfstate file
		srcFile := localTFDir + utils.PathSep + "terraform.tfstate"
		destFile := localTFDir + utils.PathSep + backupFolder + utils.PathSep + "terraform.tfstate_bkp"
		err = utils.Copy(srcFile, destFile)
		if err != nil {
			logger.Failed("ERROR:  %v", err)
			return err
		}

		// remove resources from .tf/state file
		files, err := ioutil.ReadDir(".")
		if err != nil {
			logger.Failed("ERROR:  %v", err)
		}

		var toBeWritten *hclwrite.File
		for _, file := range files {
			match, _ := regexp.MatchString("^[A-Za-z-_]+\\.tf$", file.Name())
			if match {
				info, err := os.Lstat(file.Name())
				if err != nil {
					logger.Failed("ERROR:  %v", err)
					return err
				}
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
							// take Backup of .TF file.
							if backupCount == 1 {
								srcFile = localTFDir + utils.PathSep + file.Name()
								destFile = localTFDir + utils.PathSep + backupFolder + utils.PathSep + file.Name() + "_bkp"
								err = utils.Copy(srcFile, destFile)
								if err != nil {
									logger.Failed("ERROR:  %v", err)
									return err
								}
							}

							// remove resource block from .tf file
							toBeWritten.Body().RemoveBlock(block)
							// remove resource from terraform.tfstate file
							err = tfplugin.TerraformStateRemove(localTFDir, resTypeAndName, randomID, planTimeOut)
							if err != nil {
								logger.Say("INFO: Failed to remove resource '%v' from .tf file.", resTypeAndName)
								return err
							}
							logger.Say("INFO: Removed resource '%v' from .tf/state file sucessfully.", resTypeAndName)
							continue
						}
					}
				}
				// write back the update content to .tf file
				if backupCount > 0 {
					err = ioutil.WriteFile(file.Name(), toBeWritten.Bytes(), info.Mode())
					if err != nil {
						logger.Failed("ERROR:  %v", err)
						return nil
					}
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
func getDriftResources(ctx context.Context, services []string, localTFDir, randomID string) ([]string, error) {
	logger := utils.GetLogger(ctx)
	var driftResourceList []string

	// generate terraform plan
	planText := "-out=" + localTFDir + utils.PathSep + base.PlanTextFile
	err := tfplugin.TerraformPlan(localTFDir, planText, planTimeOut, randomID)
	if err != nil {
		logger.Failed("ERROR:  Failed to generate terraform plan : %v", err)
		return nil, err
	}

	// convert plan file to json format
	planJSON := localTFDir + utils.PathSep + base.PlanTextFile + " > " + localTFDir + utils.PathSep + base.PlanJSONFile
	err = tfplugin.GenerateTerraformPlanJson(localTFDir, planJSON, planTimeOut, randomID)
	if err != nil {
		logger.Failed("ERROR:  Failed to convert terraform plan file to json : %v", err)
		return nil, err
	}

	// unmarshal plan json file
	var planIntf tfjson.Plan
	planData, err := ReadFile(ctx, localTFDir+utils.PathSep+base.PlanJSONFile)
	planErr := planIntf.UnmarshalJSON(planData)
	if planErr != nil {
		logger.Failed("ERROR:  Failed to unmarshal the plan file", zap.Error(planErr))
		return nil, err
	}

	// get service resources
	serviceResources, err := getServiceResources()
	if err != nil {
		logger.Failed("ERROR:  Failed to convert terraform plan file to json : %v", err)
		return nil, err
	}

	// get drift resources
	for _, resource := range planIntf.ResourceChanges {
		for _, action := range resource.Change.Actions {
			if action == tfjson.ActionCreate {
				// check resource Type is part of imported service resources
				if isServiceResource(services, resource.Type, serviceResources) {
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
	r, _ := regexp.Compile(`^[A-Za-z-_]+\\.tf$`)

	// check is there any other resource which has got reference of the removing resource
	// if yes, Prepare the list and ask user to fix the reference issue before delete operation
	var toBeWritten *hclwrite.File
	files, err := ioutil.ReadDir(".")
	if err != nil {
		logger.Failed("ERROR:  %v", err)
		return err
	}

	// prepare resource reference list
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

				// validate resource is not part of drift resource
				if !isValueInList(driftResources, resTypeAndName, false) {
					// validate non-drift resource is dependent on drift resources
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
		return fmt.Errorf("failed to remove the drift resource '%v' reference from .tf file.", driftResources)
	}

	return nil
}

// isServiceResource ..
func isServiceResource(services []string, resourceType string, serviceMap map[string][]string) bool {
	if serviceMap != nil {
		for _, service := range services {
			if serviceMap[service] != nil {
				for _, rType := range serviceMap[service] {
					if resourceType == rType {
						return true
					}
				}
			}
		}
	}

	return false
}
