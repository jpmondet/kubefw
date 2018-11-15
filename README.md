# KubeFW

Lightweight implementation of Kubernetes Network Policies (translating them into iptables)

## Quick Usage

### In a Kubernetes cluster (most common usage)

#### As a DaemonSet

`kubectl apply -f https://raw.githubusercontent.com/jpmondet/kubefw/master/k8s-deploy.yaml`

#### As a binary on Workers

`kubefw --kubeconfig ~/.kube/config`

### Outside a Kubernetes cluster

This usage can be used for whatever reason

For example, to allow remote troubleshooting of all the iptables across all nodes generated from Network Policies on a single machine

Or even use the Network Policies "logic" on machines completely unrelated to Kubernetes ([Still a WIP])

#### On a remote machine or a machine not hosting pods

`kubefw --kubeconfig ~/.kube/config --standalone`

or if you don't wanna use a standalone binary

``docker start `docker create --net=host --rm --privileged -v ~/.kube/config:/kconfig -v /var/lib/kubernetes/:/var/lib/kubernetes/ -ti jpmondet/kubefw:0.1 /kubefw --kubeconfig /kconfig --standalone` ``

(On the last command, `/var/lib/kubernetes` is only needed if your kubeconfig is pointing to certificates that are in that directory)

## Targets

1. Simplified version of Kube-router just for Netpols (useful when using
    a simple cni and/or to decouple networking & security)
2. Should be runnable on any type of platform (not only k8s) but still based on the netpols of a k8s cluster (can be useful for tbshooting and to span the security policies outside the cluster).
3. Allow an usage entirely decoupled from a k8s cluster.

Heavily based on Kube-router codebase (but should diverge over time).
Btw, you should take a look a their awesome project repo [kube-router](https://github.com/cloudnativelabs/kube-router)

## Current WIPs

* ~~Working version as a Daemonset inside a Kubernetes Cluster~~
* ~~Withdraw the node-hostname dependency to be able to get all the iptables on 1 standalone machine~~
* Try to translate netpols to get enhanced meaning for a remote machine
* Could maybe also implement the iptables on different linux namespaces on the remote
    machine (to avoid overriding/collisions)
      Reminder : `sudo ip netns exec NS1 iptables`