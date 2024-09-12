variable "aws_region" {
  description = "The AWS region to deploy the resources"
}

variable "ecr_repository_name" {
  description = "The name of the ECR repository"
}

variable "oidc_provider_url" {
  description = "OIDC provider URL for GitHub Actions"
  default     = "https://token.actions.githubusercontent.com"
}

variable "thumbprint" {
  description = "The thumbprint for GitHub OIDC provider"
  default     = "6938fd4d98bab03faadb97b34396831e3780aea1"
}
