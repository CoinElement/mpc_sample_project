package services

import (
	"mpc_sample_project/models"

	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type MpcService struct {
	log *logrus.Entry
	db  *models.DB
}

func NewMpcService(log *logrus.Logger, db *models.DB) *MpcService {
	return &MpcService{
		log: log.WithField("services", "mpc"),
		db:  db,
	}
}

func (ms *MpcService) Start() error {
	// 开始
	uuid4 := uuid.NewV4().String()    //获取instanceId的唯一值
	config, err := models.GetConfig() //获取ipAddresses
	if err != nil {
		ms.log.Error("failed to get config")
		return err
	}

	var mpcs []models.Mpc
	for _, ip := range config.IPAddress {
		mpcs = append(mpcs, models.Mpc{InstanceId: uuid4, IpAddress: ip, Status: "Pending"})
	}

	if err := ms.db.CreateMpcs(mpcs); err != nil {
		ms.log.Error("failed to create Mpcs")
		return err
	}

	return nil
}

func ReceiveNotification(instanceId, ip string, coefficient int, prevAddress, nextAddress string) error {
	// 收到参与计算请求
	// (POST readyOrNot to ip/commit)
	return nil
}

func ReceiveCommitment(instanceId, ip string, DB *models.DB) error {
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

func ReceiveResult(instanceId, ip string, prevResult int, DB *models.DB) error {
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
