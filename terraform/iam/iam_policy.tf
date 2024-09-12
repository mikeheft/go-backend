data "aws_caller_identity" "current" {
  # No configuration required, this retrieves information about the current caller
}

resource "aws_iam_policy" "github_actions_ecr_policy" {
  name        = "gh-actions-ecr-policy"
  description = "Policy for GitHub Actions to interact with ECR"

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Action = [
          "ecr:BatchCheckLayerAvailability",
          "ecr:BatchGetImage",
          "ecr:CompleteLayerUpload",
          "ecr:GetDownloadUrlForLayer",
          "ecr:InitiateLayerUpload",
          "ecr:PutImage",
          "ecr:UploadLayerPart"
        ],
        Resource = "arn:aws:ecr:${var.aws_region}:${data.aws_caller_identity.current.account_id}:repository/${var.ecr_repository_name}"
      },
      {
        Effect   = "Allow",
        Action   = "ecr:GetAuthorizationToken",
        Resource = "*"
      },
      {
        Effect   = "Allow",
        Action   = "sts:TagSession",
        Resource = "*"
      }
    ]
  })
}
