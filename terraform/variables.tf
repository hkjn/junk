variable "k8s_token" {}

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

variable "vpn_range" {
  default = "212.47.239.127/32"
}

variable "private_subnets" {
  default = "10.0.40.0/26,10.0.40.64/26,10.0.40.128/26"
}

variable "azs" {
  default = "eu-west-1a,eu-west-1b,eu-west-1c"
}

variable "master_ami" {
  # Current CoreOS Alpha HVM AMI from https://coreos.com/dist/aws/aws-alpha.json.
  default = "ami-a09caac6"
}

variable "worker_ami" {
  # Current CoreOS Alpha HVM AMI from https://coreos.com/dist/aws/aws-alpha.json.
  default = "ami-a09caac6"
}
