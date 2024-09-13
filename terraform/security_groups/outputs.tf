output "rds_security_group_ids" {
  value = module.rds.aws_security_group.rds_sg.id
}
