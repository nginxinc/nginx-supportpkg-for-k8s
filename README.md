# nginx-supportpkg-for-k8s

A kubectl plugin designed to collect diagnostics information on any NGINX product running on k8s. 

## Supported products

Currently, [NIC](https://github.com/nginxinc/kubernetes-ingress) is the only supported product.

## Features

Depending on the product, the plugin might collect some or all of the following global and namespace-specific information:

- k8s version, nodes information and CRDs
- pods logs
- list of pods, events, configmaps, services, deployments, statefulsets, replicasets and leases
- k8s metrics
- helm deployments
- `nginx -T` output from NGINX pods

The plugin DOES NOT collect secrets or coredumps.

## Installation

### Building from source
Clone the repo and run `make install`. This will build the binary and copy it on `/usr/local/bin/`.

Verify that the plugin is properly found by `kubectl`:

```
$ kubectl plugin list
The following compatible plugins are available:

/usr/local/bin/kubectl-nginx_supportpkg
```

### Downloading the binary

Navigate to the [releases](https://github.com/nginxinc/nginx-supportpkg-for-k8s/releases) section and download the asset for your operating system and architecture from the most recent version. 

Decompress the tarball and copy the binary somewhere in your `$PATH`. Make sure it is recognized by `kubectl`:

```
$ kubectl plugin list
The following compatible plugins are available:

/path/to/plugin/kubectl-nginx_supportpkg
```

## Usage

The plugin is invoked via `kubectl nginx-supportpkg` and has two required flags:

* `-n` or `--namespace` indicates the namespace(s) where the product is running.
* `-p` or `--product` indicates the product to collect information from.


```
$ kubectl nginx-supportpkg -n default -n nginx-ingress-0 -p nic
Running job pod-list... OK
Running job collect-pods-logs... OK
Running job events-list... OK
Running job configmap-list... OK
Running job service-list... OK
Running job deployment-list... OK
Running job statefulset-list... OK
Running job replicaset-list... OK
Running job lease-list... OK
Running job k8s-version... OK
Running job crd-info... OK
Running job nodes-info... OK
Running job metrics-information... OK
Running job helm-info... OK
Running job helm-deployments... OK
Supportpkg successfully generated: nic-supportpkg-1711384966.tar.gz

```