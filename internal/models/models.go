package models

import (
	"encoding/json"
	"time"
)

// VehicleMetadata representa os metadados do veículo
type VehicleMetadata struct {
	Hash        string          `gorm:"primaryKey;size:64" json:"hash"`
	VehicleData json.RawMessage `gorm:"type:jsonb;not null" json:"vehicle_data"`
	CreatedAt   time.Time       `gorm:"default:now()" json:"created_at"`
	UpdatedAt   time.Time       `gorm:"default:now()" json:"updated_at"`
}

// ContractRegistry representa o registro do contrato
type ContractRegistry struct {
	RegConId     string           `gorm:"primaryKey;size:50" json:"reg_con_id"`
	MetadataHash string           `gorm:"size:64;index" json:"metadata_hash"`
	BlockchainTx string           `gorm:"size:66" json:"blockchain_tx"`
	Status       string           `gorm:"size:20;default:active;index" json:"status"`
	CreatedAt    time.Time        `gorm:"default:now();index" json:"created_at"`
	Metadata     *VehicleMetadata `gorm:"foreignKey:MetadataHash;references:Hash" json:"metadata,omitempty"`
}

// VehicleData representa a estrutura dos dados do veículo
type VehicleData struct {
	Make           string         `json:"make"`
	Model          string         `json:"model"`
	Year           int            `json:"year"`
	VIN            string         `json:"vin"`
	Color          string         `json:"color"`
	Engine         string         `json:"engine"`
	Transmission   string         `json:"transmission"`
	Mileage        int            `json:"mileage"`
	Price          float64        `json:"price"`
	FinancingTerms FinancingTerms `json:"financingTerms"`
}

// FinancingTerms representa os termos de financiamento
type FinancingTerms struct {
	DownPayment  float64 `json:"downPayment"`
	LoanAmount   float64 `json:"loanAmount"`
	InterestRate float64 `json:"interestRate"`
	TermMonths   int     `json:"termMonths"`
}

// ContractRecord representa os dados on-chain do contrato
type ContractRecord struct {
	RegConId       string `json:"regConId"`
	NumeroContrato string `json:"numeroContrato"`
	DataContrato   string `json:"dataContrato"`
	MetadataHash   string `json:"metadataHash"`
	Timestamp      uint64 `json:"timestamp"`
	RegisteredBy   string `json:"registeredBy"`
	Active         bool   `json:"active"`
}

// CompleteContractData representa a resposta completa com dados on-chain e off-chain
type CompleteContractData struct {
	Success bool `json:"success"`
	Data    struct {
		OnChain  ContractRecord `json:"onChain"`
		OffChain VehicleData    `json:"offChain"`
	} `json:"data"`
}
