package controller

import (
	apadana "github.com/vayzur/apadana/pkg/client"
)

type Spasaka struct {
	apadanaClient *apadana.Client
}

func NewSpasaka(apadanaClient *apadana.Client) *Spasaka {
	return &Spasaka{
		apadanaClient: apadanaClient,
	}
}
