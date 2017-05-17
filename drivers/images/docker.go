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

func createImage(tag string) (string, error) {
	var out bytes.Buffer

	cmd := exec.Command("docker", "create", tag, "sh")
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", err
	} else {
		// only the 10 first digit of a docker id are needed
		return out.String()[:10], nil
	}
}

func createArchive(id string) (string, error) {
	archive := path.Join(utils.Config.Assets.Images, utils.SecureRandomString(10))
	return archive, exec.Command("docker", "export", "-o", archive, id).Run()
}

func getManifest(id string) ([]byte, error) {
	var out bytes.Buffer

	// retrieve the running config of the container, to save it in the AIF tarball
	cmd := exec.Command("docker", "inspect", "--format='{{json .Config}}'", id)
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	manifest := out.Bytes()
	return manifest[1 : len(manifest)-2], nil // remove quote at beginning and end plus line break

}

func appendManifest2Archive(manifest []byte, archive string) error {
	// open archive and append the manifest
	f, err := os.OpenFile(archive, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	tw := tar.NewWriter(f)
	hdr := &tar.Header{
		Name: "manifest.json", // TODO: put this outside
		Mode: 0644,
		Size: int64(len(manifest)),
	}
	if err = tw.WriteHeader(hdr); err != nil {
		return err
	}
	_, err = tw.Write(manifest)
	return err
}

func Docker2AIF(tag string) (archive string, err error) {
	var (
		id       string // id of the docker image used in this function
		manifest []byte
	)

	if id, err = createImage(tag); err != nil {
		return
	}
	defer exec.Command("docker", "rm", id).Run() // remove container at the end

	if archive, err = createArchive(id); err != nil {
		return
	}
	glog.Infof("Generate new image %s from docker image %s", archive, tag)

	// remove the archive, to avoid creating unfinished instances
	defer func() {
		if err == nil {
			return
		}
		// if remove fails, log the error but don't alter err
		if err := os.Remove(archive); err != nil {
			glog.Errorf("%s", err)
		}
	}()

	if manifest, err = getManifest(id); err != nil {
		return
	}
	err = appendManifest2Archive(manifest, archive)
	return
}