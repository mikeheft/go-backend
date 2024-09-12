terraform {
  backend "s3" {
    bucket  = "mikeheft-simplebank"
    key     = "mikeheft-simplebank/terraform.tfstate"
    region  = "us-west-2"
    encrypt = true
  }
}
