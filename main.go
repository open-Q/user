package main

import (
	common "github.com/open-Q/common/golang"
)

func main() {
	service, _, err := common.NewService("./.contract/contract.json")
	if err != nil {
		panic(err)
	}
	service.Run()
}
