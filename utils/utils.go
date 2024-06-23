package utils

import (
	"math/rand"
	"os"
	"path"
	"time"
)

func Env(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func RunEnv() string {
	return Env("ENV", "dev")
}

func generateRandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	// 为了确保每次程序运行时都能产生不同的随机数序列，
	// 我们需要设置一个变化的种子值。这里使用当前时间作为种子值。
	rand.New(rand.NewSource(time.Now().UnixNano()))

	// 创建一个长度为n的空字符串，用于存放最终生成的随机字符串
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
func GetId() string {
	filePath := path.Join(Env("CACHE_PATH", "./"), "id.data")
	data, err := os.ReadFile(filePath)
	id, _ := os.Hostname()
	if err != nil {
		if os.IsNotExist(err) {
			id = generateRandomString(10)
			_ = os.WriteFile(filePath, []byte(id), 0755)
		} else {

			return id
		}
	} else {
		id = string(data)
	}
	return id
}
