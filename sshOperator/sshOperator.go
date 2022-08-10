// sshOperator
package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

func defineParameters() (*hcloud.Client, hcloud.SSHKeyCreateOpts) {
	Name := os.Args[2]
	PublicKey := os.Args[3] //perhaps creating an ssh publickey in go?
	Labels := make(map[string]string)
	client := hcloud.NewClient(hcloud.WithToken(os.Getenv("API_TOKEN")))
	SSHKeyCreateOpts := hcloud.SSHKeyCreateOpts{Name: Name, PublicKey: PublicKey, Labels: Labels}

	return client, SSHKeyCreateOpts
}

func createKey() {
	client, SSHKeyCreateOpts := defineParameters()                                   //set parameters for SSHKey
	SSHKey, _, error := client.SSHKey.Create(context.Background(), SSHKeyCreateOpts) //create SSHKey
	if error != nil {
		fmt.Printf("%v", error) //Print out SSHKey
	}
	fmt.Printf("%v\n", SSHKey)

}

func deleteKey() {
	name := os.Args[2]
	fmt.Printf("Api Key: '%s', name: '%s'\n", os.Getenv("API_TOKEN"), name)
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

func main() {
	delKey := flag.Bool("d", false, "delete a key")
	crtKey := flag.Bool("c", false, "create a key")
	flag.Parse()
	switch {
	case *delKey && !*crtKey:
		deleteKey()
	case *crtKey && !*delKey:
		createKey()
	default:
		fmt.Println("Please enter the mode exactly once. You entered delete:%v create:%v", crtKey, delKey)
	}
}
