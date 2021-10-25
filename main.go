package main

import (
	"DrGo/dr"
	"DrGo/theme"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/Unknwon/goconfig"
	"github.com/mitchellh/go-homedir"
	"log"
	"os"
)

//config 读取配置文件路径并自动创建文件
func config(path map[string]string) {
	if _, err := os.Stat(path["dir"]); err == nil { //目录存在文件不存在
		if _, err = os.Stat(path["dir"] + path["filename"]); err != nil {
			_, err = os.Create(path["dir"] + path["filename"])
			if err != nil {
				log.Fatalln("出错！目录存在但是在创建文件过程中发生错误！", err)

			}
		}
	} else { //目录不存在
		err := os.MkdirAll(path["dir"], 0711)
		if err != nil { //
			log.Fatalln("出错！目录不存在且创建目录失败！", err)
		} else {
			if _, err = os.Stat(path["dir"] + path["filename"]); err != nil {
				_, err = os.Create(path["dir"] + path["filename"])
				if err != nil {
					log.Fatalln("出错！目录创建成功但是在创建文件过程中发生错误！", err)
				}
			}
		}
	}
}

func main() {
	//windows 隐藏客户端
	//go build -ldflags -H=windowsgui
	dir, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}
	var path = make(map[string]string, 2)
	path["dir"] = dir + "\\DrGo\\"
	path["filename"] = "config.ini"
	//检测配置文件是否存在 不存在则创建
	config(path)
	//加载配置文件
	c, err := goconfig.LoadConfigFile(path["dir"] + path["filename"])
	if err != nil {
		log.Println("加载配置文件出错")
		return
	}
	//创建窗体
	myApp := app.New()
	myApp.Settings().SetTheme(&theme.MyTheme{})
	myWindow := myApp.NewWindow("DrGo")

	account := widget.NewEntry()
	pwd := widget.NewPasswordEntry()
	account.PlaceHolder = "请输入账号"
	pwd.PlaceHolder = "请输入密码"
	userNameConf, err := c.GetValue("dr", "username")
	if userNameConf != "" && err == nil {
		account.Text = userNameConf
	}
	pwdConf, err := c.GetValue("dr", "pwd")
	if pwdConf != "" && err == nil {
		pwd.Text = pwdConf
		pwd.Wrapping = fyne.TextTruncate
		pwd.Password = true
	}
	textBalance := widget.NewLabel("UnKnown")

	//GUI不支持中文所以这里采用英文
	var form *widget.Form
	form = &widget.Form{
		Items: []*widget.FormItem{
			{Text: "账号", Widget: account},
		},
		SubmitText: "登录",
		OnSubmit: func() {
			if account.Text == "" || pwd.Text == "" {
				dialog.ShowInformation("Error", "Please Enter Username And Password", myWindow)
				return
			}
			log.Println("Username submitted:", account.Text)
			log.Println("Password submitted:", pwd.Text)
			result := dr.Login(account.Text, pwd.Text)
			if result {
				c.SetValue("dr", "username", account.Text)
				c.SetValue("dr", "pwd", pwd.Text)
				err := goconfig.SaveConfigFile(c, path["dir"]+path["filename"])
				if err != nil {
					log.Println("存储账号信息出错", err)
				}
				//获取当前余额
				balance := dr.GetBalance()
				textBalance.SetText(balance)
				dialog.ShowInformation("Success", "Login Success", myWindow)
			}
		},
		CancelText: "注销",
		OnCancel: func() {
			result := dr.Logout()
			if result {
				//注销成功
				textBalance.SetText("Unknown")
				dialog.ShowInformation("Success", "Logout Success", myWindow)
			} else {
				dialog.ShowInformation("Error", "Login Error", myWindow)
			}

		},
	}
	form.Append("密码", pwd)
	form.Append("余额", textBalance)
	myWindow.CenterOnScreen()
	myWindow.Resize(fyne.NewSize(309, 132))
	myWindow.SetContent(form)
	myWindow.ShowAndRun()

}
