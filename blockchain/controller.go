package task

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

//Controller ...
type Controller struct {
	blockchain     *Blockchain
	currentNodeURL string
}

//ResponseToSend ...
type ResponseToSend struct {
	Note string
}

//Index GET /
func (ctrl *Controller) Index(c *gin.Context) {
	c.JSON(200, nil)
}

//GetBlockchain GET /blockchain
func (ctrl *Controller) GetBlockchain(c *gin.Context) {
	c.JSON(200, ctrl.blockchain)
	return
}

//RegisterBet POST /task
func (ctrl *Controller) RegisterBet(c *gin.Context) {
	var task Bet
	err := c.Bind(&Bet)
	if  err != nil { // unmarshall body contents as a type Candidate
			log.Fatalln("Error RegisterBet unmarshalling data", err)
			c.JSON(422, err)
			return
	}

	success := ctrl.blockchain.RegisterBet(task) // registers the task into the blockchain
	if !success {
		w.WriteHeader(http.StatusInternalServerError)
		c.JSON(500, sucess)
		return
	}

	var resp ResponseToSend
	resp.Note = "Bet created and broadcast successfully."
	c.JSON(200, resp)
	return
}

//RegisterAndBroadcastBet POST /task/broadcast
func (ctrl *Controller) RegisterAndBroadcastBet(c *gin.Context) {
	var task Bet
	err := c.Bind(&task)
	if err != nil { // unmarshall body contents as a type Candidate
			log.Fatalln("Error RegisterBet unmarshalling data", err)
			c.JSON(422, err)
			return
	}

	success := ctrl.blockchain.RegisterBet(task) // registers the task into the blockchain
	if !success {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// broadcast
	for _, node := range ctrl.blockchain.NetworkNodes {
		if node != ctrl.currentNodeURL {
			// call /register-node in node
			MakePostCall(node+"/task", body)
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	var resp ResponseToSend
	resp.Note = "Bet created and broadcast successfully."
	data, _ := json.Marshal(resp)
	w.Write(data)
}

//Mine GET /mine
func (ctrl *Controller) Mine(c *gin.Context) {
	lastBlock := ctrl.blockchain.GetLastBlock()
	previousBlockHash := lastBlock.Hash
	currentBlockData := BlockData{Index: strconv.Itoa(lastBlock.Index - 1), Bets: ctrl.blockchain.PendingBets}
	currentBlockDataAsByteArray, _ := json.Marshal(currentBlockData)
	currentBlockDataAsStr := base64.URLEncoding.EncodeToString(currentBlockDataAsByteArray)

	nonce := ctrl.blockchain.ProofOfWork(previousBlockHash, currentBlockDataAsStr)
	blockHash := ctrl.blockchain.HashBlock(previousBlockHash, currentBlockDataAsStr, nonce)
	newBlock := ctrl.blockchain.CreateNewBlock(nonce, previousBlockHash, blockHash)
	blockToBroadcast, _ := json.Marshal(newBlock)

	for _, node := range ctrl.blockchain.NetworkNodes {
		if node != ctrl.currentNodeURL {
			// call /receive-new-block in node
			MakePostCall(node+"/receive-new-block", blockToBroadcast)
		}
	}
	var resp ResponseToSend
	resp.Note = "New block mined and broadcast successfully."
	c.JSON(200, resp)
	return
}

//RegisterNode POST /register-node
func (ctrl *Controller) RegisterNode(c *gin.Context) {
	if err := r.Body.Close(); err != nil {
		log.Fatalln("Error RegisterNode", err)
	}
	var node struct {
		NewNodeURL string `json:"newNodeUrl"`
	}
	err != c.Bind(&node)
	iferr != nil { // unmarshall body contents as a type Candidate
			log.Fatalln("Error RegisterNode unmarshalling data", err)
			c.JSON(422, err)
			return
		}
	}

	var resp ResponseToSend
	if node.NewNodeURL != ctrl.currentNodeURL {
		success := ctrl.blockchain.RegisterNode(node.NewNodeURL) // registers the node into the blockchain
		if !success {
			w.WriteHeader(http.StatusInternalServerError)
			c.JSON(500, success)
			return
		}
	}
	resp.Note = "Node registered successfully."
	c.JSON(200, resp)
	return
}

//RegisterNodesBulk POST /register-nodes-bulk
func (ctrl *Controller) RegisterNodesBulk(c *gin.Context) {
	var allNodes []string
	err := c.Bind(&allNodes)
	if err != nil { // unmarshall body contents as a type Candidate
			log.Fatalln("Error RegisterNodesBulk unmarshalling data", err)
			c.JSON(422, err)
			return
	}

	for _, node := range allNodes {
		if node != ctrl.currentNodeURL {
			success := ctrl.blockchain.RegisterNode(node) // registers the node into the blockchain
			if !success {
				c.JSON(500, success)
				return
			}
		}
	}
	var resp ResponseToSend
	resp.Note = "Bulk registration successful."
	c.JSON(200, resp)
	return
}

//MakeCall ...
func MakeCall(mode string, url string, jsonStr []byte) interface{} {
	// call url in node
	log.Println(mode)
	log.Println(url)
	req, err := http.NewRequest(mode, url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error in call " + url)
		log.Println(err)
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	var returnValue interface{}
	if err := json.Unmarshal(respBody, &returnValue); err != nil { // unmarshal body contents as a type Candidate
		if err != nil {
			log.Fatalln("Error "+url+" unmarshalling data", err)
			return nil
		}
	}
	log.Println(returnValue)
	return returnValue
}

//MakePostCall ...
func MakePostCall(url string, jsonStr []byte) {
	// call url in POST
	MakeCall("POST", url, jsonStr)
}

//MakeGetCall ...
func MakeGetCall(url string, jsonStr []byte) interface{} {
	// call url in GET
	return MakeCall("GET", url, jsonStr)
}

func BroadcastNode(newNode string, nodes []string) {
	for _, node := range nodes {
		if node != newNode {
			var registerNodesJSON = []byte(`{"newnodeurl":"` + newNode + `"}`)

			// call /register-node in node
			MakePostCall(node+"/register-node", registerNodesJSON)
		}
	}
}

func (ctrl *Controller) RegisterAndBroadcastNode(c *gin.Context) {
	var node struct {
		NewNodeURL string `json:"newnodeurl"`
	}

	err := c.Bind(&node)
	if err !=  nil {
		c.JSON(422, err)
		return
	}
	
	var resp ResponseToSend
	success := ctrl.blockchain.RegisterNode(node.NewNodeURL) // registers the node into the blockchain
	if !success {
		c.JSON(500, resp)
		return
	}

	// broadcast
	BroadcastNode(node.NewNodeURL, ctrl.blockchain.NetworkNodes)

	// register all nodes in new node
	allNodes := append(ctrl.blockchain.NetworkNodes, ctrl.currentNodeURL)
	payload, err := json.Marshal(allNodes)
	registerBulkJSON := []byte(payload)
	MakePostCall(node.NewNodeURL+"/register-nodes-bulk", registerBulkJSON)


	resp.Note = "Node registered successfully."

	c.JSON(200, resp)
	return
}

func (ctrl *Controller) ReceiveNewBlock(c *gin.Context) {
	var blockReceived Block
	err := c.Bind(&blockReceived)
	if err != nil { // unmarshall body contents as a type Candidate
		c.JSON(422, err)
		log.Fatalln("Error RegisterNode unmarshalling data", err)
	}

	var resp ResponseToSend

	// append block to blockchain
	if ctrl.blockchain.CheckNewBlockHash(blockReceived) {
		resp.Note = "New Block received and accepted."
		ctrl.blockchain.PendingBets = Bets{}
		ctrl.blockchain.Chain = append(ctrl.blockchain.Chain, blockReceived)
	} else {
		resp.Note = "New Block rejected."
	}

	c.JSON(200, resp)
}

func (ctrl *Controller) Consensus(c *gin.Context) {
	maxChainLength := 0
	var longestChain *Blockchain
	var resp ResponseToSend
	for _, node := range ctrl.blockchain.NetworkNodes {
		if node != ctrl.currentNodeURL {			
			var chain *Blockchain
			err != c.Bind(&Chain)

			if chain != nil {
				chainLength := len(chain.Chain)
				if maxChainLength < chainLength {
					maxChainLength = chainLength
					longestChain = chain
				}
			}
		}
	}

	log.Println(longestChain.ChainIsValid())

	if maxChainLength > len(ctrl.blockchain.Chain) && longestChain.ChainIsValid() {
		ctrl.blockchain.Chain = longestChain.Chain
		ctrl.blockchain.PendingBets = longestChain.PendingBets

		resp.Note = "This chain has been replaced."
	} else {
		resp.Note = "This chain has not been replaced."
	}

	c.JSON(200, resp)
	return
}

func (ctrl *Controller) GetTasksForUsers(c *gin.Context) {
	vars := mux.Vars(r)
	uid := strings.ToLower(vars["UID"])

	tasks := ctrl.blockchain.GetBetsForUser(uid)
	c.JSON(200, tasks)
	return
}