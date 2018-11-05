# [WIP]

## Targets : 
  * simplified version of Kube-router just for Netpols (useful when using
    a simple cni and/or to decouple networking & security)
  * should be runnable on any type of platform (not only k8s) but still based on the netpols of a k8s cluster (can be useful for tbshooting and to span the security policies outside the cluster).

Heavily based on Kube-router codebase (but should diverge over time).
Btw, you should take a look a their awesome project repo : github.com/cloudnativelabs/kube-router

## TODOs : 
  * Withdraw the node-hostname dependency to be able to get all the iptables on 1 machine
  * Try to translate netpols to get meaning for a remote machine (especially
    the "bridged" entries or the namespace rules)
  * Could also implement the iptables on different linux namespaces on the remote
    machine (to avoid overriding/collisions)
      Reminder : `sudo ip netns exec NS1 iptables`
