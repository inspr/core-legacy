# Multibroker Sidecar demo  

## Description:
This demo creates a Pod in the Cluster which contains a simple Inspr Node that runs multibroker sidecar test.  
The test consists in initializing a dApp client that will write and read messages, through the Load Balancer sidecar, into/from the Kafka Broker Sidecar and then into/from the Kafka Message Broker.  
**To check if the test worked properly, the demo should print a number that represents a counter for the messages sent (each message adds +1 to the counter).**

## How to run it:  
To run this demo, it must be assured that Kafka is installed in the cluster. Then, Kafka must be installed in Insprd so Channels can use the Kafka Broker:  
```
insprctl cluster config kafka yamls/kafka.yaml
```  

Also, the image originated from the Dockerfile in this folder also must already be created (run the Makefile to do so).  

Once the previous steps are concluded, you can proceed to create the Inspr structures in Insprd:  
```
insprctl apply -f yamls/01.type.yaml
insprctl apply -f yamls/02.channel.yaml
insprctl apply -f yamls/03.app.yaml
```
