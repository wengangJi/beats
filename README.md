## Customized modification based on file beat.
......


### changelogs，
* [Add multiplexing out. Multiple outs share the same pipeline]
* [增加基于k8s label的event数据处理功能，动态根据k8s 服务状态生成topic , 输出k8s deployment、daemonset、statfuleset、job、configmap、pod等对象数据]
* [增加output支持种类]


### 配置文件
* filebeat.yml

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
