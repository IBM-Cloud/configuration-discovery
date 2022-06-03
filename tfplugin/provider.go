package tfplugin

import (
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/zclconf/go-cty/cty"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	tfplugin "github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/providers"
	"github.com/hashicorp/terraform/version"
)

// DefaultDataDir is the default directory for storing local data.
const DefaultDataDir = ".terraform"

// DefaultPluginVendorDir is the location in the config directory to look for
// user-added plugin binaries. Terraform only reads from this path if it
// exists, it is never created by terraform.
const DefaultPluginVendorDirV12 = "terraform.d/plugins/" + pluginMachineName

// pluginMachineName is the directory name used in new plugin paths.
const pluginMachineName = runtime.GOOS + "_" + runtime.GOARCH

type ProviderWrapper struct {
	Provider     *tfplugin.GRPCProvider
	client       *plugin.Client
	rpcClient    plugin.ClientProtocol
	providerName string
	config       cty.Value
	schema       *providers.GetSchemaResponse
	retryCount   int
	retrySleepMs int
}

func NewProviderWrapper(providerName string, providerConfig cty.Value, verbose bool, options ...map[string]int) (*ProviderWrapper, error) {
	p := &ProviderWrapper{retryCount: 5, retrySleepMs: 300}
	p.providerName = providerName
	p.config = providerConfig

	if len(options) > 0 {
		retryCount, hasOption := options[0]["retryCount"]
		if hasOption {
			p.retryCount = retryCount
		}
		retrySleepMs, hasOption := options[0]["retrySleepMs"]
		if hasOption {
			p.retrySleepMs = retrySleepMs
		}
	}

	err := p.initProvider(verbose)

	return p, err
}

func (p *ProviderWrapper) Kill() {
	p.client.Kill()
}

func (p *ProviderWrapper) GetSchema() *providers.GetSchemaResponse {
	if p.schema == nil {
		r := p.Provider.GetSchema()
		p.schema = &r
	}
	return p.schema
}

func (p *ProviderWrapper) initProvider(verbose bool) error {
	providerFilePath, err := getProviderFileName(p.providerName)
	if err != nil {
		return err
	}
	options := hclog.LoggerOptions{
		Name:   "plugin",
		Level:  hclog.Error,
		Output: os.Stdout,
	}
	if verbose {
		options.Level = hclog.Trace
	}
	logger := hclog.New(&options)
	p.client = plugin.NewClient(
		&plugin.ClientConfig{
			Cmd:              exec.Command(providerFilePath),
			HandshakeConfig:  tfplugin.Handshake,
			VersionedPlugins: tfplugin.VersionedPlugins,
			Managed:          true,
			Logger:           logger,
			AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
			AutoMTLS:         true,
		})
	p.rpcClient, err = p.client.Client()
	if err != nil {
		return err
	}
	raw, err := p.rpcClient.Dispense(tfplugin.ProviderPluginName)
	if err != nil {
		return err
	}

	p.Provider = raw.(*tfplugin.GRPCProvider)

	config, err := p.GetSchema().Provider.Block.CoerceValue(p.config)
	if err != nil {
		return err
	}
	p.Provider.Configure(providers.ConfigureRequest{
		TerraformVersion: version.Version,
		Config:           config,
	})

	return nil
}

func getProviderFileName(providerName string) (string, error) {
	defaultDataDir := os.Getenv("TF_DATA_DIR")
	if defaultDataDir == "" {
		defaultDataDir = DefaultDataDir
	}
	providerFilePath, err := getProviderFileNameV13andV14(defaultDataDir, providerName)
	if err != nil || providerFilePath == "" {
		providerFilePath, err = getProviderFileNameV13andV14(os.Getenv("HOME")+string(os.PathSeparator)+
			".terraform.d", providerName)
	}
	if err != nil || providerFilePath == "" {
		return getProviderFileNameV12(providerName)
	}
	return providerFilePath, nil
}

func getProviderFileNameV13andV14(prefix, providerName string) (string, error) {
	// Read terraform v14 file path
	registryDir := prefix + string(os.PathSeparator) + "providers" + string(os.PathSeparator) +
		"registry.terraform.io"
	providerDirs, err := ioutil.ReadDir(registryDir)
	if err != nil {
		// Read terraform v13 file path
		registryDir = prefix + string(os.PathSeparator) + "plugins" + string(os.PathSeparator) +
			"registry.terraform.io"
		providerDirs, err = ioutil.ReadDir(registryDir)
		if err != nil {
			return "", err
		}
	}
	providerFilePath := ""
	for _, providerDir := range providerDirs {
		pluginPath := registryDir + string(os.PathSeparator) + providerDir.Name() +
			string(os.PathSeparator) + providerName
		dirs, err := ioutil.ReadDir(pluginPath)
		if err != nil {
			continue
		}
		for _, dir := range dirs {
			if !dir.IsDir() {
				continue
			}
			for _, dir := range dirs {
				fullPluginPath := pluginPath + string(os.PathSeparator) + dir.Name() +
					string(os.PathSeparator) + runtime.GOOS + "_" + runtime.GOARCH
				files, err := ioutil.ReadDir(fullPluginPath)
				if err == nil {
					for _, file := range files {
						if strings.HasPrefix(file.Name(), "terraform-provider-"+providerName) {
							providerFilePath = fullPluginPath + string(os.PathSeparator) + file.Name()
						}
					}
				}
			}
		}
	}
	return providerFilePath, nil
}

func getProviderFileNameV12(providerName string) (string, error) {
	defaultDataDir := os.Getenv("TF_DATA_DIR")
	if defaultDataDir == "" {
		defaultDataDir = DefaultDataDir
	}
	pluginPath := defaultDataDir + string(os.PathSeparator) + "plugins" + string(os.PathSeparator) + runtime.GOOS + "_" + runtime.GOARCH
	files, err := ioutil.ReadDir(pluginPath)
	if err != nil {
		pluginPath = os.Getenv("HOME") + string(os.PathSeparator) + "." + DefaultPluginVendorDirV12
		files, err = ioutil.ReadDir(pluginPath)
		if err != nil {
			return "", err
		}
	}
	providerFilePath := ""
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if strings.HasPrefix(file.Name(), "terraform-provider-"+providerName) {
			providerFilePath = pluginPath + string(os.PathSeparator) + file.Name()
		}
	}
	return providerFilePath, nil
}
