module "rds" {
  source         = "./rds"
  private_vpc_id = var.private_vpc_id
}
