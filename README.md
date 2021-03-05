# macvlan-ipam-etcd
a macvlan ipam with etcd

## 背景
解决macvlan使用hostlocal分配ip的ip冲突导致各种个样的bug

## 如何构建

```golang

go build 

```


## 目前简单如何使用

```shell

mv macvlan-ipam-etcd  /opt/cni/bin/

```


### 测试

```text
echo '{ "cniVersion": "0.3.1", "name": "examplenet", "ipam": { "name": "myetcd-ipam", "type": "macvlan-ipam-etcd", "etcdConfig": { "etcdURL": "http://127.0.0.1:2379"}, "ranges": [ [{"subnet": "203.0.113.0/24"}]] } }' | CNI_COMMAND=ADD CNI_CONTAINERID=example CNI_NETNS=/dev/null CNI_IFNAME=dummy0 CNI_PATH=. ./macvlan-ipam-etcd 
```

### 利用libcni调用macvlan教程

```text 


cat >/etc/cni/net.d/10-macvlan-etcd.conf  <<EOF
{
        "name": "mymacvlan",
        "type": "macvlan",
        "master": "ens192",
        "ipam": {
                "name": "ipam-etcd",
                "type": "macvlan-ipam-etcd",
                "etcdConfig": {
                        "etcdURL": "http://127.0.0.1:2379"
                },
                "subnet": "10.22.0.0/16",
                "rangeStart": "10.22.0.2",
                "rangeEnd": "10.22.0.254",
                "routes": [{
                        "dst": "0.0.0.0/0"
                }]
        }
}
EOF


[root@node cnitool]# export CNI_PATH=/root/plugins/bin/
[root@node cnitool]# ip netns delete a &&ip netns add a && ./cnitool add mymacvlan /var/run/netns/a
2021-03-05 04:00:35.789322 I | Error retrieving last reserved ip: Can not find last reserved ip
{"level":"warn","ts":"2021-03-05T04:00:35.791-0500","caller":"clientv3/retry_interceptor.go:62","msg":"retrying of unary invoker failed","target":"endpoint://client-03e5e010-bbb2-45c5-aa5d-b895eb93129e/127.0.0.1:2379","attempt":0,"error":"rpc error: code = InvalidArgument desc = etcdserver: key is not provided"}
{
    "cniVersion": "0.2.0",
    "ip4": {
        "ip": "10.22.0.2/16",
        "gateway": "10.22.0.1",
        "routes": [
            {
                "dst": "0.0.0.0/0"
            }
        ]
    },
    "dns": {}
}[root@node cnitool]# ping 10.22.0.2
PING 10.22.0.2 (10.22.0.2) 56(84) bytes of data.
64 bytes from 10.22.0.2: icmp_seq=1 ttl=64 time=0.087 ms
64 bytes from 10.22.0.2: icmp_seq=2 ttl=64 time=0.099 ms
^C
--- 10.22.0.2 ping statistics ---
2 packets transmitted, 2 received, 0% packet loss, time 1015ms
rtt min/avg/max/mdev = 0.087/0.093/0.099/0.006 ms

[root@node cnitool]# nsenter --net=/var/run/netns/a bash
[root@node cnitool]# ip a
1: lo: <LOOPBACK> mtu 65536 qdisc noop state DOWN group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
2: eth0@if2: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP group default 
    link/ether e6:97:d2:de:e9:c8 brd ff:ff:ff:ff:ff:ff link-netnsid 0
    inet 10.22.0.2/16 brd 10.22.255.255 scope global eth0
       valid_lft forever preferred_lft forever
    inet6 fe80::e497:d2ff:fede:e9c8/64 scope link 
       valid_lft forever preferred_lft forever


```

### 后续支持

使用Daemonset的方式，部署到容器云。
# macvlan-ipam-etcd
a macvlan ipam with etcd

## 背景
解决macvlan使用hostlocal分配ip的ip冲突导致各种个样的bug

## 如何构建

```golang

go build 

```


## 目前简单如何使用

```shell

mv macvlan-ipam-etcd  /opt/cni/bin/

```


### 测试

```text
echo '{ "cniVersion": "0.3.1", "name": "examplenet", "ipam": { "name": "myetcd-ipam", "type": "macvlan-ipam-etcd", "etcdConfig": { "etcdURL": "http://127.0.0.1:2379"}, "ranges": [ [{"subnet": "203.0.113.0/24"}]] } }' | CNI_COMMAND=ADD CNI_CONTAINERID=example CNI_NETNS=/dev/null CNI_IFNAME=dummy0 CNI_PATH=. ./macvlan-ipam-etcd 
```

### 利用libcni调用macvlan教程

```text 


cat >/etc/cni/net.d/10-macvlan-etcd.conf  <<EOF
{
        "name": "mymacvlan",
        "type": "macvlan",
        "master": "ens192",
        "ipam": {
                "name": "ipam-etcd",
                "type": "macvlan-ipam-etcd",
                "etcdConfig": {
                        "etcdURL": "http://127.0.0.1:2379"
                },
                "subnet": "10.22.0.0/16",
                "rangeStart": "10.22.0.2",
                "rangeEnd": "10.22.0.254",
                "routes": [{
                        "dst": "0.0.0.0/0"
                }]
        }
}
EOF


[root@node cnitool]# export CNI_PATH=/root/plugins/bin/
[root@node cnitool]# ip netns delete a &&ip netns add a && ./cnitool add mymacvlan /var/run/netns/a
2021-03-05 04:00:35.789322 I | Error retrieving last reserved ip: Can not find last reserved ip
{"level":"warn","ts":"2021-03-05T04:00:35.791-0500","caller":"clientv3/retry_interceptor.go:62","msg":"retrying of unary invoker failed","target":"endpoint://client-03e5e010-bbb2-45c5-aa5d-b895eb93129e/127.0.0.1:2379","attempt":0,"error":"rpc error: code = InvalidArgument desc = etcdserver: key is not provided"}
{
    "cniVersion": "0.2.0",
    "ip4": {
        "ip": "10.22.0.2/16",
        "gateway": "10.22.0.1",
        "routes": [
            {
                "dst": "0.0.0.0/0"
            }
        ]
    },
    "dns": {}
}[root@node cnitool]# ping 10.22.0.2
PING 10.22.0.2 (10.22.0.2) 56(84) bytes of data.
64 bytes from 10.22.0.2: icmp_seq=1 ttl=64 time=0.087 ms
64 bytes from 10.22.0.2: icmp_seq=2 ttl=64 time=0.099 ms
^C
--- 10.22.0.2 ping statistics ---
2 packets transmitted, 2 received, 0% packet loss, time 1015ms
rtt min/avg/max/mdev = 0.087/0.093/0.099/0.006 ms

[root@node cnitool]# nsenter --net=/var/run/netns/a bash
[root@node cnitool]# ip a
1: lo: <LOOPBACK> mtu 65536 qdisc noop state DOWN group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
2: eth0@if2: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP group default 
    link/ether e6:97:d2:de:e9:c8 brd ff:ff:ff:ff:ff:ff link-netnsid 0
    inet 10.22.0.2/16 brd 10.22.255.255 scope global eth0
       valid_lft forever preferred_lft forever
    inet6 fe80::e497:d2ff:fede:e9c8/64 scope link 
       valid_lft forever preferred_lft forever


```

### 后续支持

使用Daemonset的方式，部署到容器云。

