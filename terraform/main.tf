module "iam_github_actions" {
  source              = "./iam"
  aws_region          = var.aws_region
  ecr_repository_name = var.ecr_repository_name
}

module "ecr" {
  source              = "./ecr"
  ecr_repository_name = var.ecr_repository_name
}

module "rds" {
  source                 = "./rds"
  rds_security_group_ids = module.security_groups.rds_security_group_ids
  rds_subnet_group_name  = module.security_groups.rds_subnet_group_name
}

module "security_groups" {
  source             = "./security_groups"
  private_vpc_id     = module.vpc.private_vpc_id
  private_subnet_ids = module.vpc.private_subnet_ids
}

module "vpc" {
  source = "./vpc"
}
