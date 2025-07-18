package services

import (
	"vfinance-api/internal/blockchain"
	"vfinance-api/internal/models"

	"gorm.io/gorm"
)

type ContractService struct {
	db               *gorm.DB
	blockchainClient *blockchain.Client
	metadataService  *MetadataService
}

func NewContractService(db *gorm.DB, blockchainClient *blockchain.Client, metadataService *MetadataService) *ContractService {
	return &ContractService{
		db:               db,
		blockchainClient: blockchainClient,
		metadataService:  metadataService,
	}
}

func (s *ContractService) GetCompleteContract(regConId string) (*models.CompleteContractData, error) {
	// Buscar dados on-chain
	onChainData, err := s.blockchainClient.GetContract(regConId)
	if err != nil {
		return nil, err
	}

	// Buscar metadados off-chain
	offChainData, err := s.metadataService.GetMetadata(onChainData.MetadataHash)
	if err != nil {
		return nil, err
	}

	response := &models.CompleteContractData{
		Success: true,
	}
	response.Data.OnChain = *onChainData
	response.Data.OffChain = *offChainData

	return response, nil
}

func (s *ContractService) GetActiveContracts(offset, limit uint64) ([]string, error) {
	return s.blockchainClient.GetActiveContracts(offset, limit)
}

func (s *ContractService) GetContractByHash(hash string) (*models.CompleteContractData, error) {
	// Buscar no banco local primeiro
	var registry models.ContractRegistry
	if err := s.db.First(&registry, "metadata_hash = ?", hash).Error; err != nil {
		return nil, err
	}

	return s.GetCompleteContract(registry.RegConId)
}

func (s *ContractService) GetStats() (map[string]interface{}, error) {
	totalContracts, err := s.blockchainClient.GetTotalContracts()
	if err != nil {
		return nil, err
	}

	var activeCount int64
	s.db.Model(&models.ContractRegistry{}).Where("status = ?", "active").Count(&activeCount)

	return map[string]interface{}{
		"totalContracts":  totalContracts,
		"activeContracts": activeCount,
	}, nil
}

func (s *ContractService) RegisterContract(regConId, numeroContrato, dataContrato string, vehicleData models.VehicleData) (*models.ContractRegistrationResponse, error) {
	// Registrar contrato no blockchain
	txHash, metadataHash, err := s.blockchainClient.RegisterContract(regConId, numeroContrato, dataContrato)
	if err != nil {
		return nil, err
	}

	// Armazenar metadados no banco de dados
	err = s.metadataService.StoreMetadata(metadataHash, vehicleData)
	if err != nil {
		return nil, err
	}

	// Criar registro no banco local
	registry := models.ContractRegistry{
		RegConId:     regConId,
		MetadataHash: metadataHash,
		BlockchainTx: txHash,
		Status:       "active",
	}

	err = s.db.Create(&registry).Error
	if err != nil {
		return nil, err
	}

	return &models.ContractRegistrationResponse{
		Success:      true,
		Message:      "Contrato registrado com sucesso",
		RegConId:     regConId,
		MetadataHash: metadataHash,
		TxHash:       txHash,
	}, nil
}
