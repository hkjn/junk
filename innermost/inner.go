// Package inner imports some external packages.
package innermost

import "github.com/golang/glog"

func Inner() {
	glog.Infof("Inner: Woo!")
}
