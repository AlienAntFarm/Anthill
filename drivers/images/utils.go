package images

import (
	"github.com/golang/glog"
	"os"
)

func RemoveOnFail(archive string, err error) {
	if err == nil {
		return
	}
	glog.Infof("removing archive %s due to %s", archive, err)
	// if remove fails, log the error but don't alter err
	if err := os.Remove(archive); err != nil {
		glog.Errorf("%s", err)
	}
}
