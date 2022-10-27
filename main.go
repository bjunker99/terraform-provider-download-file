package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"terraform-provider-download-file/download-file"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return download_file.Provider()
		},
	})
}
