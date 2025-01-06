# timelite(WIP)

timelite 是一个基于 Prometheus TSDB 的嵌入式时序数据库以及查询引擎，支持存储 Prometheus 格式的数据以及 PromQL 查询，并在 terminal 中提供类 Grafana 的交互式查询。

## demo
<video src="demo/demo.mp4" controls></video>

## 特性

- 支持 PromQL 查询，并提供交互式查询
- 使用 terminal 进行展示，无需配置外部可视化工具
- 支持存储 Panels 配置

## 安装

```bash
go install github.com/timelite/timelite
```

## 使用
Commands:
```bash
help [<command>...]
    Show help.

query [<flags>]
    Query the data in the timeline.
    example:
    timelite query --config-file="config.json" --storage-path="./storage/tsdb"

tsdb mock
    Mock the tsdb data.

tsdb run [<flags>]
    Run the tsdb server.
    example:
    timelite tsdb run --config-file="config.json" --storage-path="./storage/tsdb" --host="0.0.0.0:9090"
```
