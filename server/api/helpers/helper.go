package helpers

import (
	"github.com/corecollectives/mist/api/handlers/dockerdeploy"
)

func SetHelper(d *dockerdeploy.Deployer) *dockerdeploy.Deployer {
	return d
}
