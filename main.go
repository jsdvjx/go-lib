package main

import "github.com/jsdvjx/go-lib/update"

func main() {
	update := &update.Update{
		Source:  "https://application-repo.oss-cn-guangzhou.aliyuncs.com/processors/",
		Target:  "/usr/bin/Processor",
		Service: "service",
		Name:    "Processor",
	}
	update.Watch()
	select {}
}
