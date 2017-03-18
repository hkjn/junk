provider "aws" {
  region = "${var.region}"
}

#
# Network resources.
#

resource "aws_vpc" "main" {
  cidr_block           = "${var.cidr_block}"
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags {
    Name          = "${var.env}_main"
    orchestration = "terraform"
  }
}

resource "aws_internet_gateway" "public" {
  vpc_id = "${aws_vpc.main.id}"

  tags {
    Name          = "${var.env}_igw"
    orchestration = "terraform"
  }
}

resource "aws_route_table" "tf_public_routes" {
  vpc_id = "${aws_vpc.main.id}"

  tags {
    Name          = "${var.env}_public_subnet_route_table"
    orchestration = "terraform"
  }
}

# Add route to public gateway
resource "aws_route" "tf_public_gateway_route" {
  route_table_id         = "${aws_route_table.tf_public_routes.id}"
  depends_on             = ["aws_route_table.tf_public_routes"]
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = "${aws_internet_gateway.public.id}"
}

resource "aws_subnet" "tf_subnets" {
  vpc_id            = "${aws_vpc.main.id}"
  cidr_block        = "${element(split(",", var.private_subnets), count.index)}"
  availability_zone = "${element(split(",", var.azs), count.index)}"
  count             = "${length(compact(split(",", var.private_subnets)))}"

  # We enable public IPs for all instances by default to be able to SSH to the workers.
  map_public_ip_on_launch = true

  tags {
    Name          = "${var.env}_tf_subnet_${count.index}"
    orchestration = "terraform"
  }
}

resource "aws_route_table_association" "tf_public_routes" {
  count          = "${length(compact(split(",", var.private_subnets)))}"
  subnet_id      = "${element(aws_subnet.tf_subnets.*.id, count.index)}"
  route_table_id = "${aws_route_table.tf_public_routes.id}"
}

data "template_file" "worker_init" {
  template = "${file("${path.module}/worker.yml")}"

  vars = {
    master_ip = "${aws_instance.master.private_ip}"
  }
}

resource "aws_security_group" "allow_ping" {
  name   = "Allow ping"
  vpc_id = "${aws_vpc.main.id}"

  ingress {
    # Allow ICMP ECHO:
    # https://www.iana.org/assignments/icmp-parameters/icmp-parameters.xhtml#icmp-parameters-codes-8
    from_port = 8

    to_port     = 0
    protocol    = "icmp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags {
    Name          = "allow_ping"
    orchestration = "terraform"
  }
}

resource "aws_security_group" "allow_https" {
  name   = "Allow HTTPS"
  vpc_id = "${aws_vpc.main.id}"

  # Allow TLS/HTTPS from VPN IP only.
  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["${var.vpn_range}"]
  }

  tags {
    Name          = "allow_https"
    orchestration = "terraform"
  }
}

resource "aws_security_group" "allow_ssh" {
  name   = "Allow ssh"
  vpc_id = "${aws_vpc.main.id}"

  # Allow SSH from VPN IP only.
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["${var.vpn_range}"]
  }

  tags {
    Name          = "allow_ssh"
    orchestration = "terraform"
  }
}

resource "aws_security_group" "allow_internal" {
  name   = "Allow all traffic within subnets"
  vpc_id = "${aws_vpc.main.id}"

  # Allow all traffic from inside subnets.
  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = "${split(",", var.private_subnets)}"
  }

  tags {
    Name          = "allow_internal"
    orchestration = "terraform"
  }
}

# Allow all outbound traffic.
resource "aws_security_group" "allow_outbound" {
  name   = "Allow outbound"
  vpc_id = "${aws_vpc.main.id}"

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags {
    Name          = "allow_outbound"
    orchestration = "terraform"
  }
}

#
# EC2 instances and associated resources.
#

data "template_file" "master_init" {
  template = "${file("${path.module}/master.yml")}"

  vars = {
    service_account_key_pem = "${base64encode(file("${path.module}/.certs/service-account-key.pem"))}"
    k8s_apiserver_pem       = "${base64encode(file("${path.module}/.certs/k8s-apiserver.pem"))}"
    k8s_apiserver_key_pem   = "${base64encode(file("${path.module}/.certs/k8s-apiserver-key.pem"))}"
  }
}

resource "aws_instance" "master" {
  key_name          = "${var.ssh_key}"
  ami               = "${var.master_ami}"
  subnet_id         = "${element(aws_subnet.tf_subnets.*.id, 0)}"
  availability_zone = "eu-west-1a"
  instance_type     = "t2.medium"

  vpc_security_group_ids = [
    "${aws_security_group.allow_ssh.id}",
    "${aws_security_group.allow_https.id}",
    "${aws_security_group.allow_ping.id}",
    "${aws_security_group.allow_internal.id}",
    "${aws_security_group.allow_outbound.id}",
  ]

  tags {
    Name          = "master"
    orchestration = "terraform"
  }

  user_data = "${data.template_file.master_init.rendered}"
}

resource "aws_eip" "master" {
  instance = "${aws_instance.master.id}"
}

resource "aws_instance" "worker1" {
  key_name          = "${var.ssh_key}"
  ami               = "${var.worker_ami}"
  instance_type     = "t2.small"
  availability_zone = "eu-west-1b"
  subnet_id         = "${element(aws_subnet.tf_subnets.*.id, 1)}"

  vpc_security_group_ids = [
    "${aws_security_group.allow_ssh.id}",
    "${aws_security_group.allow_ping.id}",
    "${aws_security_group.allow_internal.id}",
    "${aws_security_group.allow_outbound.id}",
  ]

  tags {
    Name          = "worker_1"
    orchestration = "terraform"
  }

  user_data = "${data.template_file.worker_init.rendered}"
}

resource "aws_instance" "worker2" {
  key_name          = "${var.ssh_key}"
  ami               = "${var.worker_ami}"
  instance_type     = "t2.small"
  availability_zone = "eu-west-1c"
  subnet_id         = "${element(aws_subnet.tf_subnets.*.id, 2)}"

  vpc_security_group_ids = [
    "${aws_security_group.allow_ssh.id}",
    "${aws_security_group.allow_ping.id}",
    "${aws_security_group.allow_internal.id}",
    "${aws_security_group.allow_outbound.id}",
  ]

  tags {
    Name          = "worker_2"
    orchestration = "terraform"
  }

  user_data = "${data.template_file.worker_init.rendered}"
}
