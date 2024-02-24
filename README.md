# pushgateway-cleaner

[![license](https://img.shields.io/github/license/abatilo/pushgateway-cleaner.svg)](https://github.com/abatilo/pushgateway-cleaner/blob/main/LICENSE)
[![release](https://img.shields.io/github/release/abatilo/pushgateway-cleaner.svg)](https://github.com/abatilo/pushgateway-cleaner/releases/latest)
[![GitHub release date](https://img.shields.io/github/release-date/abatilo/pushgateway-cleaner.svg)](https://github.com/abatilo/pushgateway-cleaner/releases)

Do you have a Prometheus pushgateway? Were you also surprised by the fact that
pushgateway metrics [will never be
deleted](https://github.com/prometheus/pushgateway/issues/19#issuecomment-225566114)?
That's what this project is for.

This application will use the built-in `push_time_seconds` metric that
pushgateway keeps track of, and will delete metric groups that have been around
longer than you specify.

## Configuration

These are the available flags for `pushgateway-cleaner`
```
⇒ ./pushgateway-cleaner --help
Usage of ./pushgateway-cleaner:
  -h, --help                     Show help
      --pushgateway-url string   Pushgateway URL (default "http://localhost:9091")
      --sync-period duration     How often to check for old metrics (default 3m0s)
      --ttl duration             How old metrics are allowed to be before being deleted (default 10m0s)
  -v, --verbose                  Set logging verbosity level to debug
```

## Usage

We recommend running this as a sidecar container next to your pushgateway
deployment. For example, here's an addition to the
`prometheus-community/prometheus-pushgateway` helm chart.

```yaml
extraContainers:
  - name: pushgateway-cleaner
    args:
      - --verbose
    image: ghcr.io/abatilo/pushgateway-cleaner:latest
    resources:
        requests:
          cpu: 20m
          memory: 4Mi
        limits:
          memory: 16Mi
```

## Installation

Docker containers are built for both `linux/amd64` and `linux/arm64`. Check out
available versions
[here](https://github.com/abatilo/pushgateway-cleaner/pkgs/container/pushgateway-cleaner)
