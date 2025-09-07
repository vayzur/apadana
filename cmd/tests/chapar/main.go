package main

import (
	"fmt"
	"log"
	"time"

	corev1 "github.com/vayzur/apadana/pkg/api/core/v1"

	apadana "github.com/vayzur/apadana/pkg/client"
)

func main() {
	clusterAddress := "https://sub.domain.tld:10200"
	clusterToken := "cluster-shared-token"
	apadanaClient := apadana.New(clusterAddress, clusterToken, time.Second*5)

	n := &corev1.Node{
		Metadata: corev1.NodeMetadata{
			Name: "test",
			Labels: map[string]string{
				"region":   "eu",
				"country":  "france",
				"provider": "ovh",
				"purpose":  "gaming",
			},
			Annotations: map[string]string{
				"fakeHost": "www.speedtest.net",
				"sni":      "gate.domain.tld",
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
