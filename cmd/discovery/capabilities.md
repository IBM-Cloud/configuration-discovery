# Capabilities

## Green-field import scenario

In software engineering, Green-field refers to importing existing infrastructure on a clean or blank state.

Let us assume that the following infrastructures were created through IBM Cloud console:
- Virtual Private Cloud resource
- COS resource
- IAM access group & policy resource 

Going forward, if you want to adopt the Terraform tool to manage your existing cloud infrastructure, there are many benefits. This is a perfect green-field scenario, you can import the existing infrastructure where it generates both Terraform state file & configuration files to a folder. Once it generates the Terraform code, you can manage all your existing cloud infrastructure lifecycle operations through Terraform.

In this case, the configuration discovery tool will import the state from the existing infrastructure, and next time a "terraform apply" command is run, the Terraform will consider the resources in its state. This means any changes made to configuration will be picked up as modifications, rather than addition of new infrastructure.

## Upgrading the Terraform state file to the latest Terraform versions

