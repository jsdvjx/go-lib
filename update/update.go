package update

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type Update struct {
	Source  string
	Target  string
	Service string
	Name    string
}

func (update *Update) NeedUpdate() bool {
	remoteVersion := update.GetRemoteVersion()
	if remoteVersion == "" {
		return false
	}
	return update.CurrentVersion() != remoteVersion
}
func (update *Update) localPath() string {
	return fmt.Sprintf("/usr/bin/%s", update.Name)
}
func (update *Update) SaveVersion(version string) error {
	file, err := os.Create(update.versionLocalPath())
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)
	_, err = file.Write([]byte(version))
	if err != nil {
		return err
	}
	return nil
}

func (update *Update) Do() error {
	if !update.NeedUpdate() {
		return nil
	}
	local := update.binaryLocalPath()
	err := update.download(local)
	if err != nil {
		return err
	}
	err = os.Rename(local, update.localPath())
	if err != nil {
		return err
	}
	shell := fmt.Sprintf(`
service %s stop
cp %s %s
chmod +x %s
service %s start
`, update.Service, local, update.Target, update.Target, update.Service)
	f, _ := os.Create("/tmp/update.sh")
	_, _ = f.Write([]byte(shell))
	_ = f.Chmod(0777)
	cmd := exec.Command("/bin/bash", "-c", "nohup /tmp/update.sh &>/tmp/update.log &")
	return cmd.Run()
}

func (update *Update) url() string {
	return fmt.Sprintf("%s/%s-%s", update.Source, update.Name, env())
}
func (update *Update) versionLocalPath() string {
	return fmt.Sprintf("%s-%s-version", update.Name, env())
}
func (update *Update) binaryLocalPath() string {
	return fmt.Sprintf("%s-%s", update.Name, env())
}
func env() string {
	res := os.Getenv("ENV")
	if res == "" {
		return "dev"
	}
	return res
}

func (update *Update) GetRemoteVersion() string {
	client := &http.Client{}
	req, _ := http.NewRequest("HEAD", update.Source, nil)
	resp, _ := client.Do(req)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	if resp.StatusCode == 200 {
		return strings.Trim(resp.Header.Get("ETag"), "\"")
	}
	return ""
}

func (update *Update) download(local string) error {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", update.Source, nil)
	resp, _ := client.Do(req)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	if resp.StatusCode == 200 {
		file, err := os.Create(local)
		version := strings.Trim(resp.Header.Get("ETag"), "\"")
		if err != nil {
			return err
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {

			}

		}(file)
		_, err = io.Copy(file, resp.Body)
		if err != nil {
			return err
		}
		err = update.SaveVersion(version)
		if err != nil {
			return err
		}
		return nil
	} else {
		logrus.Error("download failed[%d]", resp.StatusCode)
		return fmt.Errorf("download failed[%d]", resp.StatusCode)
	}
}

func (update *Update) CurrentVersion() string {
	file, err := os.Open(update.versionLocalPath())
	if err != nil {
		return ""
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)
	buff := make([]byte, 1024)
	_, err = file.Read(buff)
	if err != nil {
		return ""
	}
	return string(buff)
}
