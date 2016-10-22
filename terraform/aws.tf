provider "aws" {
  region                   = "${var.region}"
  shared_credentials_file  = "${var.creds_file}"
  profile                  = "${var.profile}"
}

data "template_file" "worker_init" {
  template = "${file("${path.module}/templates/worker.sh.tpl")}"

  vars = {
    k8s_token = "${var.k8s_token}"
    k8s_master_ip = "${aws_instance.tf_k8s_master.private_ip}"
  }
}

data "template_file" "master_init" {
  template = "${file("${path.module}/templates/master.sh.tpl")}"

  vars = {
    k8s_token = "${var.k8s_token}"
  }
}

resource "aws_instance" "tf_k8s_master" {
  key_name          = "${var.ssh_key}"
  ami               = "${var.master_ami}"
  subnet_id         = "${element(aws_subnet.tf_subnets.*.id, 0)}"
  availability_zone = "eu-west-1a"
  instance_type     = "t2.medium"
  vpc_security_group_ids = [
    "${aws_security_group.allow_ssh.id}",
    "${aws_security_group.allow_ping.id}",
    "${aws_security_group.allow_internal.id}",
    "${aws_security_group.allow_outbound.id}",
    "${aws_security_group.allow_api.id}",
  ]
  tags {
    Name = "tf_k8s_master"
  }
  user_data       = "${data.template_file.master_init.rendered}"
}

resource "aws_eip" "master_eip" {
  instance = "${aws_instance.tf_k8s_master.id}"
}

resource "aws_instance" "tf_k8s_worker_1" {
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
    Name = "tf_k8s_worker_1"
  }
  user_data       = "${data.template_file.worker_init.rendered}"
}

resource "aws_instance" "tf_k8s_worker_2" {
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
    Name = "tf_k8s_worker_2"
  }
  user_data       = "${data.template_file.worker_init.rendered}"
}