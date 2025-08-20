package main

import (
	"context"
	"log"

	satrapv1 "github.com/vayzur/apadana/pkg/api/satrap/v1"

	xray "github.com/vayzur/apadana/pkg/satrap/xray/client"
)

func main() {
	xrayCli, _ := xray.New("127.0.0.1:8080")
	inbound := `
	    {
	        "listen": null,
	        "port": 10800,
	        "protocol": "vless",
	        "settings": {
	            "clients": [
	                {
	                    "id": "apadana"
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
	        "tag": "proxy0",
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
	`

	if err := xrayCli.AddInbound(context.Background(), []byte(inbound)); err != nil {
		log.Printf("failed: %v", err)
	}

	ua := satrapv1.VlessUser{
		BaseUser: satrapv1.BaseUser{
			Email: "xyz@domain.tld",
		},
		ID:   "hello-world",
		Flow: "",
	}

	if err := xrayCli.AddUser(context.Background(), "proxy0", ua); err != nil {
		log.Printf("failed: %v", err)
	}

	if err := xrayCli.RemoveUser(context.Background(), "proxy0", "xyz@domain.tld"); err != nil {
		log.Printf("failed: %v", err)
	}

	// httpCli := httputil.New(time.Second * 2)

	// var conf v1.InboundConfig
	// if err := json.Unmarshal([]byte(inbound), &conf); err != nil {
	// 	panic(err)
	// }

	// conf.Tag = "proxy0"
	// conf.Port = 10900

	// status, resp, err := cli.Do(http.MethodPost, "http://127.0.0.1:10100/api/v1/inbounds", "token", conf)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("Response: %s - Status: %d\n", resp, status)

	// status, resp, err = cli.Do(http.MethodDelete, "http://127.0.0.1:10100/api/v1/inbounds/proxy0", "token", nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("Response: %s - Status: %d\n", resp, status)
}
