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