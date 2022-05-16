#!/bin/bash

OS="`uname`"
case $OS in
  'Darwin') 
    ### Install IBM Terraform Provider
    curl -LO https://github.com/IBM-Cloud/terraform-provider-ibm/releases/download/$(curl -s https://api.github.com/repos/IBM-Cloud/terraform-provider-ibm/releases/latest | grep tag_name | cut -d '"' -f 4)/terraform-provider-ibm_$(curl -s https://api.github.com/repos/IBM-Cloud/terraform-provider-ibm/releases/latest | grep tag_name | cut -d '"' -f 4 | cut -c 2-)_darwin_amd64.zip
    unzip terraform-provider-ibm_$(curl -s https://api.github.com/repos/IBM-Cloud/terraform-provider-ibm/releases/latest | grep tag_name | cut -d '"' -f 4 | cut -c 2-)_darwin_amd64.zip
    mkdir -p $HOME/.terraform.d/plugins/darwin_amd64
    sudo mv terraform-provider-ibm_$(curl -s https://api.github.com/repos/IBM-Cloud/terraform-provider-ibm/releases/latest | grep tag_name | cut -d '"' -f 4) $HOME/.terraform.d/plugins/darwin_amd64

    ### Install Terraformer
    curl -LO https://github.com/GoogleCloudPlatform/terraformer/releases/download/$(curl -s https://api.github.com/repos/GoogleCloudPlatform/terraformer/releases/latest | grep tag_name | cut -d '"' -f 4)/terraformer-ibm-darwin-amd64
    chmod +x terraformer-ibm-darwin-amd64
    sudo mv terraformer-ibm-darwin-amd64 /usr/local/bin/terraformer

    ### Install Configuration Discovery
    curl -LO https://github.com/IBM-Cloud/configuration-discovery/releases/download/$(curl -s https://api.github.com/repos/IBM-Cloud/configuration-discovery/releases | grep tag_name | cut -d '"' -f 4)/configuration_discovery_darwin_amd64
    chmod +x configuration_discovery_darwin_amd64
    sudo mv configuration_discovery_darwin_amd64 /usr/local/bin/discovery

    ;;	
  'Linux')
    ### Install IBM Terraform Provider
    curl -LO https://github.com/IBM-Cloud/terraform-provider-ibm/releases/download/$(curl -s https://api.github.com/repos/IBM-Cloud/terraform-provider-ibm/releases/latest | grep tag_name | cut -d '"' -f 4)/terraform-provider-ibm_$(curl -s https://api.github.com/repos/IBM-Cloud/terraform-provider-ibm/releases/latest | grep tag_name | cut -d '"' -f 4 | cut -c 2-)_linux_amd64.zip
    unzip terraform-provider-ibm_$(curl -s https://api.github.com/repos/IBM-Cloud/terraform-provider-ibm/releases/latest | grep tag_name | cut -d '"' -f 4 | cut -c 2-)_linux_amd64.zip
    mkdir -p $HOME/.terraform.d/plugins/linux_amd64
    sudo mv terraform-provider-ibm_$(curl -s https://api.github.com/repos/IBM-Cloud/terraform-provider-ibm/releases/latest | grep tag_name | cut -d '"' -f 4) $HOME/.terraform.d/plugins/linux_amd64

    ### Install Terraformer
    curl -LO https://github.com/GoogleCloudPlatform/terraformer/releases/download/$(curl -s https://api.github.com/repos/GoogleCloudPlatform/terraformer/releases/latest | grep tag_name | cut -d '"' -f 4)/terraformer-ibm-linux-amd64
	  chmod +x terraformer-ibm-linux-amd64
    sudo mv terraformer-ibm-linux-amd64 /usr/local/bin/terraformer   

    ### Install Configuration Discovery
    curl -LO https://github.com/IBM-Cloud/configuration-discovery/releases/download/$(curl -s https://api.github.com/repos/IBM-Cloud/configuration-discovery/releases | grep tag_name | cut -d '"' -f 4)/configuration_discovery_linux_amd64
    chmod +x configuration_discovery_linux_amd64
    sudo mv configuration_discovery_linux_amd64 /usr/local/bin/discovery

	;;
  'WindowsNT')
    OS='Windows'
    ;;
esac