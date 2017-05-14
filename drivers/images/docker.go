package images

import (
	"archive/tar"
	"bytes"
	"github.com/alienantfarm/anthive/utils"
	"github.com/golang/glog"
	"os"
	"os/exec"
	"path"
)

func Docker2AIF(tag string) (archive string, err error) {
	var out bytes.Buffer

	cmd := exec.Command("docker", "create", tag, "sh")
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return
	}
	id := out.String()[:10] // only the 10 first digit of a docker id are needed

	defer exec.Command("docker", "rm", id).Run() // remove container at the end

	archive = path.Join(utils.Config.Assets.Images, utils.SecureRandomString(10))
	glog.Infof("Generate new image %s from docker image %s", archive, tag)
	cmd = exec.Command("docker", "export", "-o", archive, id)
	err = cmd.Run()
	if err != nil {
		return
	}
	defer func() {
		if err == nil {
			return
		}
		// if remove fails, log the error but don't alter err
		if err := os.Remove(archive); err != nil {
			glog.Errorf("%s", err)
		}
	}()

	// retrieve the running config of the container, to save it in the AIF tarball
	out.Reset()
	cmd = exec.Command("docker", "inspect", "--format='{{json .Config}}'", id)
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return
	}
	// open archive and append the manifest
	var f *os.File
	f, err = os.OpenFile(archive, os.O_RDWR, 0644)
	if err != nil {
		return
	}
	manifest := out.Bytes()
	manifest = manifest[1 : len(manifest)-2] // remove quote at beginning and end plus line break
	tw := tar.NewWriter(f)
	err = tw.WriteHeader(&tar.Header{
		Name: "manifest.json", // TODO: put this outside
		Mode: 0644,
		Size: int64(len(manifest)),
	})
	_, err = tw.Write(manifest)
	if err != nil {
		return
	}

	return
}
