#
# Default Terraform variables.
#
variable "creds_file" {
  default = "/.aws/credentials"
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

variable public_ranges {
  default = "10.0.40.0/26,10.0.40.64/26,10.0.40.128/26"
}

variable azs {
  default = "eu-west-1a,eu-west-1b,eu-west-1c"
}

variable master_ami {
  # CoreOS Alpha AMI
  # default = "ami-82b5f5f1"
  # Ubuntu 14.04 LTS AMI
  # default = "ami-ed82e39e"
  # ubuntu/images/hvm-ssd/ubuntu-xenial-16.04-amd64-server-20160721:
  default = "ami-1967056a"
}

variable worker_ami {
  # CoreOS Alpha AMI
  # default = "ami-82b5f5f1"
  # Ubuntu 14.04 LTS AMI
  default = "ami-ed82e39e"
  # ubuntu/images/hvm-ssd/ubuntu-xenial-16.04-amd64-server-20160721:
  default = "ami-1967056a"
}

variable k8s_token {
}