package task

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Task struct {
	Url   string
	Token string
}
type ApiResponse[T any] struct {
	Data    Entity[T] `json:"data,omitempty"` // 使用指针来模拟可空性，omitempty表示如果字段为nil，则不会序列化到JSON中
	Success bool      `json:"success"`
	Message string    `json:"message"`
}
type ApiBaseResponse[T any] struct {
	Data    T      `json:"data,omitempty"` // 使用指针来模拟可空性，omitempty表示如果字段为nil，则不会序列化到JSON中
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// Entity 定义了一个泛型的任务实体结构
type Entity[T any] struct {
	ID           int64   `json:"id"`
	Type         string  `json:"type"`
	Data         string  `json:"data"`
	Consumer     *string `json:"consumer,omitempty"` // 可空字段，omitempty表示如果字段为空，则不会序列化到JSON中
	Timeout      int     `json:"timeout,omitempty"`  // Go中无法直接指定默认值，需要在逻辑中处理
	MaxRetry     int     `json:"max_retry,omitempty"`
	Result       *string `json:"result,omitempty"`
	ResultStatus *string `json:"result_status,omitempty"`
	Sort         int     `json:"sort,omitempty"`
	CreatedAt    *string `json:"created_at"`
	StartedAt    *string `json:"started_at,omitempty"`
	FinishedAt   *string `json:"finished_at,omitempty"`
	Param        T       `json:"param,omitempty"`
}
type CountItem struct {
	Type  string `json:"type"`
	Count int    `json:"count"`
}

func Fetch[T any](t *Task, tType string, consumer string) (*ApiResponse[T], error) {
	client := &http.Client{}
	u, err := url.Parse(t.Url + "task/pull")
	if err != nil {
		return nil, err
	}
	query := u.Query()
	query.Set("type", tType)
	query.Set("consumer", consumer)
	u.RawQuery = query.Encode()
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+t.Token)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 300 || resp.StatusCode < 200 {
		return nil, fmt.Errorf("http status code: %d,%s", resp.StatusCode, string(bs))
	}
	if len(bs) < 1 {
		return nil, fmt.Errorf("empty response")
	}
	var apiResp ApiResponse[T]
	err = json.Unmarshal(bs, &apiResp)
	if err != nil {
		return nil, err
	}
	return &apiResp, nil
}
func Counts(t *Task) ([]CountItem, error) {
	client := &http.Client{}
	u, err := url.Parse(t.Url + "task/count")
	if err != nil {
		return nil, err
	}
	query := u.Query()
	u.RawQuery = query.Encode()
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+t.Token)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 300 || resp.StatusCode < 200 {
		return nil, fmt.Errorf("http status code: %d,%s", resp.StatusCode, string(bs))
	}
	if len(bs) < 1 {
		return nil, fmt.Errorf("empty response")
	}
	var apiResp ApiBaseResponse[[]CountItem]
	err = json.Unmarshal(bs, &apiResp)
	if err != nil {
		return nil, err
	}
	return apiResp.Data, nil
}

func (t *Task) Complete(url string, status string, id int64) error {
	client := &http.Client{}

	result := make(map[string]any)
	result["id"] = id
	result["status"] = status
	result["result"] = url
	data, err := json.Marshal(result)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", t.Url+"task/complete", bytes.NewBuffer(data))
	req.Header.Set("Authorization", "Bearer "+t.Token)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("http status code: %d/%s", res.StatusCode, data)
	}
	return nil
}
