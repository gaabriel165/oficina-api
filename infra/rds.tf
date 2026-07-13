resource "random_password" "db" {
  length  = 20
  special = false
}

resource "aws_db_subnet_group" "this" {
  name       = "${var.project_name}-db"
  subnet_ids = module.vpc.private_subnets
}

resource "aws_security_group" "db" {
  name        = "${var.project_name}-db"
  description = "Allow PostgreSQL access from the EKS worker nodes"
  vpc_id      = module.vpc.vpc_id

  ingress {
    description     = "PostgreSQL from EKS nodes"
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    security_groups = [module.eks.node_security_group_id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_db_instance" "this" {
  identifier     = "${var.project_name}-db"
  engine         = "postgres"
  engine_version = var.db_engine_version
  instance_class = var.db_instance_class

  db_name  = var.db_name
  username = var.db_username
  password = random_password.db.result

  allocated_storage = 20
  storage_type      = "gp3"
  storage_encrypted = true

  db_subnet_group_name   = aws_db_subnet_group.this.name
  vpc_security_group_ids = [aws_security_group.db.id]
  multi_az               = false
  publicly_accessible    = false

  skip_final_snapshot = true
  deletion_protection = false
}
