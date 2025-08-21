package controller

import "github.com/vayzur/apadana/pkg/service"

type Spasaka struct {
	nodeService *service.NodeSerivce
}

func NewSpasaka(nodeService *service.NodeSerivce) *Spasaka {
	return &Spasaka{
		nodeService: nodeService,
	}
}
