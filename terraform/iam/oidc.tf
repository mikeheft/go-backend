resource "aws_iam_openid_connect_provider" "github_oidc" {
  client_id_list  = ["sts.amazonaws.com"]
  thumbprint_list = [var.thumbprint]
  url             = var.oidc_provider_url
}
