package models

type Client struct {
	ID             uint   `json:"id" gorm:"primary_key"`
	InstanceId     string `json:"instanceId"`
	PrevAddress    string `json:"prevAddress"`
	NextAddress    string `json:"nextAddress"`
	ReceivedData   int    `json:"receivedData"`
	CalculatedData int    `json:"calculatedData"`
	Coefficient    int    `json:"coefficient"`
}

func (db DB) CreateClient(clients []Client) error {
	return db.DB.Create(clients).Error
}
