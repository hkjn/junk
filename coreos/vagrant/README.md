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
0. SSH to the VMs using `vagrant ssh [name]`, e.g:

`$ vagrant ssh core-01`

1. Clone this repository inside the VM:

`$ git clone https://github.com/hkjn/junk`

2. Start some CoreOS services:
```
$ cd junk/vagrant/services
$ fleetctl start db-test.service db-test-discovery-vagrant.service
$ fleetctl start api-test.service api-test-discovery-vagrant.service
$ fleetctl start web-test.service web-test-discovery-vagrant.service
```

4. At this point you've launched a DB container (MySQL), an API
container (small Go app that exposes a JSON interface around the DB)
and a web container (HTTP server that shows the results from the API
as HTML).

We can check that things worked with `fleetctl`:
```
$ fleetctl list-units
UNIT                                    MACHINE                         ACTIVE  SUB
api-test-discovery-vagrant.service      5aa6ca00.../172.17.8.101        active  running
api-test.service                        5aa6ca00.../172.17.8.101        active  running
datadog.service                         5aa6ca00.../172.17.8.101        active  running
db-test-discovery-vagrant.service       5aa6ca00.../172.17.8.101        active  running
db-test.service                         5aa6ca00.../172.17.8.101        active  running
web-test-discovery-vagrant.service      5aa6ca00.../172.17.8.101        active  running
web-test.service                        5aa6ca00.../172.17.8.101        active  running
```

We can check that the services that are running registered their address for discovery in etcd:
```
$ etcdctl ls /services --recursive
/services/db
/services/db/test
/services/api
/services/api/test
/services/web
/services/web/test
```

We can query the API layer ourselves, using the value for its address that's registered in etcd:
```
$ curl http://$(etcdctl get /services/api/test)/monkeys
[{"id":1,"name":"Janelle","birthdate":"2009-03-15T14:01:43Z"},{"id":2,"name":"Billy","birthdate":"2008-01-15T11:00:15Z"}]
```

We can check out some logs:
```
$ fleetctl journal web-test
-- Logs begin at Tue 2015-01-20 13:59:56 UTC, end at Tue 2015-01-20 16:12:29 UTC. --
Jan 20 16:09:02 core-01 docker[9203]: 6b8db8be1b58: Pulling metadata
Jan 20 16:09:03 core-01 docker[9203]: 6b8db8be1b58: Pulling fs layer
Jan 20 16:09:06 core-01 docker[9203]: 6b8db8be1b58: Download complete
Jan 20 16:09:06 core-01 docker[9203]: 9e2acfc65780: Pulling metadata
Jan 20 16:09:07 core-01 docker[9203]: 9e2acfc65780: Pulling fs layer
Jan 20 16:09:11 core-01 docker[9203]: 9e2acfc65780: Download complete
Jan 20 16:09:11 core-01 docker[9203]: 9e2acfc65780: Download complete
Jan 20 16:09:11 core-01 docker[9203]: Status: Downloaded newer image for hkjn/coreosweb:latest
Jan 20 16:09:11 core-01 systemd[1]: Started Test web server.
Jan 20 16:09:11 core-01 bash[9276]: [web.v1d8d6af] web layer for stage "test" binding to :9000..
```

TODO
====
- Figure out why hostname resolution isn't working on vagrant VMs; as a
  result `foo-discovery.service` files need a cumbersome sequence of
  commands to figure out the `eth0` IP, so I've forked out
  `foo-discovery-vagrant.service` files since this works fine on
  GCE. Filed https://github.com/coreos/coreos-vagrant/issues/196 to get
  an answer from CoreOS folks. Alternatively, this incantation should
  work to dig out the IP for the appropriate network interface for all
  machine types, if we need to live with this state:
  ```
  $ ip addr list | grep 'inet.*scope global' | head -n 1 | cut -d ' ' -f 6 | cut -d/ -f1
  ```