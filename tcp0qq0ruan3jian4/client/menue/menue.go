package menue

import (
	"fmt"
	"go7tcp0liao2tian1app/20200619tcp0qq0ruan3jian4/client/origin"
	"go7tcp0liao2tian1app/20200619tcp0qq0ruan3jian4/client/utils"
	"net"
	"time"
)
//实质功能菜单界面，选择响应功能就是向server端发送系统消息
func MenueContent() {
	fmt.Println("菜单---请选择序号")
	fmt.Print("3、好友列表\n" + "4、群列表\n" + "5、查加新友\n" +
		"6、查加群\n" + "7、新建群\n" + "8、下线\n" + "9、服务器所有在线用户\n")
	var n string
	fmt.Scanln(&n)
	switch n {
	case "y":
		friendApplyer := <-origin.ChanCommon
		utils.WritePkg("yes", friendApplyer, origin.OneConn)
		MenueContent()
	case "n":
		friendApplyer := <-origin.ChanCommon
		utils.WritePkg("no", friendApplyer, origin.OneConn)
		MenueContent()
	case "admit":
		groupApplyer := <-origin.ChanCommon
		groupName := <-origin.ChanCommon
		utils.WritePkg("admit", groupApplyer+"@"+groupName, origin.OneConn)
		MenueContent()
	case "refuse":
		groupApplyer := <-origin.ChanCommon
		groupName := <-origin.ChanCommon
		utils.WritePkg("refuse", groupApplyer+"@"+groupName, origin.OneConn)
		MenueContent()
	case "3": //让服务端向客户端显示自己的好友列表
		utils.WritePkg("friendsList", "nothing", origin.OneConn)
	case "4": //让服务端向客户端显示自己所在的全部群的列表
		utils.WritePkg("groupList", "nothing", origin.OneConn)
	case "5":
		var friendName string
		fmt.Println("请输入你要查加的新友的名字")
		fmt.Scanln(&friendName)
		ExploreFriend(friendName, origin.OneConn)
		time.Sleep(2 * time.Second) //等待服务端发来"已发送添加好友申请，等待通过，您可先进行其他操作"再显示下面的递归，当然可以用管道来确保万无一失
		MenueContent()              //todo:这里有递归，其实可以用for循环？？？只是这一步有必要再次显示菜单，其他case运行后并无此必要
	case "6":
		fmt.Println("请输入你想查加群的群名称")
		var groupName string
		fmt.Scanln(&groupName)
		utils.WritePkg("addGroup", groupName, origin.OneConn)
		time.Sleep(2 * time.Second) //等待服务端发来"加群申请等待通过,你可先进行其他操作"再显示下面的菜单，当然可以用管道来确保万无一失
		MenueContent()
	case "7":
		fmt.Println("请给新建群起一个群名")
		var groupName string
		fmt.Scanln(&groupName)
		utils.WritePkg("setupGroup", groupName, origin.OneConn)
	case "8":
		utils.WritePkg("downLine", "nothing", origin.OneConn)

	case "9":
		utils.WritePkg("allOnlineUsers", "nothing", origin.OneConn)

	default:
		fmt.Println("菜单栏接收序号指令失败，请合法输入")
		MenueContent()
	}
}

func ExploreFriend(friendName string, conn net.Conn) {
	utils.WritePkg("exploreFriend", friendName, conn)
}