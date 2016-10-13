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
    cidr_blocks = ["163.172.150.69/32"]
  }

  tags {
    Name = "allow_ssh"
  }
}

resource "aws_security_group" "allow_kube" {
  name = "Allow kubernetes bootstrap traffic"
  vpc_id = "${aws_vpc.tf_vpc.id}"

	# Allow all traffic from inside subnets.
  ingress {
    from_port = 0
    to_port = 0
    protocol = "-1"
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
