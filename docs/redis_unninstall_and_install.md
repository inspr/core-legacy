## Redis Install

To install Redis on the cluster, you need to follow the steps below:

1. Go to cmd/uid_provider/k8s_redisdep folder
2. Apply both deployment and service:
    - `kubectl apply -f uidp-dep.yaml`
    - `kubectl apply -f uidp-svc.yaml`
3. Wait for all pods to be created (make sure of this). For that, you can look in k9s (in the pods tab) and see that `redis-cluster-5` is in the Running state (this is because we currently create 6 pods for redis - note that the index starts at zero)
4. Run the command (remember to type "yes" after running it):
    - `kubectl exec -it redis-cluster-0 -- redis-cli --cluster create --cluster-replicas 1 $(kubectl get pods -l app=redis-cluster -o jsonpath='{range.items[*]}{.status.podIP}:6379 ' | sed -e 's/ :6379/ 127.0.0.1:6379/')`

## Redis Unninstall

To unninstall Redis from the cluster, you need to follow the steps below:

1. Go to cmd/uid_provider/k8s_redisdep folder
2. Delete both deployment and service from the cluster:
    - `kubectl delete -f uidp-dep.yaml`
    - `kubectl delete -f uidp-svc.yaml`
3. Delete the Redis `statefulset` and `pvc` (persistent volume claim) from the cluster. You can use k9s for that.
