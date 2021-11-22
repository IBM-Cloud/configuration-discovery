package tfplugin

type IBMProvider struct { //nolint
	Provider
}

func (p *IBMProvider) Init(args []string) error {
	return nil
}

func (p *IBMProvider) GetName() string {
	return "ibm"
}

func NewIbmProvider() ProviderGenerator {
	return &IBMProvider{}
}

func Import(provider ProviderGenerator, args []string) (*ProviderWrapper, error) {
	providerWrapper, err := initOptionsAndWrapper(provider, args)
	if err != nil {
		return nil, err
	}

	defer providerWrapper.Kill()
	return providerWrapper, err
}

func initOptionsAndWrapper(provider ProviderGenerator, args []string) (*ProviderWrapper, error) {
	providerWrapper, err := NewProviderWrapper(provider.GetName(), provider.GetConfig(), false, map[string]int{"retryCount": 1, "retrySleepMs": 10})
	if err != nil {
		return nil, err
	}

	return providerWrapper, nil
}
