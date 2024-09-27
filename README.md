# System Upgrade Controller

## Deprecated

> [!WARNING]
> This project is deprecated due to the upstream project [merging the changes](https://github.com/rancher/system-upgrade-controller/pull/328) that were made here.

## Upstream

This is an opinionated fork of the upstream project at https://github.com/rancher/system-upgrade-controller.

## Introduction

This project aims to provide a general-purpose, Kubernetes-native upgrade controller (for nodes).
It introduces a new CRD, the **Plan**, for defining any and all of your upgrade policies/requirements.
A **Plan** is an outstanding intent to mutate nodes in your cluster.
For up-to-date details on defining a plan please review [v1/types.go](pkg/apis/upgrade.cattle.io/v1/types.go).

![diagram](doc/architecture.png "The Controller manages Plans by selecting Nodes to run upgrade Jobs on.
 A Plan defines which Nodes are eligible for upgrade by specifying a label selector.
 When a Job has run to completion successfully the Controller will label the Node
 on which it ran according to the Plan that was applied by the Job.")

### Presentations and Recordings

#### April 14, 2020

[CNCF Member Webinar: Declarative Host Upgrades From Within Kubernetes](https://www.cncf.io/webinars/declarative-host-upgrades-from-within-kubernetes/)
- [Slides](https://www.cncf.io/wp-content/uploads/2020/08/CNCF-Webinar-System-Upgrade-Controller-1.pdf)
- [Video](https://www.youtube.com/watch?v=uHF6C0GKjlA)

#### March 4, 2020

[Rancher Online Meetup: Automating K3s Cluster Upgrades](https://info.rancher.com/online-meetup-automating-k3s-cluster-upgrades)
- [Video](https://www.youtube.com/watch?v=UsPV8cZX8BY)

### Considerations

Purporting to support general-purpose node upgrades (essentially, arbitrary mutations) this controller attempts
minimal imposition of opinion. Our design constraints, such as they are:

- content delivery via container image a.k.a. container command pattern
- operator-overridable command(s)
- a very privileged job/pod/container:
  - host IPC, NET, and PID
  - CAP_SYS_BOOT
  - host root file-system mounted at `/host` (read/write)
- optional opt-in/opt-out via node labels
- optional cordon/drain a la `kubectl`

_Additionally, one should take care when defining upgrades by ensuring that such are idempotent--**there be dragons**._

## Deploying

Take a look at [kubesearch.dev](https://kubesearch.dev/#system-upgrade-controller) for a list of Helm charts and other deployment options.

### Example Plans

Using this for Talos and with the examples below will only work if `node-feature-discovery` is installed and configured with the system source like so:

```
worker:
  config:
    core:
      sources: ["pci", "system", "usb"]
```

- [examples/kubernetes.yaml](examples/kubernetes.yaml)
  - Demonstrates upgrading Kubernetes on Talos Linux.
- [examples/talos.yaml](examples/talos.yaml)
  - Demonstrates upgrading Talos Linux OS.

## Building

```shell script
go build -o bin/system-upgrade-controller
```

## Running

Use `./bin/system-upgrade-controller`.

## License
Copyright (c) 2019-2022 [Rancher Labs, Inc.](http://rancher.com)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

[http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
