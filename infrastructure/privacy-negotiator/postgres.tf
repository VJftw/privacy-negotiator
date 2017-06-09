resource "aws_db_instance" "db" {
  allocated_storage    = 10
  storage_type         = "gp2"
  engine               = "postgres"
  instance_class       = "db.t2.micro"
  name                 = "privneg_${var.environment}"
  username             = "${var.postgres_user}"
  password             = "${var.postgres_password}"
  vpc_security_group_ids = ["${aws_security_group.db.id}"]
  db_subnet_group_name = "${aws_db_subnet_group.db.id}"
  // skip_final_snapshot = true
  final_snapshot_identifier = "priv-neg-${var.environment}"
}

resource "aws_db_subnet_group" "db" {
  name       = "privneg"
  subnet_ids = [
    "${data.aws_subnet.app_cluster.0.id}",
    "${data.aws_subnet.app_cluster.1.id}",
    "${data.aws_subnet.app_cluster.2.id}",
    ]

  tags {
    Name = "privneg-db-subnet-grp"
  }
}

resource "aws_security_group" "db" {

  name        = "${var.environment}-db-sg"
  vpc_id      = "${data.aws_vpc.app_cluster.id}"

  # MySQL access from the VPC
  ingress {
    from_port   = 3306
    to_port     = 3306
    protocol    = "tcp"
    cidr_blocks = ["10.0.0.0/16"]
  }

}
