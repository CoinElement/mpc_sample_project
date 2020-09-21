package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type FormNotification struct {
	InstanceId  string `json:"instance_id"`
	PrevAddress string `json:"prev_address"`
	SequenceId  int    `json:"sequence_id"`
	Coefficient int64  `json:"coefficient"`
	NextAddress string `json:"next_address"`
}

type FormCommitment struct {
	InstanceId string `json:"instance_id"`
	Ready      bool   `json:"ready"`
	SequenceId int    `json:"sequence_id"`
	Secret     string `json:"secret"`
}

type FormResult struct {
	InstanceId     string `json:"instance_id"` // 感觉其实没有必要，作为身份验证的辅助依据？
	FromSequenceId int    `json:"from_sequence_id"`
	Data           int64  `json:"data"` // 上一家的 result 或者自己是第一家时的 noise
}

func PostNotification(ip string, notification FormNotification) error {
	bytesData, _ := json.Marshal(notification)
	resp, err := http.Post("http://"+ip+":8080/notification", "application/json;charset=utf-8", bytes.NewBuffer([]byte(bytesData)))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	fmt.Println(content)
	if err != nil {
		return err
	}
	return nil
}

func PostCommitment(ip string, commitment FormCommitment) error {
	bytesData, _ := json.Marshal(commitment)
	resp, err := http.Post("http://"+ip+":8080/commit", "application/json;charset=utf-8", bytes.NewBuffer([]byte(bytesData)))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	fmt.Println(content)
	if err != nil {
		return err
	}
	return nil
}

func PostResult(ip string, result FormResult) error {
	bytesData, _ := json.Marshal(result)
	resp, err := http.Post("http://"+ip+":8080/result", "application/json;charset=utf-8", bytes.NewBuffer([]byte(bytesData)))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	fmt.Println(content)
	if err != nil {
		return err
	}
	return nil
}
