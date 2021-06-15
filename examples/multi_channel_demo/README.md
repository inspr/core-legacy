# Multi Channel Demo App  

## Description

This demo creates an app that contains two nodes, **chtest** and **chcheck**, that exchange messages through four channels with each other.

**chtest** sends test messages on channels: **testch1**, **testch2** and **testch3**.

**chcheck** reads these test messages and replies to them on **checkch**, specifying witch channel it read from. 


## How to run it  

To run this demo, have a running insprctl cluster.
Run `make` and run `insprctl apply -k inspr` inside this directory.