Dev

```
rm -f /data/tmp/kube-filebeat.registry && go run main.go --path.config=/data/tmp/
```

filebeat.yml

```
filebeat.config.inputs:
  enabled: true
  path: /usr/share/filebeat/inputs.d/*.yml
  reload.enabled: true
  reload.period: 10s

filebeat.inputs:
- type: log
  paths:
  - /var/log/container/*.log
processors:
- decode_json_fields:
    fields:
    - message
    target: ''
- add_kubernetes_metadata:
    in_cluster: true
    include_labels:
    - app
    - project-app
    - project-ns
    include_annotations:
    - project.cloud/controller-kind
    matchers:
    - logs_path:
        logs_path: /var/log/containers
- add_k8s_labels_metadata:
- drop_fields:
    fields:
    - message
    - beat
    - input_type
    - prospector
    - input
    - host
    - type
    - kubernetes
    - offset

output.kafka:
  hosts: ["kafka集群地址"]
  worker: 5
  version: "0.11.0.0"
  topic: "%[topic]" #topic 自动根据k8s服务规则生成
  partition.round_robin:
    group_events: 1
    reachable_only: true
  required_acks: 1
  keep_alive: 10
  compression: snappy
  max_message_bytes: 1048576
  bulk_max_size: 2048
  channel_buffer_size: 256

output.qlog:
  container: true #topic 自动根据k8s服务规则生成
  idc: "机房"
  timeOut: 10
  minConnNum: 1
  compression: gzip
  cluster: "qlog下游存储集群地址"
  zkServer: [""]
  hosts: ["qlog集群地址"]

output.qlog:
  container: false
  idc: "机房"
  timeOut: 10
  topic: "指定topic"
  minConnNum: 1
  compression: gzip
  cluster: "qlog下游存储集群地址"
  zkServer: [""]
  hosts: ["qlog集群地址"]

output.logstash:
  hosts:[""]

output.others:
  hosts:[""]

```
```
注：
1、去除registry本地注册文件配置
2、去除input config配置文件
3、增加out多路输出，可同时支持多个out共享同一个pipeline数据通道输出
4、增加收割机processes处理类型，针对容器环境生成动态生成对象topic
5、增加qlog配置参数，包括topic自动创建、zk集群地址，超时，最小连接参数配置等。后续增加压缩参数
6、容器日志自动关联topic，非容器日志指定topic名称，自动创建topic
```
