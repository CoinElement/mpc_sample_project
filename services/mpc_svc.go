package services

import (
	"errors"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"math/rand"
	"mpc_sample_project/models"
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
		return err
	}
	var mpcs []models.Client
	for index, address := range config.IPAddress {
		//向每个 IP 发送 Notification
		mpcs = append(mpcs, models.Client{
			InstanceId: instance_id,
			IpAddress:  address,
			Status:     "INVITED",
		})
		notification := models.FormNotification{
			InstanceId:  instance_id,
			PrevAddress: config.IPAddress[index-1],
			Coefficient: index * 2,
			NextAddress: config.IPAddress[index+1],
			SequenceId:  index,
		}
		ms.log.Debug(notification)
		err := models.PostNotification(address, notification)
		if err != nil {
			return nil
		}
	}
	// 保存到数据库
	if ms.db.CreateClients(mpcs) != nil {
		return nil
	}
	return nil
}

func (ms *MpcService) ReceiveNotification(clientIp string, notification models.FormNotification) error {
	config, err := models.GetConfig() // 获取配置
	if err != nil {
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
		return nil
	}
	return nil
}

func (ms *MpcService) ReceiveCommitment(clientIp string, commitment models.FormCommitment) error {
	needContinue, err := ms.db.SetClientStatusByCommitment(commitment)
	if err != nil {
		return err
	}
	if !needContinue {
		err := ms.db.SetInstanceStatus(commitment.InstanceId, "TERMINATED") // 有一家拒绝直接设置 Terminated
		if err != nil {
			return err
		}
	}
	instance, err := ms.db.GetInstanceById(commitment.InstanceId)
	if err != nil {
		return err
	}
	mpcs, err := ms.db.GetClientsByInstanceId(commitment.InstanceId)
	if err != nil {
		return nil
	}
	// 当所有的 party_i 都发送了 Commitment 之后将 noise 发给第一家 (POST firstParty/result)
	if isReady(instance, mpcs) { // 调用方法判断是否应该开始计算
		// 生成 noise 和 form_result
		if err := ms.db.SetInstanceStatus(instance.InstanceId, "CALCULATING"); err != nil {
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

func isReady(instance *models.Instance, mpcs []models.Client) bool {
	// 超时直接返回 false
	// Instance 不为 Terminated
	if instance.Status == "TERMINATED" {
		return false
	}
	for _, mpc := range mpcs {
		// 有一家拒绝了请求
		if mpc.Status == "REFUSED" {
			return false
		}
	}
	return true
}
