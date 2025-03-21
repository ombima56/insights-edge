package models

// Contract represents a smart contract from the database
type Contract struct {
	ID              int    `json:"id"`
	ContractAddress string `json:"contract_address"`
	ContractName    string `json:"contract_name"`
	ContractABI     string `json:"contract_abi"`
	ContractBytecode string `json:"contract_bytecode"`
	DeployedAt      string `json:"deployed_at"`
	OwnerUserID     int    `json:"owner_user_id"`
	Network         string `json:"network"`
}

// ContractEvent represents an event from a smart contract
type ContractEvent struct {
	ID             int64  `json:"id"`
	ContractID     int    `json:"contract_id"`
	EventName      string `json:"event_name"`
	EventData      string `json:"event_data"`
	BlockNumber    int64  `json:"block_number"`
	TransactionHash string `json:"transaction_hash"`
	LogIndex       int    `json:"log_index"`
	Timestamp      string `json:"timestamp"`
}
