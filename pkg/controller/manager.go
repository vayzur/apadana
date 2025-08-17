package controller

import "github.com/vayzur/apadana/pkg/service"

type ControllerManager struct {
	nodeService    *service.NodeSerivce
	inboundService *service.InboundService
}

func NewControllerManager(nodeService *service.NodeSerivce, inboundService *service.InboundService) *ControllerManager {
	return &ControllerManager{
		nodeService:    nodeService,
		inboundService: inboundService,
	}
}
