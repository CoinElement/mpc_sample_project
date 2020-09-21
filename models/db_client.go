package models

// 发起请求的 IP 列表 R/W
type Client struct {
	//ID         uint   `gorm:"AUTO_INCREMENT"`
	InstanceId string `gorm:"primaryKey"`
	SequenceId int    `gorm:"primaryKey"`
	IpAddress  string // 暂不考虑 ip unreachable
	Status     string // INVITED | ACCEPTED| REFUSED
}

func (db DB) CreateClients(client Client) error {
	return db.DB.Create(client).Error
}

func (db DB) GetUncommittedClients(instanceId string) ([]Client, error) {
	clients := make([]Client, 0)
	err := db.DB.Where("committed = ? AND instance_id = ?", false, instanceId).Find(&clients).Error
	if err != nil {
		return nil, err
	}
	return clients, nil
}

func (db DB) GetClientsByInstanceId(instanceId string) ([]Client, error) {
	clients := make([]Client, 0)
	err := db.DB.Where("instance_id = ?", instanceId).Find(&clients).Error
	if err != nil {
		return nil, err
	}
	return clients, nil
}

func (db DB) SetClientStatusByCommitment(commitment FormCommitment) (bool, error) {
	client := Client{}

	if err := db.DB.Where(`"instance_id" = ? AND "sequence_id" = ?`, commitment.InstanceId, commitment.SequenceId).First(&client).Error; err != nil {
		return false, err
	}

	var str string
	need_continue := true
	if commitment.Ready {
		str = "ACCEPTED"
	} else {
		str = "REFUSED"
		need_continue = false
	}
	client.Status = str
	return need_continue, db.DB.Model(&client).Updates(client).Error
}
