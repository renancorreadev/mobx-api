package handlers

import (
	"net/http"
	"vfinance-api/internal/models"
	"vfinance-api/internal/services"

	"github.com/gin-gonic/gin"
)

type MetadataHandler struct {
	metadataService *services.MetadataService
}

func NewMetadataHandler(metadataService *services.MetadataService) *MetadataHandler {
	return &MetadataHandler{metadataService: metadataService}
}

func (h *MetadataHandler) StoreMetadata(c *gin.Context) {
	hash := c.Param("hash")
	if hash == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Hash é obrigatório"})
		return
	}

	var vehicleData models.VehicleData
	if err := c.ShouldBindJSON(&vehicleData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.metadataService.StoreMetadata(hash, vehicleData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Metadados armazenados com sucesso"})
}

func (h *MetadataHandler) GetMetadata(c *gin.Context) {
	hash := c.Param("hash")
	if hash == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Hash é obrigatório"})
		return
	}

	vehicleData, err := h.metadataService.GetMetadata(hash)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Metadados não encontrados"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": vehicleData})
}

func (h *MetadataHandler) UpdateMetadata(c *gin.Context) {
	hash := c.Param("hash")
	if hash == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Hash é obrigatório"})
		return
	}

	var vehicleData models.VehicleData
	if err := c.ShouldBindJSON(&vehicleData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.metadataService.UpdateMetadata(hash, vehicleData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Metadados atualizados com sucesso"})
}

func (h *MetadataHandler) DeleteMetadata(c *gin.Context) {
	hash := c.Param("hash")
	if hash == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Hash é obrigatório"})
		return
	}

	if err := h.metadataService.DeleteMetadata(hash); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Metadados removidos com sucesso"})
}
