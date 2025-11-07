package queue

import "github.com/corecollectives/mist/api/handlers/dockerdeploy"

func InitQueue(d *dockerdeploy.Deployer) *Queue {
	q := NewQueue(5, d)
	return q
}
