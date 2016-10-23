// Configure the Google Cloud provider
provider "google" {
  credentials = "${file(var.gcloud_credentials)}"
  project     = "${var.gcloud_project}"
  region      = "${var.region}"
}

resource "google_dns_managed_zone" "testzone" {
  name     = "tf-testzone"
  dns_name = "${var.gcloud_dns_zone}"
}

resource "google_dns_record_set" "master" {
  name = "m1.${google_dns_managed_zone.testzone.dns_name}"
  type = "A"
  ttl  = 150

  managed_zone = "${google_dns_managed_zone.testzone.name}"

  rrdatas = ["${aws_eip.master_eip.public_ip}"]
}


