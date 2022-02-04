# Multi Channel Demo App  

## Description

This demo creates two dApps, **chtest** and **chcheck**, that exchange messages with each other through four channels.

**chtest** sends test messages on channels: **testch1**, **testch2** and **testch3**.

**chcheck** reads these test messages and replies to them on **checkch**, specifying which channel it read from. 


## How to run it  

To run this demo, make sure you have Insprd and Kafka Message Broker installed in your cluster (the Channels communicate through Kafka).  
Install Kafka in Insprd by running:
```
insprctl brokers kafka yamls/kafka.yaml
```

Then build the applications that are going to run in each dApp by executing the Makefile (type `make` in the terminal from within this example's folder).  

And finally apply the Inspr structures in Insprd:
```
insprctl apply -f yamls/01.cht.yaml
insprctl apply -k yamls/channels
insprctl apply -k yamls/nodes
```

Now it's all set! If everything worked fine two new deployments should've been created in `inspr-apps` namespace in your Kubernetes cluster, and you can check it's pods to verify that the correct messages are being exchanged between **chtest** and **chcheck**.