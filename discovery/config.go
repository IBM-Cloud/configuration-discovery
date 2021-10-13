package discovery

import (
	"context"
	"io/ioutil"
	"strings"
	"time"

	"github.com/IBM-Cloud/configuration-discovery/terraformwrapper"
	"github.com/IBM-Cloud/configuration-discovery/utils"
)

// DiscoveryImport ..
func DiscoveryImport(ctx context.Context, services, tags string, compact bool, randomID, discoveryDir string) error {
	logger := utils.GetLogger(ctx)
	logger.Say("# let's import the resources (%s):\n", services)
	// Import the terraform resources & state files.

	err := terraformwrapper.TerraformerImport(discoveryDir, services, tags, compact, planTimeOut, randomID)
	if err != nil {
		return err
	}

	logger.Say("# Terraformer generated TF & state file!")

	// validate terraform exported files
	isTerraform13, err := ValidateExportedFiles(discoveryDir)
	if err != nil {
		return err
	}

	if isTerraform13 {
		logger.Say("# let's do some IaC! First, update the provider.tf file to support TF 0.13+ :")
		err = UpdateProviderFile(ctx, discoveryDir, randomID, planTimeOut)
		if err != nil {
			return err
		}
	}

	//Run terraform init commnd
	logger.Say("# let us do terrafrom initialization on exported terraform files:")
	err = terraformwrapper.TerraformInit(discoveryDir, planTimeOut, randomID)
	if err != nil {
		return err
	}

	//Run terraform refresh commnd on the generated state file
	logger.Say("# and finally compare what we imported with what we currently have in real:")
	err = terraformwrapper.TerraformRefresh(discoveryDir, planTimeOut, randomID)
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
	err = terraformwrapper.TerraformReplaceProvider(discoveryDir, randomID, planTimeOut)
	if err != nil {
		return err
	}

	return nil
}
