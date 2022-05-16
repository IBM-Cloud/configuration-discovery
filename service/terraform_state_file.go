package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

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
	tfData := TerraformState{}

	logger := utils.GetLogger(ctx)

	tfFile, err := ioutil.ReadFile(terraformerStateFile)
	if err != nil {
		logger.Failed("ERROR:  %v", err)
		os.Exit(1)
	}

	err = json.Unmarshal([]byte(tfFile), &tfData)
	if err != nil {
		logger.Failed("ERROR:  %v", err)
		os.Exit(1)
	}

	for i := 0; i < len(tfData.Modules); i++ {
		resource := tfplugin.Resource{}
		for k := range tfData.Modules[i].Resources {
			resource.Name = k
			resource.Type = tfData.Modules[i].Resources[k].ResourceType
			for p := range tfData.Modules[i].Resources[k].Primary {
				if p == "attributes" {
					resource.ID = tfData.Modules[i].Resources[k].Primary[p].ID
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
	tfData := TerraformState{}
	logger := utils.GetLogger(ctx)

	tfFile, err := ioutil.ReadFile(terraformStateFile)
	if err != nil {
		logger.Failed("ERROR:  %v", err)
		os.Exit(1)
	}

	err = json.Unmarshal([]byte(tfFile), &tfData)
	if err != nil {
		logger.Failed("ERROR:  %v", err)
		os.Exit(1)
	}

	for i := 0; i < len(tfData.Resources); i++ {
		//Don't process the mode type with 'data' value
		if tfData.Resources[i].Mode == "data" {
			continue
		}

		resource := tfplugin.Resource{
			Type:        tfData.Resources[i].Type,
			Name:        tfData.Resources[i].Name,
			TypeAndName: tfData.Resources[i].Type + "." + tfData.Resources[i].Name,
			Attributes:  tfData.Resources[i].Instances[0].Attributes,
		}

		for k := 0; k < len(tfData.Resources[i].Instances); k++ {
			var key string
			resourceId := fmt.Sprintf("%v", tfData.Resources[i].Instances[k].Attributes["id"])
			resource.ID = bytes.NewBuffer([]byte(resourceId)).String()
			if tfData.Resources[i].Instances[k].DependsOn != nil {
				resource.DependsOn = tfData.Resources[i].Instances[k].DependsOn
			}

			if repoType == "discovery" {
				key = resource.TypeAndName
				resource.TypeAndID = resource.Type + "." + resource.ID
			} else {
				key = resource.Type + "." + resource.ID
				resource.TypeAndID = key
			}
			resource.Index = i
			resources[key] = resource
		}
	}

	logger.Say("INFO: Total (%d) resource in (%s).\n", len(resources), terraformStateFile)
	return resources
}
