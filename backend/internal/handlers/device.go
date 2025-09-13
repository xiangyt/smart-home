package handlers

import (
	"net/http"
	"smart-home-backend/internal/models"
	"smart-home-backend/internal/services"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type DeviceHandler struct {
	deviceService *services.DeviceService
}

func NewDeviceHandler(deviceService *services.DeviceService) *DeviceHandler {
	return &DeviceHandler{
		deviceService: deviceService,
	}
}

// GetDevices 获取所有设备（使用缓存）
func (h *DeviceHandler) GetDevices(c *gin.Context) {
	// 添加缓存控制参数
	useCache := c.Query("cache") != "false"

	devices, err := h.deviceService.GetDevicesWithCache(useCache)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get devices",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   devices,
		"count":  len(devices),
		"cached": useCache,
	})
}

// GetDeviceState 获取设备状态
func (h *DeviceHandler) GetDeviceState(c *gin.Context) {
	entityID := c.Param("entity_id")

	device, err := h.deviceService.GetDeviceByEntityID(entityID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Device not found",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, device)
}

// ControlDevice 控制设备
func (h *DeviceHandler) ControlDevice(c *gin.Context) {
	entityID := c.Param("entity_id")

	var request struct {
		Service string                 `json:"service" binding:"required"`
		Data    map[string]interface{} `json:"data"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	err := h.deviceService.ControlDevice(entityID, request.Service, request.Data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to control device",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Device controlled successfully",
	})
}

// GetDeviceHistory 获取设备历史数据
func (h *DeviceHandler) GetDeviceHistory(c *gin.Context) {
	entityID := c.Param("entity_id")

	// 解析时间参数
	var startTime, endTime *time.Time

	if startTimeStr := c.Query("start_time"); startTimeStr != "" {
		if t, err := time.Parse(time.RFC3339, startTimeStr); err == nil {
			startTime = &t
		}
	}

	if endTimeStr := c.Query("end_time"); endTimeStr != "" {
		if t, err := time.Parse(time.RFC3339, endTimeStr); err == nil {
			endTime = &t
		}
	}

	// 如果没有指定开始时间，默认获取最近24小时的数据
	if startTime == nil {
		t := time.Now().Add(-24 * time.Hour)
		startTime = &t
	}

	history, err := h.deviceService.GetDeviceHistory(entityID, startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get device history",
			"message": err.Error(),
		})
		return
	}

	// 添加分页支持
	page := 1
	pageSize := 50

	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	start := (page - 1) * pageSize
	end := start + pageSize

	if start >= len(history) {
		history = []models.DeviceHistory{}
	} else {
		if end > len(history) {
			end = len(history)
		}
		history = history[start:end]
	}

	c.JSON(http.StatusOK, gin.H{
		"data": history,
		"pagination": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     len(history),
		},
	})
}
