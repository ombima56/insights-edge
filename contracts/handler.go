package contracts

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/ombima56/insights-edge/models"
)

// SmartContractHandler handles interactions with blockchain contracts
type SmartContractHandler struct {
	DB *sql.DB
}

// NewSmartContractHandler creates a new smart contract handler
func NewSmartContractHandler(db *sql.DB) *SmartContractHandler {
	return &SmartContractHandler{
		DB: db,
	}
}

// GetContractByName retrieves a contract by its name
func (h *SmartContractHandler) GetContractByName(name string) (*models.Contract, error) {
	var contract models.Contract
	err := h.DB.QueryRow(
		"SELECT id, contract_address, contract_name, contract_abi, contract_bytecode, deployed_at, owner_user_id, network FROM smart_contracts WHERE contract_name = ?",
		name,
	).Scan(
		&contract.ID,
		&contract.ContractAddress,
		&contract.ContractName,
		&contract.ContractABI,
		&contract.ContractBytecode,
		&contract.DeployedAt,
		&contract.OwnerUserID,
		&contract.Network,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get contract: %v", err)
	}
	return &contract, nil
}

// LogContractEvent records an event from a smart contract
func (h *SmartContractHandler) LogContractEvent(contractID int, eventName string, eventData map[string]interface{}, blockNumber int64, txHash string, logIndex int) error {
	// Serialize event data to JSON
	eventDataJSON, err := json.Marshal(eventData)
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %v", err)
	}
	
	// Insert event into the database
	_, err = h.DB.Exec(
		"INSERT INTO contract_events (contract_id, event_name, event_data, block_number, transaction_hash, log_index) VALUES (?, ?, ?, ?, ?, ?)",
		contractID,
		eventName,
		string(eventDataJSON),
		blockNumber,
		txHash,
		logIndex,
	)
	if err != nil {
		return fmt.Errorf("failed to log contract event: %v", err)
	}
	
	return nil
}
