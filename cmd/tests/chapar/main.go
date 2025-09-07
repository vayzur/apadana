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
			Labels: map[string]string{
				"region":  "EU",
				"country": "germany",
			},
			Annotations: map[string]string{
				"fakeHost": "www.speedtest.net",
				"sni":      "gate.domain.tld",
			},
		},
		Status: corev1.NodeStatus{
			Addresses: []corev1.NodeAddress{
				{
					Type: corev1.InternalAddress,
					Host: "127.0.0.1",
				},
				{
					Type: corev1.ExternalAddress,
					Host: "sub.domain.tld",
				},
			},
		},
		Spec: corev1.NodeSpec{
			Token: "satrap-shared-token",
		},
	}
	node, err := apadanaClient.CreateNode(n)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(node)
}
