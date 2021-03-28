package blockchain

type Task struct {
	UID   string `json:"playername"`
	TaskID      string `json:"matchid"`
	Reward int  	`json:"teamonescore"`
}

type Tasks []Task

type Block struct {
	Index             int    `json:"index"`
	Timestamp         int64  `json:"timestamp"`
	Tasks              Tasks `json:"tasks"`
	Nonce             int    `json:"nonce"`
	Hash              string `json:"hash"`
	PreviousBlockHash string `json:"previousblockhash"`
}

type Blocks []Block

type Blockchain struct {
	Chain        Blocks   `json:"chain"`
	PendingTasks  Tasks     `json:"pending_bets"`
	NetworkNodes []string `json:"network_nodes"`
}

type BlockData struct {
	Index string
	Tasks Tasks
}