## urlshortner run on kuber (create pod):
kubectl create -f devops/urlshortner-pod.yaml

### get all pods info:
kubectl get pod

### get a specific pod ifo
kubectl get pod urls-pod 

### get detailed pod info 
kubectl describe pod urls-pod     

### delete pod:
kubectl delete pod urls-pod  

### attach to a pod (see stdout)
kubectl attach urls-pod -i

### portforward a port on pod to a port on local host
kubectl port-forward urls-pod 8082:8081
forward port 8081 on pod to port 8082 on localhost. 
?? WRONG ?? : port on pod (8081) must be added to container exposed ports in pod yaml file (ports: - containerPort: 8081) section

### execute command on pod
executes ls command on pod
kubectl exec urls-pod -- ls

### execute shell on pod 
kubectl exec -it urls-pod -- ash 
if container have bash, we can use bash instead of ash

### see pod logs
kubectl logs urls-pod