# Tunnel

Tunnel 是一个网络隧道程序，可以在多个网络之间灵活的使用多种方式设置网络代理。

# 场景：A 通过 B 使用 Socks5 代理上网


Server端配置文件socks5.json

```
{
    "Sources": [
        {
            "Source": {
                "Category": "Socks5",
                "Location": "127.0.0.1:8080",
                "Protocol": "tcp"
            }
        }
    ]
}
```

Server端执行：

```
tunnel ./config/socks5.json
```

# 场景：安全加密

A（10.1.1.1）通过B（10.1.1.2）代理上网，A和B之间通过加密隧道相连：

A端配置：
```
{
    "Sources": [
        {
            "Source": {
                "Category": "Socks5",
                "Location": "127.0.0.1:8080",
                "Protocol": "tcp"
            }
        },
        {
            "Source": {
                "Category": "ConnectTunnel",
                "Location": "10.1.1.1:8081",
                "Protocol": "tcp",
                "SecPath" : "certs"
            },
            "Router": {
                "Domains": [
                    ".*.baidu.com",
                    "*.google.com",
                    ".*"
                ]
            }
        }
    ]
}
```

B端配置：
```
{
    "Sources": [
        {
            "Source": {
                "Category": "ListenTunnel",
                "Location": "10.1.1.2:8081",
                "Protocol": "tcp",
                "SecPath" : "certs"
            }
        }
    ]
}
```
 
# 场景：多个网络间互通

A（10.1.1.1）连接 C（10.1.1.3）
B（10.1.1.2）连接 C（10.1.1.3）
A通过B访问域名 abc.com

A端配置：
```
{
    "Sources": [
        {
            "Source": {
                "Category": "HTTPProxy",
                "Location": "127.0.0.1:12000",
                "Protocol": "tcp"
            }
        },
        {
            "Source": {
                "Category": "ConnectTunnel",
                "Location": "10.1.1.3:8080",
                "Protocol": "tcp"
            },
            "Router": {
                "Domains": [
                    "abc.com"
                ]
            }
        }
    ]
}
```

B端配置：
```
{
    "Sources": [
        {
            "Source": {
                "Category": "ConnectTunnel",
                "Location": "10.1.1.3:8080",
                "Protocol": "tcp"
            }
        }
    ],
    "Default": {
                "Domains": [
					"abc.com"
                ]
    }
}
```

C端配置：

```
{
    "Sources": [
        {
            "Source": {
                "Category": "ListenTunnel",
                "Location": "10.1.1.3:8080",
                "Protocol": "tcp"
            }
        }
    ]
}
```

这样，通过A本地的端口12000可以使用HTTP代理访问abc.com，访问数据会通过B转发。
