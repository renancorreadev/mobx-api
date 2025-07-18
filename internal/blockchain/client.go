package blockchain

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"math/big"
	"strings"

	"vfinance-api/internal/models"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Client struct {
	ethClient       *ethclient.Client
	contractAddress common.Address
	privateKey      *ecdsa.PrivateKey
	contractABI     abi.ABI
}

const contractABI = `[
        {
            "type": "constructor",
            "inputs": [
                {
                    "name": "initialOwner",
                    "type": "address",
                    "internalType": "address"
                },
                {
                    "name": "_apiServerUrl",
                    "type": "string",
                    "internalType": "string"
                },
                {
                    "name": "_apiServerAddress",
                    "type": "address",
                    "internalType": "address"
                }
            ],
            "stateMutability": "nonpayable"
        },
        {
            "type": "function",
            "name": "apiServerAddress",
            "inputs": [],
            "outputs": [
                {
                    "name": "",
                    "type": "address",
                    "internalType": "address"
                }
            ],
            "stateMutability": "view"
        },
        {
            "type": "function",
            "name": "apiServerUrl",
            "inputs": [],
            "outputs": [
                {
                    "name": "",
                    "type": "string",
                    "internalType": "string"
                }
            ],
            "stateMutability": "view"
        },
        {
            "type": "function",
            "name": "contractExists",
            "inputs": [
                {
                    "name": "",
                    "type": "string",
                    "internalType": "string"
                }
            ],
            "outputs": [
                {
                    "name": "",
                    "type": "bool",
                    "internalType": "bool"
                }
            ],
            "stateMutability": "view"
        },
        {
            "type": "function",
            "name": "contractIds",
            "inputs": [
                {
                    "name": "",
                    "type": "uint256",
                    "internalType": "uint256"
                }
            ],
            "outputs": [
                {
                    "name": "",
                    "type": "string",
                    "internalType": "string"
                }
            ],
            "stateMutability": "view"
        },
        {
            "type": "function",
            "name": "contracts",
            "inputs": [
                {
                    "name": "",
                    "type": "string",
                    "internalType": "string"
                }
            ],
            "outputs": [
                {
                    "name": "regConId",
                    "type": "string",
                    "internalType": "string"
                },
                {
                    "name": "numeroContrato",
                    "type": "string",
                    "internalType": "string"
                },
                {
                    "name": "dataContrato",
                    "type": "string",
                    "internalType": "string"
                },
                {
                    "name": "metadataHash",
                    "type": "bytes32",
                    "internalType": "bytes32"
                },
                {
                    "name": "timestamp",
                    "type": "uint256",
                    "internalType": "uint256"
                },
                {
                    "name": "registeredBy",
                    "type": "address",
                    "internalType": "address"
                },
                {
                    "name": "active",
                    "type": "bool",
                    "internalType": "bool"
                }
            ],
            "stateMutability": "view"
        },
        {
            "type": "function",
            "name": "doesContractExist",
            "inputs": [
                {
                    "name": "regConId",
                    "type": "string",
                    "internalType": "string"
                }
            ],
            "outputs": [
                {
                    "name": "",
                    "type": "bool",
                    "internalType": "bool"
                }
            ],
            "stateMutability": "view"
        },
        {
            "type": "function",
            "name": "doesHashExist",
            "inputs": [
                {
                    "name": "metadataHash",
                    "type": "bytes32",
                    "internalType": "bytes32"
                }
            ],
            "outputs": [
                {
                    "name": "",
                    "type": "bool",
                    "internalType": "bool"
                }
            ],
            "stateMutability": "view"
        },
        {
            "type": "function",
            "name": "getActiveContracts",
            "inputs": [
                {
                    "name": "offset",
                    "type": "uint256",
                    "internalType": "uint256"
                },
                {
                    "name": "limit",
                    "type": "uint256",
                    "internalType": "uint256"
                }
            ],
            "outputs": [
                {
                    "name": "",
                    "type": "string[]",
                    "internalType": "string[]"
                }
            ],
            "stateMutability": "view"
        },
        {
            "type": "function",
            "name": "getContract",
            "inputs": [
                {
                    "name": "regConId",
                    "type": "string",
                    "internalType": "string"
                }
            ],
            "outputs": [
                {
                    "name": "",
                    "type": "tuple",
                    "internalType": "struct VFinanceRegistry.ContractRecord",
                    "components": [
                        {
                            "name": "regConId",
                            "type": "string",
                            "internalType": "string"
                        },
                        {
                            "name": "numeroContrato",
                            "type": "string",
                            "internalType": "string"
                        },
                        {
                            "name": "dataContrato",
                            "type": "string",
                            "internalType": "string"
                        },
                        {
                            "name": "metadataHash",
                            "type": "bytes32",
                            "internalType": "bytes32"
                        },
                        {
                            "name": "timestamp",
                            "type": "uint256",
                            "internalType": "uint256"
                        },
                        {
                            "name": "registeredBy",
                            "type": "address",
                            "internalType": "address"
                        },
                        {
                            "name": "active",
                            "type": "bool",
                            "internalType": "bool"
                        }
                    ]
                }
            ],
            "stateMutability": "view"
        },
        {
            "type": "function",
            "name": "getContractByHash",
            "inputs": [
                {
                    "name": "metadataHash",
                    "type": "bytes32",
                    "internalType": "bytes32"
                }
            ],
            "outputs": [
                {
                    "name": "",
                    "type": "tuple",
                    "internalType": "struct VFinanceRegistry.ContractRecord",
                    "components": [
                        {
                            "name": "regConId",
                            "type": "string",
                            "internalType": "string"
                        },
                        {
                            "name": "numeroContrato",
                            "type": "string",
                            "internalType": "string"
                        },
                        {
                            "name": "dataContrato",
                            "type": "string",
                            "internalType": "string"
                        },
                        {
                            "name": "metadataHash",
                            "type": "bytes32",
                            "internalType": "bytes32"
                        },
                        {
                            "name": "timestamp",
                            "type": "uint256",
                            "internalType": "uint256"
                        },
                        {
                            "name": "registeredBy",
                            "type": "address",
                            "internalType": "address"
                        },
                        {
                            "name": "active",
                            "type": "bool",
                            "internalType": "bool"
                        }
                    ]
                }
            ],
            "stateMutability": "view"
        },
        {
            "type": "function",
            "name": "getContractIdByIndex",
            "inputs": [
                {
                    "name": "index",
                    "type": "uint256",
                    "internalType": "uint256"
                }
            ],
            "outputs": [
                {
                    "name": "",
                    "type": "string",
                    "internalType": "string"
                }
            ],
            "stateMutability": "view"
        },
        {
            "type": "function",
            "name": "getHashByRegConId",
            "inputs": [
                {
                    "name": "regConId",
                    "type": "string",
                    "internalType": "string"
                }
            ],
            "outputs": [
                {
                    "name": "",
                    "type": "bytes32",
                    "internalType": "bytes32"
                }
            ],
            "stateMutability": "view"
        },
        {
            "type": "function",
            "name": "getMetadataUrl",
            "inputs": [
                {
                    "name": "regConId",
                    "type": "string",
                    "internalType": "string"
                }
            ],
            "outputs": [
                {
                    "name": "",
                    "type": "string",
                    "internalType": "string"
                }
            ],
            "stateMutability": "view"
        },
        {
            "type": "function",
            "name": "getTotalContracts",
            "inputs": [],
            "outputs": [
                {
                    "name": "",
                    "type": "uint256",
                    "internalType": "uint256"
                }
            ],
            "stateMutability": "view"
        },
        {
            "type": "function",
            "name": "hashToRegConId",
            "inputs": [
                {
                    "name": "",
                    "type": "bytes32",
                    "internalType": "bytes32"
                }
            ],
            "outputs": [
                {
                    "name": "",
                    "type": "string",
                    "internalType": "string"
                }
            ],
            "stateMutability": "view"
        },
        {
            "type": "function",
            "name": "owner",
            "inputs": [],
            "outputs": [
                {
                    "name": "",
                    "type": "address",
                    "internalType": "address"
                }
            ],
            "stateMutability": "view"
        },
        {
            "type": "function",
            "name": "regConIdToHash",
            "inputs": [
                {
                    "name": "",
                    "type": "string",
                    "internalType": "string"
                }
            ],
            "outputs": [
                {
                    "name": "",
                    "type": "bytes32",
                    "internalType": "bytes32"
                }
            ],
            "stateMutability": "view"
        },
        {
            "type": "function",
            "name": "registerContract",
            "inputs": [
                {
                    "name": "regConId",
                    "type": "string",
                    "internalType": "string"
                },
                {
                    "name": "numeroContrato",
                    "type": "string",
                    "internalType": "string"
                },
                {
                    "name": "dataContrato",
                    "type": "string",
                    "internalType": "string"
                }
            ],
            "outputs": [
                {
                    "name": "metadataHash",
                    "type": "bytes32",
                    "internalType": "bytes32"
                }
            ],
            "stateMutability": "nonpayable"
        },
        {
            "type": "function",
            "name": "tokenURI",
            "inputs": [
                {
                    "name": "metadataHash",
                    "type": "bytes32",
                    "internalType": "bytes32"
                }
            ],
            "outputs": [
                {
                    "name": "",
                    "type": "string",
                    "internalType": "string"
                }
            ],
            "stateMutability": "view"
        },
        {
            "type": "function",
            "name": "totalContracts",
            "inputs": [],
            "outputs": [
                {
                    "name": "",
                    "type": "uint256",
                    "internalType": "uint256"
                }
            ],
            "stateMutability": "view"
        },
        {
            "type": "function",
            "name": "updateApiServer",
            "inputs": [
                {
                    "name": "newApiServerUrl",
                    "type": "string",
                    "internalType": "string"
                },
                {
                    "name": "newApiServerAddress",
                    "type": "address",
                    "internalType": "address"
                }
            ],
            "outputs": [],
            "stateMutability": "nonpayable"
        },
        {
            "type": "function",
            "name": "updateContractStatus",
            "inputs": [
                {
                    "name": "regConId",
                    "type": "string",
                    "internalType": "string"
                },
                {
                    "name": "active",
                    "type": "bool",
                    "internalType": "bool"
                }
            ],
            "outputs": [],
            "stateMutability": "nonpayable"
        },
        {
            "type": "function",
            "name": "updateMetadataHash",
            "inputs": [
                {
                    "name": "regConId",
                    "type": "string",
                    "internalType": "string"
                },
                {
                    "name": "newMetadataHash",
                    "type": "bytes32",
                    "internalType": "bytes32"
                }
            ],
            "outputs": [],
            "stateMutability": "nonpayable"
        },
        {
            "type": "function",
            "name": "updateMetadataWithNewHash",
            "inputs": [
                {
                    "name": "regConId",
                    "type": "string",
                    "internalType": "string"
                },
                {
                    "name": "newDataContrato",
                    "type": "string",
                    "internalType": "string"
                }
            ],
            "outputs": [
                {
                    "name": "newMetadataHash",
                    "type": "bytes32",
                    "internalType": "bytes32"
                }
            ],
            "stateMutability": "nonpayable"
        },
        {
            "type": "event",
            "name": "ApiServerUpdated",
            "inputs": [
                {
                    "name": "oldUrl",
                    "type": "string",
                    "indexed": false,
                    "internalType": "string"
                },
                {
                    "name": "newUrl",
                    "type": "string",
                    "indexed": false,
                    "internalType": "string"
                },
                {
                    "name": "oldAddress",
                    "type": "address",
                    "indexed": false,
                    "internalType": "address"
                },
                {
                    "name": "newAddress",
                    "type": "address",
                    "indexed": false,
                    "internalType": "address"
                }
            ],
            "anonymous": false
        },
        {
            "type": "event",
            "name": "ContractRegistered",
            "inputs": [
                {
                    "name": "regConId",
                    "type": "string",
                    "indexed": true,
                    "internalType": "string"
                },
                {
                    "name": "numeroContrato",
                    "type": "string",
                    "indexed": true,
                    "internalType": "string"
                },
                {
                    "name": "metadataHash",
                    "type": "bytes32",
                    "indexed": true,
                    "internalType": "bytes32"
                },
                {
                    "name": "registeredBy",
                    "type": "address",
                    "indexed": false,
                    "internalType": "address"
                },
                {
                    "name": "timestamp",
                    "type": "uint256",
                    "indexed": false,
                    "internalType": "uint256"
                }
            ],
            "anonymous": false
        },
        {
            "type": "event",
            "name": "ContractStatusUpdated",
            "inputs": [
                {
                    "name": "regConId",
                    "type": "string",
                    "indexed": true,
                    "internalType": "string"
                },
                {
                    "name": "active",
                    "type": "bool",
                    "indexed": false,
                    "internalType": "bool"
                },
                {
                    "name": "timestamp",
                    "type": "uint256",
                    "indexed": false,
                    "internalType": "uint256"
                }
            ],
            "anonymous": false
        },
        {
            "type": "event",
            "name": "MetadataHashUpdated",
            "inputs": [
                {
                    "name": "regConId",
                    "type": "string",
                    "indexed": true,
                    "internalType": "string"
                },
                {
                    "name": "oldHash",
                    "type": "bytes32",
                    "indexed": false,
                    "internalType": "bytes32"
                },
                {
                    "name": "newHash",
                    "type": "bytes32",
                    "indexed": false,
                    "internalType": "bytes32"
                },
                {
                    "name": "timestamp",
                    "type": "uint256",
                    "indexed": false,
                    "internalType": "uint256"
                }
            ],
            "anonymous": false
        },
        {
            "type": "error",
            "name": "ContractAlreadyExists",
            "inputs": [
                {
                    "name": "regConId",
                    "type": "string",
                    "internalType": "string"
                }
            ]
        },
        {
            "type": "error",
            "name": "ContractNotFound",
            "inputs": [
                {
                    "name": "regConId",
                    "type": "string",
                    "internalType": "string"
                }
            ]
        },
        {
            "type": "error",
            "name": "HashAlreadyUsed",
            "inputs": [
                {
                    "name": "hash",
                    "type": "bytes32",
                    "internalType": "bytes32"
                }
            ]
        },
        {
            "type": "error",
            "name": "InvalidInput",
            "inputs": [
                {
                    "name": "paramName",
                    "type": "string",
                    "internalType": "string"
                }
            ]
        },
        {
            "type": "error",
            "name": "UnauthorizedAccess",
            "inputs": []
        }
    ]`

func NewClient(rpcURL, contractAddr, privateKeyHex string) (*Client, error) {
	ethClient, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(privateKeyHex, "0x"))
	if err != nil {
		return nil, err
	}

	contractAbi, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		return nil, err
	}

	return &Client{
		ethClient:       ethClient,
		contractAddress: common.HexToAddress(contractAddr),
		privateKey:      privateKey,
		contractABI:     contractAbi,
	}, nil
}

func (c *Client) GetContract(regConId string) (*models.ContractRecord, error) {
	// Preparar dados para chamada
	data, err := c.contractABI.Pack("getContract", regConId)
	if err != nil {
		return nil, err
	}

	// Fazer chamada ao contrato
	result, err := c.ethClient.CallContract(context.Background(), ethereum.CallMsg{
		To:   &c.contractAddress,
		Data: data,
	}, nil)
	if err != nil {
		return nil, err
	}

	// Decodificar resultado
	var contractData struct {
		RegConId       string
		NumeroContrato string
		DataContrato   string
		MetadataHash   [32]byte
		Timestamp      *big.Int
		RegisteredBy   common.Address
		Active         bool
	}

	err = c.contractABI.UnpackIntoInterface(&contractData, "getContract", result)
	if err != nil {
		return nil, err
	}

	// Mapear para modelo
	contractRecord := &models.ContractRecord{
		RegConId:       contractData.RegConId,
		NumeroContrato: contractData.NumeroContrato,
		DataContrato:   contractData.DataContrato,
		MetadataHash:   hex.EncodeToString(contractData.MetadataHash[:]),
		Timestamp:      contractData.Timestamp.Uint64(),
		RegisteredBy:   contractData.RegisteredBy.Hex(),
		Active:         contractData.Active,
	}

	return contractRecord, nil
}

func (c *Client) RegisterContract(regConId, numeroContrato, dataContrato string) (string, string, error) {
	auth, err := bind.NewKeyedTransactorWithChainID(c.privateKey, big.NewInt(1337))
	if err != nil {
		return "", "", err
	}

	// Preparar dados para transação
	data, err := c.contractABI.Pack("registerContract", regConId, numeroContrato, dataContrato)
	if err != nil {
		return "", "", err
	}

	// Estimar gas
	gasLimit, err := c.ethClient.EstimateGas(context.Background(), ethereum.CallMsg{
		To:   &c.contractAddress,
		Data: data,
	})
	if err != nil {
		return "", "", err
	}

	auth.GasLimit = gasLimit

	// Criar transação
	tx := types.NewTransaction(
		auth.Nonce.Uint64(),
		c.contractAddress,
		auth.Value,
		gasLimit,
		auth.GasPrice,
		data,
	)

	// Assinar transação
	signedTx, err := auth.Signer(auth.From, tx)
	if err != nil {
		return "", "", err
	}

	// Enviar transação
	err = c.ethClient.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", "", err
	}

	// Aguardar confirmação da transação
	receipt, err := bind.WaitMined(context.Background(), c.ethClient, signedTx)
	if err != nil {
		return signedTx.Hash().Hex(), "", err
	}

	// Decodificar logs para obter o hash gerado
	if len(receipt.Logs) > 0 {
		for _, log := range receipt.Logs {
			if log.Topics[0] == c.contractABI.Events["ContractRegistered"].ID {
				// Decodificar o evento para obter o metadataHash
				var event struct {
					RegConId       string
					NumeroContrato string
					MetadataHash   [32]byte
					RegisteredBy   common.Address
					Timestamp      *big.Int
				}

				err = c.contractABI.UnpackIntoInterface(&event, "ContractRegistered", log.Data)
				if err == nil {
					return signedTx.Hash().Hex(), hex.EncodeToString(event.MetadataHash[:]), nil
				}
			}
		}
	}

	// Se não conseguir decodificar o evento, buscar o hash diretamente
	metadataHash, err := c.GetHashByRegConId(regConId)
	if err != nil {
		return signedTx.Hash().Hex(), "", err
	}

	return signedTx.Hash().Hex(), metadataHash, nil
}

func (c *Client) GetHashByRegConId(regConId string) (string, error) {
	// Preparar dados para chamada
	data, err := c.contractABI.Pack("getHashByRegConId", regConId)
	if err != nil {
		return "", err
	}

	// Fazer chamada ao contrato
	result, err := c.ethClient.CallContract(context.Background(), ethereum.CallMsg{
		To:   &c.contractAddress,
		Data: data,
	}, nil)
	if err != nil {
		return "", err
	}

	// Decodificar resultado
	var hash [32]byte
	err = c.contractABI.UnpackIntoInterface(&hash, "getHashByRegConId", result)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hash[:]), nil
}

func (c *Client) GetActiveContracts(offset, limit uint64) ([]string, error) {
	// Preparar dados para chamada
	data, err := c.contractABI.Pack("getActiveContracts", big.NewInt(int64(offset)), big.NewInt(int64(limit)))
	if err != nil {
		return nil, err
	}

	// Fazer chamada ao contrato
	result, err := c.ethClient.CallContract(context.Background(), ethereum.CallMsg{
		To:   &c.contractAddress,
		Data: data,
	}, nil)
	if err != nil {
		return nil, err
	}

	// Decodificar resultado
	var contracts []string
	err = c.contractABI.UnpackIntoInterface(&contracts, "getActiveContracts", result)
	if err != nil {
		return nil, err
	}

	return contracts, nil
}

func (c *Client) GetTotalContracts() (uint64, error) {
	// Preparar dados para chamada
	data, err := c.contractABI.Pack("getTotalContracts")
	if err != nil {
		return 0, err
	}

	// Fazer chamada ao contrato
	result, err := c.ethClient.CallContract(context.Background(), ethereum.CallMsg{
		To:   &c.contractAddress,
		Data: data,
	}, nil)
	if err != nil {
		return 0, err
	}

	// Decodificar resultado
	var total *big.Int
	err = c.contractABI.UnpackIntoInterface(&total, "getTotalContracts", result)
	if err != nil {
		return 0, err
	}

	return total.Uint64(), nil
}

func (c *Client) DoesContractExist(regConId string) (bool, error) {
	data, err := c.contractABI.Pack("doesContractExist", regConId)
	if err != nil {
		return false, err
	}

	result, err := c.ethClient.CallContract(context.Background(), ethereum.CallMsg{
		To:   &c.contractAddress,
		Data: data,
	}, nil)
	if err != nil {
		return false, err
	}

	var exists bool
	err = c.contractABI.UnpackIntoInterface(&exists, "doesContractExist", result)
	if err != nil {
		return false, err
	}

	return exists, nil
}
