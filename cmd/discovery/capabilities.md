
## Capabilities

### Green-field import scenario 

In software engineering, green-field refers to importing existing infrastructure on a clean or blank state.


Let's assume that the following infrastructures were created through IBM Cloud console:
Virtual Private Cloud resource
IBM Cloud Object Storage resource
IAM access group and policy resource 


Going forward, if you want to adopt the Terraform tool to manage your existing cloud infrastructure, there are many benefits. This is a perfect green-field scenario — you can import the existing infrastructure where it generates both Terraform state file and configuration files to a folder. Once it generates the Terraform code, you can manage all your existing cloud infrastructure lifecycle operations through Terraform.


In this case, the configuration discovery tool will import the state from the existing infrastructure, and next time a `terraform apply` command is run, the Terraform will consider the resources in its state. This means any changes made to configuration will be picked up as modifications, rather than addition of new infrastructure.


### Upgrading the Terraform state file to the latest Terraform versions

Terraformer will import the state from existing infrastructure and exports a terraform configuration and state files to folder.
Even though the import was successful, but still Terraformer exported state file are not useful. By default Terraformer imports the sate file with 0.12 and If you have installed Terraform with 0.13 or latter version, and running the subsequent Terraform commands on the exported folder will throw a error due to this version incompatibility. In this case, the configuration discovery tool will export the terraform files and migrate the terraform state file from 0.12 to 0.13+ version compatible.
