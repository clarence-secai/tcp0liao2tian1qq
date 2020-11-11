package register0login

import (
	"fmt"
	"go7tcp0liao2tian1app/20200619tcp0qq0ruan3jian4/client/origin"
	"go7tcp0liao2tian1app/20200619tcp0qq0ruan3jian4/client/utils"
	"go7tcp0liao2tian1app/20200619tcp0qq0ruan3jian4/common"
	"net"
)

//登录和注册的界面
func MenuePreface() {
	var n int
	for { //用于避免客户端乱输入除了1、2之外的其他东西
		fmt.Println("请选择序号")
		fmt.Print("1、注册\n" + "2、登录\n")
		fmt.Scanln(&n)
		if n == 1 || n == 2 {
			break
		} else {
			fmt.Println("请合法输入")
			continue
		}
	}
	switch n {
	case 1:
		userName, passWord := MenuRegistLogin()
		//进行注册
		regist(userName, passWord, origin.OneConn)
	case 2:
		userName, passWord := MenuRegistLogin()
		//进行登录
		login(userName, passWord, origin.OneConn)
	}
}

//获取用户输入的登录或注册信息
func MenuRegistLogin() (string, int64) {
	fmt.Println("请输入用户名")
	var userName string
	fmt.Scanln(&userName)
	fmt.Println("请输入用户密码")
	var passWord int64
	fmt.Scanln(&passWord)
	return userName, passWord
}

//发送登录消息给server端
func login(name string, word int64, conn net.Conn) {
	user := common.NewUser(name, word)
	head := "loginMessage"
	//todo:用utils.go文件中WritePkg函数发送信息包
	utils.WritePkg(head, user, conn) //这里的user是指针类型
}

//发送注册消息给server端
func regist(userName string, passWord int64, conn net.Conn) {
	user := common.NewUser(userName, passWord) //model.go文件中用户结构体，信息包的body
	//todo:用utils.go文件中WritePkg函数发送信息包
	utils.WritePkg("registMessage", user, conn) //注意这里user是指针类型
}
