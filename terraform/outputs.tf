output "master_ip" {
  value = "${aws_eip.master_eip.public_ip}"
}

output "worker1_ip" {
  value = "${aws_instance.worker_1.public_ip}"
}

output "worker2_ip" {
  value = "${aws_instance.worker_2.public_ip}"
}

output "master1_dns" {
  value = "m1.${var.gcloud_dns}"
}

output "worker1_dns" {
  value = "w1.${var.gcloud_dns}"
}

output "worker2_dns" {
  value = "w2.${var.gcloud_dns}"
}
