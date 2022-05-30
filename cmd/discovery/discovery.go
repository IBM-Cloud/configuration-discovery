package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/IBM-Cloud/configuration-discovery/service"
	"github.com/IBM-Cloud/configuration-discovery/utils"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
)

const (
	releasesLink = "https://github.com/anilkumarnagaraj/terraform-provider-ibm-api/releases"
)

var (
	ui          terminal.UI
	planTimeOut = 60 * time.Minute
	confDir     string
	goctx       context.Context
)

func init() {
	// Bluemix terminal UI
	ui = terminal.NewStdUI()
	terminal.InitColorSupport()

	goctx = context.WithValue(context.Background(), utils.ContextKeyLogger, ui)

	confDir = os.Getenv("DISCOVERY_CONFIG_DIR")
	if confDir == "" {
		var err error
		confDir, err = os.Getwd()
		if err != nil {
			ui.Failed("Couldn't get DISCOVERY_CONFIG_DIR %v", err)
		}
	}
}

// todo - can users set this directly
// if msg.LOGLEVEL != "" {
// 	os.Setenv("TF_LOG", msg.LOGLEVEL)
// } user can set this directly

func main() {

	app := cli.NewApp()
	app.Name = "IBM Cloud Discovery CLI"
	app.HelpName = "discovery"
	app.Usage = "Lets you create state file and TF Config from Resources in your cloud account. " +
		"For the green field and brown field imports of config and statefile, " +
		"and all terraformer related"
	app.Writer = ui.Writer()
	app.ErrWriter = ui.Writer()
	app.Version = cliBuild

	// we create our commands
	app.Commands = []cli.Command{
		{
			Category:    "discovery",
			Name:        "version",
			Description: "Version",
			Usage: fmt.Sprint(
				"discovery",
				" version",
			),
			Action: actForVersion,
			OnUsageError: func(ctx *cli.Context, err error, isSub bool) error {
				return cli.ShowCommandHelp(ctx, ctx.Args().First())
			},
		},
		{
			Category:    "discovery",
			Name:        "metronome",
			Description: "Tick every 2 seconds for given time in mins",
			Usage: fmt.Sprint(
				"discovery",
				" metronome",
				" [--at AT]",
				" [--for FOR]",
			),
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "for",
					Usage: "for how much time in minutes",
					Value: 2,
				},
				cli.IntFlag{
					Name:  "at",
					Usage: "tick at which interval in seconds",
					Value: 2,
				},
			},
			Action: func(c *cli.Context) error {

				start_time := time.Now()

				secs := c.Int("at")
				mins := c.Int("for")
				for time.Since(start_time) < time.Duration(mins)*time.Minute {
					ui.Print("Ticking at ", time.Now())
					time.Sleep(time.Duration(secs) * time.Second)
				}
				return nil
			},
			OnUsageError: func(ctx *cli.Context, err error, isSub bool) error {
				return cli.ShowCommandHelp(ctx, ctx.Args().First())
			},
		},
		{
			Category: "discovery",
			Name:     "config",
			Aliases:  []string{"configure"},
			Usage: fmt.Sprint(
				"discovery",
				" config",
				// " [--git_url GIT_URL]",
				" [--config_name CONFIG_NAME]",
			),
			// Description: "Clone and create a local configuration in an empty repo and run terraform " +
			// 	"init. Clones in to a directory (printed with name repo_name) " +
			// 	"If git_url is not passed, terraformer and terraform will be " +
			// 	"installed and terraform init is done in the config_dir. " +
			// 	"Set TF_LOG like you set for terraform for debug logs in your env ",
			Description: "Create a local configuration directory for importing the terrraform configuration. " +
				"config_dir is read from the env variable DISCOVERY_CONFIG_DIR. " +
				" If not set, current folder will be config_dir",
			Flags: []cli.Flag{
				// cli.StringFlag{
				// 	Name:  "git_url",
				// 	Usage: "The git url to get the configuration from. If empty, config_dir should have tf files.",
				// },
				cli.StringFlag{
					Name: "config_name",
					// Usage: "If git_url is passed, Must be an empty existing folder. A folder to operate in. " +
					// 	"If git_url is not passed, this folder should have tf files already. In this case, " +
					// 	"empty means current folder. Can be used to download terraformer and terraform.",
					Usage: "Name of the folder in config_dir, to which to download the terraform configuration. ",
					Value: "",
				},
			},

			Action: func(c *cli.Context) error {

				gitURL := c.String("git_url")
				configName := c.String("config_name")

				ui.Say("config dir: %s", confDir)
				if configName != "" {
					ui.Say("config folder name: %s", configName)
				}

				if gitURL != "" { // todo: @srikar - remove these in brownfield
					ui.Say("git url: %s", gitURL)
					ui.Failed("git_url not supported yet. Can clone once brownfield is supported")
				}

				var repoName string
				var err error
				if gitURL == "" {
					// ui.Say("EMPTY GIT URL: No git_url given, skipping to tf init")
				} else {

					ui.Say("Will clone git repo", gitURL)

					_, repoName, err = service.CloneRepo(service.ConfigRequest{
						GitURL: gitURL,
					})
					if err != nil {
						ui.Failed("Eror Cloning repo..err : %v\n", err)
						return err
					}
					ui.Say("\n config name: ", repoName)
				}

				if configName == "" {
					b := make([]byte, 10)
					rand.Read(b)
					randomID := fmt.Sprintf("%x", b)
					configName = "discovery" + randomID // todo: @srikar - change to time based
				}

				configRepoFolder, _ := utils.Filepathjoin(confDir, configName)

				if _, err := os.Stat(configRepoFolder); os.IsNotExist(err) {
					ui.Say("\ncreating configRepoFolder %s", configRepoFolder)
					err = os.MkdirAll(configRepoFolder, os.ModePerm)
					if err != nil {
						ui.Failed("Couldn't create %s, error: %v", err)
						return err
					}
				} else {
					isEmpty, err := utils.IsFolderEmpty(configRepoFolder)
					if err != nil {
						ui.Warn("Couldn't open dir %s, err: %v", configRepoFolder, err)
					}
					if !isEmpty {
						ui.Failed("Folder %s should be empty", configRepoFolder)
						return fmt.Errorf("config_name folder should be empty")
					}
					// ui.Say("\nRunning terraform init in %s", configRepoFolder)
					// err = utils.TerraformInit(configRepoFolder, &planTimeOut, "")
					// if err != nil {
					// 	ui.Failed("TF INIT ERROR: %v", err)
					// 	return err
					// }
				}
				ui.Ok()
				return nil
			},
			OnUsageError: func(ctx *cli.Context, err error, isSub bool) error {
				ui.Failed("ERROR: " + err.Error())
				return cli.ShowCommandHelp(ctx, ctx.Command.Name)
			},
		},
		{
			Category: "discovery",
			Name:     "import",
			Usage: fmt.Sprint(
				"discovery",
				" import",
				" --services SERVICES_TO_IMPORT", // ibm_is_instance
				" [--tags TAGS]",
				" [--config_name CONFIG_NAME]",
				" [--compact]",
				" [--merge]",
			),
			Description: "Import TF config for resources in your ibm cloud account. " +
				"Import all the resources for this service. Imports config and statefile. " +
				"If a statefile is already present, merging will be done. ",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "services",
					Usage: "The IBM service(s) to import the resources from. Comma separated. " +
						"'discovery version' to see all available services",
				},
				cli.StringFlag{
					Name: "config_name",
					Usage: "Folder inside config_dir, where to import the config. " +
						"A folder with prefix discovery is created inside the config_dir, if not given. " +
						"If this folder has some tf config already, merge flag has to be given. Imported" +
						"tf config will be merge to existing config then.",
					Value: "",
				},
				cli.StringFlag{
					Name:  "tags",
					Usage: "Tags in the format a:b,c:d",
				},
				cli.BoolFlag{
					Name: "compact",
					Usage: "Use --compact to generate all the terraform code into one single file. " +
						"If not passed, a file is created for each resource",
				},
				cli.BoolFlag{
					Name:  "merge",
					Usage: "Use --merge to import and merge with config/statefile in folder config_name ",
				},
			},

			Action: func(c *cli.Context) error {

				if !c.IsSet("services") {
					ui.Failed("services flag not set")
					return fmt.Errorf("services flag not set")
				}

				services := c.String("services")
				configName := c.String("config_name")
				isCompact := c.Bool("compact")
				isBrownField := c.Bool("merge")
				tags := c.String("tags")

				ui.Say("config_directory is %s", confDir)

				if key := os.Getenv("IC_API_KEY"); key == "" {
					ui.Warn("IC_API_KEY not exported")
				}

				b := make([]byte, 10) // todo: @srikar - decrease the lenghth
				rand.Read(b)
				randomID := fmt.Sprintf("%x", b) // todo: @srikar - change to time based
				var localTFDir string

				if configName == "" {
					if !isBrownField {
						configName = "discovery" + randomID
					}
				}

				greenFieldImportDir, _ := utils.Filepathjoin(confDir, configName)
				if !isBrownField {
					if _, err := os.Stat(greenFieldImportDir); os.IsNotExist(err) {
						ui.Say("\ncreating Folder %s for generating config", greenFieldImportDir)
						err = os.MkdirAll(greenFieldImportDir, os.ModePerm)
						if err != nil {
							ui.Failed("Couldn't create %s, error: %v", err)
							return err
						}
					} else {
						isEmpty, err := utils.IsFolderEmpty(greenFieldImportDir)
						if err != nil {
							ui.Warn("Couldn't open dir %s, err: %v", greenFieldImportDir, err)
						}
						if !isEmpty && !isBrownField {
							ui.Failed("Folder %s should be empty", greenFieldImportDir)
							return fmt.Errorf("config_name folder should be empty or pass --merge option to merge with existing config")
						}
					}
				} else {
					localTFDir = greenFieldImportDir
				}

				// todo: @srikar - where is opts used
				opts := []string{}
				var tagsChanged string

				if services != "" {
					opts = append(opts, "--resources="+services)
				}
				if tags != "" {
					ui.Say("Tags provided: %s", tags)
					splittedTags := strings.Split(tags, ",")
					ui.Say("Split tags: %v ", splittedTags)
					if len(splittedTags) > 0 {
						for _, v := range splittedTags {
							tag := strings.SplitN(v, ":", 2)
							if len(tag) == 2 {
								opts = append(opts, fmt.Sprintf("--%s=%s",
									strings.TrimSpace(strings.ToLower(tag[0])), tag[1]))
								tagsChanged += fmt.Sprintf("--%s=%s",
									strings.TrimSpace(strings.ToLower(tag[0])), tag[1])
							}
						}
					}
				}

				if isCompact {
					opts = append(opts, "--compact")
					ui.Say("terraformer cmd options passed: %v", opts)
				}

				ui.Say("Importing resources from ibm cloud")
				if !isBrownField {
					ui.Say("folder is greenFieldImportDir", greenFieldImportDir)
					err := service.DiscoveryImport(goctx, services, tagsChanged, isCompact, "", greenFieldImportDir)
					if err != nil {
						ui.Failed("Error in Importing resources: %v", err)
						return err
					}
				} else {
					ui.Say("folder is localTFDir", localTFDir)
					//Brown field scenario
					// todo:  - don't we support merging just the tf files
					if _, err := os.Stat(localTFDir + utils.PathSep + "terraform.tfstate"); os.IsNotExist(err) {
						ui.Say("No merging needed bcz terraform statefile doesn't already exist in local repo.", localTFDir)
						ui.Say("Done. Exiting")
						ui.Ok()
						return nil
					}

					//Create discovery directory
					importDir := localTFDir + utils.PathSep + "generated" + randomID
					if _, err := os.Stat(importDir); os.IsNotExist(err) { // todo: @srikar - always true
						ui.Say("\ncreating Folder %s for generating config", importDir)
						err = os.MkdirAll(importDir, os.ModePerm)
						if err != nil {
							ui.Failed("Couldn't create %s, error: %v", err)
							return err
						}
					} else {
						// change the randomid and try again or fail  // todo: @srikar -
					}

					serviceList := strings.Split(services, ",")
					if len(serviceList) == 0 {
						ui.Failed("No services were given") // todo: @srikar - what is this error @anil
						return fmt.Errorf("no services were given")
					}

					// Import the terraform resources & state files.
					err := service.DiscoveryImport(goctx, services, tagsChanged, isCompact, "", importDir)
					if err != nil {
						ui.Failed("Error with importing: %v", err)
						return err
					}

					ui.Say("Imported resources from ibm cloud at " + importDir)
					// if _, err := os.Stat(importDir); os.IsNotExist(err) {  // todo: @srikar - check with Anil
					if isEmpty, err := utils.IsFolderEmpty(importDir); err != nil && isEmpty {
						ui.Say("No configuration files after importing!!!")
						// todo: @srikar - exit here
					} else {
						ui.Say("Import successful. Imported into %s\n", importDir)
					}

					//Compare & merge local & remote .tf/state file
					terraformStateFile := localTFDir + utils.PathSep + "terraform.tfstate"
					terraformerStateFile := importDir + utils.PathSep + "terraform.tfstate"
					ui.Say("# Comparing and merging statefiles local %s and remote %s\n",
						terraformStateFile, terraformerStateFile)

					defer func() {
						//Cleanup the discovery generated files
						err = service.CleanUpDiscoveryFiles(localTFDir, importDir)
						if err != nil {
							ui.Warn("# there was a problem in cleaning up.", err)
						}
					}()

					//Delete drift resources from local .tf/state file
					err = service.DeleteDriftResources(goctx, serviceList, localTFDir, randomID, planTimeOut)
					if err != nil {
						ui.Warn("# Couldn't delete resource from .tf/state files", err)
						return err
					}

					//Merge resources from remote to local .tf/state file
					err = service.MergeResources(goctx, terraformerStateFile, terraformStateFile,
						importDir, localTFDir, randomID, planTimeOut)
					if err != nil {
						ui.Warn("# Couldn't merge .tf/state files", err)
						return err
					}

					ui.Say("Backend action: file state here finally", terraformStateFile)
				}

				ui.Say("Successful import")
				ui.Ok()

				return nil
			},
			OnUsageError: func(ctx *cli.Context, err error, isSub bool) error {
				ui.Failed("ERROR: " + err.Error())
				return cli.ShowCommandHelp(ctx, ctx.Command.Name)
			},
		},
	}

	// start our application
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
