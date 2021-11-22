package tfplugin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/hashicorp/hcl/hcl/ast"
	hclPrinter "github.com/hashicorp/hcl/hcl/printer"
	hclParser "github.com/hashicorp/hcl/json/parser"
	"github.com/hashicorp/terraform/terraform"
)

// Copy code from https://github.com/kubernetes/kops project with few changes for support many provider and heredoc

const safeChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"

var unsafeChars = regexp.MustCompile(`[^0-9A-Za-z_]`)

// sanitizer fixes up an invalid HCL AST, as produced by the HCL parser for JSON
type astSanitizer struct{}

// HCL Config ..
type Resource struct {
	ID                  string
	ResourceIndex       int
	ResourceType        string
	ResourceName        string
	ResourceTypeAndName string
	ResourceTypeAndID   string
	DependsOn           []string `json:",omitempty"`
	Provider            string
	Attributes          map[string]interface{}            `json:",omitempty"`
	Outputs             map[string]*terraform.OutputState `json:",omitempty"`
}

// output prints creates b printable HCL output and returns it.
func (v *astSanitizer) visit(n interface{}) error {
	if n == nil {
		return nil
	}
	switch t := n.(type) {
	case *ast.File:
		v.visit(t.Node)
	case *ast.ObjectList:
		var index int
		for {
			if index == len(t.Items) {
				break
			}
			v.visit(t.Items[index])
			index++
		}
	case *ast.ObjectKey:
	case *ast.ObjectItem:
		v.visitObjectItem(t)
	case *ast.LiteralType:
	case *ast.ListType:
	case *ast.ObjectType:
		v.visit(t.List)
	default:
		fmt.Printf(" unknown type: %T\n", n)
	}
	return nil
}

func (v *astSanitizer) visitObjectItem(o *ast.ObjectItem) {
	for i, k := range o.Keys {
		if i == 0 {
			text := k.Token.Text
			if text != "" && text[0] == '"' && text[len(text)-1] == '"' {
				v := text[1 : len(text)-1]
				safe := true
				for _, c := range v {
					if !strings.ContainsRune(safeChars, c) {
						safe = false
						break
					}
				}
				if safe {
					k.Token.Text = v
				}
			}
		}
	}

	// A hack so that Assign.IsValid is true, so that the printer will output =
	o.Assign.Line = 1

	v.visit(o.Val)
}

func Print(data interface{}, mapsObjects map[string]struct{}, format string) ([]byte, error) {
	return hclPrint(data, mapsObjects)
}

func hclPrint(data interface{}, mapsObjects map[string]struct{}) ([]byte, error) {
	dataJSONBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Println(string(dataJSONBytes))
		return []byte{}, fmt.Errorf("error marshalling terraform data to json: %v", err)
	}
	dataJSON := string(dataJSONBytes)
	nodes, err := hclParser.Parse([]byte(dataJSON))
	if err != nil {
		return []byte{}, fmt.Errorf("error parsing terraform json: %v", err)
	}
	var sanitizer astSanitizer
	sanitizer.visit(nodes)

	var b bytes.Buffer
	err = hclPrinter.Fprint(&b, nodes)
	if err != nil {
		return nil, fmt.Errorf("error writing HCL: %v", err)
	}
	s := b.String()

	// Remove extra whitespace...
	s = strings.ReplaceAll(s, "\n\n", "\n")

	// ...but leave whitespace between resources
	s = strings.ReplaceAll(s, "}\nresource", "}\n\nresource")

	// Apply Terraform style (alignment etc.)
	formatted, err := hclPrinter.Format([]byte(s))
	if err != nil {
		return nil, err
	}

	return formatted, nil
}

// Print hcl file from TerraformResource + provider
func HclPrintResource(resources []Resource, providerData map[string]interface{}, output string) ([]byte, error) {
	resourcesByType := map[string]map[string]interface{}{}
	mapsObjects := map[string]struct{}{}
	indexRe := regexp.MustCompile(`\.[0-9]+`)
	for _, res := range resources {
		r := resourcesByType[res.ResourceType]
		if r == nil {
			r = make(map[string]interface{})
			resourcesByType[res.ResourceType] = r
		}

		if r[res.ResourceName] != nil {
			log.Printf("[ERR]: duplicate resource found: %s.%s", res.ResourceType, res.ResourceName)
			continue
		}

		if len(res.DependsOn) > 0 {
			res.Attributes["depends_on"] = res.DependsOn
		}
		r[res.ResourceName] = res.Attributes

		for k := range res.Attributes {
			if strings.HasSuffix(k, ".%") {
				key := strings.TrimSuffix(k, ".%")
				mapsObjects[indexRe.ReplaceAllString(key, "")] = struct{}{}
			}
		}
	}

	data := map[string]interface{}{}
	if len(resourcesByType) > 0 {
		data["resource"] = resourcesByType
	}
	if len(providerData) > 0 {
		data["provider"] = providerData
	}
	var err error
	hclBytes, err := Print(data, mapsObjects, output)
	if err != nil {
		return []byte{}, err
	}
	return hclBytes, nil
}

func PrintHcl(hclData []byte, path string) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.Write([]byte(hclData)); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

}
