vagrant + coreos example
=====

Adapted from http://github.com/coreos/coreos-vagrant.

Setup
====
0. Clone this repository:
`$ git clone https://github.com/hkjn/junk`

1. Install dependencies:
- VirtualBox 4.3.10 or greater.
- Vagrant 1.6 or greater.

2. (Optional) Edit `config.rb` to set the number of CoreOS instances (currently 1) or any other options.

3. Bring up the Vagrant VMs (using the default VirtualBox provider):
`$ vagrant up`

Usage
====
0. SSH to the VMs using `vagrant ssh [name]`:
`$ vagrant ssh core-01`

1. Clone this repository inside the VM:
`$ git clone https://github.com/hkjn/junk`

2. Start some CoreOS services:
```$ cd junk/coreos/services
$ fleetctl start db-test*.service
$ fleetctl start api-test*.service
$ fleetctl start web-test*.service
```