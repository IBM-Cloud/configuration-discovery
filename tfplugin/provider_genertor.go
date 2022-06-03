package tfplugin

import (
	"github.com/zclconf/go-cty/cty"
)

type ProviderGenerator interface {
	GetName() string
	GetConfig() cty.Value
}

type Provider struct {
	Config cty.Value
}

func (p *Provider) GetConfig() cty.Value {
	return p.Config
}

func (p *Provider) GetName() string {
	panic("implement me")
}
