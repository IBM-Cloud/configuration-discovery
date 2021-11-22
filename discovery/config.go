package discovery

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/IBM-Cloud/configuration-discovery/tfplugin"
	"github.com/IBM-Cloud/configuration-discovery/utils"
	"github.com/tidwall/sjson"
)

// ReadTerraformerStateFile ..
// TF 0.12 compatible
func ReadTerraformerStateFile(ctx context.Context, terraformerStateFile string) ResourceList {
	var rList ResourceList
	tfData := TerraformState{}

	logger := utils.GetLogger(ctx)

	tfFile, err := ioutil.ReadFile(terraformerStateFile)
	if err != nil {
		logger.Failed("Error: %v", err)
		os.Exit(1)
	}

	err = json.Unmarshal([]byte(tfFile), &tfData)
	if err != nil {
		logger.Failed("Error: %v", err)
		os.Exit(1)
	}

	for i := 0; i < len(tfData.Modules); i++ {
		rData := tfplugin.Resource{}
		for k := range tfData.Modules[i].Resources {
			rData.ResourceName = k
			rData.ResourceType = tfData.Modules[i].Resources[k].ResourceType
			for p := range tfData.Modules[i].Resources[k].Primary {
				if p == "attributes" {
					rData.ID = tfData.Modules[i].Resources[k].Primary[p].ID
				}
			}
			rList = append(rList, rData)
		}
	}

	logger.Say("Total (%d) resource in (%s).\n", len(rList), terraformerStateFile)
	return rList
}

// ReadTerraformStateFile ..
// TF 0.13+ compatible
func ReadTerraformStateFile(ctx context.Context, terraformStateFile, repoType string) map[string]interface{} {
	rIDs := make(map[string]interface{})
	tfData := TerraformState{}
	logger := utils.GetLogger(ctx)

	tfFile, err := ioutil.ReadFile(terraformStateFile)
	if err != nil {
		logger.Failed("Error: %v", err)
		os.Exit(1)
	}

	err = json.Unmarshal([]byte(tfFile), &tfData)
	if err != nil {
		logger.Failed("Error: %v", err)
		os.Exit(1)
	}

	for i := 0; i < len(tfData.Resources); i++ {
		rData := tfplugin.Resource{}
		var key string
		//Don't process the mode type with 'data' value
		if tfData.Resources[i].Mode == "data" {
			continue
		}

		rData.ResourceName = tfData.Resources[i].ResourceName
		rData.ResourceType = tfData.Resources[i].ResourceType
		rData.Attributes = tfData.Resources[i].Instances[0].Attributes
		for k := 0; k < len(tfData.Resources[i].Instances); k++ {
			resourceId := fmt.Sprintf("%v", tfData.Resources[i].Instances[k].Attributes["id"])
			rData.ID = bytes.NewBuffer([]byte(resourceId)).String()
			if tfData.Resources[i].Instances[k].DependsOn != nil {
				rData.DependsOn = tfData.Resources[i].Instances[k].DependsOn
			}

			if repoType == "discovery" {
				key = rData.ResourceType + "." + rData.ResourceName
			} else {
				key = rData.ResourceType + "." + rData.ID
			}
			rData.ResourceIndex = i
			rIDs[key] = rData
		}
	}

	logger.Say("Total (%d) resource in (%s).\n", len(rIDs), terraformStateFile)
	return rIDs
}

// DiscoveryImport ..
//  // todo: opts []string is needed to be taken as arg
func DiscoveryImport(ctx context.Context, services, tags string, compact bool, randomID, discoveryDir string) error {
	logger := utils.GetLogger(ctx)
	logger.Say("# let's import the resources (%s) 2/6:\n", services)
	// Import the terraform resources & state files.

	err := tfplugin.TerraformerImport(discoveryDir, services, tags, compact, planTimeOut, randomID)
	if err != nil {
		return err
	}

	logger.Say("# Writing HCL Done!")
	logger.Say("# Writing TFState Done!")

	//Check terraform version compatible
	logger.Say("# now, we can do some infra as code ! First, update the IBM Terraform provider to support TF 0.13 [3/6]:")
	err = UpdateProviderFile(ctx, discoveryDir, randomID, planTimeOut)
	if err != nil {
		return err
	}

	//Run terraform init commnd
	logger.Say("# we need to init our Terraform project [4/6]:")
	err = tfplugin.TerraformInit(discoveryDir, planTimeOut, randomID)
	if err != nil {
		return err
	}

	//Run terraform refresh commnd on the generated state file
	logger.Say("# and finally compare what we imported with what we currently have [5/6]:")
	err = tfplugin.TerraformRefresh(discoveryDir, planTimeOut, randomID)
	if err != nil {
		return err
	}

	return nil
}

// UpdateProviderFile ..
func UpdateProviderFile(ctx context.Context, discoveryDir, randomID string, timeout time.Duration) error {
	providerTF := discoveryDir + "/provider.tf"
	input, err := ioutil.ReadFile(providerTF)
	if err != nil {
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
		return err
	}

	//Replace provider path in state file
	err = tfplugin.TerraformReplaceProvider(discoveryDir, randomID, planTimeOut)
	if err != nil {
		return err
	}
	return nil
}

// MergeStateFile ..
func MergeStateFile(ctx context.Context, configRepoMap, discoveryRepoMap map[string]interface{}, src, dest, discoveryDir, configDir, randomID string, timeout time.Duration) error {
	provider := tfplugin.NewIbmProvider()
	providerWrapper, err := tfplugin.Import(provider, []string{})
	if err != nil {
		log.Fatalln("Could not create IBM Cloud provider schema object:", err)
	}

	//Read discovery state file
	terraformStateFileData, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	statefilecontent := string(terraformStateFileData)
	var addResourceList []tfplugin.Resource
	hclConf := []byte("\n")

	//Loop through each discovery repo resource with config repo resource
	for _, discoveryResource := range discoveryRepoMap {
		//Discovery resource
		discovery_resource := discoveryResource.(tfplugin.Resource).ResourceType + "." + discoveryResource.(tfplugin.Resource).ID

		//Check discovery resource exist in config repo.
		//If resource not exist, Move the discovery resource to config repo
		if configRepoMap[discovery_resource] == nil {
			resource := discoveryResource.(tfplugin.Resource)
			resource.ResourceName = discoveryResource.(tfplugin.Resource).ResourceName
			resource.ResourceTypeAndName = discoveryResource.(tfplugin.Resource).ResourceType + "." + discoveryResource.(tfplugin.Resource).ResourceName
			resource = RemoveComputedAttributes(resource, providerWrapper)

			//Check discovery resource has got depends_on attribute
			//If depends_on attribute exist in discovery resource, Get the depends_on resource name from config repo & update in discovery state file.
			if discoveryResource.(tfplugin.Resource).DependsOn != nil {
				var dependsOn []string

				for i, d := range discoveryResource.(tfplugin.Resource).DependsOn {
					configParentResource := discoveryRepoMap[d].(tfplugin.Resource).ResourceType + "." + discoveryRepoMap[d].(tfplugin.Resource).ID

					//Get parent resource from config repo
					if configRepoMap[configParentResource] != nil {
						//Get depends_on resource name from config repo to update in discovery state file
						configParentResource = configRepoMap[configParentResource].(tfplugin.Resource).ResourceType + "." + configRepoMap[configParentResource].(tfplugin.Resource).ResourceName
						dependsOn = append(dependsOn, configParentResource)

						//Update depends_on parameter in discovery state file content
						statefilecontent, err = sjson.Set(statefilecontent, "resources."+strconv.Itoa(discoveryResource.(tfplugin.Resource).ResourceIndex)+".instances.0.dependencies."+strconv.Itoa(i), configParentResource)
						if err != nil {
							return err
						}
					}
				}
				if len(dependsOn) > 0 {
					resource.DependsOn = dependsOn
				}
			}
			addResourceList = append(addResourceList, resource)
		}
	}

	//Copy the state file content changes to discovery repo state file
	if len(statefilecontent) > 0 {
		err = ioutil.WriteFile(src, []byte(statefilecontent), 0644)
		if err != nil {
			return err
		}
	}

	//Move resource from discovery repo to config repo
	if len(addResourceList) > 0 {
		for _, resource := range addResourceList {
			err = tfplugin.TerraformMoveResource(discoveryDir, src, dest, resource.ResourceTypeAndName, planTimeOut, randomID)
			if err != nil {
				return err
			}
		}

		//Print HCL
		providerData := map[string]interface{}{}
		data, err := tfplugin.HclPrintResource(addResourceList, providerData, "hcl")
		if err != nil {
			log.Println("Error in building resource ::", err)
		}

		hclConf = append(hclConf, string(data)...)
		tfplugin.PrintHcl(hclConf, configDir+"/conf_discovery.tf")
		log.Printf("\n\n# Discovery service successfuly moved (%v) resources from (%s) to (%s).", len(addResourceList), src, dest)
	} else {
		log.Printf("\n\n# Discovery service didn't find any resource to move from (%s) to (%s).", src, dest)
	}

	return nil
}

func RemoveComputedAttributes(resource tfplugin.Resource, providerWrapper *tfplugin.ProviderWrapper) tfplugin.Resource {
	//Get computed attributes
	readOnlyAttributes := []string{}
	obj := providerWrapper.GetSchema().ResourceTypes[resource.ResourceType]
	readOnlyAttributes = append(readOnlyAttributes, "id")
	for k, v := range obj.Block.Attributes {
		if !v.Optional && !v.Required {
			readOnlyAttributes = append(readOnlyAttributes, k)
		}
	}

	//Remove computed attributes
	for key, value := range resource.Attributes {
		switch t := value.(type) {
		case interface{}:
			v := reflect.ValueOf(t)
			if v.Kind() != reflect.Bool && (v.Len() == 0 || utils.Contains(readOnlyAttributes, key)) {
				delete(resource.Attributes, key)
			}
		default:
			if value == nil || utils.Contains(readOnlyAttributes, key) {
				delete(resource.Attributes, key)
			}
		}
	}

	return resource
}
