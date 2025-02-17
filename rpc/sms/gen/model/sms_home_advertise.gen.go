// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameSmsHomeAdvertise = "sms_home_advertise"

// SmsHomeAdvertise 首页轮播广告表
type SmsHomeAdvertise struct {
	ID         int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name       string    `gorm:"column:name;not null;comment:名称" json:"name"`                        // 名称
	Type       int32     `gorm:"column:type;not null;comment:轮播位置：0->PC首页轮播；1->app首页轮播" json:"type"` // 轮播位置：0->PC首页轮播；1->app首页轮播
	Pic        string    `gorm:"column:pic;not null;comment:图片地址" json:"pic"`                        // 图片地址
	StartTime  time.Time `gorm:"column:start_time;not null;comment:开始时间" json:"start_time"`          // 开始时间
	EndTime    time.Time `gorm:"column:end_time;not null;comment:结束时间" json:"end_time"`              // 结束时间
	Status     int32     `gorm:"column:status;not null;comment:上下线状态：0->下线；1->上线" json:"status"`     // 上下线状态：0->下线；1->上线
	ClickCount int32     `gorm:"column:click_count;not null;comment:点击数" json:"click_count"`         // 点击数
	OrderCount int32     `gorm:"column:order_count;not null;comment:下单数" json:"order_count"`         // 下单数
	URL        string    `gorm:"column:url;not null;comment:链接地址" json:"url"`                        // 链接地址
	Note       string    `gorm:"column:note;not null;comment:备注" json:"note"`                        // 备注
	Sort       int32     `gorm:"column:sort;not null;comment:排序" json:"sort"`                        // 排序
}

// TableName SmsHomeAdvertise's table name
func (*SmsHomeAdvertise) TableName() string {
	return TableNameSmsHomeAdvertise
}
