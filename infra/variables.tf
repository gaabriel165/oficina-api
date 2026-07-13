variable "region" {
  description = "AWS region to provision into"
  type        = string
  default     = "us-east-1"
}

variable "project_name" {
  description = "Name prefix applied to all resources"
  type        = string
  default     = "oficina-api"
}

variable "cluster_name" {
  description = "EKS cluster name"
  type        = string
  default     = "oficina-api-eks"
}

variable "cluster_version" {
  description = "Kubernetes version for the EKS control plane"
  type        = string
  default     = "1.31"
}

variable "vpc_cidr" {
  description = "CIDR block for the VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "node_instance_type" {
  description = "EC2 instance type for the EKS managed node group"
  type        = string
  default     = "t3.small"
}

variable "node_desired_size" {
  description = "Desired number of worker nodes"
  type        = number
  default     = 2
}

variable "node_min_size" {
  description = "Minimum number of worker nodes"
  type        = number
  default     = 2
}

variable "node_max_size" {
  description = "Maximum number of worker nodes"
  type        = number
  default     = 4
}

variable "db_instance_class" {
  description = "RDS instance class"
  type        = string
  default     = "db.t3.micro"
}

variable "db_engine_version" {
  description = "PostgreSQL engine version for RDS"
  type        = string
  default     = "16.4"
}

variable "db_name" {
  description = "Initial database name"
  type        = string
  default     = "oficina_mecanica"
}

variable "db_username" {
  description = "Master username for RDS"
  type        = string
  default     = "postgres"
}

variable "github_repo" {
  description = "owner/repo allowed to assume the CI role via GitHub OIDC"
  type        = string
  default     = "gaabriel165/oficina-api"
}
