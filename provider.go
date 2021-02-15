package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type Config struct {
	AssumeRoleARN               string
	AssumeRoleSessionName       string
}

func Provider() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"assume_role": assumeRoleSchema(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"deployment_ecs": dataSourceSimpliaEcsCurrentDeployment(),
		},
	}

	provider.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		terraformVersion := provider.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}

		config := Config{}

		if l, ok := d.Get("assume_role").([]interface{}); ok && len(l) > 0 && l[0] != nil {
			m := l[0].(map[string]interface{})

			if v, ok := m["role_arn"].(string); ok && v != "" {
				config.AssumeRoleARN = v
			}

			if v, ok := m["session_name"].(string); ok && v != "" {
				config.AssumeRoleSessionName = v
			}
		}

		return config, nil
	}

	return provider
}

func assumeRoleSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"role_arn": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Amazon Resource Name of an IAM Role to assume prior to making API calls.",
				},
				"session_name": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Identifier for the assumed role session.",
				},
			},
		},
	}
}
