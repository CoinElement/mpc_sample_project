package services

import (
	"github.com/sirupsen/logrus"
)

type MpcService struct {
	log *logrus.Entry
}

func NewMpcService(log *logrus.Logger) *MpcService {
	return &MpcService{
		log: log.WithField("services", "Mpc"),
	}
}

func Start() error {
	// 开始
	return nil
}

func ReceiveNotification(instanceId, ip string, coefficient int, prevAddress, nextAddress string) error {
	// 收到参与计算请求
	// (POST readyOrNot to ip/commit)
	return nil
}

func ReceiveCommitment(instanceId, ip string) error {
	// 作为 K 时收到 Commitment
	// 当所有的 party_i 都发送了 Commitment 之后将 noise 发给第一家 (POST firstParty/result)
	return nil
}

func ReceiveResult(instanceId, ip string, prevResult int) error {
	// 收到 noise 或者 prevResult
	// 根据 instanceId 判断自己的身份
	// 如果自己是k，结束
	// 对比 ip 和 prev_address，一致则继续，不一致则退出
	// 如果自己是 i， 根据系数和prevResult， 计算出this_result, 发给下一家 (POST nextAddress/result)
	return nil
}
