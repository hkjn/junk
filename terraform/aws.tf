provider "aws" {
  region = "${var.region}"
}

data "template_file" "worker_init" {
  template = "${file("${path.module}/worker.yml")}"

  vars = {
    master_ip = "${aws_instance.master_1.private_ip}"
  }
}

data "template_file" "master_init" {
  template = "${file("${path.module}/master.yml")}"

  vars = {
    service_account_key_pem = "${base64encode(file("${path.module}/.certs/service-account-key.pem"))}"
    k8s_apiserver_pem = "${base64encode(file("${path.module}/.certs/k8s-apiserver.pem"))}"
    k8s_apiserver_key_pem = "${base64encode(file("${path.module}/.certs/k8s-apiserver-key.pem"))}"
  }
}

resource "aws_instance" "master_1" {
  key_name          = "${var.ssh_key}"
  ami               = "${var.master_ami}"
  subnet_id         = "${element(aws_subnet.tf_subnets.*.id, 0)}"
  availability_zone = "eu-west-1a"
  instance_type     = "t2.medium"

  vpc_security_group_ids = [
    "${aws_security_group.allow_ssh.id}",
    "${aws_security_group.allow_tls.id}",
    "${aws_security_group.allow_ping.id}",
    "${aws_security_group.allow_internal.id}",
    "${aws_security_group.allow_outbound.id}",
  ]

  tags {
    Name          = "master_1"
    orchestration = "terraform"
  }

  user_data = "${data.template_file.master_init.rendered}"
}

resource "aws_eip" "master_eip" {
  instance = "${aws_instance.master_1.id}"
}

resource "aws_instance" "worker_1" {
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

resource "aws_instance" "worker_2" {
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
