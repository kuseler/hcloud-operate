// sshOperator
package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

func createKey(name, publicKey string) {
	labels := make(map[string]string)
	client := hcloud.NewClient(hcloud.WithToken(os.Getenv("API_TOKEN")))
	SSHKeyCreateOpts := hcloud.SSHKeyCreateOpts{Name: name, PublicKey: publicKey, Labels: labels}
	SSHKey, _, error := client.SSHKey.Create(context.Background(), SSHKeyCreateOpts) //create SSHKey
	if error != nil {
		fmt.Printf("%v", error) //Print out SSHKey
	}
	fmt.Printf("%v\n", SSHKey)

}

func deleteKey(name string) {
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
	// flag things
	delKey := flag.Bool("deletesshkey", false, "delete a key")
	crtKey := flag.Bool("createsshkey", false, "create a key")
	flag.Parse()
	switch {
	case *delKey && !*crtKey:
		deleteKey(os.Args[2])
	case *crtKey && !*delKey:
		name := os.Args[2]
		PublicKey := os.Args[3]
		createKey(name, PublicKey)
	default:
		fmt.Printf("Please enter the mode exactly once. You entered delete:%v create:%v\n", *crtKey, *delKey)
	}
}
