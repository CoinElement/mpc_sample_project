package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"mpc_sample_project/services"
	"net/http"
)

type MpcController struct {
	log *logrus.Entry
	ms  *services.MpcService
}

func NewMpcController(log *logrus.Logger, ms *services.MpcService) *MpcController {
	return &MpcController{
		log: log.WithField("controller", "match"),
		ms:  ms,
	}
}

type FormNotification struct {
	InstanceId  string `json:"instance_id"`
	PrevAddress string `json:"prev_address"`
	Coefficient int    `json:"coefficient"`
	NextAddress string `json:"next_address"`
}

type FormCommitment struct {
	InstanceId string `json:"instance_id"`
	Ready      bool   `json:"ready"`
	SequenceId string `json:"sequence_id"`
	Secret     string `json:"secret"`
}

type FormResult struct {
	InstanceId     string `json:"instance_id"` // 感觉其实没有必要，作为身份验证的辅助依据？
	FromSequenceId string `json:"from_sequence_id"`
	Data           int    `json:"data"` // 上一家的 result 或者自己是第一家时的 noise
}

func (mc *MpcController) HandleNotification(c *gin.Context) {
	// 收到参与计算的请求
	result := FormNotification{}
	mc.log.Debug("From IP:" + c.ClientIP())
	if err := c.ShouldBindJSON(&result); err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{
				"msg":   "failure",
				"error": err.Error(),
			})
	}
}

func (mc *MpcController) HandleCommitment(c *gin.Context) {
	// 作为 K 时收到 FormCommitment
	commitment := FormCommitment{}
	mc.log.Debug("FromIP: " + c.ClientIP())
	if err := c.ShouldBindJSON(&commitment); err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{
				"msg":   "failure",
				"error": err.Error(),
			},
		)
		return
	}
}

func (mc *MpcController) HandleResult(c *gin.Context) {
	// 作为 K 或者 I 时收到 FormResult
	result := FormResult{}
	mc.log.Debug("From IP: " + c.ClientIP())
	if err := c.ShouldBindJSON(&result); err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{
				"msg":   "failure",
				"error": err.Error(),
			},
		)
		return
	}
}
