package sender

import (
	"asyncSender/pkg/logger"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type SendError struct {
	Code        int
	Description string
}

type Sender struct {
	ID        int    `json:"id"`
	IsBot     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
	token     string
	client    *http.Client //будем переиспользовать клиент чтобы не создавать новый
}

type SendMessageParams struct {
	Url         string `json:"Url"`
	Data        string `json:"json data"`
	RequestType string `json:"Request type"`
}

func InitSender() (*Sender, error) {
	client := &http.Client{
		Timeout: time.Second * 15,
	}

	var result struct {
		Ok     bool   `json:"ok"`
		Result Sender `json:"result"`
	}

	result.Result.client = client

	return &result.Result, nil
}

func (b *Sender) SendMessage(params SendMessageParams) error {
	if !inArray(params.RequestType, []string{"POST", "GET"}) {
		return &SendError{
			Code:        400,
			Description: "Unsupported protocol scheme",
		}
	}

	req, err := http.NewRequest(params.RequestType, params.Url, bytes.NewBuffer([]byte(params.Data)))
	if err != nil {
		logger.Errorf("Ошибка при создании запроса: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("boobl", "booblgoom")

	resp, err := b.client.Do(req)
	if err != nil {
		logger.Errorf("Ошибка при выполнении запроса: %v", err)
	}

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	logger.Info("respBody: " + string(respBody))

	var result struct {
		Status         string `json:"status"`
		ErrorMessage   string `json:"errorMessage"`
		ErrorId        string `json:"errorId"`
		HttpStatusCode int    `json:"httpStatusCode"`
	}

	err = json.Unmarshal(respBody, &result)

	if err != nil {
		return err
	}

	if result.Status != "Success" {
		return &SendError{
			Code:        resp.StatusCode,
			Description: result.ErrorMessage,
		}
	}

	return nil
}

func (e *SendError) Error() string {
	return fmt.Sprintf("Ошибка %d: %s", e.Code, e.Description)
}

func inArray(val string, array []string) bool {
	for _, item := range array {
		if item == val {
			return true
		}
	}
	return false
}
