
module "iam_github_actions" {
  source              = "./iam"
  aws_region          = var.aws_region
  ecr_repository_name = var.ecr_repository_name
}

module "ecr" {
  source              = "./ecr"
  ecr_repository_name = var.ecr_repository_name
}
