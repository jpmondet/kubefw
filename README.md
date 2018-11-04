# [WIP]

## Targets : 
  * simplified version of Kube-router just for Netpols (useful when using
    a simple cni and/or to decouple networking & security)
  * should be usable on any type of platform (not only k8s) but still based on k8s netpols to implement the same k8s security everywhere.

Heavily based on Kube-router code thus keeping the same Licence/Authors.
Btw, you should take a look a their repo : github.com/cloudnativelabs/kube-router

## TODOs : 
  * Withdraw the node-hostname dependency to be able to get all the iptables on 1 machine
  * Try to translate netpols to get meaning for a remote machine (especially
    the "bridged" entries)
  * Could also implement the iptables on different linux namespaces on the remote
    machine (to avoid overriding/collisions)
      Reminder : `sudo ip netns exec NS1 iptables`
