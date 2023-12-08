package imagedistributor

import (
	"github.com/sealerio/sealer/pkg/dfs"
	"github.com/sealerio/sealer/pkg/imageengine"
	v1 "github.com/sealerio/sealer/types/api/v1"
)

type dfsMounter struct {
	deployer    dfs.Deployer
	imageEngine imageengine.Interface
}

func NewDfsMounter(d dfs.Deployer, imageEngine imageengine.Interface) Mounter {
	return &dfsMounter{
		d,
		imageEngine,
	}
}

func (d *dfsMounter) Mount(imageName string, platform v1.Platform, dest string) (string, string, string, error) {
	//TODO implement me
	panic("implement me")
}

func (d *dfsMounter) Umount(dir, containerID string) error {
	//TODO implement me
	panic("implement me")
}
