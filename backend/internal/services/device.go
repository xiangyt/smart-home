package services

import (
	"encoding/json"
	"smart-home-backend/internal/models"
	"strings"
	"time"

	"gorm.io/gorm"
)

type DeviceService struct {
	db        *gorm.DB
	haService *HomeAssistantService
	cache     map[string]interface{}
	cacheTime time.Time
}

func NewDeviceService(db *gorm.DB, haService *HomeAssistantService) *DeviceService {
	return &DeviceService{
		db:        db,
		haService: haService,
		cache:     make(map[string]interface{}),
		cacheTime: time.Time{},
	}
}

// SyncDevicesFromHA 从Home Assistant同步设备数据
func (s *DeviceService) SyncDevicesFromHA() error {
	haDevices, err := s.haService.GetStates()
	if err != nil {
		return err
	}

	for _, haDevice := range haDevices {
		// 过滤掉不需要的实体
		if s.shouldSkipEntity(haDevice.EntityID) {
			continue
		}

		device := s.convertHADeviceToDevice(haDevice)

		// 更新或创建设备记录
		var existingDevice models.Device
		result := s.db.Where("entity_id = ?", device.EntityID).First(&existingDevice)

		if result.Error == gorm.ErrRecordNotFound {
			// 创建新设备
			if err := s.db.Create(&device).Error; err != nil {
				continue // 跳过错误，继续处理其他设备
			}
		} else {
			// 更新现有设备
			device.ID = existingDevice.ID
			device.CreatedAt = existingDevice.CreatedAt
			if err := s.db.Save(&device).Error; err != nil {
				continue
			}
		}

		// 记录历史数据
		s.recordDeviceHistory(device.EntityID, device.State)
	}

	return nil
}

// GetDevices 获取所有设备（快速版本，不同步）
func (s *DeviceService) GetDevices() ([]models.Device, error) {
	var devices []models.Device
	err := s.db.Find(&devices).Error
	return devices, err
}

// GetDeviceByEntityID 根据实体ID获取设备
func (s *DeviceService) GetDeviceByEntityID(entityID string) (*models.Device, error) {
	var device models.Device
	err := s.db.Where("entity_id = ?", entityID).First(&device).Error
	if err != nil {
		return nil, err
	}
	return &device, nil
}

// ControlDevice 控制设备
func (s *DeviceService) ControlDevice(entityID, service string, serviceData map[string]interface{}) error {
	// 解析域名
	parts := strings.Split(entityID, ".")
	if len(parts) != 2 {
		return gorm.ErrInvalidData
	}
	domain := parts[0]

	// 调用Home Assistant服务
	err := s.haService.CallService(domain, service, entityID, serviceData)
	if err != nil {
		return err
	}

	// 同步设备状态
	go func() {
		time.Sleep(1 * time.Second) // 等待HA状态更新
		s.SyncDevicesFromHA()
	}()

	return nil
}

// GetDeviceHistory 获取设备历史数据
func (s *DeviceService) GetDeviceHistory(entityID string, startTime, endTime *time.Time) ([]models.DeviceHistory, error) {
	query := s.db.Where("entity_id = ?", entityID)

	if startTime != nil {
		query = query.Where("timestamp >= ?", *startTime)
	}

	if endTime != nil {
		query = query.Where("timestamp <= ?", *endTime)
	}

	var history []models.DeviceHistory
	err := query.Order("timestamp DESC").Limit(100).Find(&history).Error
	return history, err
}

// convertHADeviceToDevice 转换HA设备到本地设备模型
func (s *DeviceService) convertHADeviceToDevice(haDevice models.HADevice) models.Device {
	attributesJSON, _ := json.Marshal(haDevice.Attributes)

	lastChanged, _ := time.Parse(time.RFC3339, haDevice.LastChanged)
	lastUpdated, _ := time.Parse(time.RFC3339, haDevice.LastUpdated)

	// 从attributes中获取友好名称
	friendlyName := haDevice.EntityID
	if name, ok := haDevice.Attributes["friendly_name"].(string); ok {
		friendlyName = name
	}

	// 获取设备类别
	deviceClass := ""
	if class, ok := haDevice.Attributes["device_class"].(string); ok {
		deviceClass = class
	}

	// 获取单位
	unitOfMeasurement := ""
	if unit, ok := haDevice.Attributes["unit_of_measurement"].(string); ok {
		unitOfMeasurement = unit
	}

	// 解析域名
	domain := strings.Split(haDevice.EntityID, ".")[0]

	return models.Device{
		EntityID:          haDevice.EntityID,
		FriendlyName:      friendlyName,
		State:             haDevice.State,
		Domain:            domain,
		DeviceClass:       deviceClass,
		UnitOfMeasurement: unitOfMeasurement,
		Attributes:        string(attributesJSON),
		LastChanged:       lastChanged,
		LastUpdated:       lastUpdated,
		UpdatedAt:         time.Now(),
	}
}

// shouldSkipEntity 判断是否应该跳过某个实体
func (s *DeviceService) shouldSkipEntity(entityID string) bool {
	// 跳过一些系统实体
	skipPrefixes := []string{
		"sun.",
		"zone.",
		"automation.",
		"script.",
		"group.",
		"device_tracker.",
		"person.",
	}

	for _, prefix := range skipPrefixes {
		if strings.HasPrefix(entityID, prefix) {
			return true
		}
	}

	return false
}

// recordDeviceHistory 记录设备历史
func (s *DeviceService) recordDeviceHistory(entityID, state string) {
	history := models.DeviceHistory{
		EntityID:  entityID,
		State:     state,
		Timestamp: time.Now(),
	}
	s.db.Create(&history)
}

// GetDevicesWithCache 获取设备（带缓存）
func (s *DeviceService) GetDevicesWithCache(useCache bool) ([]models.Device, error) {
	cacheKey := "devices"
	cacheDuration := 30 * time.Second // 缓存30秒

	// 检查缓存
	if useCache && time.Since(s.cacheTime) < cacheDuration {
		if cachedDevices, exists := s.cache[cacheKey]; exists {
			if devices, ok := cachedDevices.([]models.Device); ok {
				return devices, nil
			}
		}
	}

	// 从数据库获取
	devices, err := s.GetDevices()
	if err != nil {
		return nil, err
	}

	// 更新缓存
	s.cache[cacheKey] = devices
	s.cacheTime = time.Now()

	return devices, nil
}
