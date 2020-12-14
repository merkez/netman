# netman

netman is a simple tool which applies `ifup` command to  network interfaces which are down. 
Used in the virtual machine which is integrated to Network Analysis Platform. 

## How it works 

There is a service file called [netman.service](.github/scripts/netman.service). It needs to be placed 
under `/lib/systemd/system/` in the virtual machine. (Ubuntu 18.04)

The service file indicates that netman will run once when system booted then exits. (No further run). 
In order to run service file, released version of netman should be deployed into the virtual machine
from [releases page](https://github.com/mrturkmenhub/netman/releases)