variable "aws_region" {
  description = "The AWS region to deploy the resources"
  default     = "us-west-2"
}

variable "bucket_name" {
  description = "unique name of bucket"
  default     = "mikeheft-simplebank"
}

variable "ecr_repository_name" {
  description = "The name of the ECR repository"
  default     = "simplebank"
}
