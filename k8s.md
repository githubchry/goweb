概念理解

[](https://juejin.cn/post/6919639635509379085)

[从0到1使用Kubernetes系列——K8s入门](https://zhuanlan.zhihu.com/p/43266412)

[Kubernetes，2020 快速入门](https://zhuanlan.zhihu.com/p/100644716)

# 安装

[官方教程](https://kubernetes.io/zh/docs/setup/production-environment/tools/kubeadm/install-kubeadm/)

[Kubernetes 最新版本安装过程和注意事项](https://blog.csdn.net/isea533/article/details/86769125)

[kubeadm部署单master节点](https://blog.csdn.net/weixin_40585721/article/details/109545699)

## 前置条件



### 允许 iptables 检查桥接流量 

```
加载br_netfilter模块
sudo modprobe br_netfilter

验证是否已加载
lsmod | grep br_netfilte

cat <<EOF | sudo tee /etc/modules-load.d/k8s.conf
br_netfilter
EOF

添加网桥过滤及地址转发
cat <<EOF | sudo tee /etc/sysctl.d/k8s.conf
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables = 1
EOF

重加载网桥过滤文件
sudo sysctl --system

```

### 禁用swap

```
临时关闭
sudo swapoff -a

永久关闭
sudo vim /etc/fstab 
使用#注释掉以下行:
#/swap.img      none    swap    sw      0       0
```



### docker安装

详阅docker.md

### 修改docker配置

```
sudo vim /etc/docker/daemon.json
{
        "registry-mirrors":["https://9cpn8tt6.mirror.aliyuncs.com", "https://registry.docker-cn.com"],
        "log-driver":"json-file",
        "log-opts": {"max-size":"100m", "max-file":"3"},
        "exec-opts": ["native.cgroupdriver=systemd"]
}


重启docker服务
sudo systemctl daemon-reload
sudo systemctl restart docker
```

这一步主要解决`kubeadm init `时产生的警告:

```shell
[WARNING IsDockerSystemdCheck]: detected "cgroupfs" as the Docker cgroup driver. The recommended driver is "systemd". Please follow the guide at https://kubernetes.io/docs/setup/cri/
```







## kubeadm安装

```
0. 依赖
sudo apt-get update && apt-get install -y curl apt-transport-https

1. 添加阿里云Kubernetes镜像apt key
curl -s https://mirrors.aliyun.com/kubernetes/apt/doc/apt-key.gpg | sudo apt-key add -

2.添加阿里云Kubernetes镜像源
sudo cat <<EOF | sudo tee /etc/apt/sources.list.d/kubernetes.list
deb https://mirrors.aliyun.com/kubernetes/apt/ kubernetes-xenial main
EOF

3. 更新源并安装Kubernetes 
sudo apt-get update
sudo apt-get install -y kubelet kubeadm kubectl

4. 阻止软件自动更新
sudo apt-mark hold -y kubelet kubeadm kubectl

5. 启用 kubectl 自动补齐
5.1 补齐脚本依赖于 bash-completion 软件包, 一般系统已经安装了
sudo apt-get install bash-completion
source /usr/share/bash-completion/bash_completion
5.2 在 ~/.bashrc 文件中源引自动补齐脚本
sudo echo 'source <(kubectl completion bash)' >>~/.bashrc
```

PS: 以上为Linux操作, Windows下安装docker for windows时就已经附带了kubectl

### 说明
 `kubelet:` 运行在cluster所有节点上,负责启动 `Pod` 和容器
 `kubeadm:` 用于初始化cluster
 `kubectl:` kubenetes命令行工具，通过kubectl可以部署和管理应用，查看各种资源，创建，删除和更新组件



**至此 k8s安装完成, kubelet 现在每隔几秒就会重启，因为它陷入了一个等待 kubeadm 指令的死循环。**



# 创建集群

master: 控制节点

worker: 工作节点

一些命令备忘

```
列出所需镜像列表
kubeadm config images list

拉取镜像到本地
kubeadm config images pull
```

## 初始化Master部署

根据[](https://zhuanlan.zhihu.com/p/46341911)配置`--pod-network-cidr`参数

```shell
sudo kubeadm init \
  --apiserver-advertise-address=10.2.3.50 \
  --image-repository registry.aliyuncs.com/google_containers \
  --kubernetes-version=v1.20.4
  
  
  

--apiserver-advertise-address指明用 Master 的哪个 interface 与 Cluster 的其他节点通信。如果 Master 有多个 interface，建议明确指定，如果不指定，kubeadm 会自动选择有默认网关的 interface。
--image-repository指定初始化需要的镜像源从阿里云镜像仓库拉取, 默认的gcr.io会被墙
--pod-network-cidr是指配置节点中的pod的可用IP地址，此为内部IP, 根据网络插件指定不同段的地址(https://zhuanlan.zhihu.com/p/46341911)
--kubernetes-version如果不指定会自动在线获取最新版本号, 结果就是会被墙
```

#### 附kubeadm init参数说明

```shell
--apiserver-advertise-address string   设置 apiserver 绑定的 IP.
--apiserver-bind-port int32            设置apiserver 监听的端口. (默认 6443)
--apiserver-cert-extra-sans strings    api证书中指定额外的Subject Alternative Names (SANs) 可以是IP 也可以是DNS名称。 证书是和SAN绑定的。
--cert-dir string                      证书存放的目录 (默认 "/etc/kubernetes/pki")
--certificate-key string               kubeadm-cert secret 中 用于加密 control-plane 证书的key
--config string                        kubeadm 配置文件的路径.
--cri-socket string                    CRI socket 文件路径，如果为空 kubeadm 将自动发现相关的socket文件; 只有当机器中存在多个 CRI  socket 或者 存在非标准 CRI socket 时才指定.
--dry-run                              测试，并不真正执行;输出运行后的结果.
--feature-gates string                 指定启用哪些额外的feature 使用 key=value 对的形式。
-h, --help                                 帮助文档
--ignore-preflight-errors strings      忽略前置检查错误，被忽略的错误将被显示为警告. 例子: 'IsPrivilegedUser,Swap'. Value 'all' ignores errors from all checks.
--image-repository string              选择拉取 control plane images 的镜像repo (default "k8s.gcr.io")
--kubernetes-version string            选择K8S版本. (default "stable-1")
--node-name string                     指定node的名称，默认使用 node 的 hostname.
--pod-network-cidr string              指定 pod 的网络， control plane 会自动将 网络发布到其他节点的node，让其上启动的容器使用此网络
--service-cidr string                  指定service 的IP 范围. (default "10.96.0.0/12")
--service-dns-domain string            指定 service 的 dns 后缀, e.g. "myorg.internal". (default "cluster.local")
--skip-certificate-key-print           不打印 control-plane 用于加密证书的key.
--skip-phases strings                  跳过指定的阶段（phase）
--skip-token-print                     不打印 kubeadm init 生成的 default bootstrap token 
--token string                         指定 node 和control plane 之间，简历双向认证的token ，格式为 [a-z0-9]{6}\.[a-z0-9]{16} - e.g. abcdef.0123456789abcdef
--token-ttl duration                   token 自动删除的时间间隔。 (e.g. 1s, 2m, 3h). 如果设置为 '0', token 永不过期 (default 24h0m0s)
--upload-certs                         上传 control-plane 证书到 kubeadm-certs Secret.
```



------

初始化成功后可以看到以下打印:

```shell
Your Kubernetes control-plane has initialized successfully!

To start using your cluster, you need to run the following as a regular user:

  mkdir -p $HOME/.kube
  sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
  sudo chown $(id -u):$(id -g) $HOME/.kube/config

Alternatively, if you are the root user, you can run:

  export KUBECONFIG=/etc/kubernetes/admin.conf

You should now deploy a pod network to the cluster.
Run "kubectl apply -f [podnetwork].yaml" with one of the options listed at:
  https://kubernetes.io/docs/concepts/cluster-administration/addons/

Then you can join any number of worker nodes by running the following on each as root:

kubeadm join 10.2.3.50:6443 --token sg59kp.j7qnpz7rdhdoso0h \
    --discovery-token-ca-cert-hash sha256:8ec157b313dee0c68a7bcf25fefe97304db7160edd181d71ce6eb7609c56fd11
```

根据以上提示可划分为三个步骤:

- 配置授权信息
- 部署pod网络
- 创建worker

### 配置授权信息

如果不执行这一步, 会出现下面的问题:

```shell
$ sudo kubectl get pods
The connection to the server localhost:8080 was refused - did you specify the right host or port?
```

解决:

```shell
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config
```

有了上面配置后，后续才能使用 `kubectl` 执行命令。

```
$ sudo kubectl get pods
No resources found in default namespace.
因为当前是没有添加任何pod的
```

后面的`kubeadm join`语句用于初始化`Worker`部署

### 部署pod网络

有很多种选择,  随便选了weave作为pod网络

```
kubectl apply -f "https://cloud.weave.works/k8s/net?k8s-version=$(kubectl version | base64 | tr -d '\n')" 
```

### 部署Web UI

可选 还没弄好 不管

```
kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/v2.2.0/aio/deploy/recommended.yaml  
```



## 初始化Worker部署

根据上面的初始化Master之后得到的`kubeadm join`语句进行`Worker`初始化

该语句可以运行在其他机器(多例集群), 也可以运行在Master机器上(单例集群)

### 单例部署

kubernetes官方默认策略是worker节点运行Pod，master节点不运行Pod。

如果只是为了开发或者其他目的而需要部署单节点集群，可以通过以下的命令设置：

```
kubectl taint nodes --all node-role.kubernetes.io/master-
```

验证安装成功,  全部都是`running`状态

```shell
$ kubectl get pods -n kube-system
NAME                           READY   STATUS              RESTARTS   AGE
coredns-7f89b7bc75-9ghs8       0/1     ContainerCreating   0          73m
coredns-7f89b7bc75-sf9sd       0/1     ContainerCreating   0          73m
etcd-root                      1/1     Running             0          74m
kube-apiserver-root            1/1     Running             0          74m
kube-controller-manager-root   1/1     Running             0          74m
kube-proxy-9qbsm               1/1     Running             0          73m
kube-scheduler-root            1/1     Running             0          74m
weave-net-5n5zz                1/2     Running             0          10s

重复操作查看状态  等待直到ContainerCreating变成Running

$ kubectl get pods -n kube-system
NAMESPACE     NAME                           READY   STATUS    RESTARTS   AGE
kube-system   coredns-7f89b7bc75-9ghs8       1/1     Running   0          76m
kube-system   coredns-7f89b7bc75-sf9sd       1/1     Running   0          76m
kube-system   etcd-root                      1/1     Running   0          76m
kube-system   kube-apiserver-root            1/1     Running   0          76m
kube-system   kube-controller-manager-root   1/1     Running   0          76m
kube-system   kube-proxy-9qbsm               1/1     Running   0          76m
kube-system   kube-scheduler-root            1/1     Running   0          76m
kube-system   weave-net-5n5zz                2/2     Running   0          2m18s

查看阶段状态为Ready
$ kubectl get nodes
NAME   STATUS   ROLES                  AGE   VERSION
root   Ready    control-plane,master   88m   v1.20.4
```

检查集群健康状态

```
kubectl get cs
Warning: v1 ComponentStatus is deprecated in v1.19+
NAME                 STATUS      MESSAGE                                                                                       ERROR
scheduler            Unhealthy   Get "http://127.0.0.1:10251/healthz": dial tcp 127.0.0.1:10251: connect: connection refused
controller-manager   Unhealthy   Get "http://127.0.0.1:10252/healthz": dial tcp 127.0.0.1:10252: connect: connection refused
etcd-0               Healthy     {"health":"true"}

```

可以看到是`Unhealthy`非健康状态

```
修改以下两个文件, 找到--port=0删除或注释掉

sudo vim /etc/kubernetes/manifests/kube-scheduler.yaml

sudo vim /etc/kubernetes/manifests/kube-controller-manager.yaml

```

最多等待半分钟,  再次检查集群健康状态:

```
$ kubectl get cs
Warning: v1 ComponentStatus is deprecated in v1.19+
NAME                 STATUS    MESSAGE             ERROR
scheduler            Healthy   ok
controller-manager   Healthy   ok
etcd-0               Healthy   {"health":"true"}
```

至此 单例部署就完事了, 不用去执行`kubeadm join`语句

### 多例部署

把worker节点加入master:

```shell
sudo kubeadm join 10.2.3.50:6443 --token sg59kp.j7qnpz7rdhdoso0h \
    --discovery-token-ca-cert-hash sha256:8ec157b313dee0c68a7bcf25fefe97304db7160edd181d71ce6eb7609c56fd11
```

注意: `token`的有效期是24小时,  如果过期后想再添加一台worker, 可通过[这里的方法](https://kubernetes.io/docs/setup/independent/create-cluster-kubeadm/#join-nodes)获取和创建`token`和`discovery-token-ca-cert-hash`



## 关于反初始化

```
sudo kubeadm reset 
还需要手动删除
sudo rm -rf /etc/kubernetes/
sudo rm -rf $HOME/.kube
```



# 集群测试

## 创建一个`nginx`的`pod`：

```shell
$ kubectl create deployment nginx --image=nginx
deployment.apps/nginx created

$ kubectl expose deployment nginx --port=80 --type=NodePort
service/nginx exposed

```

## 查看`pod`和`service`

```
$ kubectl get pod,svc -o wide
I0303 14:12:47.382455   10425 request.go:655] Throttling request took 1.178735318s, request: GET:https://10.2.3.50:6443/apis/authorization.k8s.io/v1beta1?timeout=32s
NAME                         READY   STATUS    RESTARTS   AGE    IP          NODE   NOMINATED NODE   READINESS GATES
pod/nginx-6799fc88d8-mvzz4   1/1     Running   0          116s   10.32.0.7   root   <none>           <none>

NAME                 TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)        AGE    SELECTOR
service/kubernetes   ClusterIP   10.96.0.1      <none>        443/TCP        110m   <none>
service/nginx        NodePort    10.106.75.63   <none>        80:31601/TCP   102s   app=nginx

```

在网页访问http://10.2.3.50:31601/即可看到nginx的页面, 其中10.2.3.50就是机器真实网卡上的ip, 不要用10.106.75.63, 那是内部虚拟地址. 

