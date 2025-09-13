package models

import (
	"time"

	"gorm.io/gorm"
)

// Device 设备模型
type Device struct {
	ID                uint           `json:"id" gorm:"primaryKey"`
	EntityID          string         `json:"entity_id" gorm:"type:varchar(255);uniqueIndex;not null"`
	FriendlyName      string         `json:"friendly_name" gorm:"type:varchar(255)"`
	State             string         `json:"state" gorm:"type:varchar(100)"`
	Domain            string         `json:"domain" gorm:"type:varchar(50);index"`
	DeviceClass       string         `json:"device_class,omitempty" gorm:"type:varchar(100)"`
	UnitOfMeasurement string         `json:"unit_of_measurement,omitempty" gorm:"type:varchar(50)"`
	Attributes        string         `json:"attributes" gorm:"type:text"` // JSON字符串
	LastChanged       time.Time      `json:"last_changed"`
	LastUpdated       time.Time      `json:"last_updated"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `json:"-" gorm:"index"`
}

// DeviceHistory 设备历史记录
type DeviceHistory struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	EntityID  string    `json:"entity_id" gorm:"type:varchar(255);index;not null"`
	State     string    `json:"state" gorm:"type:varchar(100)"`
	Timestamp time.Time `json:"timestamp" gorm:"index"`
	CreatedAt time.Time `json:"created_at"`
}

// HADevice Home Assistant设备响应结构
type HADevice struct {
	EntityID    string                 `json:"entity_id"`
	State       string                 `json:"state"`
	Attributes  map[string]interface{} `json:"attributes"`
	LastChanged string                 `json:"last_changed"`
	LastUpdated string                 `json:"last_updated"`
	Context     map[string]interface{} `json:"context"`
}

// HAServiceCall Home Assistant服务调用结构
type HAServiceCall struct {
	Domain      string                 `json:"domain"`
	Service     string                 `json:"service"`
	Target      map[string]interface{} `json:"target,omitempty"`
	ServiceData map[string]interface{} `json:"service_data,omitempty"`
}
