output "rds_security_group_ids" {
  value = module.rds.aws_security_group.rds_sg.id
}

output "rds_subnet_group_name" {
  value = module.rds.aws_db_subnet_group.rds_subnet_group.name
}
