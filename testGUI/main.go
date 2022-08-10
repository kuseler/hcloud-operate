package main

import (
	"context"
	"fmt"
	"image/color"
	"testGUI/helper"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/hetznercloud/hcloud-go/hcloud"
)

var myApp = app.New()
var AllSSHKeys []*hcloud.SSHKey

func main() {
	mainWindow := myApp.NewWindow("SSHKey manager")
	mainWindow.Resize(fyne.NewSize(800, 500))

	title := helper.ConfigureTitle()

	buttonCreate := widget.NewButton("create a new SSHKey", func() { //create the "create SSHKey" button and set logic for it
		helper.ButtonCreateSSHKey(myApp)
	})
	buttonCreate.Resize(fyne.NewSize(200, 40))
	buttonCreate.Move(fyne.NewPos(300, 150))

	buttonDelete := widget.NewButton("delete a SSHKey", func() { //create the "delete SSHKey" button and set logic for it
		helper.ButtonDeleteSSHKey(myApp)
	})
	buttonDelete.Resize(fyne.NewSize(200, 40))
	buttonDelete.Move(fyne.NewPos(300, 250))

	buttonList := widget.NewButton("list existing Keys", func() {

		listForm := myApp.NewWindow("verification")
		listForm.Resize(fyne.NewSize(300, 200))

		inputToken := widget.NewEntry()
		inputToken.SetPlaceHolder("your API token")

		formCreate := widget.NewForm(
			widget.NewFormItem("enter your API token", inputToken),
		)

		formCreate.OnCancel = func() {
			listForm.Hide()
		}
		formCreate.OnSubmit = func() {

			token := inputToken.Text

			if token == "" {
				content := canvas.NewText("please enter a token", color.Black)
				information := widget.NewModalPopUp(content, listForm.Canvas())
				information.Show()
				go helper.Wait5SecAndClose3(listForm)

			} else {
				client := hcloud.NewClient(hcloud.WithToken(token))
				allSSHKeys, err := client.SSHKey.All(context.Background())
				if err != nil {
					content := canvas.NewText("something went wrong, please try again", color.Black)
					information := widget.NewModalPopUp(content, listForm.Canvas())
					information.Show()
					go helper.Wait5SecAndClose3(listForm)
					fmt.Print(err)

				} else {

					fmt.Print(allSSHKeys[0])
					AllSSHKeys := allSSHKeys
					fmt.Print(AllSSHKeys)

					/*var data = [][]string{[]string{"Name", "SSHKey"}, //data for myTable
						[]string{"bottom left", "bottom right"}}

					myTable := widget.NewTable(
						func() (int, int) {
							return len(data), len(data[0])
						},
						func() fyne.CanvasObject {
							return widget.NewLabel("wide content")
						},
						func(i widget.TableCellID, o fyne.CanvasObject) {
							o.(*widget.Label).SetText(data[i.Row][i.Col])
						})

					content := myTable
					listWindow.SetContent(content)
					listWindow.Show()*/
				}
			}
		}

		content := formCreate
		listForm.SetContent(content)
		listForm.Show()
	})
	buttonList.Resize(fyne.NewSize(200, 40))
	buttonList.Move(fyne.NewPos(300, 350))

	content := container.NewWithoutLayout(title, buttonCreate, buttonDelete, buttonList)

	mainWindow.SetContent(content)
	mainWindow.ShowAndRun()
}
