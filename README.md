# Terraform plugin for fargate deployment

## Usage
1) Copy `plugins` and `schemas` folders into `$HOME/.terraform.d` (or `%APPDATA%/terraform.d` in case of Windows) 
2) Run terraform command `terraform init`
3) Restart PhpStorm   


```hcl-terraform
data "simplia_ecs_current_deployment" "deployment" {
  cluster = module.cluster.cluster_name
  service = "test" # Service name
  container_name = "test" 

  default_image_digest = "default" # Fallback, when no image is deployed
  region = "eu-west-1"
}
# ...
resource "aws_ecs_task_definition" "task" {
  family = "http-redirects"
  cpu = "256"
  memory = "512"
  execution_role_arn = aws_iam_role.task_role.arn
  task_role_arn = aws_iam_role.role.arn
  requires_compatibilities = [
    "FARGATE"
  ]
  network_mode = "awsvpc"
  container_definitions = <<DEFINITION
[
  {
    "name": "test",
    "image": "ecr_url:${data.simplia_ecs_current_deployment.deployment.image_digest}",
    "essential": true,
    "environment": [],
    "portMappings": [],
    "logConfiguration": {
      "logDriver": "awslogs",
      "secretOptions": null,
      "options": {
        "awslogs-group": "${aws_cloudwatch_log_group.log.name}",
        "awslogs-region": "eu-west-1",
        "awslogs-stream-prefix": "task"
      }
    }
  }
]
DEFINITION
}
```
Note our plugin is used to determine image digest **data.simplia_ecs_current_deployment.deployment.image_digest** - it returns currently
used image, or default value in case there is none (e. g. first deployment of this service).


## Kompilace
```bash
$  go build -o plugins/terraform-provider-simplia
```

## Zdroje
* https://github.com/VladRassokhin/intellij-hcl
* https://github.com/VladRassokhin/terraform-metadata
* https://www.terraform.io/docs/extend/index.html