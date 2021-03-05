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
### 后续支持

使用Daemonset的方式，部署到容器云。

