package helper

import (
	"context"
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/hetznercloud/hcloud-go/hcloud"
)

//func for creating title-------------------------------------------------------------------------------------------------------------------------------------
func ConfigureTitle() *canvas.Text {
	title := canvas.NewText("manage your SSHKeys", color.Black) //create Title
	title.TextStyle = fyne.TextStyle{
		Bold: true,
	}
	title.TextSize = 18
	title.Resize(fyne.NewSize(200, 0))
	title.Alignment = fyne.TextAlignCenter
	title.Move(fyne.NewPos(300, 15))
	return title
}

//func for creating title end----------------------------------------------------------------------------------------------------------------------------------

//funcs for creating SSHKey------------------------------------------------------------------------------------------------------------------------------------
func ButtonCreateSSHKey(myApp fyne.App) {
	createWindow := myApp.NewWindow("create new SSHKey")
	createWindow.Resize(fyne.NewSize(300, 200))
	label := widget.NewLabel("")

	inputName := widget.NewEntry()
	inputName.SetPlaceHolder("name for your new key")

	inputPublicKey := widget.NewEntry()
	inputPublicKey.SetPlaceHolder("your public key")

	inputToken := widget.NewEntry()
	inputToken.SetPlaceHolder("your api token")

	formCreate := widget.NewForm(
		widget.NewFormItem("Name", inputName),
		widget.NewFormItem("PublicKey", inputPublicKey),
		widget.NewFormItem("ApiToken", inputToken),
	)

	formCreate.OnCancel = func() {
		createWindow.Close()
	}
	formCreate.OnSubmit = func() {

		name := inputName.Text
		publicKey := inputPublicKey.Text
		token := inputToken.Text

		if name == "" {
			content := canvas.NewText("please enter a name", color.Black)
			information := widget.NewModalPopUp(content, createWindow.Canvas())
			information.Show()
			go wait5SecAndClose(createWindow)
		} else if publicKey == "" {
			content := canvas.NewText("please enter a public key", color.Black)
			information := widget.NewModalPopUp(content, createWindow.Canvas())
			information.Show()
			go wait5SecAndClose(createWindow)
		} else if token == "" {
			content := canvas.NewText("please enter an API token", color.Black)
			information := widget.NewModalPopUp(content, createWindow.Canvas())
			information.Show()
			go wait5SecAndClose(createWindow)
		} else {

			error := createSSHKey(name, publicKey, token, createWindow)
			if error == nil {
				content := canvas.NewText("The new SSHKey was created successfully", color.Black)
				information := widget.NewModalPopUp(content, createWindow.Canvas())
				information.Show()
				go wait5SecAndClose(createWindow)
			} else {
				content := canvas.NewText("something went wrong, please try again and check your input data", color.Black)
				information := widget.NewModalPopUp(content, createWindow.Canvas())
				information.Show()
				go wait5SecAndClose(createWindow)

				fmt.Printf("%v", error)
			}
		}
	}

	createWindow.SetContent(container.NewVBox(
		formCreate, label,
	))
	createWindow.Show()
}

func createSSHKey(name string, publicKey string, token string, createWindow fyne.Window) error {
	Name := name
	PublicKey := publicKey
	Labels := make(map[string]string)
	client := hcloud.NewClient(hcloud.WithToken(token))

	SSHKeyCreateOpts := hcloud.SSHKeyCreateOpts{Name: Name, PublicKey: PublicKey, Labels: Labels}
	_, _, error := client.SSHKey.Create(context.Background(), SSHKeyCreateOpts) //create SSHKey

	return error

}

func wait5SecAndClose(createWindow fyne.Window) {
	time.Sleep(5 * time.Second)
	createWindow.Close()
}

//funcs for creating SSHKey end------------------------------------------------------------------------------------------------------------------------------------

//funcs for deleting SSHKey----------------------------------------------------------------------------------------------------------------------------------------
func ButtonDeleteSSHKey(myApp fyne.App) {
	deleteWindow := myApp.NewWindow("create new SSHKey")
	deleteWindow.Resize(fyne.NewSize(300, 200))
	label := widget.NewLabel("")

	inputName := widget.NewEntry()
	inputName.SetPlaceHolder("name of the key to delete")

	inputToken := widget.NewEntry()
	inputToken.SetPlaceHolder("your api token")

	formCreate := widget.NewForm(
		widget.NewFormItem("Name", inputName),
		widget.NewFormItem("ApiToken", inputToken),
	)

	formCreate.OnCancel = func() {
		deleteWindow.Hide()
	}
	formCreate.OnSubmit = func() {

		name := inputName.Text
		token := inputToken.Text

		if name == "" {
			content := canvas.NewText("please enter a name", color.Black)
			information := widget.NewModalPopUp(content, deleteWindow.Canvas())
			information.Show()
			go wait5SecAndClose(deleteWindow)
		} else if token == "" {
			content := canvas.NewText("please enter an API token", color.Black)
			information := widget.NewModalPopUp(content, deleteWindow.Canvas())
			information.Show()
			go wait5SecAndClose(deleteWindow)
		} else {

			sshKeys, err, client := deleteSSHKey(name, token, label, deleteWindow)

			if sshKeys == nil {
				content := canvas.NewText("this Key does not exist", color.Black)
				information := widget.NewModalPopUp(content, deleteWindow.Canvas())
				information.Show()
				go wait5SecAndClose(deleteWindow)
			} else if err != nil {
				content := canvas.NewText(err.Error(), color.Black)
				information := widget.NewModalPopUp(content, deleteWindow.Canvas())
				information.Show()
				go wait5SecAndClose(deleteWindow)
				fmt.Printf("%v", err)
			} else {
				_, err = client.SSHKey.Delete(context.Background(), sshKeys) //delete SSHKey
				content := canvas.NewText("the SSHKey was deleted successfully", color.Black)
				information := widget.NewModalPopUp(content, deleteWindow.Canvas())
				information.Show()
				go wait5SecAndClose(deleteWindow)
			}
		}
	}

	deleteWindow.SetContent(container.NewVBox(
		formCreate, label,
	))
	deleteWindow.Show()
}

func deleteSSHKey(name string, token string, label *widget.Label, deleteWindow fyne.Window) (*hcloud.SSHKey, error, *hcloud.Client) {
	namE := name
	client := hcloud.NewClient(hcloud.WithToken(token))
	sshKeys, _, err := client.SSHKey.GetByName(context.Background(), namE)
	return sshKeys, err, client
}

func wait5SecAndClose2(deleteWindow fyne.Window) {
	time.Sleep(5 * time.Second)
	deleteWindow.Close()
}

//funcs for deleting SSHKey end------------------------------------------------------------------------------------------------------------------------------------

func Wait5SecAndClose3(listForm fyne.Window) {
	time.Sleep(5 * time.Second)
	listForm.Close()
}
