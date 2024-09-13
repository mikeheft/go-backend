resource "aws_db_instance" "postgres" {
  identifier        = "simple-bank-postgres"
  allocated_storage = 20
  engine            = "postgres"
  engine_version    = 14
  instance_class    = "db.t2.micro"
  # Will be annonymized in actual prod project
  username               = "root"
  password               = "secret"
  publicly_accessible    = false
  skip_final_snapshot    = true
  vpc_security_group_ids = var.rds_security_group_ids
}
