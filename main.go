package main

import (
	"flag"

	"github.com/chesshacker/terraform-provider-confluence/confluence"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &plugin.ServeOpts{
		Debug:        debugMode,
		ProviderFunc: confluence.Provider,
	}

	plugin.Serve(opts)
}
