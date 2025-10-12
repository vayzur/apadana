package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	satrapv1 "github.com/vayzur/apadana/pkg/apis/satrap/v1"
	"github.com/xtls/xray-core/infra/conf"

	xray "github.com/vayzur/apadana/pkg/satrap/xray/client"
	xrayconfigv1 "github.com/vayzur/apadana/pkg/satrap/xray/config/v1"
)

func main() {
	xrayConf := xrayconfigv1.XrayConfig{
		Address:               "127.0.0.1",
		Port:                  8080,
		RuntimeRequestTimeout: time.Second * 5,
	}
	xrayCli, _ := xray.New(&xrayConf)

	listenAddr := "0.0.0.0"
	port := 10800
	tag := "proxy0"
	clientID := "testing"

	inbound := fmt.Sprintf(`
	    {
	        "listen": "%s",
	        "port": %d,
	        "protocol": "vless",
	        "settings": {
	            "clients": [
	                {
	                    "id": "%s"
	                }
	            ],
	            "decryption": "none",
	            "fallbacks": []
	        },
	        "streamSettings": {
	            "network": "ws",
	            "security": "none",
	            "wsSettings": {
	                "acceptProxyProtocol": false,
	                "headers": {},
	                "heartbeatPeriod": 0,
	                "host": "",
	                "path": ""
	            },
	            "sockopt": {
	                "tcpFastOpen": true,
	                "tcpCongestion": "bbr",
	                "tcpMptcp": true,
	                "tcpNoDelay": true
	            }
	        },
	        "tag": "%s",
	        "sniffing": {
	            "enabled": false,
	            "destOverride": [
	                "http",
	                "tls",
	                "quic",
	                "fakedns"
	            ],
	            "metadataOnly": false,
	            "routeOnly": false
	        },
	        "allocate": {
	            "strategy": "always",
	            "refresh": 5,
	            "concurrency": 3
	        }
	    }
	`, listenAddr, port, clientID, tag)

	conf := &conf.InboundDetourConfig{}
	if err := json.Unmarshal([]byte(inbound), conf); err != nil {
		log.Printf("unmarshal failed: %v", err)
	}

	if err := xrayCli.AddInbound(context.Background(), conf); err != nil {
		log.Printf("add inbound failed: %v", err)
	}

	email := "user@domain.tld"
	a := satrapv1.VlessAccount{
		ID:   "hello-world",
		Flow: "",
	}

	if err := xrayCli.AddUser(context.Background(), tag, email, a); err != nil {
		log.Printf("add user failed: %v", err)
	}

	if err := xrayCli.RemoveUser(context.Background(), tag, email); err != nil {
		log.Printf("remove user failed: %v", err)
	}
}
