package models

type Client struct {
	InstanceId     string `gorm:"primary_key"`
	PrevAddress    string
	NextAddress    string
	ReceivedData   string
	CalculatedData string
	Coefficient    float32
}

func (db DB) CreateClients(clients []Client) error {
	return db.DB.Create(clients).Error
}

func (db DB) GetClientInfoByInstanceId(instanceId string) (*Client, error) {
	var found Client
	err := db.DB.Where(&Client{InstanceId: instanceId}).First(&found).Error
	if err != nil {
		return nil, err
	}
	return &found, nil
}
