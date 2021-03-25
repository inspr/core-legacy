# Inspr CLI Install
This is a quick tutorial on how to install the **latest version avaliable** of Inspr's CLI.  

## Linux / Mac

Run the following command in the Terminal to install :  
```
curl -s https://storage.googleapis.com/inspr-cli/install_cli.sh | bash
```
To uninstall Inspr's CLI :
```
sudo rm /usr/local/bin/inspr
```
## Other OS
To get other Inspr CLI versions, download the binary file from the GitLab repository [here](https://gitlab.inspr.dev/inspr/core/-/releases).


## After installing

You can check if the installation was successful by typing `inspr -h` on your terminal.

It's important to remeber that if you already have a server with all the necessary helm configuration, you **have** to set the inspr **serverip** to your cluster.

First check the current config using `inspr config list`. 
This will allow you to see all the cli environment values and their current values.

To be able to make changes to your cluster the `serverip` variable has to be changed to comport the hostname currently being used to access the cluster IP. This can be done using the following command 
`inspr config serverip http://<your_domain>.com`