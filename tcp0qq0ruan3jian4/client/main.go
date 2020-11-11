package main

import (

	"go7tcp0liao2tian1app/20200619tcp0qq0ruan3jian4/client/ioMes"
	"go7tcp0liao2tian1app/20200619tcp0qq0ruan3jian4/client/origin"
	"go7tcp0liao2tian1app/20200619tcp0qq0ruan3jian4/client/register0login"
)

func main() {
	//初始化打开本main.go终端的用户与server的连接OneConn
	origin.Origin()

	//启动后的登录注册界面，向server端发送注册、登录信息
	register0login.MenuePreface()

	chanBool := make(chan bool)

	//QQ菜单界面，持续接收server端发来的消息及用户发出消息
	go ioMes.RecieveMsg(origin.OneConn, chanBool)

	<-chanBool
}






