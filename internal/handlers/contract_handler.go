package handlers

import (
	"net/http"
	"strconv"
	"vfinance-api/internal/models"
	"vfinance-api/internal/services"

	"github.com/gin-gonic/gin"
)

type ContractHandler struct {
	contractService *services.ContractService
}

func NewContractHandler(contractService *services.ContractService) *ContractHandler {
	return &ContractHandler{contractService: contractService}
}

func (h *ContractHandler) GetContract(c *gin.Context) {
	regConId := c.Param("regConId")
	if regConId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "regConId é obrigatório"})
		return
	}

	contractData, err := h.contractService.GetCompleteContract(regConId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contrato não encontrado"})
		return
	}

	c.JSON(http.StatusOK, contractData)
}

func (h *ContractHandler) GetActiveContracts(c *gin.Context) {
	offsetStr := c.DefaultQuery("offset", "0")
	limitStr := c.DefaultQuery("limit", "10")

	offset, err := strconv.ParseUint(offsetStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Offset inválido"})
		return
	}

	limit, err := strconv.ParseUint(limitStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Limit inválido"})
		return
	}

	contracts, err := h.contractService.GetActiveContracts(offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": contracts})
}

func (h *ContractHandler) GetContractByHash(c *gin.Context) {
	hash := c.Param("hash")
	if hash == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Hash é obrigatório"})
		return
	}

	contractData, err := h.contractService.GetContractByHash(hash)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contrato não encontrado"})
		return
	}

	c.JSON(http.StatusOK, contractData)
}

func (h *ContractHandler) GetStats(c *gin.Context) {
	stats, err := h.contractService.GetStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": stats})
}

func (h *ContractHandler) RegisterContract(c *gin.Context) {
	var request models.ContractRegistrationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validar campos obrigatórios
	if request.RegConId == "" || request.NumeroContrato == "" || request.DataContrato == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "regConId, numeroContrato e dataContrato são obrigatórios"})
		return
	}

	// Registrar contrato
	response, err := h.contractService.RegisterContract(
		request.RegConId,
		request.NumeroContrato,
		request.DataContrato,
		request.VehicleData,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}
