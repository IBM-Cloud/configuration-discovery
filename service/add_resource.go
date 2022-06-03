package service

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/IBM-Cloud/configuration-discovery/tfplugin"
	"github.com/IBM-Cloud/configuration-discovery/utils"
	"github.com/tidwall/sjson"
)

// DiscoveryImport ..
//  // todo: opts []string is needed to be taken as arg
func DiscoveryImport(ctx context.Context, services, tags string, compact bool, randomID, discoveryDir string) error {
	logger := utils.GetLogger(ctx)
	logger.Say("INFO:  let's import the resources (%s) 2/6:\n", services)

	// import the terraform resources & state files.
	err := tfplugin.TerraformerImport(discoveryDir, services, tags, compact, planTimeOut, randomID)
	if err != nil {
		logger.Failed("ERROR:  %v", err)
		return err
	}

	logger.Say("INFO:  Writing HCL Done!")
	logger.Say("INFO:  Writing TFState Done!")

	// check terraform version compatible
	logger.Say("INFO:  now, we can do some infra as code ! First, update the IBM Terraform provider to support TF 0.13 [3/6]:")
	err = UpdateProviderFile(ctx, discoveryDir, randomID, planTimeOut)
	if err != nil {
		logger.Failed("ERROR:  %v", err)
		return err
	}

	// run terraform init commnd
	logger.Say("INFO:  we need to init our Terraform project [4/6]:")
	err = tfplugin.TerraformInit(discoveryDir, planTimeOut, randomID)
	if err != nil {
		logger.Failed("ERROR:  %v", err)
		return err
	}

	// run terraform refresh commnd on the generated state file
	logger.Say("INFO:  and finally compare what we imported with what we currently have [5/6]:")
	err = tfplugin.TerraformRefresh(discoveryDir, planTimeOut, randomID)
	if err != nil {
		logger.Failed("ERROR:  %v", err)
		return err
	}

	return nil
}

// UpdateProviderFile ..
func UpdateProviderFile(ctx context.Context, discoveryDir, randomID string, timeout time.Duration) error {
	logger := utils.GetLogger(ctx)
	providerTF := discoveryDir + "/provider.tf"
	input, err := ioutil.ReadFile(providerTF)
	if err != nil {
		logger.Failed("ERROR:  %v", err)
		return err
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "version") {
			lines[i] = "source = \"IBM-Cloud/ibm\""
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(providerTF, []byte(output), 0644)
	if err != nil {
		logger.Failed("ERROR:  %v", err)
		return err
	}

	// replace provider path in state file
	err = tfplugin.TerraformReplaceProvider(discoveryDir, randomID, planTimeOut)
	if err != nil {
		logger.Failed("ERROR:  %v", err)
		return err
	}

	return nil
}

// MergeResources ..
func MergeResources(ctx context.Context, terraformerStateFile, terraformStateFile, discoveryDir, localDir, randomID string, timeout time.Duration) error {
	logger := utils.GetLogger(ctx)
	logger.Say("INFO:  Merge local tf/state file with discovery generated tf/state file!!")

	// initialize
	var addResourceList []tfplugin.Resource
	hclConf := []byte("\n")
	provider := tfplugin.NewIbmProvider()
	providerWrapper, err := tfplugin.Import(provider, []string{})
	if err != nil {
		log.Fatalln("ERROR: Could not create IBM Cloud provider schema object:", err)
	}

	// get resource attributes
	resourceAtrrs, err := getResourceIgnoredAttributes()
	if err != nil {
		return fmt.Errorf("ERROR: Failed to parse resource attributes yaml file.")
	}

	// read local repo state file content
	terraformStateFileData, err := ioutil.ReadFile(terraformerStateFile)
	if err != nil {
		logger.Failed("ERROR:  %v", err)
		return err
	}
	statefilecontent := string(terraformStateFileData)

	// read terraform state file from local & discovery repo directory
	localRepoMap := ReadTerraformStateFile13(ctx, terraformStateFile, "")
	discoveryRepoMap := ReadTerraformStateFile13(ctx, terraformerStateFile, "discovery")

	// loop through each discovery repo resource with local repo resource
	for _, discovery := range discoveryRepoMap {
		// discovery resource
		discoveryResource := discovery.(tfplugin.Resource)

		// check discovery resource exist in local repo.
		// if resource not exist, Move the discovery resource to local repo
		if localRepoMap[discoveryResource.TypeAndID] == nil {
			resource := tfplugin.Resource{
				Type:        discoveryResource.Type,
				Name:        discoveryResource.Name,
				TypeAndName: discoveryResource.TypeAndName,
				Attributes:  RemoveComputedAttributes(discoveryResource, resourceAtrrs, providerWrapper),
			}

			// check discovery resource has got depends_on attribute
			// if depends_on attribute exist in discovery resource, Get the depends_on resource name from local repo & update in discovery state file.
			if discoveryResource.DependsOn != nil {
				var dependsOn []string

				for i, d := range discoveryResource.DependsOn {
					localParentResource := discoveryRepoMap[d].(tfplugin.Resource).TypeAndID

					// check dependent resource exist in local repo
					if localRepoMap[localParentResource] != nil {
						// get depends_on value from local repo resource to update in discovery state file
						// if dependent resource is from module, set module name value to depends_on attribute
						if len(localRepoMap[localParentResource].(tfplugin.Resource).Module) > 0 {
							dependsOn = append(dependsOn, localRepoMap[localParentResource].(tfplugin.Resource).Module)
						} else {
							localParentResource = localRepoMap[localParentResource].(tfplugin.Resource).TypeAndName
							dependsOn = append(dependsOn, localParentResource)
						}
					} else {
						// if deendent resource not exist in local repo, set depends_on value of the discovery resource
						localParentResource = discoveryRepoMap[d].(tfplugin.Resource).TypeAndName
						dependsOn = append(dependsOn, discoveryRepoMap[d].(tfplugin.Resource).TypeAndName)
					}

					// update depends_on parameter in discovery state file content
					statefilecontent, err = sjson.Set(statefilecontent, "resources."+strconv.Itoa(discovery.(tfplugin.Resource).Index)+".instances.0.dependencies."+strconv.Itoa(i), localParentResource)
					if err != nil {
						logger.Failed("ERROR:  %v", err)
						return err
					}
				}

				// set dependsOn to resource
				if len(dependsOn) > 0 {
					resource.DependsOn = dependsOn
				}
			}

			addResourceList = append(addResourceList, resource)
		}
	}

	// copy the state file content changes to discovery repo state file
	if len(statefilecontent) > 0 {
		err = ioutil.WriteFile(terraformerStateFile, []byte(statefilecontent), 0644)
		if err != nil {
			logger.Failed("ERROR:  %v", err)
			return err
		}
	}

	// move resources from discovery repo to local repo
	if len(addResourceList) > 0 {
		for _, resource := range addResourceList {
			err = tfplugin.TerraformMoveResource(discoveryDir, terraformerStateFile, terraformStateFile, resource.TypeAndName, planTimeOut, randomID)
			if err != nil {
				logger.Failed("ERROR: Error in moving resource from state file : %v", err)
				return err
			}
		}

		// print HCL
		providerData := map[string]interface{}{}
		data, err := tfplugin.HclPrintResource(addResourceList, providerData, "hcl")
		if err != nil {
			logger.Failed("ERROR: Error in creating HCL resource ::", err)
		}

		hclConf = append(hclConf, string(data)...)
		tfplugin.PrintHcl(hclConf, localDir+"/conf_service.tf")
		logger.Say("INFO: Discovery service successfuly moved (%v) resources from (%s) to (%s).", len(addResourceList), discoveryDir, localDir)
	} else {
		logger.Say("INFO: Discovery service didn't find any resource to move from (%s) to (%s).", discoveryDir, localDir)
	}

	return nil
}

// Remove computed/optional parameter from resource
func RemoveComputedAttributes(resource tfplugin.Resource, resourceIgnoredAtrrs map[string]interface{}, providerWrapper *tfplugin.ProviderWrapper) map[string]interface{} {
	// get resource attributes from provider schema
	resourceSchema := providerWrapper.GetSchema().ResourceTypes[resource.Type]
	for k, attr := range resourceSchema.Block.Attributes {
		if !attr.Optional && !attr.Required {
			delete(resource.Attributes, k)
			continue
		}

		// remove ignored attributes
		if resourceIgnoredAtrrs[resource.Type] != nil {
			ignoreAttrs := resourceIgnoredAtrrs[resource.Type].(map[interface{}]interface{})
			if ignoreAttrs["ignore_keys"] != nil {
				for _, v := range ignoreAttrs["ignore_keys"].([]interface{}) {
					delete(resource.Attributes, v.(string))
				}
			}
		}
	}

	// remove computed attributes from resource
	for key, value := range resource.Attributes {
		switch t := value.(type) {
		case interface{}:
			v := reflect.ValueOf(t)
			if v.Kind() == reflect.Float32 || v.Kind() == reflect.Float64 {
				continue
			}
			if v.Kind() != reflect.Bool && v.Len() == 0 {
				delete(resource.Attributes, key)
			}

			if key == "timeouts" {
				delete(resource.Attributes, key)
			}

		default:
			if value == nil {
				delete(resource.Attributes, key)
			}
		}
	}
	delete(resource.Attributes, "id")

	return resource.Attributes
}
