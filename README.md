"# goTester" 
This Project is intended to test the RPC failed polling issue by simulating the job triggering mechanism of existing jobs  

Pre requisites:  
ubuntu: 20.04  
vscode  
go extension in VSCODE & install if there is any recomendation from VScode for go  
go => sudo apt install golang-go  
postman => sudo snap install postman  

Step 1: git clone https://github.com/GoPlugin/Rpc_Tester.git  
Step 2: navigate to Rpc_Tester  
Step 3: Replace your RPC URL at line 148 in Tester.go  
Step 4: execute command "go run Tester.go" in VSCode terminal => this will download all the go dependencies and run on port 5100 with harcoded jobs you can watch this on VSCode console  
Step 5: send get request from post man "http://localhost:5100/triggerJob1" or triggerJob2 or triggerJob3 => this will trigger an existing job to poll on specified rpc server (trigger multiple times to have more polls to rpc)  
