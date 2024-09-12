variable "aws_region" {
  description = "The AWS region to deploy the resources"
  default     = "us-west-2"
}

variable "ecr_repository_name" {
  description = "The name of the ECR repository"
  default     = "simplebank"
}

variable "github_org_repo" {
  description = "GitHub organization and repository in the format org/repo"
  default     = "mikeheft/*"
}

variable "branch_name" {
  description = "The branch name for GitHub Actions"
  default     = "main"
}
