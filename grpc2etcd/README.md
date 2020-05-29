# grpc etcdv3

grpc使用etcdv3作为服务发现组件, 使用新版实现Build接口方式

测试方式:

1. 进入server目录, 启动n个server

```bash
./server --port 8090 --ip 172.16.7.16 &
./server --port 8091 --ip 172.16.7.16 &
./server --port 8092 --ip 172.16.7.16 &
```

2.进入client目录, 启动client

`./client -n 10`

由于client会请求num次, 因此会看到负载均衡的效果

3. 查看监听服务地址

`etcdctl --endpoints "http://172.16.7.16:9002,http://172.16.7.16:9004,http://172.16.7.16:9006" watch /etcd/helloworld/ --prefix`

也会看到以下服务注册, 取消的效果

```bash
PUT
/etcd/helloworld/172.16.7.16:8092
172.16.7.16:8092
DELETE
/etcd/helloworld/172.16.7.16:8092
PUT
/etcd/helloworld/172.16.7.16:8091
172.16.7.16:8091
```