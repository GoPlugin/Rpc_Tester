"# goTester" 
This Project is intended to test the RPC failed polling issue by simulating the triggering mechanism of existing jobs  

Pre requisites:  
ubuntu: 20.04  
vscode  
go version gol.17.3  
postman  

Step 1: git clone https://github.com/nareshmmr/PluginRpcTester.git  
Step 2: navigate to PluginRpcTester  
Step 3: go run Tester.go => this will download all the go dependencies and run on port 5100 with harcoded jobs you can watch this on VSCode console  
Step 4: send get request from post man http://localhost:5100/triggerJob1 or triggerJob2 or triggerJob3 => this will trigger an existing job to poll on specified rpc server every 15 seconds(trigger multiple times to have more polls to rpc)  
