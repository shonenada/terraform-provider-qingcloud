package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/shonenada/terraform-provider-qingcloud/qingcloud"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: qingcloud.Provider})
}
