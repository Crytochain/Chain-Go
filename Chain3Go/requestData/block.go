// requestData project block.go
package requestData

type Block struct {
	ParentHash        string `json:"parentHash"`
	UncleHash         string `json:"sha3Uncles"`
	Coinbase          string `json:"miner"`
	StateRoot         string `json:"stateRoot"`
	TransactionsRoot  string `json:"transactionsRoot"`
	ReceiptHash       string `json:"receiptsRoot"`
	Bloom             string `json:"logsBloom"`
	Difficulty        string `json:"difficulty"`
	Size              string `json:"size"`
	TotalDifficulty   string `json:"totalDifficulty"`
	Number            string `json:"number"`
	GasLimit          string `json:"gasLimit"`
	GasUsed           string `json:"gasUsed"`
	Time              string `json:"timestamp"`
	Extra             string `json:"extraData"`
	MixDigest         string `json:"mixHash"`
	Nonce             string `json:"nonce"`
	Hash              string `json:"hash"`
	TransactionsFalse []string
	TransactionsTrue  []TransactionResponse
}

type ScsBlock struct {
	Extra            string                `json:"extraData"`
	Hash             string                `json:"hash"`
	Number           string                `json:"number"`
	ParentHash       string                `json:"parentHash"`
	ReceiptHash      string                `json:"receiptsRoot"`
	StateRoot        string                `json:"stateRoot"`
	Time             string                `json:"timestamp"`
	TransactionsRoot string                `json:"transactionsRoot"`
	TransactionsTrue []TransactionResponse `json:"transactions"`
}
