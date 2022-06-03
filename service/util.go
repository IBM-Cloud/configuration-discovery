package service

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/IBM-Cloud/configuration-discovery/utils"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

// Find takes a slice and looks for an element in it. If found it will
// return it's key, otherwise it will return -1 and a bool of false.
func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

//Remove element from list
func removeElementFromList(urlList, removeList []string) []cty.Value {
	for i := 0; i < len(urlList); i++ {
		url := urlList[i]
		for _, rem := range removeList {
			if url == rem {
				urlList = append(urlList[:i], urlList[i+1:]...)
				i-- // Important: decrease index
				continue
			}
		}
	}

	retVals := make([]cty.Value, len(urlList))
	for i, s := range urlList {
		retVals[i] = cty.StringVal(s)
	}

	return retVals
}

// Iterate over an array of strings and check if a value is equal
func isValueInList(strList []string, value string, subStr bool) bool {
	for _, s := range strList {
		if subStr == true {
			if strings.Contains(value, s) {
				return true
			}
		} else {
			if s == value {
				return true
			}
		}

	}
	return false
}

//function to check if a de-provisioned resource present in the attributes
func checkMapContainsValue(resourceList []string, attributes map[string]*hclwrite.Attribute) bool {

	//traverse through the resource attributes
	for _, value := range attributes {
		paramValue := string(hclwrite.Format(value.Expr().BuildTokens(nil).Bytes()))
		//check if attribute value is equals to de-provisioned resource
		for _, res := range resourceList {
			r, _ := regexp.Compile("\\b" + res + "\\b(.*)$")
			if r.MatchString(paramValue) {
				return true
			}
		}
	}

	//if value not found return false
	return false
}

//It will clone the git repo which contains the configuration file.
func CloneRepo(msg ConfigRequest) ([]byte, string, error) {
	var stdouterr []byte
	gitURL := msg.GitURL
	urlPath, err := url.Parse(msg.GitURL)
	if err != nil {
		return nil, "", err
	}
	baseName := filepath.Base(urlPath.Path)
	extName := filepath.Ext(urlPath.Path)
	p := baseName[:len(baseName)-len(extName)]
	if _, err := os.Stat(currentDir + utils.PathSep + p); err == nil {
		stdouterr, err = pullRepo(p)

	} else {
		cmd := exec.Command("git", "clone", gitURL)
		fmt.Println(cmd.Args)
		cmd.Dir = currentDir
		stdouterr, err = cmd.CombinedOutput()
		if err != nil {
			return nil, "", err
		}
	}
	path := currentDir + utils.PathSep + p + utils.PathSep + "terraform.tfvars"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		createFile(msg, path)
	} else {
		err = os.Remove(path)
		createFile(msg, path)
	}

	return stdouterr, p, err
}

//It will create a vars file
func createFile(msg ConfigRequest, path string) {
	// detect if file exists

	_, err := os.Stat(path)

	// create file if not exists
	if os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			return
		}
		defer file.Close()
	}

	writeFile(path, msg)
}

func pullRepo(repoName string) ([]byte, error) {
	cmd := exec.Command("git", "pull")
	fmt.Println(cmd.Args)
	cmd.Dir = currentDir + utils.PathSep + repoName
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	return stdoutStderr, err
}

func removeRepo(path, repoName string) error {
	removePath := filepath.Join(path, repoName)
	err := os.RemoveAll(removePath)
	return err
}

func writeFile(path string, msg ConfigRequest) {
	// open file using READ & WRITE permission
	var file, err = os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return
	}
	defer file.Close()

	variables := msg.VariableStore

	if variables != nil {
		for _, v := range *variables {
			_, err = file.WriteString(v.Name + " = \"" + v.Value + "\" \n")
			if err != nil {
				return
			}
		}
	}

	// save changes
	err = file.Sync()
	if err != nil {
		return
	}
}

func ValidateExportedFiles(discoveryDir string) (bool, error) {
	// check for terraform file existence
	isTfExit, err := utils.IsFileExists(discoveryDir + "/outputs.tf")
	if !isTfExit {
		return false, fmt.Errorf("ERROR: Error in importing resources.")
	}

	// Check terraform version compatible
	reTerraformversion, _ := regexp.Compile("Terraform v(.*)[\\s]")
	cmd := exec.Command("terraform", "version")
	sysOutput, err := cmd.Output()
	if err != nil {
		log.Println("ERROR: Error running command :", err)
	}

	tfVersion := string(sysOutput)
	results := reTerraformversion.FindStringSubmatch(tfVersion)
	v1, err := version.NewVersion("0.12.31")
	v2, err := version.NewVersion(results[1])
	if v2.LessThanOrEqual(v1) {
		return false, nil
	}

	return true, nil
}

func ReadFile(ctx context.Context, filepath string) ([]byte, error) {

	// Open our jsonFile
	jsonFile, err := os.Open(filepath)
	// if we os.Open returns an error then handle it
	if err != nil {
		log.Println("ERROR: failed to open the file", zap.Error(err))
		return nil, err
	}

	byteValue, readError := ioutil.ReadAll(bufio.NewReader(jsonFile))
	if readError != nil {
		log.Println("ERROR: failed to read the file content", zap.Error(readError))
		return nil, readError
	}

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'fileContent' which we defined above
	//json.Unmarshal(byteValue, &fileContent)
	//base.LogExit(methodName, logger)
	return byteValue, nil
}

// CleanUpDiscoveryFiles ..
func CleanUpDiscoveryFiles(localTFDir, discoveryDir string) error {

	//Remove discovery generated folder
	err := utils.RemoveDir(discoveryDir)
	if err != nil {
		log.Println("ERROR: Failed to remove the discovery generated folder.", zap.Error(err))
		return err
	}

	//Remove terraform.tfstate file genearted during merge
	files, err := filepath.Glob(localTFDir + "/terraform.tfstate.*.backup")
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		if err := os.Remove(f); err != nil {
			log.Println("ERROR: Failed to remove the discovery generated state files.", zap.Error(err))
			return err
		}
	}

	//Remove terraform.tfstate file genearted during merge
	files, err = filepath.Glob(localTFDir + "/tfplan.*")
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		if err := os.Remove(f); err != nil {
			log.Println("ERROR: Failed to remove the tfplan files.", zap.Error(err))
			return err
		}
	}

	return nil
}

func getServiceResources() (map[string][]string, error) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	resourceFile, err := ioutil.ReadFile(basepath + utils.PathSep + "conf/services.yml")
	if err != nil {
		return nil, err
	}

	serviceMap := make(map[string][]string)
	err = yaml.Unmarshal(resourceFile, serviceMap)
	if err != nil {
		panic(err)
	}

	return serviceMap, nil
}

func getResourceIgnoredAttributes() (map[string]interface{}, error) {
	resourceAtrrs := map[string]interface{}{}

	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	content, err := ioutil.ReadFile(basepath + utils.PathSep + "conf/resource_attributes.yml")
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}

	err = yaml.Unmarshal([]byte(content), &resourceAtrrs)
	if err != nil {
		log.Fatal(err)
	}

	return resourceAtrrs, nil
}
