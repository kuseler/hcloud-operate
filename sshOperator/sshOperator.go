// sshOperator
package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

func createServer(name string, locationIDOrName string, serverTypeName string, imageNameOrID string) {
	// Location, Image, Type, Volume(has to be created, sizing can be determined there), Networking (IPv4, IPv6, private), firewalls, additional features, ssh-key, name
	// https://pkg.go.dev/github.com/hetznercloud/hcloud-go/hcloud?utm_source=godoc#Server
	// https://docs.hetzner.com/cloud/general/locations/

	// create the PublicKey in Golang? https://gist.github.com/devinodaniel/8f9b8a4f31573f428f29ec0e884e6673
	// unique pair of SSHkey and Server, thus new SSHKey for every server

	// possible Sizing approach: if the standard Size of the serverType is smaller than the required size, create a volume with the smallest available size and link it
	client := hcloud.NewClient(hcloud.WithToken(os.Getenv("API_TOKEN")))
	// setting up server options
	serverOpts := hcloud.ServerCreateOpts{Name: name}
	serverOpts.Location, _, _ = client.Location.Get(context.Background(), locationIDOrName)
	serverOpts.ServerType, _, _ = client.ServerType.GetByName(context.Background(), serverTypeName)
	serverOpts.Image, _, _ = client.Image.Get(context.Background(), imageNameOrID)
	err := serverOpts.Validate()
	fmt.Printf("name: %v\n", serverOpts.Image)

	if err != nil {
		fmt.Printf("%v\n", err)
	}
	result, _, err := client.Server.Create(context.Background(), serverOpts)
	if err != nil {
		fmt.Printf("%v", err)
	}
	fmt.Println(result)
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
		createServer("abc", "nbg1", "cx11", "79028095")
	default:
		fmt.Printf("Please enter the mode exactly once. You entered delete:%v create:%v\n", *crtKey, *delKey)
	}
}
