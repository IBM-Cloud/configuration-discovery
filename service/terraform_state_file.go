package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/IBM-Cloud/configuration-discovery/tfplugin"
	"github.com/IBM-Cloud/configuration-discovery/utils"
)

// TerraformSate ..
type TerraformState struct {
	Resources []Resources `json:"resources"`
	Modules   []Modules   `json:"modules"`
}

// Resources ..
type Resources struct {
	Instances []Instances `json:"instances"`
	Mode      string      `json:"mode"`
	Module    string      `json:"module"`
	Type      string      `json:"type"`
	Name      string      `json:"name"`
}

// Instances ..
type Instances struct {
	Mode       string                 `json:"mode"`
	Attributes map[string]interface{} `json:"attributes"`
	DependsOn  []string               `json:"dependencies"`
}

// ResourceList ..
type ResourceList []tfplugin.Resource

// Modules ..
type Modules struct {
	Resources map[string]ResourceBody `json:"resources"`
}

// ResourceBody .
type ResourceBody struct {
	Primary      map[string]PrimaryBody `json:"primary"`
	Provider     string                 `json:"provider"`
	ResourceType string                 `json:"type"`
}

// PrimaryBody .
type PrimaryBody struct {
	ID            string   `json:"id"`
	Location      string   `json:"location"`
	AttributeData []string `json:"attributes"`
}

// ReadTerraformerStateFile ..
// TF 0.12 compatible
func ReadTerraformerStateFile(ctx context.Context, terraformerStateFile string) ResourceList {
	var rList ResourceList
	tfStateFile := TerraformState{}

	logger := utils.GetLogger(ctx)

	tfFile, err := ioutil.ReadFile(terraformerStateFile)
	if err != nil {
		logger.Failed("ERROR:  %v", err)
		os.Exit(1)
	}

	err = json.Unmarshal([]byte(tfFile), &tfStateFile)
	if err != nil {
		logger.Failed("ERROR:  %v", err)
		os.Exit(1)
	}

	for i := 0; i < len(tfStateFile.Modules); i++ {
		resource := tfplugin.Resource{}
		for k := range tfStateFile.Modules[i].Resources {
			resource.Name = k
			resource.Type = tfStateFile.Modules[i].Resources[k].ResourceType
			for p := range tfStateFile.Modules[i].Resources[k].Primary {
				if p == "attributes" {
					resource.ID = tfStateFile.Modules[i].Resources[k].Primary[p].ID
				}
			}
			rList = append(rList, resource)
		}
	}

	logger.Say("INFO: Total (%d) resource in (%s).\n", len(rList), terraformerStateFile)
	return rList
}

// ReadTerraformStateFile ..
// TF 0.13+ compatible
func ReadTerraformStateFile13(ctx context.Context, terraformStateFile, repoType string) map[string]interface{} {
	resources := make(map[string]interface{})
	tfStateFile := TerraformState{}
	logger := utils.GetLogger(ctx)

	tfFile, err := ioutil.ReadFile(terraformStateFile)
	if err != nil {
		logger.Failed("ERROR:  %v", err)
		os.Exit(1)
	}

	err = json.Unmarshal([]byte(tfFile), &tfStateFile)
	if err != nil {
		logger.Failed("ERROR:  %v", err)
		os.Exit(1)
	}

	for i := 0; i < len(tfStateFile.Resources); i++ {
		tfResource := tfStateFile.Resources[i]

		//Don't process the mode type with 'data' value
		if tfStateFile.Resources[i].Mode == "data" {
			continue
		}

		newResource := tfplugin.Resource{
			Type:        tfResource.Type,
			Name:        tfResource.Name,
			TypeAndName: tfResource.Type + "." + tfResource.Name,
			Attributes:  tfResource.Instances[0].Attributes,
			Module:      tfResource.Module,
		}

		for k := 0; k < len(tfResource.Instances); k++ {
			var key string
			var dependsOn []string

			resourceId := fmt.Sprintf("%v", tfResource.Instances[k].Attributes["id"])
			newResource.ID = bytes.NewBuffer([]byte(resourceId)).String()

			if tfResource.Instances[k].DependsOn != nil {
				for _, dependOn := range tfResource.Instances[k].DependsOn {
					if strings.HasPrefix("module.", dependOn) {
						result := strings.Split(dependOn, ".")
						dependsOn = append(dependsOn, fmt.Sprintf("%s.%s", result[0], result[1]))
					} else {
						dependsOn = append(dependsOn, dependOn)
					}
				}
				newResource.DependsOn = dependsOn
			}

			if repoType == "discovery" {
				key = newResource.TypeAndName
				newResource.TypeAndID = newResource.Type + "." + newResource.ID
			} else {
				key = newResource.Type + "." + newResource.ID
				newResource.TypeAndID = key
			}

			newResource.Index = i
			resources[key] = newResource
		}
	}

	logger.Say("INFO: Total (%d) resource in (%s).\n", len(resources), terraformStateFile)
	return resources
}
