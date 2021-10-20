#!/bin/bash

OS="`uname`"
case $OS in
  'Darwin') 
    ### Install Terraformer
    curl -LO https://github.com/GoogleCloudPlatform/terraformer/releases/download/$(curl -s https://api.github.com/repos/GoogleCloudPlatform/terraformer/releases/latest | grep tag_name | cut -d '"' -f 4)/terraformer-ibm-darwin-amd64
    chmod +x terraformer-ibm-darwin-amd64
    sudo mv terraformer-ibm-darwin-amd64 /usr/local/bin/terraformer

    ### Install Configuration Discovery
    curl -LO https://github.com/anilkumarnagaraj/configuration-discovery-1/releases/download/$(curl -s https://api.github.com/repos/anilkumarnagaraj/configuration-discovery-1/releases | grep tag_name | cut -d '"' -f 4)/discovery_darwin_amd64
    chmod +x configuration_discovery_darwin_amd64
    sudo mv discovery_darwin_amd64 /usr/local/bin/discovery

    ;;	
  'Linux')
    ### Install Terraformer
    curl -LO https://github.com/GoogleCloudPlatform/terraformer/releases/download/$(curl -s https://api.github.com/repos/GoogleCloudPlatform/terraformer/releases/latest | grep tag_name | cut -d '"' -f 4)/terraformer-ibm-linux-amd64
	chmod +x terraformer-ibm-linux-amd64
    sudo mv terraformer-ibm-linux-amd64 /usr/local/bin/terraformer   

    ### Install Configuration Discovery
    curl -LO https://github.com/anilkumarnagaraj/configuration-discovery-1/releases/download/$(curl -s https://api.github.com/repos/anilkumarnagaraj/configuration-discovery-1/releases | grep tag_name | cut -d '"' -f 4)/discovery_linux_amd64    
    chmod +x configuration_discovery_linux_amd64
    sudo mv discovery_linux_amd64 /usr/local/bin/discovery

	;;
  'WindowsNT')
    OS='Windows'
    ;;
esac