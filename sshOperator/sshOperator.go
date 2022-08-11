// sshOperator
package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

func createServer(name string) {
	//name string, image *hcloud.Image, sizing float32
	// Location, Image, Type, Volume(has to be created, sizing can be determined there), Networking (IPv4, IPv6, private), firewalls, additional features, ssh-key, name
	// https://pkg.go.dev/github.com/hetznercloud/hcloud-go/hcloud?utm_source=godoc#Server
	// https://docs.hetzner.com/cloud/general/locations/
	client := hcloud.NewClient(hcloud.WithToken(os.Getenv("API_TOKEN")))
	cl, _ := client.ServerType.All(context.Background())
	server, _ := client.Server.All(context.Background())
	fmt.Println(cl[0], server)

	serverOpts := hcloud.ServerCreateOpts{Name: name}
	serverOpts.SSHKeys
	err := serverOpts.Validate()
	if err != nil {
		fmt.Printf("%v", err)
	}
	result, response, err := client.Server.Create(context.Background(), serverOpts)
	if err != nil {
		fmt.Printf("%v", err)
	}

}

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
	switch {
	case sshKeys == nil:
		fmt.Printf("The key %v does not exist.\n", name)
	case err != nil:
		fmt.Printf("%v", err)
	default:
		_, err = client.SSHKey.Delete(context.Background(), sshKeys)
	}
}

func main() {
	// flag things
	delKey := flag.Bool("deletesshkey", false, "delete a key")
	crtKey := flag.Bool("createsshkey", false, "create a key")
	serv := flag.Bool("createServer", false, "create server")
	flag.Parse()
	switch {
	case *delKey && !*crtKey:
		deleteKey(os.Args[2])
	case *crtKey && !*delKey:
		name := os.Args[2]
		PublicKey := os.Args[3]
		createKey(name, PublicKey)
	case *serv:
		createServer()
	default:
		fmt.Printf("Please enter the mode exactly once. You entered delete:%v create:%v\n", *crtKey, *delKey)
	}
}
