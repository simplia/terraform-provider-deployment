package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		DataSourcesMap: map[string]*schema.Resource{
			"simplia_ecs_current_deployment": dataSourceSimpliaEcsCurrentDeployment(),
		},
	}
}
