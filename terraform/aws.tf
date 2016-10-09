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

resource "aws_subnet" "tf_subnet_1" {
  #vpc_id = "vpc-1f8bab7a"
  vpc_id = "${aws_vpc.tf_vpc.id}"
  cidr_block = "10.0.40.0/26"
  availability_zone = "eu-west-1a"
  tags {
    Name = "tf_subnet_1"
  }
  map_public_ip_on_launch = true
}

resource "aws_subnet" "tf_subnet_2" {
  vpc_id = "${aws_vpc.tf_vpc.id}"
  cidr_block = "10.0.40.65/26"
  availability_zone = "eu-west-1b"
  tags {
    Name = "tf_subnet_2"
  }
  map_public_ip_on_launch = true
}

resource "aws_subnet" "tf_subnet_3" {
  vpc_id = "${aws_vpc.tf_vpc.id}"
  cidr_block = "10.0.40.130/26"
  availability_zone = "eu-west-1c"
  tags {
    Name = "tf_subnet_3"
  }
  map_public_ip_on_launch = true
}

resource "aws_route_table_association" "tf_associate_route_1" {
  subnet_id = "${aws_subnet.tf_subnet_1.id}"
  route_table_id = "${aws_route_table.tf_public_routes.id}"
  #"rtb-c4bbb5a1"
}

resource "aws_route_table_association" "tf_associate_route_2" {
  subnet_id = "${aws_subnet.tf_subnet_2.id}"
  route_table_id = "${aws_route_table.tf_public_routes.id}"
}

resource "aws_route_table_association" "tf_associate_route_3" {
  subnet_id = "${aws_subnet.tf_subnet_3.id}"
  route_table_id = "${aws_route_table.tf_public_routes.id}"
}

#resource "aws_route" "tf_route_1" {
#  route_table_id = "rtb-c4bbb5a1"
#  # destination_cidr_block = "10.0.20.0/24"
#  destination_cidr_block = "0.0.0.0/0"
#  gateway_id = "igw-70355415"
#}

resource "aws_security_group" "allow_ping" {
  name = "Allow ping"
  vpc_id = "${aws_vpc.tf_vpc.id}"
  # vpc_id = "vpc-1f8bab7a"
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

resource "aws_eip" "ip" {
  instance = "${aws_instance.tf-core-1.id}"
}

resource "aws_instance" "tf-core-1" {
  ami               = "ami-82b5f5f1"
  availability_zone = "eu-west-1a"
  instance_type     = "t2.micro"
  subnet_id = "${aws_subnet.tf_subnet_1.id}"
  vpc_security_group_ids = [
    "${aws_security_group.allow_ssh.id}",
    "${aws_security_group.allow_ping.id}",
    "${aws_security_group.allow_outbound.id}",
  ]
  key_name = "hkjn-key-1"
  tags {
    Name = "tf-core-1"
  }
  # iam_instance_profile = ..
  # user_data = ""
}

resource "aws_instance" "tf-core-2" {
  ami               = "ami-82b5f5f1"
  availability_zone = "eu-west-1b"
  instance_type     = "t2.micro"
  subnet_id = "${aws_subnet.tf_subnet_2.id}"
  vpc_security_group_ids = [
    "${aws_security_group.allow_ssh.id}",
    "${aws_security_group.allow_ping.id}",
    "${aws_security_group.allow_outbound.id}",
  ]
  key_name = "hkjn-key-1"
  tags {
    Name = "tf-core-2"
  }
  associate_public_ip_address = true
}

#resource "aws_instance" "tf-core-3" {
#  ami               = "ami-82b5f5f1"
#  availability_zone = "eu-west-1c"
#  instance_type     = "t2.tiny"
#  tags {
#    Name = "tf-core-3"
#  }
#}