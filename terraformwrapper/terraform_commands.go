package terraformwrapper

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"
	"time"

	"github.com/IBM-Cloud/configuration-discovery/utils"
)

var logDir string

// HCL Config ..
type Resource struct {
	ID                  string
	ResourceIndex       int
	ResourceType        string
	ResourceName        string
	ResourceTypeAndName string
	ResourceTypeAndID   string
	DependsOn           []string               `json:",omitempty"`
	Attributes          map[string]interface{} `json:",omitempty"`
}

// TerraformInit ...
func TerraformInit(execDir string, timeout time.Duration, randomID string) error {

	return run(context.Background(), "terraform", []string{"init"}, execDir, timeout, randomID)
}

// TerraformApply ...
func TerraformApply(execDir, stateDir string, stateFileName string, timeout time.Duration, randomID string) error {
	return run(context.Background(), "terraform", []string{"apply", fmt.Sprintf("-state=%s", stateDir+utils.PathSep+stateFileName+".tfstate"), "-auto-approve"}, execDir, timeout, randomID)
}

// TerraformPlan ...
func TerraformPlan(execDir string, timeout time.Duration, randomID string) error {
	return run(context.Background(), "terraform", []string{"plan"}, execDir, timeout, randomID)
}

//TerraformRefresh ...
func TerraformRefresh(configDir string, timeout time.Duration, randomID string) error {
	return run(context.Background(), "terraform", []string{"refresh"}, configDir, timeout, randomID)
}

// TerraformDestroy ...
func TerraformDestroy(execDir, stateDir string, stateFileName string, timeout time.Duration, randomID string) error {

	return run(context.Background(), "terraform", []string{"destroy", "-force", fmt.Sprintf("-state=%s", stateDir+utils.PathSep+stateFileName+".tfstate")}, execDir, timeout, randomID)
}

// TerraformShow ...
func TerraformShow(execDir, stateDir string, stateFileName string, timeout time.Duration, randomID string) error {

	return run(context.Background(), "terraform", []string{"show", stateDir + utils.PathSep + stateFileName + ".tfstate"}, execDir, timeout, randomID)
}

//TerraformerImport ...
func TerraformerImport(configDir, resources, tags string, compact bool, timeout time.Duration, randomID string) error {

	if compact {

		return run(context.Background(), "terraformer", []string{"import", "ibm", fmt.Sprintf("--resources=%s", resources), tags, "--compact", fmt.Sprintf("-p=%s", configDir)}, configDir, timeout, randomID)
	} else {
		return run(context.Background(), "terraformer", []string{"import", "ibm", fmt.Sprintf("--resources=%s", resources), tags, fmt.Sprintf("-p=%s", configDir)}, configDir, timeout, randomID)
	}
}

//TerraformMoveResource ...
func TerraformMoveResource(configDir, srcStateFile, destStateFile, resourceName string, timeout time.Duration, randomID string) error {

	return run(context.Background(), "terraform", []string{"state", "mv", fmt.Sprintf("-state=%s", srcStateFile), fmt.Sprintf("-state-out=%s", destStateFile), resourceName, resourceName}, configDir, timeout, randomID)
}

//TerraformReplaceProvider ..
func TerraformReplaceProvider(configDir, randomID string, timeout time.Duration) error {
	//terraform state
	return run(context.Background(), "terraform", []string{"state", "replace-provider", "-auto-approve", "registry.terraform.io/-/ibm", "registry.terraform.io/ibm-cloud/ibm"}, configDir, timeout, randomID)
}

// TerraformVersion ...
func TerraformVersion(execDir string, timeout time.Duration, randomID string) error {

	return run(context.Background(), "terraform", []string{"verson"}, execDir, timeout, randomID)
}

// todo: @srikar - Make attribute function, remove too many func arguments
func run(ctx context.Context, cmdName string, args []string, execDir string, timeout time.Duration, randomID string) error {
	if timeout == 0 {
		timeout = 3 * time.Minute
	}

	ui := utils.GetLogger(ctx)

	ctx, cancel := context.WithTimeout(ctx, timeout)
	cmd := exec.CommandContext(ctx, cmdName, args...)
	defer cancel()

	cmd.Dir = execDir
	// cmd.Env = env // set any env needed
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if randomID != "" {
		stdoutFile, stderrFile, err := getLogFiles(logDir, randomID)
		if err != nil {
			return err
		}
		defer stdoutFile.Close()
		defer stderrFile.Close()

		//Write the stdout to log file
		go func() {
			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				fmt.Fprintln(stdoutFile, scanner.Text())
			}
		}()

		//Write the stderr to log file
		go func() {
			scanner := bufio.NewScanner(stderr)
			for scanner.Scan() {
				fmt.Fprintln(stderrFile, scanner.Text())

			}
		}()
	} else {
		go func() {
			stdioReader := bufio.NewReader(stdout)
			for {
				line, err := stdioReader.ReadString('\n')
				if err == nil || len(line) > 1 { // todo: @srikar - why is timestamp coming
					ui.Say(strings.TrimSpace(fmt.Sprintf("%s | %s", cmdName, line)))
				}
				if err != nil {
					break
				}
			}
		}()

		go func() {
			stderrReader := bufio.NewReader(stderr)
			for {
				line, err := stderrReader.ReadString('\n')
				if err == nil || len(line) > 1 {
					ui.Say(strings.TrimSpace(fmt.Sprintf("%s | %s", cmdName, line)))
				}
				if err != nil {
					break
				}
			}
		}()
	}

	// Start the command
	ui.Say("Starting command %s %v", cmd.Path, cmd.Args)
	err = cmd.Start()
	if err != nil {
		return err
	}

	//Wait for command to finish
	err = cmd.Wait()
	if err != nil {
		return err
	}
	return err
}

func getLogFiles(logDir, fileName string) (stdoutFile, stderrFile *os.File, err error) {
	stdoutPath := path.Join(logDir, fileName+".out")
	stderrPath := path.Join(logDir, fileName+".err")

	if _, err = os.Stat(stdoutPath); err == nil {
		stdoutFile, err = os.OpenFile(stdoutPath, os.O_APPEND|os.O_WRONLY, 0600)
	} else {
		stdoutFile, err = os.Create(stdoutPath)
	}
	if err != nil {
		return
	}

	if _, err = os.Stat(stderrPath); err == nil {
		stderrFile, err = os.OpenFile(stderrPath, os.O_APPEND|os.O_WRONLY, 0600)
	} else {
		stderrFile, err = os.Create(stderrPath)
	}
	return
}

func ReadLogFile(logID string) (stdout, stderr string, err error) {
	stdoutPath := path.Join(logDir, logID+".out")
	stderrPath := path.Join(logDir, logID+".err")

	outFile, err := ioutil.ReadFile(stdoutPath)
	if err != nil {
		return
	}
	errFile, err := ioutil.ReadFile(stderrPath)

	if err != nil {
		return
	}

	return string(outFile), string(errFile), nil
}
