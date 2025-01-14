# timelite(WIP)

timelite is an embedded time-series database and query engine based on Prometheus TSDB. It supports storing Prometheus-format data and PromQL queries, providing a Grafana-like interactive query experience in the terminal.

## demo

![demo](demo/demo.gif)

## Features

- Supports PromQL queries with interactive query interface
- Terminal-based visualization without external tools configuration
- Supports storing Panels configurations

## Installation

```bash
go install github.com/timelite/timelite
```

## Commands:
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

## Usage
1. Run the tsdb server
```bash
timelite tsdb run --config-file="config.json" --storage-path="./storage/tsdb" --host="0.0.0.0:9090"
```

2. Query the data
```bash
timelite query --config-file="config.json" --storage-path="./storage/tsdb"
```