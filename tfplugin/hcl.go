package tfplugin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
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
	ID          string
	Index       int
	Type        string
	Name        string
	TypeAndName string
	TypeAndID   string
	DependsOn   []string `json:",omitempty"`
	Provider    string
	Attributes  map[string]interface{}            `json:",omitempty"`
	Outputs     map[string]*terraform.OutputState `json:",omitempty"`
	Module      string
}

// make HCL output reproducible by sorting the AST nodes
func sortHclTree(tree interface{}) {
	switch t := tree.(type) {
	case []*ast.ObjectItem:
		sort.Slice(t, func(i, j int) bool {
			var bI, bJ bytes.Buffer
			_, _ = hclPrinter.Fprint(&bI, t[i]), hclPrinter.Fprint(&bJ, t[j])
			return bI.String() < bJ.String()
		})
	case []ast.Node:
		sort.Slice(t, func(i, j int) bool {
			var bI, bJ bytes.Buffer
			_, _ = hclPrinter.Fprint(&bI, t[i]), hclPrinter.Fprint(&bJ, t[j])
			return bI.String() < bJ.String()
		})
	default:
	}
}

// output prints creates b printable HCL output and returns it.
func (v *astSanitizer) visit(n interface{}) {
	switch t := n.(type) {
	case *ast.File:
		v.visit(t.Node)
	case *ast.ObjectList:
		var index int
		sortHclTree(t.Items)
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
		sortHclTree(t.List)
	case *ast.ObjectType:
		sortHclTree(t.List)
		v.visit(t.List)
	default:
		fmt.Printf(" unknown type: %T\n", n)
	}
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
	switch t := o.Val.(type) {
	case *ast.LiteralType: // heredoc support
		if strings.HasPrefix(t.Token.Text, `"<<`) {
			t.Token.Text = t.Token.Text[1:]
			t.Token.Text = t.Token.Text[:len(t.Token.Text)-1]
			t.Token.Text = strings.ReplaceAll(t.Token.Text, `\n`, "\n")
			t.Token.Text = strings.ReplaceAll(t.Token.Text, `\t`, "")
			t.Token.Type = 10
			// check if text json for Unquote and Indent
			jsonTest := t.Token.Text
			lines := strings.Split(jsonTest, "\n")
			jsonTest = strings.Join(lines[1:len(lines)-1], "\n")
			jsonTest = strings.ReplaceAll(jsonTest, "\\\"", "\"")
			// it's json we convert to heredoc back
			var tmp interface{} = map[string]interface{}{}
			err := json.Unmarshal([]byte(jsonTest), &tmp)
			if err != nil {
				tmp = make([]interface{}, 0)
				err = json.Unmarshal([]byte(jsonTest), &tmp)
			}
			if err == nil {
				dataJSONBytes, err := json.MarshalIndent(tmp, "", "  ")
				if err == nil {
					jsonData := strings.Split(string(dataJSONBytes), "\n")
					// first line for heredoc
					jsonData = append([]string{lines[0]}, jsonData...)
					// last line for heredoc
					jsonData = append(jsonData, lines[len(lines)-1])
					hereDoc := strings.Join(jsonData, "\n")
					t.Token.Text = hereDoc
				}
			}
		}
	case *ast.ListType:
		sortHclTree(t.List)
	default:
	}

	// A hack so that Assign.IsValid is true, so that the printer will output =
	o.Assign.Line = 1

	v.visit(o.Val)
}

func Print(data interface{}, mapsObjects map[string]struct{}, format string) ([]byte, error) {
	return hclPrint(data, mapsObjects)
}

func hclPrint(data interface{}, mapsObjects map[string]struct{}) ([]byte, error) {
	dataBytesJSON, err := jsonPrint(data)
	if err != nil {
		return dataBytesJSON, err
	}
	dataJSON := string(dataBytesJSON)
	nodes, err := hclParser.Parse([]byte(dataJSON))
	if err != nil {
		log.Println(dataJSON)
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

	// hack for support terraform 0.12
	formatted = terraform12Adjustments(formatted, mapsObjects)
	// hack for support terraform 0.13
	formatted = terraform13Adjustments(formatted)
	if err != nil {
		log.Println("Invalid HCL follows:")
		for i, line := range strings.Split(s, "\n") {
			fmt.Printf("%4d|\t%s\n", i+1, line)
		}
		return nil, fmt.Errorf("error formatting HCL: %v", err)
	}

	return formatted, nil
}

// Print hcl file from TerraformResource + provider
func HclPrintResource(resources []Resource, providerData map[string]interface{}, output string) ([]byte, error) {
	resourcesByType := map[string]map[string]interface{}{}
	mapsObjects := map[string]struct{}{}
	indexRe := regexp.MustCompile(`\.[0-9]+`)
	for _, res := range resources {
		r := resourcesByType[res.Type]
		if r == nil {
			r = make(map[string]interface{})
			resourcesByType[res.Type] = r
		}

		if r[res.Name] != nil {
			log.Printf("[ERR]: duplicate resource found: %s.%s", res.Type, res.Name)
			continue
		}

		if len(res.DependsOn) > 0 {
			res.Attributes["depends_on"] = res.DependsOn
		}
		r[res.Name] = res.Attributes

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

func terraform13Adjustments(formatted []byte) []byte {
	s := string(formatted)
	requiredProvidersRe := regexp.MustCompile("required_providers \".*\" {")
	oldRequiredProviders := "\"required_providers\""
	newRequiredProviders := "required_providers"
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		if requiredProvidersRe.MatchString(line) {
			parts := strings.Split(strings.TrimSpace(line), " ")
			provider := strings.ReplaceAll(parts[1], "\"", "")
			lines[i] = "\t" + newRequiredProviders + " {"
			lines[i+1] = "\t\t" + provider + " = {\n\t" + lines[i+1] + "\n\t\t}"
		}
		lines[i] = strings.Replace(lines[i], oldRequiredProviders, newRequiredProviders, 1)
	}
	s = strings.Join(lines, "\n")
	return []byte(s)
}

func terraform12Adjustments(formatted []byte, mapsObjects map[string]struct{}) []byte {
	singletonListFix := regexp.MustCompile(`^\s*\w+ = {`)
	singletonListFixEnd := regexp.MustCompile(`^\s*}`)

	s := string(formatted)
	old := " = {"
	newEquals := " {"
	lines := strings.Split(s, "\n")
	prefix := make([]string, 0)
	for i, line := range lines {
		if singletonListFixEnd.MatchString(line) && len(prefix) > 0 {
			prefix = prefix[:len(prefix)-1]
			continue
		}
		if !singletonListFix.MatchString(line) {
			continue
		}
		key := strings.Trim(strings.Split(line, old)[0], " ")
		prefix = append(prefix, key)
		if _, exist := mapsObjects[strings.Join(prefix, ".")]; exist {
			continue
		}
		lines[i] = strings.ReplaceAll(line, old, newEquals)
	}
	s = strings.Join(lines, "\n")
	return []byte(s)
}

var OpeningBracketRegexp = regexp.MustCompile(`.?\\<`)
var ClosingBracketRegexp = regexp.MustCompile(`.?\\>`)

func jsonPrint(data interface{}) ([]byte, error) {
	dataJSONBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Println(string(dataJSONBytes))
		return []byte{}, fmt.Errorf("error marshalling terraform data to json: %v", err)
	}
	// We don't need to escape > or <
	s := strings.ReplaceAll(string(dataJSONBytes), "\\u003c", "<")
	s = OpeningBracketRegexp.ReplaceAllStringFunc(s, escapingBackslashReplacer("<"))
	s = strings.ReplaceAll(s, "\\u003e", ">")
	s = ClosingBracketRegexp.ReplaceAllStringFunc(s, escapingBackslashReplacer(">"))
	return []byte(s), nil
}

func escapingBackslashReplacer(backslashedCharacter string) func(string) string {
	return func(match string) string {
		if strings.HasPrefix(match, "\\\\") {
			return match // Don't replace regular backslashes
		}
		return strings.Replace(match, "\\"+backslashedCharacter, backslashedCharacter, 1)
	}
}
