package main

import (
	"context"
	"fmt"
	"os"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

func main() {
	client, SSHKeyCreateOpts := defineParameters()                                   //set parameters for SSHKey
	SSHKey, _, error := client.SSHKey.Create(context.Background(), SSHKeyCreateOpts) //create SSHKey
	if error != nil {
		fmt.Printf("%v", error) //Print out SSHKey
	}
	fmt.Printf("%v", SSHKey)

}

func defineParameters() (*hcloud.Client, hcloud.SSHKeyCreateOpts) {
	Name := os.Args[1]
	PublicKey := os.Args[2] //perhaps creating an ssh publickey in go?
	Labels := make(map[string]string)
	client := hcloud.NewClient(hcloud.WithToken(os.Getenv("API_TOKEN")))
	SSHKeyCreateOpts := hcloud.SSHKeyCreateOpts{Name: Name, PublicKey: PublicKey, Labels: Labels}

	return client, SSHKeyCreateOpts
}
