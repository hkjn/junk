variable "gcloud_credentials" {
  default = ".gcp/tf-dns-editor.json"
}

variable "gcloud_project" {
  default = "henrik-jonsson"
}

variable "gcloud_region" {
  default = "europe-west1"
}

variable "gcloud_dns" {
  default = "tf.hkjn.me"
}

variable "ssh_key" {
  default = "hkjn-key-1"
}

variable "region" {
  default = "eu-west-1"
}

variable "profile" {
  default = "default"
}

variable "cidr_block" {
  default = "10.0.40.0/24"
}

variable "env" {
  default = "tf_test"
}

variable vpn_range {
  default = "163.172.150.69/32"
}

variable public_ranges {
  default = "10.0.40.0/26,10.0.40.64/26,10.0.40.128/26"
}

variable azs {
  default = "eu-west-1a,eu-west-1b,eu-west-1c"
}

variable master_ami {
  # Current CoreOS Alpha HVM AMI from https://coreos.com/os/docs/latest/booting-on-ec2.html.  # CoreOS Alpha AMI
  default = "ami-29511f5a"
}

variable worker_ami {
  # Current CoreOS Alpha HVM AMI from https://coreos.com/os/docs/latest/booting-on-ec2.html.
  default = "ami-29511f5a"
}

variable k8s_token {
}