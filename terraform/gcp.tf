// Configure the Google Cloud provider
provider "google" {
  credentials = "${file(var.gcloud_credentials)}"
  project     = "${var.gcloud_project}"
  region      = "${var.region}"
}

resource "google_dns_managed_zone" "zone" {
  name     = "tf-zone"
  dns_name = "${var.gcloud_dns}."
}

resource "google_dns_record_set" "master1" {
  name = "m1.${google_dns_managed_zone.zone.dns_name}"
  type = "A"
  ttl  = 150

  managed_zone = "${google_dns_managed_zone.zone.name}"

  rrdatas = ["${aws_eip.master_eip.public_ip}"]
}

resource "google_dns_record_set" "worker1" {
  name = "w1.${google_dns_managed_zone.zone.dns_name}"
  type = "A"
  ttl  = 150

  managed_zone = "${google_dns_managed_zone.zone.name}"
  rrdatas = ["${aws_instance.worker_1.public_ip}"]
}


resource "google_dns_record_set" "worker2" {
  name = "w2.${google_dns_managed_zone.zone.dns_name}"
  type = "A"
  ttl  = 150

  managed_zone = "${google_dns_managed_zone.zone.name}"
  rrdatas = ["${aws_instance.worker_2.public_ip}"]
}


