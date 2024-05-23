package update

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type Update struct {
	source  string
	target  string
	service string
	name    string
}

func New(source, target, service, name string) *Update {
	return &Update{
		source:  source,
		target:  target,
		service: service,
		name:    name,
	}
}

func (update *Update) NeedUpdate() bool {
	remoteVersion := update.GetRemoteVersion()
	if remoteVersion == "" {
		return false
	}
	return update.CurrentVersion() != remoteVersion
}
func (update *Update) localPath() string {
	return fmt.Sprintf("/usr/bin/%s", update.name)
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
	local := update.binaryLocalPath()
	err := update.download(local)
	if err != nil {
		return err
	}
	err = os.Rename(local, update.localPath())
	if err != nil {
		return err
	}
	command := fmt.Sprintf("service %s restart", update.service)
	cmd := exec.Command("/bin/bash", "-c", command)
	return cmd.Run()
}

func (update *Update) url() string {
	return fmt.Sprintf("%s/%s-%s", update.source, update.name, env())
}
func (update *Update) versionLocalPath() string {
	return fmt.Sprintf("%s-%s-version", update.name, env())
}
func (update *Update) binaryLocalPath() string {
	return fmt.Sprintf("%s-%s", update.name, env())
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
	req, _ := http.NewRequest("HEAD", update.source, nil)
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
	req, _ := http.NewRequest("GET", update.source, nil)
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
