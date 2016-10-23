#
# Plan for AWS security groups.
#

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
    cidr_blocks = ["${var.vpn_range}"]
  }

  tags {
    Name = "allow_ssh"
  }
}

resource "aws_security_group" "allow_api" {
  name = "Allow traffic to apiserver"
  vpc_id = "${aws_vpc.tf_vpc.id}"

  # Allow apiserver requests from VPN IP only.
  ingress {
    from_port = 443
    to_port = 443
    protocol = "tcp"
    cidr_blocks = ["${var.vpn_range}"]
  }

  tags {
    Name = "allow_api"
  }
}

resource "aws_security_group" "allow_internal" {
  name = "Allow all traffic within subnets"
  vpc_id = "${aws_vpc.tf_vpc.id}"

	# Allow all traffic from inside subnets.
  ingress {
    from_port = 0
    to_port = 0
    protocol = "-1"
    cidr_blocks = "${split(",", var.public_ranges)}"
  }

  tags {
    Name = "allow_internal"
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
