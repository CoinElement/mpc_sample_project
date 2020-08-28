package services

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"mpc_sample_project/models"
)

type MpcService struct {
	log *logrus.Entry
	db  *models.DB
}

func NewMpcService(log *logrus.Logger) *MpcService {
	return &MpcService{
		log: log.WithField("services", "Mpc"),
	}
}

func (ms *MpcService) Start() error {
	// 开始
	instance_id := uuid.New().String()
	config, err := models.GetConfig()
	if err != nil {
		return err
	}
	var mpcs []models.Mpc
	for index, address := range config.IPAddress {
		mpcs = append(mpcs, models.Mpc{
			InstanceId: instance_id,
			IpAddress:  address,
			Status:     "INVITED",
		})
		notification := models.FormNotification{
			InstanceId:  instance_id,
			PrevAddress: config.IPAddress[index-1],
			Coefficient: index*2,
			NextAddress: config.IPAddress[index+1],
			SequenceId: index,
		}
		err := models.PostNotification(address, notification)
		if err != nil {
			return nil
		}
	}
	if ms.db.CreateMpcs(mpcs) != nil {
		return nil
	}
	return nil
}

func (ms *MpcService) ReceiveNotification(instanceId, ip string, coefficient int, prevAddress, nextAddress string) error {
	// 收到参与计算请求
	// (POST readyOrNot to ip/commit)
	// 暂时全部接受
	form := models.FormCommitment{
		InstanceId: instanceId,
		Ready: ShouldAcceptCommitment(),
		SequenceId:
	}
	if ShouldAcceptCommitment() {
		err := models.PostCommitment(ip)
	}
	return nil
}

func (ms *MpcService) ReceiveCommitment(instanceId, ip string, DB *models.DB) error {
	// 作为 K 时收到 Commitment
	err := DB.SetCommitment(instanceId, ip, true)
	if err != nil {
		return err
	}
	CheckCommitments, err := DB.GetInfoByInstanceId(instanceId)
	if err != nil {
		return err
	}
	for _, commitment := range CheckCommitments {
		if commitment.Committed == false {
			return nil //存在未验证
		}
	} //至此所有的 party_i 都发送了 Commitment
	// 当所有的 party_i 都发送了 Commitment 之后将 noise 发给第一家 (POST firstParty/result)
	return nil
}

func (ms *MpcService) ReceiveResult(instanceId, ip string, prevResult int, DB *models.DB) error {
	// 收到 noise 或者 prevResult
	SelectResults, err := DB.GetInfoByInstanceId(instanceId) // 根据 instanceId 判断自己的身份
	if err != nil {
		return err
	}
	if len(SelectResults) != 0 { // 如果自己是k，结束
		return nil
	}
	if ip != prev_address { // 对比 ip 和 prev_address，一致则继续，不一致则退出
		return nil
	}
	result := prevResult
	// 如果自己是 i， 根据系数和prevResult， 计算出this_result, 发给下一家 (POST nextAddress/result)
	return nil
}

func ShouldAcceptCommitment() bool {
	return true
}
