package main

import (
	"bytes"
	"fmt"

	"github.com/urfave/cli"
)

var cliBuild = "0.1.1"
var buildDate string

type versionCheckInfo struct {
	outdated bool
	latest   string
}

// Should give out versions of discovery, terraform, terraformer Should also give a sentence if the discovery version is not the latest
// Should also give the available service names in the version information
// Version should throw error - if tf and tfr don’t meet recommended versions. Or every command should through this error
// Add json option for version

func versionCheckFunc() (versionCheckInfo, error) {
	// Call and Wait for the result to come through
	// Or have a chechpoint
	return versionCheckInfo{
		outdated: false,    // info.Outdated,
		latest:   cliBuild, // info.CurrentVersion,
		// Alerts:   alerts,
	}, nil
}

func actForVersion(ctx *cli.Context) error {
	var outdated bool
	var latest string
	var versionString bytes.Buffer

	fmt.Fprintf(&versionString, "Discovery v%s", cliBuild)

	// todo: @srikar - print the versions of terraform and terraformer

	// Check the latest version
	info, err := versionCheckFunc()
	if err != nil {
		ui.Failed(fmt.Sprintf(
			"\nError checking latest version: %s", err))
	}
	if info.outdated {
		outdated = true
		latest = info.latest
	}

	ui.Print(versionString.String())
	ui.Print(fmt.Sprintf("on %s", platform))
	ui.Print(fmt.Sprintf("built at %s", buildDate))

	ui.Print("List of IBM Cloud resources that can be imported:")
	PrintResources()

	if outdated {
		ui.Print(fmt.Sprintf(
			"\nYour discovery executable is out of date! The latest version\n"+
				"is %s. You can update by downloading from %s",
			latest, releasesLink))
	}

	return nil
}

// PrintActionLogs .
func PrintResources() {
	header := []string{"services", "resources"}
	table := ui.Table(header)
	table.Add("ibm_kp", "ibm_resource_instance")
	table.Add("", "ibm_kms_key")
	table.Add("", "")
	table.Add("ibm_cos", "ibm_resource_instance")
	table.Add("", "ibm_cos_bucket")
	table.Add("", "")
	table.Add("ibm_iam", "ibm_iam_user_policy")
	table.Add("", "ibm_iam_access_group")
	table.Add("", "ibm_iam_access_group_members")
	table.Add("", "ibm_iam_access_group_policy")
	table.Add("", "ibm_iam_access_group_dynamic_rule")
	table.Add("", "")
	table.Add("ibm_container_vpc_cluster", "ibm_container_vpc_cluster")
	table.Add("", "ibm_container_vpc_worker_pool")
	table.Add("", "")
	table.Add("ibm_database_etcd", "ibm_database")
	table.Add("", "")
	table.Add("ibm_database_mongo", "ibm_database")
	table.Add("", "")
	table.Add("ibm_database_postgresql", "ibm_database")
	table.Add("", "")
	table.Add("ibm_database_rabbitmq", "ibm_database")
	table.Add("", "")
	table.Add("ibm_database_redis", "ibm_database")
	table.Add("", "")
	table.Add("ibm_is_instance_group", "ibm_is_instance_group")
	table.Add("", "ibm_is_instance_group_manager")
	table.Add("", "ibm_is_instance_group_manager_policy")
	table.Add("", "")
	table.Add("ibm_cis", "ibm_cis")
	table.Add("", "ibm_cis_dns_record")
	table.Add("", "ibm_cis_firewall")
	table.Add("", "ibm_cis_domain_settings")
	table.Add("", "ibm_cis_global_load_balancer")
	table.Add("", "ibm_cis_edge_functions_action")
	table.Add("", "ibm_cis_edge_functions_trigger")
	table.Add("", "ibm_cis_healthcheck")
	table.Add("", "ibm_cis_rate_limit")
	table.Add("", "")
	table.Add("ibm_is_vpc", "ibm_is_vpc")
	table.Add("", "ibm_is_vpc_address_prefix")
	table.Add("", "ibm_is_vpc_route")
	table.Add("", "ibm_is_vpc_routing_table")
	table.Add("", "ibm_is_vpc_routing_table_route")
	table.Add("", "ibm_is_subnet")
	table.Add("", "ibm_is_instance")
	table.Add("", "")
	table.Add("ibm_is_security_group", "ibm_is_security_group_rule")
	table.Add("", "ibm_is_network_acl")
	table.Add("", "ibm_is_public_gateway")
	table.Add("", "ibm_is_volume")
	table.Add("", "")
	table.Add("ibm_is_vpn_gateway", "ibm_is_vpn_gateway_connections")
	table.Add("", "")
	table.Add("ibm_is_lb", "ibm_is_lb_pool")
	table.Add("", "ibm_is_lb_pool_member")
	table.Add("", "ibm_is_lb_listener")
	table.Add("", "ibm_is_lb_listener_policy")
	table.Add("", "ibm_is_lb_listener_policy_rule")
	table.Add("", "ibm_is_floating_ip")
	table.Add("", "ibm_is_flow_log")
	table.Add("", "ibm_is_ike_policy")
	table.Add("", "ibm_is_image")
	table.Add("", "ibm_is_instance_template")
	table.Add("", "ibm_is_ipsec_policy")
	table.Add("", "ibm_is_ssh_key")
	table.Add("", "")
	table.Add("ibm_function", "ibm_function_package")
	table.Add("", "ibm_function_action")
	table.Add("", "ibm_function_rule")
	table.Add("", "ibm_function_trigger")
	table.Add("", "")
	table.Add("ibm_private_dns", "ibm_resource_instance")
	table.Add("", "ibm_dns_zone")
	table.Add("", "ibm_dns_resource_record")
	table.Add("", "ibm_dns_permitted_network")
	table.Add("", "ibm_dns_glb_monitor")
	table.Add("", "ibm_dns_glb_pool")
	table.Add("", "ibm_dns_glb")
	table.Add("ibm_satellite", "ibm_satellite_location")
	table.Add("", "ibm_satellite_host")
	table.Print()
}
