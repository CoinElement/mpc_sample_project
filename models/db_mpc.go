package models

// 参与计算的信息 R/W
type Mpc struct {
	ID             uint64 `gorm:"primary_key"`
	InstanceId     string
	SequenceId     int
	PrevAddress    string
	NextAddress    string
	ReceivedData   string
	CalculatedData string
	Coefficient    int64
	Status         string //PENDING | CALCULATING | FINISHED
}

func (db DB) CreateMpcs(mpc Mpc) error {
	return db.DB.Create(&mpc).Error
}

func (db DB) GetMpcInfoByInstanceId(instanceId string) (*Mpc, error) {
	var found Mpc
	err := db.DB.Where(`"instance_id" = ?`, instanceId).First(&found).Error
	if err != nil {
		return nil, err
	}
	return &found, nil
}

func (db DB) SetMpcStatus(instanceId, status string) error {
	mpc := Mpc{}
	return db.DB.Model(&mpc).Where(`"instance_id" = ?`, instanceId).Updates(Mpc{Status: status}).Error
}
