package sender

import (
	"asyncSender/pkg/logger"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

type SendError struct {
	Code        int
	Description string
}

type Bot struct {
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

const telegramAPI = "https://api.telegram.org/bot"

func InitSender() (*Bot, error) {
	client := &http.Client{
		Timeout: time.Second * 15,
	}

	var result struct {
		Ok     bool `json:"ok"`
		Result Bot  `json:"result"`
	}

	result.Result.client = client

	return &result.Result, nil
}

func (b *Bot) SendMessage(params SendMessageParams) error {
	var resp *http.Response
	var err error

	switch params.RequestType {
	case "POST":
		resp, err = b.client.Post(params.Url, "application/json", bytes.NewBuffer([]byte(params.Data)))
		fmt.Println("POST")
	default:
		resp, err = b.client.Get(params.Url)
		fmt.Println("GET")
	}

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	logger.Info(respBody)

	var result struct {
		Ok          bool   `json:"ok"`
		Description string `json:"description"`
	}

	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return err
	}

	if !result.Ok {
		return &SendError{
			Code:        resp.StatusCode,
			Description: result.Description,
		}
	}

	return nil
}

func ParseRetryAfter(err *SendError) (int, error) {
	if err.Code == 429 {
		matches := regexp.MustCompile(`retry after (\d+)`).FindStringSubmatch(err.Description)
		retryAfter, err := strconv.Atoi(matches[1])
		if err != nil {
			return 0, fmt.Errorf("failed to convert retry after to int: %w", err)
		}
		return retryAfter, nil
	}

	return 0, nil
}

func Pluralize(n int, singular, plural1, plural2 string) string {
	n = n % 100
	if n > 10 && n < 20 {
		return plural2
	}

	n = n % 10
	if n == 1 {
		return singular
	}

	if n > 1 && n < 5 {
		return plural1
	}

	return plural2
}

func (e *SendError) Error() string {
	return fmt.Sprintf("Ошибка %d: %s", e.Code, e.Description)
}
