// sshOperator
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

func createServer(name string, locationIDOrName string, serverTypeName string, imageNameOrID string, userdata string, publicKey string) {
	// unique pair of SSHkey and Server, thus new SSHKey for every server
	client := hcloud.NewClient(hcloud.WithToken(os.Getenv("API_TOKEN")))
	// setting up server options
	serverOpts := hcloud.ServerCreateOpts{Name: name}
	serverOpts.Location, _, _ = client.Location.Get(context.Background(), locationIDOrName)
	serverOpts.ServerType, _, _ = client.ServerType.Get(context.Background(), serverTypeName)
	serverOpts.Image, _, _ = client.Image.Get(context.Background(), imageNameOrID)
	serverOpts.UserData = userdata
	// validation of server options
	err := serverOpts.Validate()
	if err != nil {
		fmt.Printf("Error during validation: %v\n", err)
		return
	}
	// if the other server options are correct, we create the publicKey
	createKey(name+"-Key", publicKey)
	publicKeySSH, _, _ := client.SSHKey.Get(context.Background(), name+"-Key")
	serverOpts.SSHKeys = append(serverOpts.SSHKeys, publicKeySSH)

	result, _, err := client.Server.Create(context.Background(), serverOpts)
	if err != nil {
		fmt.Printf("Error during creation: %v", err)
		fmt.Println(result)
	}
}

func testCreationTime() {
	start := time.Now()
	createServer("abc", "flk1", "cx11", "ubuntu-20.04", "", "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC1ORh6h8PpZ57zzx0rYBS/WjRu7ObAws6dSN+xQ5zcC1VZo2H/yJdcuyUU8HObkRZHRBTaMEbh3W3nnWj1PggeO7BQxUsLhtuSneI8FvIodbmYsyvAigReyv5pxfj9N0o06oCvkDP/kFTgidcAt1kUvBcSQfT97KltGYo4i+zVt6U+YCaeHOZTz7R11tHaOeh7b7A4z2olwcrhrfzq+s55WumvH0sM+Ohfh6Xo0FYgoO/G4XCLeymdYPbAA1JU96qarHF0sFBTv0zdCNl/grK2im4D4giSCjsYdxU9xFYLgsj8QIBZeAvQ7RSZTtlgh1IKsBvuQHBTwOzlVsb3YzJFVOI053TnMinhrJjJCtIWJYpVCW6QNNkMnCtiU+SAD0PKdX0uFF4Gy5/9K2m4PfPgyvtrjusPEGgkt3+BeKgbZHhoX8efktVBaj/aph0PUum3VkSPfBbduISsypl2cXCIOeTshBg3zPQxptK9qepMF1DY8JkRgQNSjcjPWy0MrLlAaG/UiUvgeFXhr6Hi5paIZ9bzSv1V66MNHvlxW3HXj4LtQjbZnDFfLo/pK+fMjSwW4ZDewgvYPrevMFvxEansEPbAIPvd0SYCjbRyOdSRH7hNH1bOapxiSZTD1Ja1P4umbRe1RXyRBgx02T7sAKvqJkUqpkgwDbowi6TxdTEXuQ== kimi@kimiarch")
	elapsed := time.Since(start)
	fmt.Printf("Creation of Server took %s", elapsed)
	//deleteServer("abc")
}

/*
Ausprobieren:
Server erstellen, userdata:
Cloudinit angucken
über cloudinit Programme erstellen
ubuntu 20.04
über userdaten Programme installieren
Ziel: Docker und ZSH installieren
überprüfen
Was ist schneller? Docker über cloudinit zsh oder Docker über snapshot
*/

func deleteServer(name string) {
	client := hcloud.NewClient(hcloud.WithToken(os.Getenv("API_TOKEN")))
	serverToBeDeleted, _, err := client.Server.Get(context.Background(), name)
	if err != nil {
		fmt.Println(err)
	}
	deleteKey(name + "-Key")
	client.Server.Delete(context.Background(), serverToBeDeleted)
}

func createKey(name, publicKey string) {
	labels := make(map[string]string)
	client := hcloud.NewClient(hcloud.WithToken(os.Getenv("API_TOKEN")))
	SSHKeyCreateOpts := hcloud.SSHKeyCreateOpts{Name: name, PublicKey: publicKey, Labels: labels}
	_, _, error := client.SSHKey.Create(context.Background(), SSHKeyCreateOpts) //create SSHKey
	if error != nil {
		fmt.Printf("Error while creating sshKey: %v", error) //Print out SSHKey
	}

}

func deleteKey(name string) {
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
	del := flag.Bool("del", false, "delete server")
	flag.Parse()
	switch {
	case *delKey && !*crtKey:
		deleteKey(os.Args[2])
	case *crtKey && !*delKey:
		name := os.Args[2]
		PublicKey := os.Args[3]
		createKey(name, PublicKey)
	case *serv:
		testCreationTime()
	case *del:
		deleteServer("abc")
	default:
		fmt.Printf("Please enter the mode exactly once. You entered delete:%v create:%v\n", *crtKey, *delKey)
	}
}
