package models

type Setting struct {
	ID            int64           `gorm:"primary_key;AUTO_INCREMENT;unique;column:id" json:"id"`                               //协议ID
	ConfigName    string          `json:"ConfigName" form:"ConfigName" gorm:"unique;not null;size:128" binding:"required"`     //配置名称类型
	ConfigValue   string          `json:"ConfigValue" form:"ConfigValue" gorm:"not null" binding:"required"`                   //配置值
}


func (setting *Setting) CreateSetting() error {
	if err := db.Create(setting).Error; err != nil {
		return err
	}
	return nil
}
