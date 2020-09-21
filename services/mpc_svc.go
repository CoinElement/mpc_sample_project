package services

import (
	"errors"
	"math/rand"
	"mpc_sample_project/models"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type MpcService struct {
	log *logrus.Entry
	db  *models.DB
}

func NewMpcService(log *logrus.Logger, db *models.DB) *MpcService {
	return &MpcService{
		log: log.WithField("services", "Client"),
		db:  db,
	}
}

func (ms *MpcService) Start() error {
	// 开始
	instance_id := uuid.New().String()
	config, err := models.GetConfig()
	if err != nil {
		ms.log.Error("failed to get config when start")
		return err
	}

	for index, address := range config.IPAddress {

		if index == 0 {
			//向instance数据库插入一条新数据
			instance := models.Instance{
				InstanceId: instance_id,
				FirstIp:    address,
				Function:   "prevResult+coefficient*data",
				Status:     "PENDING",
				StartTime:  time.Now(),
			}
			if err := ms.db.CreateInstances(instance); err != nil {
				ms.log.Error("failed to insert new information into instance")
				return err
			}
		}
		//向每个 IP 发送 Notification
		client := models.Client{
			InstanceId: instance_id,
			SequenceId: index,
			IpAddress:  address,
			Status:     "INVITED",
		}
		// 插入到数据库client
		if err := ms.db.CreateClients(client); err != nil {
			ms.log.Error("failed to insert data into Client")
			return err
		}

		notification := models.FormNotification{}
		if index == 0 { //第一个
			notification.InstanceId = instance_id
			notification.PrevAddress = config.SelfIp
			notification.Coefficient = 10
			notification.NextAddress = config.IPAddress[index+1]
			notification.SequenceId = index
		} else if index == len(config.IPAddress)-1 { //最后一个
			notification.InstanceId = instance_id
			notification.PrevAddress = config.IPAddress[index-1]
			notification.Coefficient = 10
			notification.NextAddress = config.SelfIp
			notification.SequenceId = index
		} else {
			notification.InstanceId = instance_id
			notification.PrevAddress = config.IPAddress[index-1]
			notification.Coefficient = 10
			notification.NextAddress = config.IPAddress[index+1]
			notification.SequenceId = index
		}

		ms.log.Debug(notification)
		err := models.PostNotification(address, notification)
		if err != nil {
			ms.log.Error("failed to post notification")
			return err
		}
	}

	return nil
}

func (ms *MpcService) ReceiveNotification(clientIp string, notification models.FormNotification) error {
	config, err := models.GetConfig() // 获取配置
	if err != nil {
		ms.log.Error("failed to get config when receive notification")
		return err
	}
	// 收到参与计算请求
	// (POST readyOrNot to ip/commit)
	// 暂时全部接受
	commitment := models.FormCommitment{
		InstanceId: notification.InstanceId,
		Ready:      ShouldAcceptCommitment(),
		SequenceId: notification.SequenceId,
		Secret:     config.Secret,
	}
	if err := models.PostCommitment(clientIp, commitment); err != nil {
		ms.log.Error("failed to post commitment")
		return err
	}

	mpc := models.Mpc{
		InstanceId:  notification.InstanceId,
		SequenceId:  notification.SequenceId,
		PrevAddress: notification.PrevAddress,
		NextAddress: notification.NextAddress,
		Coefficient: notification.Coefficient,
		Status:      "PENDING",
	}

	if err := ms.db.CreateMpcs(mpc); err != nil {
		ms.log.Error("failed to insert new information into mpc")
		return err
	}
	return nil
}

func (ms *MpcService) ReceiveCommitment(clientIp string, commitment models.FormCommitment) error {
	needContinue, err := ms.db.SetClientStatusByCommitment(commitment)
	if err != nil {
		ms.log.Error("failed to set client status by commitment")
		return err
	}
	if !needContinue {
		err := ms.db.SetInstanceStatus(commitment.InstanceId, "TERMINATED") // 有一家拒绝直接设置 Terminated
		if err != nil {
			ms.log.Error("failed to set instance status of TERMINATED")
			return err
		}
		return nil
	}
	instance, err := ms.db.GetInstanceById(commitment.InstanceId)
	if err != nil {
		ms.log.Error("failed to get instance by id")
		return err
	}
	clients, err := ms.db.GetClientsByInstanceId(commitment.InstanceId)
	if err != nil {
		ms.log.Error("failed to get clients by instance id")
		return nil
	}
	// 当所有的 party_i 都发送了 Commitment 之后将 noise 发给第一家 (POST firstParty/result)
	if isReady(instance, clients) { // 调用方法判断是否应该开始计算
		// 生成 noise 和 form_result
		if err := ms.db.SetInstanceStatus(instance.InstanceId, "CALCULATING"); err != nil {
			ms.log.Error("failed to set instance status of CALCULATING")
			return err
		}
		noise := rand.Int63()
		result := models.FormResult{
			InstanceId:     commitment.InstanceId,
			FromSequenceId: 0,     //本机的SequenceId
			Data:           noise, //本机计算出来的结果
		}
		err := models.PostResult(instance.FirstIp, result) // Post 给第一家
		if err != nil {
			ms.log.Error("failed to post noise to first party")
			return err
		}
	}
	return nil
}

func (ms *MpcService) ReceiveResult(clientIp string, result models.FormResult) error {
	// 收到 noise 或者 prevResult
	instance, err_instance := ms.db.GetInstanceById(result.InstanceId)
	if err_instance == nil {
		//自己是发起者
		if err := ms.db.SetFinalResult(instance.InstanceId, result.Data); err != nil {
			return err
		} else {
			return nil
		}
	} else if !errors.Is(err_instance, gorm.ErrRecordNotFound) {
		return err_instance
	}
	mpc, err_mpc := ms.db.GetMpcInfoByInstanceId(result.InstanceId)
	if err_mpc == nil {
		// 如果自己是 参与者， 根据系数和prevResult， 计算出this_result, 发给下一家 (POST nextAddress/result)
		config, err := models.GetConfig()
		if err != nil {
			ms.log.Error("failed to get config when receive result")
			return err
		}
		if mpc.PrevAddress != clientIp {
			ms.log.Error("From IP:" + clientIp + "PrevIp:" + mpc.PrevAddress)
		} else {
			ms.log.Debug("IP Confirmed")
		}
		nextData := (mpc.Coefficient*config.Data + result.Data) % (2 ^ 64) //计算 Data
		nextResult := models.FormResult{
			InstanceId:     result.InstanceId,
			FromSequenceId: mpc.SequenceId,
			Data:           nextData,
		}
		if err := models.PostResult(mpc.NextAddress, nextResult); err != nil {
			return err
		}
	} else if !errors.Is(err_mpc, gorm.ErrRecordNotFound) {
		return err_mpc
	}

	return nil
}

func ShouldAcceptCommitment() bool {
	return true
}

func isReady(instance *models.Instance, clients []models.Client) bool {
	// 超时直接返回 false
	// Instance 不为 Terminated
	if instance.Status == "TERMINATED" {
		return false
	}
	for _, client := range clients {
		// 有一家拒绝了请求
		if client.Status == "REFUSED" {
			return false
		}
	}
	return true
}
