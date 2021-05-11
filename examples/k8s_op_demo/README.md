# K8s Operator demo  

## Description:
This demo creates channels/types/dApps in memory an then applies these creations inside the running Cluster.

## How to run it:  
To run this demo, it must be assured that the Cluster is up and running, aswell as the image originated from the Dockerfile in this folder also must already be created and pushed.  
If the above is valid, just run the following command to create given structures inside the Cluster:  

`kubectl apply -f examples/k8s_op_demo/testpod.yaml`  

**PS:**  
If desired, the main.go in this folder can be edited to create channels/types/dApps with different information. Be advised that by doing so, it is necessary to rebuild and push the Docker Image.
