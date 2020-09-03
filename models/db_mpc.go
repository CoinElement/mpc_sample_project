package models

// 参与计算的信息 R/W
type Mpc struct {
	ID             uint
	InstanceId     string `gorm:"primary_key"`
	SequenceId     int    `gorm:"primaryKey;autoIncrement:false"`
	PrevAddress    string
	NextAddress    string
	ReceivedData   string
	CalculatedData string
	Coefficient    int64
	Status         string //PENDING | FINISHED
}

func (db DB) CreateMpcs(clients []Mpc) error {
	return db.DB.Create(clients).Error
}

func (db DB) GetMpcInfoByInstanceId(instanceId string) (*Mpc, error) {
	var found Mpc
	err := db.DB.Where(&Mpc{InstanceId: instanceId}).First(&found).Error
	if err != nil {
		return nil, err
	}
	return &found, nil
}

func (db DB) SetMpcStatus(instanceId, status string) error {
	return db.DB.Where(&Mpc{InstanceId: instanceId}).Updates(&Mpc{Status: status}).Error
}
