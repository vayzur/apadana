package main

import (
	"fmt"
	"log"
	"time"

	corev1 "github.com/vayzur/apadana/pkg/apis/core/v1"
	metav1 "github.com/vayzur/apadana/pkg/apis/meta/v1"

	apadana "github.com/vayzur/apadana/pkg/client"
)

func main() {
	clusterAddress := "https://sub.domain.tld:10200"
	clusterToken := "cluster-shared-token"
	apadanaClient := apadana.New(clusterAddress, clusterToken, time.Second*5)

	n := &corev1.Node{
		Metadata: metav1.ObjectMeta{
			Name: "test",
			Labels: map[string]string{
				"name":     "my-first-apadana-node",
				"region":   "eu",
				"country":  "france",
				"provider": "ovh",
				"purpose":  "gaming",
			},
			Annotations: map[string]string{
				"sni": "gate.domain.tld",
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
