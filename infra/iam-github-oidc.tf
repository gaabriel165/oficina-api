module "github_oidc_provider" {
  source  = "terraform-aws-modules/iam/aws//modules/iam-github-oidc-provider"
  version = "~> 5.44"
}

resource "aws_iam_policy" "eks_describe" {
  name        = "${var.project_name}-eks-describe"
  description = "Allow the CI role to describe the EKS cluster for kubeconfig"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect   = "Allow"
        Action   = ["eks:DescribeCluster", "eks:ListClusters"]
        Resource = "*"
      }
    ]
  })
}

module "github_actions_role" {
  source  = "terraform-aws-modules/iam/aws//modules/iam-github-oidc-role"
  version = "~> 5.44"

  name     = "${var.project_name}-github-actions"
  subjects = ["repo:${var.github_repo}:*"]

  policies = {
    ecr = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryPowerUser"
    eks = aws_iam_policy.eks_describe.arn
  }

  depends_on = [module.github_oidc_provider]
}
