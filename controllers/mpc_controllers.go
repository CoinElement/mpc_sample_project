package controllers

import (
	"mpc_sample_project/models"
	"mpc_sample_project/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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

func (mc *MpcController) HandlePing(c *gin.Context) {
	mc.log.Info("handling ping")
	c.JSON(
		http.StatusOK,
		gin.H{
			"msg": "pong",
		},
	)
}

func (mc *MpcController) HandleStart(c *gin.Context) {
	mc.log.Info("starting")
	err := mc.ms.Start()
	if err != nil {
		mc.log.Error("failed to handle start")
		c.JSON(http.StatusInternalServerError,
			gin.H{
				"msg":   "error",
				"error": err.Error(),
			},
		)
	}
	c.JSON(http.StatusOK,
		gin.H{
			"msg": "OK",
		},
	)
}

func (mc *MpcController) HandleNotification(c *gin.Context) {
	// 收到参与计算的请求
	notification := models.FormNotification{}
	mc.log.Debug("From IP:" + c.ClientIP())
	if err := c.ShouldBindJSON(&notification); err != nil {
		mc.log.Error("failed to bind json of notification")
		c.JSON(http.StatusBadRequest,
			gin.H{
				"msg":   "failure",
				"error": err.Error(),
			},
		)
		return
	}

	err := mc.ms.ReceiveNotification(
		c.ClientIP(),
		notification,
	)
	if err != nil {
		mc.log.Error("failed to receive notification")
		c.JSON(http.StatusInternalServerError,
			gin.H{
				"msg":   "error",
				"error": err.Error(),
			},
		)
		return
	}
	c.JSON(http.StatusOK,
		gin.H{
			"msg": "OK",
		},
	)
}

func (mc *MpcController) HandleCommitment(c *gin.Context) {
	// 作为 K 时收到 FormCommitment
	commitment := models.FormCommitment{}
	mc.log.Debug("FromIP: " + c.ClientIP())
	if err := c.ShouldBindJSON(&commitment); err != nil {
		mc.log.Error("failed to bind json of commitment")
		c.JSON(http.StatusBadRequest,
			gin.H{
				"msg":   "failure",
				"error": err.Error(),
			},
		)
		return
	}

	err := mc.ms.ReceiveCommitment(c.ClientIP(), commitment)
	if err != nil {
		mc.log.Error("failed to receive commitment")
		c.JSON(http.StatusInternalServerError,
			gin.H{
				"msg":   "error",
				"error": err.Error(),
			},
		)
		return
	}
	c.JSON(http.StatusOK,
		gin.H{
			"msg": "OK",
		},
	)
}

func (mc *MpcController) HandleResult(c *gin.Context) {
	// 作为 K 或者 I 时收到 FormResult
	result := models.FormResult{}
	mc.log.Debug("From IP: " + c.ClientIP())
	if err := c.ShouldBindJSON(&result); err != nil {
		mc.log.Error("failed to bind json of result")
		c.JSON(http.StatusBadRequest,
			gin.H{
				"msg":   "failure",
				"error": err.Error(),
			},
		)
		return
	}

	err := mc.ms.ReceiveResult(c.ClientIP(), result)
	if err != nil {
		mc.log.Error("failed to receive result")
		c.JSON(http.StatusInternalServerError,
			gin.H{
				"msg":   "error",
				"error": err.Error(),
			},
		)
		return
	}
	c.JSON(http.StatusOK,
		gin.H{
			"msg": "OK",
		},
	)
}
