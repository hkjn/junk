// Configure the Google Cloud provider
provider "google" {
  credentials = "${file(var.google_credentials)}"
  project     = "${var.google_project}"
  region      = "${var.region}"
}

resource "google_dns_managed_zone" "testzone" {
  name     = "tf-testzone"
  dns_name = "tf.hkjn.me."
}

resource "google_dns_record_set" "master" {
  name = "m1.${google_dns_managed_zone.testzone.dns_name}"
  type = "A"
  ttl  = 300

  managed_zone = "${google_dns_managed_zone.testzone.name}"

  rrdatas = ["${aws_eip.master_eip.public_ip}"]
}


