package images

import (
	"bytes"
	"github.com/alienantfarm/anthive/utils"
	"github.com/alienantfarm/anthive/utils/structs"
	"github.com/golang/glog"
	"os/exec"
	"path"
	"path/filepath"
)

type DockerManifest struct {
	Hostname   string
	User       string
	Env        []string
	Cmd        []string
	WorkingDir string
}

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

func getManifest(id string) (*DockerManifest, error) {
	var out bytes.Buffer
	var dm = &DockerManifest{}

	// retrieve the running config of the container, to save it in the AIF tarball
	cmd := exec.Command("docker", "inspect", "--format='{{json .Config}}'", id)
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	if err := utils.UnmarshalJSON(out.Bytes()[1:out.Len()-2], dm); err != nil {
		return nil, err
	}
	return dm, nil
}

func Docker2AIF(tag string) (image *structs.Image, err error) {
	var (
		id       string
		archive  string
		manifest *DockerManifest
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
	defer func() { RemoveOnFail(archive, err) }()

	if manifest, err = getManifest(id); err != nil {
		return
	}
	image = &structs.Image{
		Archive:  filepath.Base(archive),
		Hostname: manifest.Hostname,
		Cmd:      manifest.Cmd,
		Cwd:      manifest.WorkingDir,
		Env:      manifest.Env,
	}
	return
}
