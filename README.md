# KubeFW

Lightweight implementation of Kubernetes Network Policies (translating them into iptables)

## Release state

* Ready and Stable for Kubernetes (tested heavily as a DaemonSet).

## Quick Usage

### In a Kubernetes cluster (most common usage)

#### As a DaemonSet

`kubectl apply -f https://raw.githubusercontent.com/jpmondet/kubefw/master/k8s-deploy.yaml`

#### As a binary on Workers

Assuming your kubeconfig is in `~/.kube/config` on the Worker and the current
user is privileged (it is needed to use iptables/ipsets):

`kubefw --kubeconfig ~/.kube/config`

or if you don't wanna use a standalone binary

``docker start `docker create --net=host --rm --privileged -v ~/.kube/config:/kconfig -v /var/lib/kubernetes/:/var/lib/kubernetes/ -ti jpmondet/kubefw:0.1.2 /kubefw --kubeconfig /kconfig` ``

(Mounting `/var/lib/kubernetes` is only needed if your kubeconfig is pointing to certificates that are in that directory)

#### And after that ? 

You are ready to use Network Policies as defined by Kubernetes.

If you need example, you should take a look at this nice [cheatsheet/learning repo](https://github.com/ahmetb/kubernetes-network-policy-recipes) 

A quick example of a network policy that denies all traffic to an application
would be : 

```yaml
kind: NetworkPolicy
apiVersion: networking.k8s.io/v1
metadata:
  name: nginx-deny-all
spec:
  podSelector:
    matchLabels:
      run: nginx
  ingress: []
```

## More about this project 

This project may be useful when you are using a simple CNI plugin and/or to decouple Networking & Security (which should be a best practice in my point of view ;-) )

Deploying KubeFW, even with default parameters will NOT mess up with your network plugin nor with your /etc/cni directory nor anything like that. This is
safe.

It WILL add `iptables` to your nodes depending on your Network Policies but will NOT override the `iptables` of the node itself (if it had any)

There is another usage for servers outside Kubernetes but it's still a **WIP** (see [standalone_support](https://github.com/jpmondet/kubefw/tree/standalone_support) branch)

This is **heavily** based on Kube-router codebase (but lighter and what's left is diverging over time).

Btw, you should take a look a their awesome project repo [kube-router](https://github.com/cloudnativelabs/kube-router) if you need a fully fledge Network Plugin.

