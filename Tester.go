package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type filterQuery struct {
	BlockHash *common.Hash     // used by eth_getLogs, return logs only from block with this hash
	FromBlock string           // beginning of the queried range, nil means genesis block
	ToBlock   string           // end of the range, nil means latest block
	Addresses []common.Address // restricts matches to events created by specific contracts

	// The Topic list restricts matches to particular event topics. Each event has a list
	// of topics. Topics matches a prefix of that list. An empty element slice matches any
	// topic. Non-empty elements represent an alternative that matches any of the
	// contained topics.
	//
	// Examples:
	// {} or nil          matches any topic list
	// {{A}}              matches topic A in first position
	// {{}, {B}}          matches any topic in first position AND B in second position
	// {{A}, {B}}         matches topic A in first position AND B in second position
	// {{A, B}, {C, D}}   matches topic (A OR B) in first position AND (C OR D) in second position
	Topics [][]common.Hash
}

type JsonrpcMessage struct {
	Version string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
	Error   *interface{}    `json:"error,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
}

func main() {
	//http.HandleFunc("/", handleRoot)
	http.HandleFunc("/createJobs", handleCreateJobs)
	http.HandleFunc("/triggerJob1", triggerJob1)
	http.HandleFunc("/triggerJob2", triggerJob2)
	http.HandleFunc("/triggerJob3", triggerJob3)
	fmt.Println("Listening on port 5100")
	http.ListenAndServe(":5100", nil)

}

type Jobs struct {
	jobId      string
	ocaAddress []string
	fromBlock  string
	toBlock    string
}

var job1 Jobs
var job2 Jobs
var job3 Jobs
var jobArray []Jobs
var FQA []filterQuery
var triggerCount int

func handleCreateJobs(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Creating Jobs")
	job1.jobId = "6ba16cd62b8b4114b95449e35800cba0"
	job1.ocaAddress = []string{"0x291D46AE055Dd592B8CC7DdBF232cc206CEfa975"}
	job1.fromBlock = "0x35e8e0c"
	job1.toBlock = "latest"
	q1 := createEvmFilterQuery(job1.jobId, job1.ocaAddress)
	q1.FromBlock = job1.fromBlock
	q1.ToBlock = job1.toBlock

	job2.jobId = "e21f90744a6f4b44bc00086d30717777"
	job2.ocaAddress = []string{"0x58d4BbD5F75c87ECa079a31d3Bd6d50e737aC0f5"}
	job2.fromBlock = "0x35e8e0c"
	job2.toBlock = "latest"
	q2 := createEvmFilterQuery(job1.jobId, job1.ocaAddress)
	q2.FromBlock = job1.fromBlock
	q2.ToBlock = job1.toBlock

	job3.jobId = "09d1e90f9eaa48ad8042f9e255b21984"
	job3.ocaAddress = []string{"0xa778D1539b8fb94662Fae3D33b9d06D48E75021C"}
	job3.fromBlock = "0x35e96FC"
	job1.toBlock = "latest"
	q3 := createEvmFilterQuery(job1.jobId, job1.ocaAddress)
	q3.FromBlock = job1.fromBlock
	q3.ToBlock = job1.toBlock

	jobArray = append(jobArray, job1)
	jobArray = append(jobArray, job2)
	jobArray = append(jobArray, job3)
	FQA = append(FQA, *q1)
	FQA = append(FQA, *q2)
	FQA = append(FQA, *q3)
	fmt.Println("Created filter queries")
	// for _, filterquerys := range FQA {
	// 	fmt.Println("FQA address:", filterquerys.Addresses)
	// 	fmt.Println("FQA topics:", filterquerys.Topics)
	// 	fmt.Println("FQA fromBlock:", filterquerys.FromBlock)
	// 	fmt.Println("FQA toBlock:", filterquerys.ToBlock)
	// }

}
func triggerJob1(w http.ResponseWriter, r *http.Request) {
	triggerCount++
	go triggerJob(FQA[0], triggerCount)
}
func triggerJob2(w http.ResponseWriter, r *http.Request) {
	triggerCount++
	go triggerJob(FQA[1], triggerCount)
}
func triggerJob3(w http.ResponseWriter, r *http.Request) {
	triggerCount++
	go triggerJob(FQA[2], triggerCount)
}

func triggerJob(q filterQuery, triggerCount int) {
	for {
		time.Sleep(15)
		filterBytes, err := json.Marshal(q)
		if err != nil {
			//return nil
			fmt.Println("Marshalling error")
		}
		msg := JsonrpcMessage{
			Version: "2.0",
			ID:      json.RawMessage(`1`),
		}
		msg.Method = "eth_getLogs"
		msg.Params = json.RawMessage(`[` + string(filterBytes) + `]`)
		bytes, err := json.Marshal(msg)
		time.Sleep(1)
		fmt.Println("Polling on job trigger number", triggerCount)
		url := "http://18.217.184.54:8989"
		resp, _ := sendPostRequest(url, bytes)
		//response,_ := ioutil.ReadAll(resp.body)
		var responseJSON map[string]interface{}
		json.Unmarshal(resp, &responseJSON)
		fmt.Println("Response :", responseJSON["result"])
	}
}

// func handleRoot(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Welcome to the server!")
// 	jobid := "c3d9861a75b945888e14b37e406cd85f"
// 	ocaAddress := []string{"0xA847a7b737e2414Fc6BEef7A1eF05aE446206B52"}
// 	q := createEvmFilterQuery(jobid, ocaAddress)
// 	q.FromBlock = "0x2a922c1"
// 	q.ToBlock = "latest"
// 	fmt.Println("value in fromBlock", q.FromBlock)
// 	fmt.Println("value in ToBlock", q.ToBlock)
// 	fmt.Println("value in Addresses", q.Addresses)
// 	fmt.Println("value in Topics", q.Topics)

// 	filterBytes, err := json.Marshal(q)
// 	if err != nil {
// 		//return nil
// 	}

// 	msg := JsonrpcMessage{
// 		Version: "2.0",
// 		ID:      json.RawMessage(`1`),
// 	}
// 	msg.Method = "eth_getLogs"
// 	msg.Params = json.RawMessage(`[` + string(filterBytes) + `]`)
// 	bytes, err := json.Marshal(msg)
// 	//sendRpcRequest(bytes)
// 	for {
// 		time.Sleep(1)
// 		fmt.Println("Polling")
// 		url := "https://3.16.227.1262"
// 		resp, _ := sendPostRequest(url, bytes)
// 		//response,_ := ioutil.ReadAll(resp.body)
// 		var responseJSON map[string]interface{}
// 		json.Unmarshal(resp, &responseJSON)

// 		fmt.Println("Response :", responseJSON["result"])
// 	}
// }

func sendPostRequest(url string, body []byte) ([]byte, error) {
	time.Sleep(2 * time.Second)
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	r, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	//defer logger.ErrorIfCalling(r.Body.Close)

	if r.StatusCode < 200 || r.StatusCode >= 400 {
		return nil, errors.New("got unexpected status code")
	}

	return ioutil.ReadAll(r.Body)
}

func createEvmFilterQuery(jobid string, strAddresses []string) *filterQuery {
	var addresses []common.Address
	for _, a := range strAddresses {
		b := strings.Replace(a, "xdc", "0x", 1)
		addresses = append(addresses, common.HexToAddress(b))
	}

	var (
		// RunLogTopic20190207withoutIndexes was the new RunRequest filter topic as of 2019-01-28,
		// after renaming Solidity variables, moving data version, and removing the cast of requestId to uint256
		RunLogTopic20190207withoutIndexes = utils.MustHash("OracleRequest(bytes32,address,bytes32,uint256,address,bytes4,uint256,uint256,bytes)")
	)
	topics := [][]common.Hash{{
		RunLogTopic20190207withoutIndexes,
	}, {
		StringToBytes32(jobid),
	}}

	return &filterQuery{
		Addresses: addresses,
		Topics:    topics,
	}
}

func StringToBytes32(str string) common.Hash {
	value := common.RightPadBytes([]byte(str), utils.EVMWordByteLen)
	hx := utils.RemoveHexPrefix(hexutil.Encode(value))

	if len(hx) > utils.EVMWordHexLen {
		hx = hx[:utils.EVMWordHexLen]
	}

	hxStr := utils.AddHexPrefix(hx)
	return common.HexToHash(hxStr)
}
