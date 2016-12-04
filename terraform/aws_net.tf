#

# Plan for AWS network resources.

#

resource "aws_vpc" "tf_vpc" {
  cidr_block           = "${var.cidr_block}"
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags {
    Name          = "${var.env}_vpc"
    orchestration = "terraform"
  }
}

resource "aws_internet_gateway" "tf_igw" {
  vpc_id = "${aws_vpc.tf_vpc.id}"

  tags {
    Name          = "${var.env}_igw"
    orchestration = "terraform"
  }
}

resource "aws_route_table" "tf_public_routes" {
  vpc_id = "${aws_vpc.tf_vpc.id}"

  tags {
    Name          = "${var.env}_public_subnet_route_table"
    orchestration = "terraform"
  }
}

# add route to public gateway
resource "aws_route" "tf_public_gateway_route" {
  route_table_id         = "${aws_route_table.tf_public_routes.id}"
  depends_on             = ["aws_route_table.tf_public_routes"]
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = "${aws_internet_gateway.tf_igw.id}"
}

resource "aws_subnet" "tf_subnets" {
  vpc_id            = "${aws_vpc.tf_vpc.id}"
  cidr_block        = "${element(split(",", var.public_ranges), count.index)}"
  availability_zone = "${element(split(",", var.azs), count.index)}"
  count             = "${length(compact(split(",", var.public_ranges)))}"

  tags {
    Name          = "${var.env}_tf_subnet_${count.index}"
    orchestration = "terraform"
  }

  map_public_ip_on_launch = true
}

resource "aws_route_table_association" "tf_public_routes" {
  count          = "${length(compact(split(",", var.public_ranges)))}"
  subnet_id      = "${element(aws_subnet.tf_subnets.*.id, count.index)}"
  route_table_id = "${aws_route_table.tf_public_routes.id}"
}
