package models

import "time"

// 每轮计算 R/W
type Instance struct {
	ID          uint64 `gorm:"primary_key"`
	InstanceId  string
	FirstIp     string
	FinalResult int64
	Function    string
	Status      string // PENDING | CALCULATING | FINISHED | TERMINATED
	StartTime   time.Time
}

func (db *DB) CreateInstances(instance Instance) error {
	return db.DB.Create(&instance).Error
}

func (db *DB) SetFinalResult(instanceId string, finalResult int64) error {
	instance := Instance{}
	return db.DB.Model(&instance).Where(`"instance_id" = ?`, instanceId).Updates(Instance{FinalResult: finalResult, Status: "FINISHED"}).Error
}

func (db *DB) SetInstanceStatus(instanceId, status string) error {
	instance := Instance{}
	return db.DB.Model(&instance).Where(`"instance_id" = ?`, instanceId).Updates(Instance{Status: status}).Error
}

func (db *DB) GetInstanceById(instanceId string) (*Instance, error) {
	var instance Instance
	err := db.DB.Where(`"instance_id" = ?`, instanceId).First(&instance).Error
	return &instance, err
}
