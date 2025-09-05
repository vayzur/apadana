package main

import (
	"fmt"
	"log"
	"time"

	corev1 "github.com/vayzur/apadana/pkg/api/core/v1"

	apadana "github.com/vayzur/apadana/pkg/client"
)

func main() {
	apadanaClient := apadana.New("http://127.0.0.1:10200", "cluster-shared-token", time.Second*5)

	n := &corev1.Node{
		Metadata: corev1.NodeMetadata{
			Name: "test",
		},
		Status: corev1.NodeStatus{
			Capacity: corev1.NodeCapacity{
				MaxInbounds: 110,
			},
		},
		Address: "http://127.0.0.1:10100",
		Token:   "satrap-shared-token",
	}
	node, err := apadanaClient.CreateNode(n)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(node)
}
