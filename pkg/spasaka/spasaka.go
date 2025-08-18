package spasaka

import "github.com/vayzur/apadana/pkg/service"

type Spasaka struct {
	nodeService    *service.NodeSerivce
	inboundService *service.InboundService
}

func NewSpasaka(nodeService *service.NodeSerivce, inboundService *service.InboundService) *Spasaka {
	return &Spasaka{
		nodeService:    nodeService,
		inboundService: inboundService,
	}
}
