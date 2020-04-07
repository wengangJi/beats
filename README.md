## Customized modification based on file beat.
......


### changelogs，
* [Add multiplexing out. Multiple outs share the same pipeline]
* [Add event handler based on K8S label. It outputs k8s deployment、daemonset、statfuleset、job、configmap、pod and generate k8s topic dynamically]
* [Add multiple output objects]


### config
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
  hosts: ["kafka culuster"]
  worker: 5
  version: "0.11.0.0"
  topic: "%[topic]" #topic rule
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
  container: true #topic rule
  idc: "idc"
  timeOut: 10
  minConnNum: 1
  compression: gzip
  cluster: "qlog culster"
  zkServer: [""]
  hosts: ["qlog cluster"]

output.qlog:
  container: false
  idc: "idc"
  timeOut: 10
  topic: "topic"
  minConnNum: 1
  compression: gzip
  cluster: "qlog cluster"
  zkServer: [""]
  hosts: ["qlog cluster"]

output.logstash:
  hosts:[""]

output.others:
  hosts:[""]
