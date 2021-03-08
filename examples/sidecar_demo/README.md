# Sidecar demo  

## Description:
This demo creates a Pod in the Cluster which contains a simple Inspr Node that runs Kafka sidecar test.  
The test consists in initializing a dApp client that will write and read messages, through the sidecar, into/from Kafka Message Broker.  
**To check if the test worked properly, the demo should print the message "Mensagem Lida: Boa tarde amigos".**

## How to run it:  
To run this demo, it must be assured that Kafka and it's sidecar are already deployed and running in the cluster where the demo will be executed.  
Also, the image originated from the Dockerfile in this folder also must already be created.  
If the above is valid, just run the following command:  

`kubectl apply -f examples/sidecar_demo/testpod.yaml` 
