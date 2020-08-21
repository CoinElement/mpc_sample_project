package models

type Mpc struct {
	InstanceId string `gorm:"primary_key"`
	IpAddress  string `gorm:"primary_key"`
	Committed  bool
}

func (db DB) CreateMpcs(mpcs []Mpc) error {
	return db.DB.Create(mpcs).Error
}

func (db DB) GetUncommittedIps(instanceId string) ([]Mpc, error) {
	mpcs := make([]Mpc, 0)
	err := db.DB.Where("committed = ? AND instance_id = ?", false, instanceId).Find(&mpcs).Error
	if err != nil {
		return nil, err
	}
	return mpcs, nil
}

func (db DB) GetInfoByInstanceId(instanceId string) ([]Mpc, error) {
	mpcs := make([]Mpc, 0)
	err := db.DB.Where("instance_id = ?", instanceId).Find(&mpcs).Error
	if err != nil {
		return nil, err
	}
	return mpcs, nil
}

func (db DB) SetCommitment(instanceId, ipAddress string, committed bool) error {
	return db.DB.Model(Mpc{}).Where("instance_id = ? AND ip_address = ?", instanceId, ipAddress).Update("committed", committed).Error
}
