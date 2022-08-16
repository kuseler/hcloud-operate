// sshOperator
package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

func do_keyscan(ip string) {
	//ssh-keyscan -t rsa github.com >> ~/.ssh/known_hosts;
	app := "ssh-keyscan"
	arg1 := "-t"
	arg2 := "rsa"
	arg3 := ip
	cmd := exec.Command(app, arg1, arg2, arg3)
	output, err := cmd.CombinedOutput()
	out := strings.Fields(string(output))
	fmt.Printf("output: %v", out)
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + string(output))
		return
	}
	path := "/home/kimi/.ssh/known_hosts"
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	if _, err = f.WriteString(strings.Join(out[len(out)-3:], " ")); err != nil {
		panic(err)
	}
}
func raw_connect(host string, port string) {
	timeout := time.Second
	for {
		conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
		if err != nil {
			fmt.Println("Connecting error:", err)
		}
		fmt.Println(conn)
		if conn != nil {
			defer conn.Close()
			fmt.Println("Opened", net.JoinHostPort(host, port))
			break
		}
	}
}

func doconn(ip string) {
	raw_connect(ip, "22")
	do_keyscan(ip)
	hostKeyCallback, err := knownhosts.New("/home/kimi/.ssh/known_hosts")
	if err != nil {
		fmt.Println("hostkeyerror: ", err)
	}
	key, err := ioutil.ReadFile("/home/kimi/.ssh/id_rsa")
	if err != nil {
		fmt.Printf("unable to read private key: %v", err)
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		fmt.Printf("unable to parse private key: %v", err)
	}

	config := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			// Use the PublicKeys method for remote authentication.
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: hostKeyCallback,
	}

	config.SetDefaults()
	for {
		for {
			_, err := ssh.Dial("tcp", ip+":22", config)
			time.Sleep(100 * time.Millisecond)
			if err != nil {
				break
			}
		}
		client, err := ssh.Dial("tcp", ip+":22", config)
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer client.Close()
		session, _ := client.NewSession()
		answer, err := session.Output("docker version")
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(answer)
		if answer != nil {
			break
		}

	}
}

func createServer(name string, locationIDOrName string, serverTypeName string, imageNameOrID string, userdata string, publicKey string) {
	// unique pair of SSHkey and Server, thus new SSHKey for every server
	client := hcloud.NewClient(hcloud.WithToken(os.Getenv("API_TOKEN")))
	// setting up server options
	serverOpts := hcloud.ServerCreateOpts{Name: name}
	serverOpts.Location, _, _ = client.Location.Get(context.Background(), locationIDOrName)
	serverOpts.ServerType, _, _ = client.ServerType.Get(context.Background(), serverTypeName)
	serverOpts.Image, _, _ = client.Image.Get(context.Background(), imageNameOrID)
	if userdata != "" {
		serverOpts.UserData = userdata
	}
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
	for {
		server, _, _ := client.Server.Get(context.Background(), "abc")
		if server.Status == "running" {
			break
			fmt.Println("breaking loop", server.Status)
		}
		time.Sleep(10 * time.Millisecond)
		fmt.Println(server.PublicNet.IPv4)

	}
	server, _, _ := client.Server.Get(context.Background(), "abc")
	doconn(server.PublicNet.IPv4.IP.String())
}

func testCreationTime() {
	start := time.Now()
	cloudconfig, _ := ioutil.ReadFile("cloudinit.yaml")
	fmt.Println(cloudconfig)
	createServer("abc", "flk1", "cx11", "ubuntu-20.04", string(cloudconfig), "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC1ORh6h8PpZ57zzx0rYBS/WjRu7ObAws6dSN+xQ5zcC1VZo2H/yJdcuyUU8HObkRZHRBTaMEbh3W3nnWj1PggeO7BQxUsLhtuSneI8FvIodbmYsyvAigReyv5pxfj9N0o06oCvkDP/kFTgidcAt1kUvBcSQfT97KltGYo4i+zVt6U+YCaeHOZTz7R11tHaOeh7b7A4z2olwcrhrfzq+s55WumvH0sM+Ohfh6Xo0FYgoO/G4XCLeymdYPbAA1JU96qarHF0sFBTv0zdCNl/grK2im4D4giSCjsYdxU9xFYLgsj8QIBZeAvQ7RSZTtlgh1IKsBvuQHBTwOzlVsb3YzJFVOI053TnMinhrJjJCtIWJYpVCW6QNNkMnCtiU+SAD0PKdX0uFF4Gy5/9K2m4PfPgyvtrjusPEGgkt3+BeKgbZHhoX8efktVBaj/aph0PUum3VkSPfBbduISsypl2cXCIOeTshBg3zPQxptK9qepMF1DY8JkRgQNSjcjPWy0MrLlAaG/UiUvgeFXhr6Hi5paIZ9bzSv1V66MNHvlxW3HXj4LtQjbZnDFfLo/pK+fMjSwW4ZDewgvYPrevMFvxEansEPbAIPvd0SYCjbRyOdSRH7hNH1bOapxiSZTD1Ja1P4umbRe1RXyRBgx02T7sAKvqJkUqpkgwDbowi6TxdTEXuQ== kimi@kimiarch")
	elapsed := time.Since(start)
	fmt.Printf("Creation of Server took %s\n", elapsed)
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
		for i := 0; i < 1; i++ {
			testCreationTime()
		}
	case *del:
		deleteServer("abc")
	default:
		fmt.Printf("Please enter the mode exactly once. You entered delete:%v create:%v\n", *crtKey, *delKey)
	}
}
