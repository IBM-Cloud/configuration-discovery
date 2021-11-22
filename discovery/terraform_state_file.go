package discovery

import "github.com/IBM-Cloud/configuration-discovery/tfplugin"

// TerraformSate ..
type TerraformState struct {
	Resources []Resources `json:"resources"`
	Modules   []Modules   `json:"modules"`
}

// Resources ..
type Resources struct {
	Instances    []Instances `json:"instances"`
	Mode         string      `json:"mode"`
	ResourceType string      `json:"type"`
	ResourceName string      `json:"name"`
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
