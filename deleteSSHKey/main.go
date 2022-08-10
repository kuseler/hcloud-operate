package main

import (
	"context"
	"fmt"
	"os"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

func main() {
	name := os.Args[1]
	fmt.Printf("Api Key: '%s', name: '%s'\n", os.Getenv("API_TOKEN"), os.Args[1])
	client := hcloud.NewClient(hcloud.WithToken(os.Getenv("API_TOKEN")))
	sshKeys, _, err := client.SSHKey.GetByName(context.Background(), name)
	if sshKeys == nil {
		fmt.Printf("The key %v does not exist.\n", name)
	} else if err != nil {
		fmt.Printf("%v", err)
	} else {
		_, err = client.SSHKey.Delete(context.Background(), sshKeys)
	}
}
