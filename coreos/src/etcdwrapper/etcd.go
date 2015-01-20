// Package etcdwrapper provides some shared tools for interacting with etcd.
package etcdwrapper

import (
	"fmt"

	"github.com/coreos/go-etcd/etcd"
	"github.com/golang/glog"
)

var (
	etcdPeers = []string{
		"http://172.17.42.1:4001", // on GCE / most others
		"http://10.1.42.1:4001",   // on Vagrant
	}
	client *etcd.Client
)

// Read returns the simple string value at path from etcd.
func Read(path string) (string, error) {
	if client == nil {
		client = etcd.NewClient(etcdPeers)
	}
	r, err := client.Get(path, false, false)
	if err != nil {
		return "", fmt.Errorf("failed to read etcd path %s from peers %v: %v", path, etcdPeers, err)
	}
	v := r.Node.Value
	glog.V(2).Infof("read value %q from %s\n", v, path)
	return v, nil
}
