package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceSimpliaEcsCurrentDeployment() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSimpliaEcsCurrentDeploymentRead,

		Schema: map[string]*schema.Schema{
			"cluster": {
				Type:     schema.TypeString,
				Required: true,
			},
			"container_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"default_image_digest": {
				Type:     schema.TypeString,
				Required: true,
			},
			"region": {
				Type:     schema.TypeString,
				Required: true,
			},
			"service": {
				Type:     schema.TypeString,
				Required: true,
			},

			// Computed values.
			"image_digest": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_found": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceSimpliaEcsCurrentDeploymentRead(d *schema.ResourceData, meta interface{}) error {
	awsSession, _ := session.NewSession()
	ecsClient := ecs.New(awsSession, aws.NewConfig().WithRegion(d.Get("region").(string)))

	serviceName := d.Get("service").(string)
	describeServicesInput := &ecs.DescribeServicesInput{
		Cluster:  aws.String(d.Get("cluster").(string)),
		Services: []*string{&serviceName},
	}

	var currentTaskDefinition string
	services, servicesError := ecsClient.DescribeServices(describeServicesInput)

	if servicesError != nil || len(services.Failures) > 0 {
		if(servicesError != nil) {
			log.Println("[SIMPLIA] %s", servicesError.Error())
		}

		d.SetId("default")
		d.Set("image_digest", d.Get("default_image_digest"))
		d.Set("image_found", false)

		return nil
	}

	for _, serviceDefinition := range services.Services {
		currentTaskDefinition = *serviceDefinition.TaskDefinition
	}

	describeTaskDefinitionInput := &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: aws.String(currentTaskDefinition),
	}
	taskDefinitionOutput, taskDefinitionError := ecsClient.DescribeTaskDefinition(describeTaskDefinitionInput)

	if taskDefinitionError != nil {
		log.Println("[SIMPLIA] %s", servicesError.Error())

		d.SetId("default")
		d.Set("image_digest", d.Get("default_image_digest"))
		d.Set("image_found", false)

		return nil
	}

	taskDefinition := *taskDefinitionOutput.TaskDefinition
	for _, containerDefinition := range taskDefinition.ContainerDefinitions {
		if aws.StringValue(containerDefinition.Name) != d.Get("container_name").(string) {
			continue
		}

		d.SetId(aws.StringValue(taskDefinition.TaskDefinitionArn))
		image := aws.StringValue(containerDefinition.Image)
		if strings.Contains(image, ":") {
			d.Set("image_digest", strings.Split(image, ":")[1])
		}
		d.Set("image_found", true)
	}

	if d.Id() == "" {
		return fmt.Errorf("Task definition or container with name %q not found.", d.Get("container_name").(string), d.Get("task_definition").(string))
	}

	return nil
}