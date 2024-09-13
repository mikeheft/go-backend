module "rds" {
  source             = "./rds"
  private_vpc_id     = var.private_vpc_id
  private_subnet_ids = var.private_subnet_ids
}
