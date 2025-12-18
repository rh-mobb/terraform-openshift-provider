// Package main provides the entry point for the Terraform provider binary.
package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/redhat/terraform-provider-openshift-operator/internal/provider"
	"github.com/redhat/terraform-provider-openshift-operator/internal/version"
)

func main() {
	var debugMode bool
	var showVersion bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.BoolVar(&showVersion, "version", false, "show version information")
	flag.Parse()

	if showVersion {
		log.Printf("terraform-provider-openshift %s\n", version.String())
		return
	}

	opts := &plugin.ServeOpts{
		ProviderFunc: provider.New,
	}

	if debugMode {
		//nolint:staticcheck // plugin.Debug is the standard way to debug providers in SDK v2
		if err := plugin.Debug(context.Background(), "registry.terraform.io/rh-mobb/openshift", opts); err != nil {
			log.Fatal(err.Error())
		}
		return
	}

	plugin.Serve(opts)
}
