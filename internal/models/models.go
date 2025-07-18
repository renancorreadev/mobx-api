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

// VehicleData representa a estrutura dos dados do contrato de financiamento
type VehicleData struct {
	RegConId                              string `json:"regConId"`
	NumeroContrato                        string `json:"numeroContrato"`
	DataContrato                          string `json:"dataContrato"`
	CnpjAgenteFinanceiro                  string `json:"cnpjAgenteFinanceiro"`
	NomeAgenteFinanceiro                  string `json:"nomeAgenteFinanceiro"`
	EnderecoAgenteFinanceiro              string `json:"enderecoAgenteFinanceiro"`
	NumeroEnderecoAgenteFinanceiro        string `json:"numeroEnderecoAgenteFinanceiro"`
	ComplementoEnderecoAgenteFinanceiro   string `json:"complementoEnderecoAgenteFinanceiro"`
	BairroEnderecoAgenteFinanceiro        string `json:"bairroEnderecoAgenteFinanceiro"`
	NomeMunicipioEnderecoAgenteFinanceiro string `json:"nomeMunicipioEnderecoAgenteFinanceiro"`
	UfEnderecoAgenteFinanceiro            string `json:"ufEnderecoAgenteFinanceiro"`
	CepEnderecoAgenteFinanceiro           string `json:"cepEnderecoAgenteFinanceiro"`
	TelefoneAgenteFinanceiro              string `json:"telefoneAgenteFinanceiro"`
	EmailAgenteFinanceiro                 string `json:"emailAgenteFinanceiro"`
	CpfCnpjProprietario                   string `json:"cpfCnpjProprietario"`
	NomeProprietario                      string `json:"nomeProprietario"`
	EnderecoProprietario                  string `json:"enderecoProprietario"`
	NumeroEnderecoProprietario            string `json:"numeroEnderecoProprietario"`
	BairroEnderecoProprietario            string `json:"bairroEnderecoProprietario"`
	NomeMunicipioProprietario             string `json:"nomeMunicipioProprietario"`
	UfEnderecoProprietario                string `json:"ufEnderecoProprietario"`
	CepEnderecoProprietario               string `json:"cepEnderecoProprietario"`
	TelefoneProprietario                  string `json:"telefoneProprietario"`
	EmailProprietario                     string `json:"emailProprietario"`
	VeiculoZeroKm                         bool   `json:"veiculoZeroKm"`
	ChassiVeiculo                         string `json:"chassiVeiculo"`
	ChassiRemarcadoVeiculo                string `json:"chassiRemarcadoVeiculo"`
	PlacaVeiculo                          string `json:"placaVeiculo"`
	TipoPlacaVeiculo                      string `json:"tipoPlacaVeiculo"`
	UfPlacaVeiculo                        string `json:"ufPlacaVeiculo"`
	RenavamVeiculo                        string `json:"renavamVeiculo"`
	AnoFabricacaoVeiculo                  string `json:"anoFabricacaoVeiculo"`
	AnoModeloVeiculo                      string `json:"anoModeloVeiculo"`
	NumeroRestricaoVeiculo                string `json:"numeroRestricaoVeiculo"`
	EspecieVeiculo                        string `json:"especieVeiculo"`
	MarcaVeiculo                          string `json:"marcaVeiculo"`
	ModeloVeiculo                         string `json:"modeloVeiculo"`
	TipoRestricacaoContrato               string `json:"tipoRestricacaoContrato"`
	UfRegistroContrato                    string `json:"ufRegistroContrato"`
	CnpjResponsavelPeloRegistro           string `json:"cnpjResponsavelPeloRegistro"`
	ValorTotalContrato                    string `json:"valorTotalContrato"`
	ValorParcelaContrato                  string `json:"valorParcelaContrato"`
	QuantidadeParcelasContrato            string `json:"quantidadeParcelasContrato"`
	TaxaJurosMesContrato                  string `json:"taxaJurosMesContrato"`
	TaxaJurosMesAnoContrato               string `json:"taxaJurosMesAnoContrato"`
	PossuiJurosMultaContrato              string `json:"possuiJurosMultaContrato"`
	TaxaJurosMultaContrato                string `json:"taxaJurosMultaContrato"`
	PossuiJurosMoraDiaContrato            string `json:"possuiJurosMoraDiaContrato"`
	TaxaJurosMoraDiaContrato              string `json:"taxaJurosMoraDiaContrato"`
	ValorCustoRegistroContrato            string `json:"valorCustoRegistroContrato"`
	ValorIofContrato                      string `json:"valorIofContrato"`
	DataVencimentoPrimeiraParcelaContrato string `json:"dataVencimentoPrimeiraParcelaContrato"`
	DataVencimentoUltimaParcelaContrato   string `json:"dataVencimentoUltimaParcelaContrato"`
	DataLiberacaoCreditoContrato          string `json:"dataLiberacaoCreditoContrato"`
	CidadeLiberacaoCreditoContrato        string `json:"cidadeLiberacaoCreditoContrato"`
	UfLiberacaoCreditoContrato            string `json:"ufLiberacaoCreditoContrato"`
	IndiceCorrecaoContrato                string `json:"indiceCorrecaoContrato"`
	NumeroGrupoConsorcioContrato          string `json:"numeroGrupoConsorcioContrato"`
	NumeroCotaConsorcioContrato           string `json:"numeroCotaConsorcioContrato"`
	IndicativoPenalidadeContrato          string `json:"indicativoPenalidadeContrato"`
	PenalidadeContrato                    string `json:"penalidadeContrato"`
	IndicativoComissaoContrato            string `json:"indicativoComissaoContrato"`
	ComissaoContrato                      string `json:"comissaoContrato"`
	CategoriaVeiculo                      string `json:"categoriaVeiculo"`
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

// ContractRegistrationRequest representa a requisição para registrar um contrato
type ContractRegistrationRequest struct {
	RegConId       string      `json:"regConId" binding:"required"`
	NumeroContrato string      `json:"numeroContrato" binding:"required"`
	DataContrato   string      `json:"dataContrato" binding:"required"`
	VehicleData    VehicleData `json:"vehicleData" binding:"required"`
}

// ContractRegistrationResponse representa a resposta do registro de contrato
type ContractRegistrationResponse struct {
	Success      bool   `json:"success"`
	Message      string `json:"message"`
	RegConId     string `json:"regConId"`
	MetadataHash string `json:"metadataHash"`
	TxHash       string `json:"txHash"`
}
