provider "aws" {
  region                   = "${var.region}"
  shared_credentials_file  = "${var.creds_file}"
  profile                  = "${var.profile}"
}

resource "aws_vpc" "tf_vpc" {
  cidr_block = "${var.cidr_block}"
  enable_dns_hostnames = true
  enable_dns_support = true
  tags {
    Name = "${var.env}_vpc"
  }
}

resource "aws_internet_gateway" "tf_igw" {
  vpc_id = "${aws_vpc.tf_vpc.id}"
  tags {
    Name = "${var.env}_igw"
  }
}

resource "aws_route_table" "tf_public_routes" {
  vpc_id = "${aws_vpc.tf_vpc.id}"
  tags {
    Name = "${var.env}_public_subnet_route_table"
  }
}

# add a public gateway to the public route table
resource "aws_route" "tf_public_gateway_route" {
  route_table_id = "${aws_route_table.tf_public_routes.id}"
  depends_on = ["aws_route_table.tf_public_routes"]
  destination_cidr_block = "0.0.0.0/0"
  gateway_id = "${aws_internet_gateway.tf_igw.id}"
}

resource "aws_subnet" "tf_subnets" {
  vpc_id = "${aws_vpc.tf_vpc.id}"
  cidr_block = "${element(split(",", var.public_ranges), count.index)}"
  availability_zone = "${element(split(",", var.azs), count.index)}"
  count = "${length(compact(split(",", var.public_ranges)))}"
  tags {
    Name = "${var.env}_tf_subnet_${count.index}"
  }
  map_public_ip_on_launch = true
}

resource "aws_route_table_association" "tf_public_routes" {
  count = "${length(compact(split(",", var.public_ranges)))}"
  subnet_id = "${element(aws_subnet.tf_subnets.*.id, count.index)}"
  route_table_id = "${aws_route_table.tf_public_routes.id}"
}

resource "aws_security_group" "allow_ping" {
  name = "Allow ping"
  vpc_id = "${aws_vpc.tf_vpc.id}"
  ingress {
    # Allow ICMP ECHO:
    # https://www.iana.org/assignments/icmp-parameters/icmp-parameters.xhtml#icmp-parameters-codes-8
    from_port = 8
    to_port = 0
    protocol = "icmp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  tags {
    Name = "allow_ping"
  }
}

resource "aws_security_group" "allow_ssh" {
  name = "Allow ssh"
  vpc_id = "${aws_vpc.tf_vpc.id}"

  # Allow SSH from VPN IP only.
  ingress {
    from_port = 22
    to_port = 22
    protocol = "tcp"
    cidr_blocks = ["163.172.150.69/32"]
  }

  tags {
    Name = "allow_ssh"
  }
}

resource "aws_security_group" "allow_kube" {
  name = "Allow kubernetes bootstrap traffic"
  vpc_id = "${aws_vpc.tf_vpc.id}"

  # Allow SSH from VPN IP only.
  ingress {
    from_port = 9898
    to_port = 9898
    protocol = "tcp"
    cidr_blocks = "${split(",", var.public_ranges)}"
  }

  tags {
    Name = "allow_kube"
  }
}

# Allow all outbound traffic.
resource "aws_security_group" "allow_outbound" {
  name = "Allow outbound"
  vpc_id = "${aws_vpc.tf_vpc.id}"

  egress {
    from_port = 0
    to_port = 0
    protocol = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
  tags {
    Name = "allow_outbound"
  }
}

data "template_file" "worker_init" {
  template = "${file("${path.module}/templates/worker.sh.tpl")}"

  vars = {
    k8s_token = "${var.k8s_token}"
    k8s_master_ip = "${aws_eip.master_eip.public_ip}"
  }
}

data "template_file" "master_init" {
  template = "${file("${path.module}/templates/master.sh.tpl")}"

  vars = {
    k8s_token = "${var.k8s_token}"
  }
}

resource "aws_instance" "tf_k8s_master" {
  key_name          = "hkjn-key-1"
  ami               = "${var.master_ami}"
  subnet_id         = "${element(aws_subnet.tf_subnets.*.id, 0)}"
  availability_zone = "eu-west-1a"
  instance_type     = "t2.medium"
  vpc_security_group_ids = [
    "${aws_security_group.allow_ssh.id}",
    "${aws_security_group.allow_ping.id}",
    "${aws_security_group.allow_kube.id}",
    "${aws_security_group.allow_outbound.id}",
  ]
  tags {
    Name = "tf_k8s_master"
  }
  user_data       = "${data.template_file.master_init.rendered}"
}

resource "aws_instance" "tf_k8s_worker_1" {
  key_name          = "hkjn-key-1"
  ami               = "${var.worker_ami}"
  instance_type     = "t2.small"
  availability_zone = "eu-west-1b"
  subnet_id         = "${element(aws_subnet.tf_subnets.*.id, 1)}"
  vpc_security_group_ids = [
    "${aws_security_group.allow_ssh.id}",
    "${aws_security_group.allow_ping.id}",
    "${aws_security_group.allow_kube.id}",
    "${aws_security_group.allow_outbound.id}",
  ]
  tags {
    Name = "tf_k8s_worker_1"
  }
  user_data       = "${data.template_file.worker_init.rendered}"
}

resource "aws_instance" "tf_k8s_worker_2" {
  key_name          = "hkjn-key-1"
  ami               = "${var.worker_ami}"
  instance_type     = "t2.small"
  availability_zone = "eu-west-1c"
  subnet_id         = "${element(aws_subnet.tf_subnets.*.id, 2)}"
  vpc_security_group_ids = [
    "${aws_security_group.allow_ssh.id}",
    "${aws_security_group.allow_ping.id}",
    "${aws_security_group.allow_kube.id}",
    "${aws_security_group.allow_outbound.id}",
  ]
  tags {
    Name = "tf_k8s_worker_2"
  }
  user_data       = "${data.template_file.worker_init.rendered}"
}

resource "aws_eip" "master_eip" {
  instance = "${aws_instance.tf_k8s_master.id}"
}
