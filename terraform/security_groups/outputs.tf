output "rds_security_group_ids" {
  value = [module.rds.rds_sg_id]
}

output "rds_subnet_group_name" {
  value = module.rds.rds_subnet_group_name
}
