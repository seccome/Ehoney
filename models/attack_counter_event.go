package models

type CounterEvent struct {
	ID    int64  `gorm:"primary_key;AUTO_INCREMENT;not null;unique;column:id" json:"id"`
	IP    string `gorm:"" form:"ip" json:"ip" gorm:"size:256" binding:"required"`
	Type  string `gorm:"" form:"type" json:"type" gorm:"size:256" binding:"required"`
	Token string `gorm:"" form:"token" json:"token" gorm:"size:256" binding:"required"`
	Info  string `gorm:"" form:"info" json:"info" gorm:"" binding:"required"`
}

func (event *CounterEvent) CreateCountEvent() error {
	if err := db.Create(event).Error; err != nil {
		return err
	}
	return nil
}

func (event *CounterEvent) GetCounterEvent(info, protocolType, ip, token string) (*CounterEvent, error) {
	var ret CounterEvent
	if err := db.Where("info = ? AND protocolType = ? AND ip = ? AND token = ?", info, protocolType, ip, token).Find(event).Scan(&ret).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}

func (event *CounterEvent) GetCounterEvents() (*[]CounterEvent, error) {
	var ret []CounterEvent
	if err := db.Find(&ret).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}

func (event *CounterEvent) GetCounterEventsByAttackIp(attackIp string) (*[]CounterEvent, error) {
	var ret []CounterEvent
	if err := db.Where("IP = ?", attackIp).Find(event).Scan(&ret).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}
