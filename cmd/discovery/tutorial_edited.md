# IBM Cloud Configuration Discovery

Let's take a scenario that you have an infrastructure containing multiple services provisioned on IBM Cloud. One of your business requirement is to migrate existing infrastructure setup to declarative approach using Terraform. Since you have basic skill of Terraform, thinking more about challenges and workaround to create an entire infrastructure to Terraform. 
You need not think about, here is a tutorial that walks your through the [Cloud Configuration Discovery](https://test.cloud.ibm.com/docs/ibm-cloud-provider-for-terraform?topic=ibm-cloud-provider-for-terraform-terraformer-intro) to simplify your requirement.
For more information, about getting started with the tool, see [blog](https://ibm.box.com/s/0ou4erd2t65ndiv1v83egfjgle699pcy).

Objectives
Following are the objectives that are covered in this tutorial.
- Usage
- Prerequisites 
- Installation and setup
- Generate infrastructure to code
- 

## Usage

The Cloud Configuration Discovery is a reverse engineering tool used to generate your cloud infrastructure to code. Recreating or coding entire provisioned resource requires minimal coding skill and less time to implement.

The Cloud configuration discovery is a command line tool that augments an open-source Google Terraformer project with a number of capabilities that are native to IBM Cloud.

Google Terraformer is an open-source project with a number of capabilities that are native to IBM Cloud. For more information, about Google Terrafomer, see [Google Terraform support for IBM Cloud](https://github.com/GoogleCloudPlatform/terraformer/blob/master/docs/ibmcloud.md)

## Prerequisites

In your local system, following prerequisites must be met to install the Cloud Configuration Discovery plugin.
- [Terraform](https://test.cloud.ibm.com/docs/ibm-cloud-provider-for-terraform?topic=ibm-cloud-provider-for-terraform-getting-started#tf_installation)

## Installation and setup

The installation can be done by following two approaches:

**Approach 1**
Install precompiled binary by running following CURL command.

```
curl -qL https://github.com/IBM-Cloud/configuration-discovery/install.sh | sh
```
Clone and install the tool to your GOPATH

```
make install-cli
```

**Approach 2**
Install by using the precompiled binary files manually by following the given steps.

- Download the binary file from the latest [release](https://github.com/IBM-Cloud/configuration-discovery/releases).
- Rename the file to **discovery**, ensure it is executable, and place the binary in the PATH. 

  **Example**
  
  ```
  mv discovery <$GOPATH/bin or any other directory in path>
  ```
- Or, If you have `GO` lang installed, run this command.

    ```
    go install github.com/IBM-Cloud/configuration-discovery/cmd/discovery
    ```
### Setup your environment

- Export the required `env vars` 
    * IC_API_KEY: Provide your IBM Cloud API key. This imports resources that your account has access.
    * DISCOVERY_CONFIG_DIR: The directory, where you want to generate and import the Terraform code. Cloud Configuration Discovery uses only this directory for all the read or write operations.
- Run `discovery help` to see all the commands as `discovery help <command>`.
   **Example**
       discovery version
       discovery config --config_name testfolder
       discovery import --services ibm_is_vpc --config_name testfolder --compact --merge
-------







### Tutorial with examples

- Run `discovery version`. This will show the version of the discovery binary and all the services and resources that can be imported. Here is the list of services and the supported resources.
    ```
        Configuration Discovery v0.0.1 on unix
        List of IBM Cloud resources that can be imported:
        services                    resources
        ibm_kp                      ibm_resource_instance
                                    ibm_kms_key

        ibm_cos                     ibm_resource_instance
                                    ibm_cos_bucket

        ibm_iam                     ibm_iam_user_policy
                                    ibm_iam_access_group
                                    ibm_iam_access_group_members
                                    ibm_iam_access_group_policy
                                    ibm_iam_access_group_dynamic_rule

        ibm_container_vpc_cluster   ibm_container_vpc_cluster
                                    ibm_container_vpc_worker_pool

        ibm_database_etcd           ibm_database

        ibm_database_mongo          ibm_database

        ibm_database_postgresql     ibm_database

        ibm_database_rabbitmq       ibm_database

        ibm_database_redis          ibm_database

        ibm_is_instance_group       ibm_is_instance_group
                                    ibm_is_instance_group_manager
                                    ibm_is_instance_group_manager_policy

        ibm_cis                     ibm_cis
                                    ibm_cis_dns_record
                                    ibm_cis_firewall
                                    ibm_cis_domain_settings
                                    ibm_cis_global_load_balancer
                                    ibm_cis_edge_functions_action
                                    ibm_cis_edge_functions_trigger
                                    ibm_cis_healthcheck
                                    ibm_cis_rate_limit

        ibm_is_vpc                  ibm_is_vpc
                                    ibm_is_vpc_address_prefix
                                    ibm_is_vpc_route
                                    ibm_is_vpc_routing_table
                                    ibm_is_vpc_routing_table_route
                                    ibm_is_subnet
                                    ibm_is_instance

        ibm_is_security_group       ibm_is_security_group_rule
                                    ibm_is_network_acl
                                    ibm_is_public_gateway
                                    ibm_is_volume

        ibm_is_vpn_gateway          ibm_is_vpn_gateway_connections

        ibm_is_lb                   ibm_is_lb_pool
                                    ibm_is_lb_pool_member
                                    ibm_is_lb_listener
                                    ibm_is_lb_listener_policy
                                    ibm_is_lb_listener_policy_rule
                                    ibm_is_floating_ip
                                    ibm_is_flow_log
                                    ibm_is_ike_policy
                                    ibm_is_image
                                    ibm_is_instance_template
                                    ibm_is_ipsec_policy
                                    ibm_is_ssh_key

        ibm_function                ibm_function_package
                                    ibm_function_action
                                    ibm_function_rule
                                    ibm_function_trigger

        ibm_private_dns             ibm_resource_instance
                                    ibm_dns_zone
                                    ibm_dns_resource_record
                                    ibm_dns_permitted_network
                                    ibm_dns_glb_monitor
                                    ibm_dns_glb_pool
                                    ibm_dns_glb

        ibm_satellite
                                    ibm_satellite_location
                                    ibm_satellite_host                           
    ```

- You have a vpc created in the IBM cloud, but now you want to maintain this vpc using terraform. 
    ```
        discovery import --services ibm_is_vpc
    ```

    For example, this imported the resources ibm_is_vpc and ibm_is_vpc_address_prefix in my account 

    ![image](images/vpc.png)

    Instead, to get all the resources in a single terraform file, run
    ```
        discovery import --services ibm_is_vpc --compact
    ```

    Here, the name of the folder inside DISCOVERY_CONFIG_DIR, is randomly generated string. To pass a proper folder_name. If this folder already exists, make sure it is empty. Merging will old terraform files is not supported yet.
    ```
        discovery import --services ibm_is_vpc --config_name vpcconfig --compact
    ```

    ![image](images/vpccompact.png)

- Import all the access_groups in your account. This is one resource `ibm_iam_access_group` of service `ibm_iam`. Run
    ```
        discovery import --services ibm_iam_access_group --config_name myaccessgroups
    ```

- Import a particular vpc. You can use tags to filter the resources. Run 
    ```
        discovery import --services ibm_vpc --tags resource_group:1edf423492d34609a7759de8ed7e675a --config_name vpctag --compact
    ```


## Validating & Re-creating the Environment

Now that we’ve captured the environment into Terraform files using the Configuration Discovery tool, you can re-create that environment in the IBM Cloud using Schematics. Let’s walk through doing that.

### Scenario

You’ve manually created an IBM Kubernetes Service and a VPC on the IBM Cloud, you’ve deployed your application, and it’s working perfectly. Now, you want to put that into production just as you have it, using Infrastructure as Code, so you’ve run the Configuration Discovery tool, and captured that environment to Terraform files. These are the files you’ve captured

```
* resources.tf -  Provides the configuration.
* terraform.tfstate  -  Holds the last-known state of the infrastructure.
* variables.tf  -  Contain values for the declared variables.
* provider.tf   - Define which providers require for terraform to install and use.
* outputs.tf - Allow you to export structured data about your infrastructure.
```

These represent the environment you saved and want to duplicated. In this case, you are ready to go to full production, so you want to create 3 environments: Dev, Staging, and Prod. To create three different environments, you’ll need to take those files and duplicate them, one set for each environment type, and edit the files and inject new input values that are unique to each environment. 

    .
    ├── ...
    ├── MyApp-Dev               
    │   ├── resources.tf    
    │   ├── terraform.tfstate
    │   ├── variables.tf    
    │   ├── provider.tf      
    │   └── outputs.tf 
    ├── MyApp-Staging               
    │   ├── resources.tf    
    │   ├── terraform.tfstate
    │   ├── variables.tf    
    │   ├── provider.tf      
    │   └── outputs.tf 
    ├── MyApp-Prod               
    │   ├── resources.tf    
    │   ├── terraform.tfstate
    │   ├── variables.tf    
    │   ├── provider.tf      
    │   └── outputs.tf      
    └── ...


To edit the files, load them in an editor, and open the `resources.tf` or `variables.tf` , and change the input values for each environment.
Once you save all three sets of files, you can store them on your disk or in a Github repository under different directories.

Now you can create these 3 environments using IBM Cloud Schematics. You will run a series of commands from Schematics 3 times, all from the CLI.

### Create Environments Commands:

- Create a schematics workspace 

```
ibmcloud schematics workspace new --file myapp-prod.json --state STATE_FILE_PATH

myapp-prod.json:
{
  "name": "cde",
  "shared_data": {
    "region": "us-south"
  },
  "type": [
    "terraform_v0.13"
  ],
  "description": "terraform workspace",
  "tags": [
    "department:HR",
    "application:compensation",
    ";"
  ]
    }
  ]
}
```

- Create archive file (.tar) 

`tar -cf myapp-prod.tar terraform.tfstate resources.tf ...`

- Upload archive file (.tar) to your Schematics workspace

`ibmcloud schematics workspace upload  --id WORKSPACE_ID --file myapp-prod.tar --template TEMPLATE_ID`

- Provision infrastructure

`ibmcloud schematics apply --id WORKSPACE_ID --var-file vars.tfvars`

Run the above commands three times on all environment folders(MyApp-Dev, MyApp-Staging, and MyApp-Prod).

Now that you’ve executed these commands, and created these environments, let’s look at them and use them.

### View environments Commands

- List the IBM Cloud resources provisioned in schematics workspace.

`ibmcloud schematics state list --id WORKSPACE_ID`

As you can see, now you have 3 sets of resources for your application MyApp - one for Dev, one for Staging, and one for Prod.
You can now deploy your application on each one, do your development, staging, and production work, confident that you have an identical environment as you originally set up and had working.

## Conclusion

So that is the process to reverse-engineer your existing cloud environment and re-deploy it as needed.

The steps are:

- Use this Configuration Discovery tool to reproduce a manually created cloud environment by capturing the Terraform “Infrastructure to Code” files
- Modify the files to create 3 different environments: Dev, Staging, Prod.
- Use the to re-create those resources using IBM Cloud Schematics, and then to check they were created the same way as you originally had set up


You now have the ability to take your own manually created environments, capture them to re-usable Terraform code files, and then re-create in the IBM Cloud for future use.

## Future enhancements

- Import will support a new flag `--merge`. This can be used to import the resources and merge with existing terraform statefile and configuration. 
- Rewrite terraform resource into terraform modules.
