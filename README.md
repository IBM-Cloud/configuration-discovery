# IBM Cloud Configuration Discovery

Use "Configuration Discovery" to import the existing Cloud resources (in your account) and its configuration settings - to auto-generate the terraform configuration file (.tf) and state file (.tfstate).  It makes it easy for you to adopt the Infrastructure-as-Code practices; it can reverse engineer the current IBM Cloud environment (that was provisioned using UI or CLI).  

## Dependencies

-   [Terraform](https://www.terraform.io/downloads.html) 0.9.3+
-   [Terraformer](https://github.com/GoogleCloudPlatform/terraformer) 0.8.15+
-   [Go](https://golang.org/doc/install) 1.15+ (to build the discovery cli)
-   [IBM Cloud Provider](https://github.com/IBM-Cloud/terraform-provider-ibm/)
-   [Mongodb](https://docs.mongodb.com/manual/installation/) v4.4.5+


## Steps to use the Configuration Discovery project
### Files

*   main.go

    This file contains the web server and handlers.

*   cmd/discovery

    Code for the executable.

### Steps to use the project as an executable

*  Build and install the executable to your GOPATH

       make install-cli

*  Export the required env vars 

       IC_API_KEY: Your ibm cloud api key. Imports resources in that account, for which user has access.
       DISCOVERY_CONFIG_DIR: Directory, where to generate and import the terraform code. discovery uses only this directory for any read/write op.

*  Example commands

       discovery config --config_name testfolder
       discovery import --services ibm_is_vpc --config_name testfolder --compact --merge


### Steps to use the project as a server 
 <!-- todo: @anil - add the swagger api link here, may be later we can host the swagger github page if needed -->

*  Start the server

        cd /go/src/github.com
        git clone git@github.ibm.com:IBMTerraform/configuration-discovery/.git
        cd configuration-discovery/
        go run main.go docs.go
        http://localhost:8080

*  Or 

       make run-mac <or make run-local for linux>

### How to run the Configuration Discovery as a docker container

*  Or run as docker container

       make docker-build
       make docker-run

    First two need mongodb service running on localhost:27017. Third needs mongodb running as docker container. To run mongodb as docker container with 27017 exposed outside. This will work for all three methods above. Run this before any of the above steps
        
        make docker-run-mongo
        


### How to run the terraform-ibmcloud-provider-api as a container
        
        cd /go/src/github.com
        git clone git@github.ibm.com:IBMTerraform/configuration-discovery/.git
        cd configuration-discovery/
        docker build -t configuration-discovery .
        docker images
        export API_IMAGE=configuration-discovery:latest
        docker-compose up --build -d
        
### Contributing to Configuration Discovery

Please have a look at the [CONTRIBUTING.md](./CONTRIBUTING.md) file for some guidance before
submitting a pull request. Thank you for your help and interest!


## Discovery executable

### Install

Head over to the [releases page](https://github.com/anilkumarnagaraj/terraform-provider-ibm-api/releases) and download the latest release and place it under the path. For example,

```
wget https://github.com/anilkumarnagaraj/terraform-provider-ibm-api/releases/download/v0.1.1/discovery_0.1.1_darwin_arm64.zip
untar discovery_0.1.1_darwin_arm64.zip
mv discovery <$GOPATH/bin or any other directory in path>
```

### Install using go

If you have go installed, run this command 
<!-- Need to verify and update this. Add the -u flag to update -->
```
go install github.com/anilkumarnagaraj/terraform-provider-ibm-api/cmd/discovery
```
### Commands

Below commands help you to import your resources into terraform configuration. For more detailed description, run `discovery help`. For help on a command, run `discovery help <command>`

| Command                           | Description                                                                                                                                                                                             |
| --------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| discovery version | IBM Cloud Discovery CLI version. Gives you other dependency information. Shows all supported importable resources                    |
| discovery config                         | Create a local configuration directory for importing the terrraform configuration.    |
| discovery import | Import TF config for resources in your ibm cloud account. Import all the resources for this service. Imports config and statefile. If a statefile is already present, merging will be done.        |

Start using [discovery executable](cmd/discovery/tutorial.md)


## Report a Issue / Feature request

-   Is something broken? Have a issue/bug to report? use the [Bug report](https://github.com/IBM-Cloud/configuration-discovery/issues/new?assignees=&labels=&template=bug_report.md&title=) link. But before raising a issue, please check the [issues list](https://github.com/IBM-Cloud/configuration-discovery/issues) to see if the issue is already raised by someone
-   Do you have a new feature or enhancement you would like to see? use the [Feature request](https://github.com/IBM-Cloud/configuration-discovery/issues/new?assignees=&labels=&template=feature_request.md&title=) link.

## License

The project is licensed under the [Apache 2.0 License](https://www.apache.org/licenses/LICENSE-2.0).
A copy is also available in the LICENSE file in this repository.