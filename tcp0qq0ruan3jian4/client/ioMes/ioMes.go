package ioMes

import (
	"bufio"
	"fmt"
	"go7tcp0liao2tian1app/20200619tcp0qq0ruan3jian4/client/menue"
	"go7tcp0liao2tian1app/20200619tcp0qq0ruan3jian4/client/origin"
	"go7tcp0liao2tian1app/20200619tcp0qq0ruan3jian4/client/register0login"
	"go7tcp0liao2tian1app/20200619tcp0qq0ruan3jian4/client/utils"
	"io"
	"net"
	"os"
	"strings"
)
//循环持续接收server端发来的系统消息或server转发来的其他用户群聊发来的聊天消息，
//展示QQ菜单界面上的功能选项(用户选择后即是向server端发送系统消息)，
//及本用户发出聊天消息
func RecieveMsg(conn net.Conn, chanBool chan bool) {
	buf := make([]byte, 1024*3)
	var response string
	for {
		n, err := conn.Read(buf)
		if err == io.EOF {
			//conn.Close() //如果服务端已关闭跟该用户的conn，就会err == io.EOF 这里客户端就也关闭，不关也没事，毕竟服务端已经关了
			//chanBool <-true  //有没有这个关系也不是很大
			break
		}
		utils.ErrExit("35行报错", err)
		response = string(buf[:n])
		//接收到服务端发来的非聊天信息
		if strings.Contains(response, "@") {
			fmt.Println(response + "(系统消息)")
			strSlice := strings.Split(response, "@")
			switch strSlice[1] {
			case "用户已存在，无需注册，请登录":
				register0login.MenuePreface()
			case "注册成功，请登录":
				register0login.MenuePreface()
			case "注册失败":
				register0login.MenuePreface()
			case "用户名或密码错误":
				register0login.MenuePreface()
			case "登录成功":
				fmt.Println("登录成功，欢迎:", strSlice[0])
				utils.WritePkg("upline", strSlice[0], conn) //strSlice[0]是成功登录的用户名
//todo:选择菜单上的功能让server发来好友列表、群聊列表、转发出加群加好友申请等
				menue.MenueContent()                              //显示完菜单，马上就又循环回去等待接收服务端的消息
			case "申请添加你为好友":
				fmt.Print("y、允许\n" + "n、不允许\n")
				fmt.Println("此为浮窗消息，请menue中处理：请选择代号")
				origin.ChanCommon <- strSlice[0]
				//var n string  //todo:数据类型的不同可以避免和MenueContent的scanln竞争，但最好还是不这样做
				//fmt.Scanln(&n)
				//if n == "y" {
				//	utils.WritePkg("yes",strSlice[0],conn)//strSlice[0]是想要添加好友的发起方的名字
				//}else if n== "n" {
				//	utils.WritePkg("no",strSlice[0],conn)//strSlice[0]是想要添加好友的发起方的名字
				//}else {
				//	fmt.Println("请合法输入，请选择序号")
				//}
			case "申请加群":
				fmt.Print("admit、允许\n" + "refuse、不允许\n")
				fmt.Println("浮窗消息请menue中处理：请选择单词")
				origin.ChanCommon <- strSlice[0] //申请加群的人
				origin.ChanCommon <- strSlice[2] //申请人想加入的群名
				//var n string
				//fmt.Scanln(&n)
				//if n == "admit" {
				//	utils.WritePkg("admit",strSlice[0]+"@"+strSlice[2],conn)//strSlice[0]是申请加群的发起方的名字
				//}else if n=="refuse" {
				//	utils.WritePkg("refuse",strSlice[0]+"@"+strSlice[2],conn)//strSlice[2]是群名
				//}else {
				//	fmt.Println("请合法输入，请选择序号")
				//}

			case "添加好友成功":
				fmt.Println("此为浮窗消息，请menue中处理：你已允许申请，打开好友列表和该好友聊天吧！")
				//MenueContent()
			case "添加好友失败":
				fmt.Println("此为浮窗消息，请menue中处理：您添加好友失败，等待对方再次添加你为好友")
			case "申请通过，快和新好友聊天吧":
				fmt.Println("此为浮窗消息，请menue中处理：申请通过，打开好友列表和该好友聊天吧！")
				//MenueContent()  和最上面的MenueContent函数的第一次递归重了，导致争抢cmd的输入
			case "申请未通过":
				fmt.Println("此为浮窗消息，请menue中处理：申请未通过")
				//MenueContent()
			case "好友列表":
				fmt.Println(strSlice[0]) //打印出好友列表
				fmt.Println("聊天开始，请以----要说的话#好友名字---的格式进行聊天")
				//这里还需想办法限制用户#一个好友列表中没有的好友来发消息的操作
				fmt.Println("输入单词menue将退出聊天室并显示菜单,再发聊天消息将发不出去")
//todo:在持续循环接收消息的同时开启发出聊天消息的协程
				go sendMsg(conn)
			case "所在全部群的列表":
				fmt.Println(strSlice[0]) //打印出所在全部群的列表
				fmt.Println("请输入你想进入群聊的群名,输入menue返回菜单")
				var groupName string
				fmt.Scanln(&groupName)
				if groupName == "menue" {
					menue.MenueContent()
					continue
				}
				utils.WritePkg("intoThisGroup", groupName, conn)
			case "该群不存在，请正确输入群列表中的群名":
				fmt.Println(strSlice[1]) //打印出这句话：该群不存在，请正确输入群列表中的群名
			case "你申请加群被允许":
				utils.WritePkg("beAdmited", strSlice[0]+"@"+strSlice[2], conn)
			case "请开始群聊":
				fmt.Println("请开始群聊，menue返回菜单，群聊消息无所需遵守的聊天格式")
				fmt.Println("如收到好友私聊消息，可按  ***#好友名字 的格式回复")
//todo:在接收消息的同时开启发聊天消息的协程
				go sendGroupMsg(strSlice[0], conn) //strSlice[0]是群聊名
			case "新建群成功":
				fmt.Println("你成功新建了一个聊天群")
				menue.MenueContent()
			}
			//接收到服务端转发来的聊天消息
		} else if strings.Contains(response, "#") {
			fmt.Println(response)
		} else {
			fmt.Println("后台发来无@的系统提示用户消息")
			fmt.Println(response)
		}

	}
}

//循环读取cmd输入端内容作为聊天消息发送出去或切换到菜单
func sendMsg(conn net.Conn) {
	reader := bufio.NewReader(os.Stdin)
	var input string
	for {
		lineByte, _, err := reader.ReadLine()
		if err == io.EOF {
			conn.Close()
			return
		}
		utils.ErrContinue("45行报错", err)
		input = string(lineByte)

		if input == "1" || input == "2" || input == "3" || input == "4" || input == "5" ||
			input == "6" || input == "7" || input == "8" || input == "9" {
			fmt.Println("您是想选择菜单吗？请确认选择")
			menue.MenueContent()
			return
			//break	//todo:这一步极为重要，按菜单提示输入指令和发送聊天消息是可以混合进行的，这里则将其分开
		}
		if input == "menue" {
			menue.MenueContent()
			return //为何不能是break？？？
			//break   //todo:这一步极为重要，按菜单提示输入指令和发送聊天消息是可以混合进行的，这里则将其分开
		}
		if input == "y" {
			friendApplyer := <-origin.ChanCommon
			utils.WritePkg("yes", friendApplyer, conn)
			continue
		}
		if input == "n" {
			friendApplyer := <-origin.ChanCommon
			utils.WritePkg("no", friendApplyer, conn)
			continue
		}
		if input == "admit" {
			groupApplyer := <-origin.ChanCommon
			groupName := <-origin.ChanCommon
			utils.WritePkg("admit", groupApplyer+"@"+groupName, conn)
			continue
		}
		if input == "refuse" {
			groupApplyer := <-origin.ChanCommon
			groupName := <-origin.ChanCommon
			utils.WritePkg("refuse", groupApplyer+"@"+groupName, conn)
			continue
		}
		if !(strings.Contains(input, "#")) {
			fmt.Println("聊天消息发送失败，请再次编辑消息，以----****#好友名字---的格式进行聊天")
			continue
		}
//todo:发出聊天消息
		utils.WritePkg("chatmessage", input, conn)
	}
}

//发送群聊消息或切换到菜单
func sendGroupMsg(groupname string, conn net.Conn) {
	reader := bufio.NewReader(os.Stdin)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			conn.Close()
			return
		}
		lineStr := string(line)
		if err != nil {
			fmt.Println("腾讯接收消息出错，请再次输入")
			continue
		}
		if lineStr == "y" {
			friendApplyer := <-origin.ChanCommon
			utils.WritePkg("yes", friendApplyer, conn)
			fmt.Println("请继续输入群聊消息")
			continue
		}
		if lineStr == "n" {
			friendApplyer := <-origin.ChanCommon
			utils.WritePkg("no", friendApplyer, conn)
			fmt.Println("请继续输入群聊消息")
			continue
		}
		if lineStr == "admit" {
			groupApplyer := <-origin.ChanCommon
			groupName := <-origin.ChanCommon
			utils.WritePkg("admit", groupApplyer+"@"+groupName, conn)
			fmt.Println("请继续输入群聊消息")
			continue
		}
		if lineStr == "refuse" {
			groupApplyer := <-origin.ChanCommon
			groupName := <-origin.ChanCommon
			utils.WritePkg("refuse", groupApplyer+"@"+groupName, conn)
			fmt.Println("请继续输入群聊消息")
			continue
		}
		if strings.Contains(lineStr, "#") {
			utils.WritePkg("chatmessage", lineStr, conn)
			fmt.Println("请继续输入群聊消息")
			continue
		}
		if lineStr == "menue" {
			utils.WritePkg("withdrawGroup", groupname, conn) //退出群聊，不再接收该群的群聊消息
			menue.MenueContent()
			break //todo:这一步极为重要，按菜单提示输入指令和发送聊天消息是无法混合进行的
		}
		utils.WritePkg("groupmessage", lineStr+"@"+groupname, conn)
	}
}