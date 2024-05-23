package config

import (
	"fmt"
	"github.com/jsdvjx/go-lib/aes"
	"gopkg.in/yaml.v2"
	"io"
	"net/http"
	"strings"
	"time"
)

func download(url string) ([]byte, string, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	resp, _ := client.Do(req)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	if resp.StatusCode == 200 {
		bs, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, "", err
		}
		if len(bs) < 1 {
			return nil, "", fmt.Errorf("empty response")
		}
		tag := strings.Trim(resp.Header.Get("ETag"), "\"")
		return bs, tag, nil
	}
	return nil, "", fmt.Errorf("download failed[%d]", resp.StatusCode)
}

func Watch[T any](url string, key string, duration time.Duration, onChange func(T)) {
	var result *T = nil
	version := ""
	go func() {
		tick := time.Tick(duration)
		for {
			select {
			case <-tick:
				bs, nVersion, err := download(url)
				bs, _ = aes.DecryptAES(bs, []byte(key))
				if err != nil {
					fmt.Println(err)
					continue
				}
				if version != nVersion {
					version = nVersion
					err = yaml.Unmarshal(bs, &result)
					if err != nil {
						return
					}
					onChange(*result)
				}
			}
		}
	}()
}

func Load[T any](url string, key string) (*T, error) {
	bs, _, err := download(url)
	if err != nil {
		return nil, err
	}
	bs, _ = aes.DecryptAES(bs, []byte(key))
	var result T
	err = yaml.Unmarshal(bs, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
