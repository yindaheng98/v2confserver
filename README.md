# v2confserver
根据vmess订阅链接自动生成v2ray配置文件，且定时更新；v2ray通过HTTP读取之

## 原理

1. 输入你的订阅地址
2. vmessconfig获取订阅地址上的内容
3. vmessping解析出`vmess://...`链接列表
4. vmessping对这些`vmess://...`链接的质量进行测量
5. vmessconfig根据质量进行排序
6. vmessconfig根据模板生成配置文件
7. v2confserver用http.Server接口让外部访问配置文件
8. 在v2ray[命令行参数](https://www.v2fly.org/guide/command.html#v2ray)中指定`v2ray -c http://...`以访问配置文件

## 用法

### 单节点模式

* 我的配置文件里面只需要一个服务器节点
* 我的vmess订阅地址是`https://my.subscribe.xyz/vmess`，里面有很多个服务器
* 我想每10分钟测量一次订阅中的所有服务器
* 让v2ray通过8080端口访问做好的配置文件
* 我的模板在`/v2ray/template.json`路径下
* 把选出的服务器outbound放在标签为`direct`的outbound前面

那么v2confserver可以这样用：
```
v2confserver single -urls "https://my.subscribe.xyz/vmess" -interval 600 -addr ":8080" -template-config-from "/v2ray/template.json" -outbound-insert-before-tag "direct"
```
于是，`v2confserver`会每10分钟测量一次订阅中的所有服务器，然后选其中通信质量最好的，根据模板生成配置文件，你可以从`localhost:8080`端口访问到这个配置文件。

比如，如果`/v2ray/template.json`路径下的配置文件长这样：
```json
{
  "inbounds": [
    {
      "port": 80,
      "listen": "0.0.0.0",
      "protocol": "http",
      "settings": {
        "udp": true
      }
    }
  ],
  "outbounds": [
    {
      "tag": "direct",
      "protocol": "freedom",
      "settings": {}
    }
  ]
}
```
那么最好的那个outbound的配置将会放在这个标签为`direct`的outbound前面，即生成一个大概长这样的配置文件：
```json
{
  "inbounds": [
    {
      "port": 80,
      "listen": "0.0.0.0",
      "protocol": "http",
      "settings": {
        "udp": true
      }
    }
  ],
  "outbounds": [
    {
      "protocol": "vmess",
      "sendThrough": null,
      "tag": "",
      "settings": {
        "vnext": [
          {
            "address": "xx.xxx.xx",
            "port": 12345,
            "users": [
              {
                "id": "X0XX0XXX-XXX0-XXXX-X0XX-0XXXXXXX00XX",
                "alterId": 0,
                "security": "auto"
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "ws",
        "security": "",
        "tlsSettings": null,
        "tcpSettings": null,
        "kcpSettings": null,
        "wsSettings": {
          "path": "/0X0X00X",
          "headers": {
            "Host": ""
          },
          "acceptProxyProtocol": false,
          "maxEarlyData": 0,
          "useBrowserForwarding": false,
          "earlyDataHeaderName": ""
        },
        "httpSettings": null,
        "dsSettings": null,
        "quicSettings": null,
        "gunSettings": null,
        "grpcSettings": null,
        "sockopt": null
      },
      "proxySettings": null,
      "mux": {
        "enabled": false,
        "concurrency": 0
      }
    },
    {
      "tag": "direct",
      "protocol": "freedom",
      "settings": {}
    }
  ]
}
```

### 负载均衡模式

* 我的配置文件里面需要4个服务器节点做负载均衡
* 我的vmess订阅地址是`https://my.subscribe.xyz/vmess`，里面有很多个服务器
* 我想每10分钟测量一次订阅中的所有服务器
* 让v2ray通过8080端口访问做好的配置文件
* 我的模板在`/v2ray/template.json`路径下
* 把选出的服务器outbound放在标签为`direct`的outbound前面
* outbound的标签按照`my-bl-%d`的格式生成
* 生成好的outbound标签列表放进标签为`my-bl`的balancer里面

那么v2confserver可以这样用：
```
v2confserver balancer -max-select 4 -urls "https://my.subscribe.xyz/vmess" -interval 600 -addr ":8080" -template-config-from "/v2ray/template.json" -outbound-insert-before-tag "direct" -tag-format "my-bl-%d" -balancer-insert-to-tag "my-bl"
```
于是，`v2confserver`会每10分钟测量一次订阅中的所有服务器，然后选其中通信质量最好的4个，根据模板生成配置文件，你可以从`localhost:8080`端口访问到这个配置文件。

比如，如果`/v2ray/template.json`路径下的配置文件长这样：
```json
{
  "inbounds": [
    {
      "port": 80,
      "listen": "0.0.0.0",
      "protocol": "http",
      "settings": {
        "udp": true
      }
    }
  ],
  "outbounds": [
    {
      "tag": "direct",
      "protocol": "freedom",
      "settings": {}
    }
  ],
  "routing": {
    "domainStrategy": "IPIfNonMatch",
    "rules": [
      {
        "type": "field",
        "balancerTag": "my-bl",
        "domainStrategy": "IPOnDemand",
        "ip": [
          "0.0.0.0/0"
        ]
      },
      {
        "type": "field",
        "outboundTag": "direct",
        "domainStrategy": "IPOnDemand",
        "ip": [
          "geoip:private",
          "geoip:cn"
        ]
      },
      {
        "type": "field",
        "outboundTag": "direct",
        "domain": [
          "geosite:private",
          "geosite:cn"
        ]
      }
    ],
    "balancers": [
      {
        "tag": "my-bl",
        "selector": [],
        "strategy": {
          "type": "random"
        }
      }
    ]
  }
}
```
那么最好的那个outbound的配置将会放在这个标签为`direct`的outbound前面，并且他们的标签将会被放入标签为`my-bl`的balancer的selector里面，即生成一个大概长这样的配置文件：
```json
{
  "inbounds": [
    {
      "port": 1080,
      "listen": "0.0.0.0",
      "protocol": "socks",
      "sniffing": {
        "enabled": true,
        "destOverride": [
          "http",
          "tls"
        ]
      },
      "settings": {
        "auth": "noauth",
        "udp": true
      }
    },
    {
      "port": 80,
      "listen": "0.0.0.0",
      "protocol": "http",
      "settings": {
        "udp": true
      }
    }
  ],
  "log": {
    "loglevel": "warning"
  },
  "outbounds": [
    {
      "protocol": "vmess",
      "sendThrough": null,
      "tag": "my-bl-0",
      "settings": {
        "vnext": [
          {
            "address": "xx.xxx.xx",
            "port": 12345,
            "users": [
              {
                "id": "X0XX0XXX-XXX0-XXXX-X0XX-0XXXXXXX00XX",
                "alterId": 0,
                "security": "auto"
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "ws",
        "security": "",
        "tlsSettings": null,
        "tcpSettings": null,
        "kcpSettings": null,
        "wsSettings": {
          "path": "/0X0X00X",
          "headers": {
            "Host": ""
          },
          "acceptProxyProtocol": false,
          "maxEarlyData": 0,
          "useBrowserForwarding": false,
          "earlyDataHeaderName": ""
        },
        "httpSettings": null,
        "dsSettings": null,
        "quicSettings": null,
        "gunSettings": null,
        "grpcSettings": null,
        "sockopt": null
      },
      "proxySettings": null,
      "mux": {
        "enabled": false,
        "concurrency": 0
      }
    },
    {
      "protocol": "vmess",
      "sendThrough": null,
      "tag": "my-bl-1",
      "settings": {
        "vnext": [
          {
            "address": "xx.xxx.xx",
            "port": 12345,
            "users": [
              {
                "id": "X0XX0XXX-XXX0-XXXX-X0XX-0XXXXXXX00XX",
                "alterId": 0,
                "security": "auto"
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "ws",
        "security": "",
        "tlsSettings": null,
        "tcpSettings": null,
        "kcpSettings": null,
        "wsSettings": {
          "path": "/0X0X00X",
          "headers": {
            "Host": ""
          },
          "acceptProxyProtocol": false,
          "maxEarlyData": 0,
          "useBrowserForwarding": false,
          "earlyDataHeaderName": ""
        },
        "httpSettings": null,
        "dsSettings": null,
        "quicSettings": null,
        "gunSettings": null,
        "grpcSettings": null,
        "sockopt": null
      },
      "proxySettings": null,
      "mux": {
        "enabled": false,
        "concurrency": 0
      }
    },
    {
      "protocol": "vmess",
      "sendThrough": null,
      "tag": "my-bl-2",
      "settings": {
        "vnext": [
          {
            "address": "xx.xxx.xx",
            "port": 12345,
            "users": [
              {
                "id": "X0XX0XXX-XXX0-XXXX-X0XX-0XXXXXXX00XX",
                "alterId": 0,
                "security": "auto"
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "ws",
        "security": "",
        "tlsSettings": null,
        "tcpSettings": null,
        "kcpSettings": null,
        "wsSettings": {
          "path": "/0X0X00X",
          "headers": {
            "Host": ""
          },
          "acceptProxyProtocol": false,
          "maxEarlyData": 0,
          "useBrowserForwarding": false,
          "earlyDataHeaderName": ""
        },
        "httpSettings": null,
        "dsSettings": null,
        "quicSettings": null,
        "gunSettings": null,
        "grpcSettings": null,
        "sockopt": null
      },
      "proxySettings": null,
      "mux": {
        "enabled": false,
        "concurrency": 0
      }
    },
    {
      "protocol": "vmess",
      "sendThrough": null,
      "tag": "my-bl-3",
      "settings": {
        "vnext": [
          {
            "address": "xx.xxx.xx",
            "port": 12345,
            "users": [
              {
                "id": "X0XX0XXX-XXX0-XXXX-X0XX-0XXXXXXX00XX",
                "alterId": 0,
                "security": "auto"
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "ws",
        "security": "",
        "tlsSettings": null,
        "tcpSettings": null,
        "kcpSettings": null,
        "wsSettings": {
          "path": "/0X0X00X",
          "headers": {
            "Host": ""
          },
          "acceptProxyProtocol": false,
          "maxEarlyData": 0,
          "useBrowserForwarding": false,
          "earlyDataHeaderName": ""
        },
        "httpSettings": null,
        "dsSettings": null,
        "quicSettings": null,
        "gunSettings": null,
        "grpcSettings": null,
        "sockopt": null
      },
      "proxySettings": null,
      "mux": {
        "enabled": false,
        "concurrency": 0
      }
    },
    {
      "tag": "direct",
      "protocol": "freedom",
      "settings": {}
    }
  ],
  "routing": {
    "balancers": [
      {
        "selector": [
          "my-bl-0",
          "my-bl-1",
          "my-bl-2",
          "my-bl-3"
        ],
        "strategy": {
          "type": "random"
        },
        "tag": "my-bl-"
      }
    ],
    "domainStrategy": "IPIfNonMatch",
    "rules": [
      {
        "type": "field",
        "balancerTag": "my-bl",
        "domainStrategy": "IPOnDemand",
        "ip": [
          "0.0.0.0/0"
        ]
      },
      {
        "type": "field",
        "outboundTag": "direct",
        "domainStrategy": "IPOnDemand",
        "ip": [
          "geoip:private",
          "geoip:cn"
        ]
      },
      {
        "type": "field",
        "outboundTag": "direct",
        "domain": [
          "geosite:private",
          "geosite:cn"
        ]
      }
    ]
  }
}
```

### 所有选项
```
v2confserver -h
INVALID ARGS: [-h]
v2confserver [balancer|single] -urls https://... -urls https://...
Usage of balancer:
  -addr value
        Address where the server listen on (default :80)
  -balancer-insert-to-tag value
        Insert the selector into the balancer whose tag is this (default vmessconfig-autogenerated-balancer)
  -interval value
        Interval for get and ping vmess outbounds (default 1800)
  -max-select value
        How many outbounds do you want to put into (default 8)
  -outbound-insert-before-tag value
        Insert outbound before the exists outbound whose tag is this (default vmessconfig-outbound-insert)
  -ping-config-allow-insecure
        allow insecure TLS connections (default false)
  -ping-config-count value
        Count. Stop after sending COUNT requests (default 4)
  -ping-config-dest value
        the test destination url, need 204 for success return (default http://www.google.com/gen_204)
  -ping-config-inteval value
        inteval seconds between pings (default 1)
  -ping-config-quit value
        fast quit on error counts (default 0)
  -ping-config-show-node
        show node location/outbound ip (default true)
  -ping-config-threads value
        How many pinging coroutines exists at the same time (default 16)
  -ping-config-timeoutsec value
        timeout seconds for each request (default 8)
  -ping-config-use-mux
        use mux outbound (default false)
  -ping-config-verbose
        verbose (debug log) (default false)
  -tag-format value
        Format of the auto-generated outbounds' tag (default vmessconfig-autogenerated-%d)
  -template-config-from value
        Where the template file is
  -template-config-to value
        Where the v2ray json config file should write to
  -urls value
        List of your subscription urls
Usage of single:
  -addr value
        Address where the server listen on (default :80)
  -interval value
        Interval for get and ping vmess outbounds (default 1800)
  -outbound-insert-before-tag value
        Insert outbound before the exists outbound whose tag is this (default vmessconfig-outbound-insert)
  -ping-config-allow-insecure
        allow insecure TLS connections (default false)
  -ping-config-count value
        Count. Stop after sending COUNT requests (default 4)
  -ping-config-dest value
        the test destination url, need 204 for success return (default http://www.google.com/gen_204)
  -ping-config-inteval value
        inteval seconds between pings (default 1)
  -ping-config-quit value
        fast quit on error counts (default 0)
  -ping-config-show-node
        show node location/outbound ip (default true)
  -ping-config-threads value
        How many pinging coroutines exists at the same time (default 16)
  -ping-config-timeoutsec value
        timeout seconds for each request (default 8)
  -ping-config-use-mux
        use mux outbound (default false)
  -ping-config-verbose
        verbose (debug log) (default false)
  -template-config-from value
        Where the template file is
  -template-config-to value
        Where the v2ray json config file should write to
  -urls value
        List of your subscription urls
```
