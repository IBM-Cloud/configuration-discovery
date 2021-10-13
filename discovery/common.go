package discovery

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/IBM-Cloud/configuration-discovery/utils"
	"github.com/hashicorp/go-version"
)

var stdouterr []byte

//It will clone the git repo which contains the configuration file.
func CloneRepo(msg ConfigRequest) ([]byte, string, error) {
	gitURL := msg.GitURL
	urlPath, err := url.Parse(msg.GitURL)
	if err != nil {
		return nil, "", err
	}
	baseName := filepath.Base(urlPath.Path)
	extName := filepath.Ext(urlPath.Path)
	p := baseName[:len(baseName)-len(extName)]
	if _, err := os.Stat(currentDir + "/" + p); err == nil {
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
	path := currentDir + "/" + p + "/terraform.tfvars"
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
	cmd.Dir = currentDir + "/" + repoName
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
		return false, fmt.Errorf("Error in importing resources.")
	}

	// Check terraform version compatible
	reTerraformversion, _ := regexp.Compile("Terraform v(.*)[\\s]")
	cmd := exec.Command("terraform", "version")
	sysOutput, err := cmd.Output()
	if err != nil {
		log.Println("Error running command :", err)
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
