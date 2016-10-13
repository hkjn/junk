output "master_ip" {
  value = "${aws_eip.master_eip.public_ip}"
}

output "worker1_ip" {
  value = "${aws_instance.tf_k8s_worker_1.public_ip}"
}

output "worker2_ip" {
  value = "${aws_instance.tf_k8s_worker_2.public_ip}"
}