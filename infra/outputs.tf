output "region" {
  description = "AWS region"
  value       = var.region
}

output "cluster_name" {
  description = "EKS cluster name"
  value       = module.eks.cluster_name
}

output "configure_kubectl" {
  description = "Command to point kubectl and Lens at the cluster"
  value       = "aws eks update-kubeconfig --name ${module.eks.cluster_name} --region ${var.region}"
}

output "ecr_repository_url" {
  description = "ECR repository URL for the application image"
  value       = aws_ecr_repository.this.repository_url
}

output "rds_endpoint" {
  description = "RDS PostgreSQL endpoint (host)"
  value       = aws_db_instance.this.address
}

output "database_url" {
  description = "Full connection string for the application"
  value       = "postgres://${var.db_username}:${random_password.db.result}@${aws_db_instance.this.address}:5432/${var.db_name}?sslmode=require"
  sensitive   = true
}

output "github_actions_role_arn" {
  description = "IAM role ARN assumed by GitHub Actions via OIDC"
  value       = module.github_actions_role.arn
}
