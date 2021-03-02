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

### 后续支持

使用Daemonset的方式，部署到容器云。

