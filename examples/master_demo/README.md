# Master demo  
## Description  
Master's demo consists in running Inspr daemon on localhost and applying operations to modify the dApp structure created.  
By running main.go whilst the server is up, it will execute some basic operations to show how Create/Get/Update/Delete methods are altering the root dApp structure (in the memory).
## How to run it:
First, the Inspr daemon server must be initialized in the locahost. To do so, run the following command in the terminal:

`go run cmd/insprd/main.go`    

Then, in another terminal instance, run:

`go run examples/master_demo/main.go`  

As a result, new dApps, Channels and Channel Types will be  created, got (as in Get), updated and then deleted, and each operation done will be printed in the terminal, along with the diff applied.  

**PS:**  
If desired, the main.go in this folder can be edited to do different operations. Also, if Inspr's CLI is installed, it is possible to run Create/Get/Update/Delete methods by using the CLI's commands.