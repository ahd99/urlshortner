## urlshortner run on kuber:
kubectl create -f devops/urlshortner-pod.yaml

### get all pods info:
kubectl get pod

### get a specific pod ifo
kubectl get pod urls-pod 

### get detailed pod info 
kubectl describe pod urls-pod     

### delete pod:
kubectl delete pod urls-pod  