package ioMes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go7tcp0liao2tian1app/20200619tcp0qq0ruan3jian4/common"
	"go7tcp0liao2tian1app/20200619tcp0qq0ruan3jian4/server/dao"
	"go7tcp0liao2tian1app/20200619tcp0qq0ruan3jian4/server/origin"
	"go7tcp0liao2tian1app/20200619tcp0qq0ruan3jian4/server/register0login"
	"go7tcp0liao2tian1app/20200619tcp0qq0ruan3jian4/server/serve"
	"go7tcp0liao2tian1app/20200619tcp0qq0ruan3jian4/server/utils"
	"io"
	"net"
	"strings"
	"time"
)

//循环读取该专属连接所属用户发来的系统消息或需转发的聊天消息，发出系统消息或转发消息
func IoWithConn(conn net.Conn) {
	buf := make([]byte, 1024*3)
	var thisUser *common.OnlineUser
	for {
		n, err := conn.Read(buf)
		if err == io.EOF {
			fmt.Println("43行，客户端不再发来消息")
			return
		}
		utils.ErrContinue("22行", err)

		var msg common.Msg
		//todo:此时只完成了第一次反序列化，msg.Head已可直接用，msg.Body还需一次反序列化
		err = json.Unmarshal(buf[:n], &msg)
		utils.ErrExit("27行报错", err)

		//todo:根据信息包头的描述信息，对信息的body部分进行多路复用器似的处理
		switch msg.Head {
		case "registMessage":
			result := register0login.DealRegist(msg.Body) //msg.Body是一个user，包含姓名、密码
			conn.Write([]byte(result))
		case "loginMessage":
			result := register0login.DealLogin(msg.Body) //msg.Body是一个user，包含姓名、密码
			conn.Write([]byte(result))    //result 为 username@登录成功"或者"ha@用户名或密码错误"
		case "upline": //一旦上线，添加进在线用户列表，该步骤会在用户成功登录后才执行 //其实也可以合并到上一个case
			var upLineName string
			json.Unmarshal([]byte(msg.Body), &upLineName)

			//todo:把和该客户端的专属conn连接、该客户的信息二者捆绑维护，便于后面
			// 按用户能找到专属该用户的长连接，依靠该长连接向该用户发消息
			origin.OnlineUsers[upLineName] = &common.OnlineUser{
				Name: upLineName,
				Conn: conn,
				Addr: conn.RemoteAddr().String(),
			}
			fmt.Println(upLineName, "上线，进入OnlineUsers的map里")

			//todo:thisUser在用户登录上线后完成实例化，方便后续想找到当前conn对应的用户姓名的用户信息
			thisUser = origin.OnlineUsers[upLineName]

			//todo:将该用户所属的全部群的自己添加进服务端维护的若干相应群的在线群成员中
			//todo:这样可以避免在一个群中广播消息每次都要去根据群成员名字遍历OnlineUsers，因为OnlineUsers用户很多，这样效率较低
			groupByte1 := dao.ShowGroupsList(upLineName)//显示当前用户所加的全部群
			groupByte2 := bytes.Split(groupByte1, []byte("@"))
			var groupSlice []string
			json.Unmarshal(groupByte2[0], &groupSlice)
			for _, groupName := range groupSlice {
				//由于所有群在OnlineGroup中都有维护着，故必定能找到，无需判断
				origin.OnlineGroup[groupName].Members = append(origin.OnlineGroup[groupName].Members, thisUser)
				//&common.OnlineUser{
				//		Name:upLineName,
				//		Conn: conn,
				//		Addr: conn.RemoteAddr().String(),
				//})
			}
			fmt.Println(upLineName, "在所在的全部群内的自己完成上线")
		case "beAdmited":
			var applyerAndGoalGroupStr string
			json.Unmarshal([]byte(msg.Body), &applyerAndGoalGroupStr)
			applyerAndGoalGroup := strings.Split(applyerAndGoalGroupStr, "@")
			//applyer := applyerAndGoalGroup[0]//此时的applyer就是thisuser自己
			appliedGroup := applyerAndGoalGroup[1]
			origin.OnlineGroup[appliedGroup].Members = append(origin.OnlineGroup[appliedGroup].Members, thisUser)
			//OnlineUsers[thisUser.Name].Conn.Write([]byte(thisUser.Name+"@你申请加群被允许@"+appliedGroup))
			origin.OnlineUsers[thisUser.Name].Conn.Write([]byte("加群申请通过，你已是群名为：" + appliedGroup + "：的群成员"))
		case "downLine":
			thisUser.Location = "" //从当前群退出，抹去标记  其实这一步没必要，后面的代码就将该用户从服务器中的用户和群列表中删去了
			downLineName := thisUser.Name

			//从所有在线用户列表中退出
			delete(origin.OnlineUsers, downLineName)

			//从服务器维护的OnlineGroup中该用户所有所在的群聊中退出在线
			groupByte1 := dao.ShowGroupsList(downLineName)
			groupByte2 := bytes.Split(groupByte1, []byte("@"))
			var groupSlice []string
			json.Unmarshal(groupByte2[0], &groupSlice)
			for _, groupName := range groupSlice {
				//由于所有群在OnlineGroup中都有维护着，故必定能找到，无需判断
				for k, v := range origin.OnlineGroup[groupName].Members {
					if v.Name == downLineName {
						origin.OnlineGroup[groupName].Members = append(origin.OnlineGroup[groupName].Members[:k], origin.OnlineGroup[groupName].Members[k+1:]...)
					}
				}
			}
			thisUser.Conn.Write([]byte("你成功下线"))
			//todo:关闭与该用户的conn
			thisUser.Conn.Close()
			return                //todo:结束与该用户的conn的交互
		case "exploreFriend":
			var friendName string
			json.Unmarshal([]byte(msg.Body), &friendName)
			result := serve.DealExploreFriend(friendName) //从数据库中查找是否有名为friendName的用户
			if result {                             //该用户是注册用户，是存在的用户
				conn.Write([]byte("已发送添加好友申请，等待通过，您可先进行其他操作"))
				//todo:开启一个背后的小协程，默默地等待被申请加好友的用户上线，告知有人想加他为好友
				go func() {
					serve.Wait2conn(friendName, 2*time.Second, thisUser.Name+"@申请添加你为好友")
				}() //存在一个问题是多条留言因为是协程，不一定按先后顺序显示给上线的客户端
			} else {
				conn.Write([]byte("你要找的用户不存在"))
			}
		case "yes":
			var friendName string
			json.Unmarshal([]byte(msg.Body), &friendName)
			result := serve.AddFriend(friendName, thisUser)         //friendName就是发起加好友的人
			serve.AddFriend(thisUser.Name, origin.OnlineUsers[friendName]) //todo:将自己也加入对方好友列表，好友添加是一个互相的过程
			conn.Write([]byte(result))                        //"添加好友失败"或"添加好友成功"
			//找到当初发起加好友者的连接，向其发送提示消息
			origin.OnlineUsers[friendName].Conn.Write([]byte("ha@申请通过，快和新好友聊天吧"))
		case "no":
			var friendName string
			json.Unmarshal([]byte(msg.Body), &friendName)
			//找到当初发起加好友者的连接，向其发送提示消息
			origin.OnlineUsers[friendName].Conn.Write([]byte("ha@申请未通过"))
		case "admit":
			var admitWhat string //strSlice[0]+"@"+strSlice[2] 即 加群申请人@群名
			json.Unmarshal([]byte(msg.Body), &admitWhat)
			splitResult := strings.Split(admitWhat, "@")
			applyer := splitResult[0]
			appliedGroup := splitResult[1]
			////群主将该用户添加进服务端维护的OnlineGroup群成员中
			//OnlineGroup[appliedGroup].Members = append(OnlineGroup[appliedGroup].Members,
			//											&common.OnlineUser{  //此时尚未给字段location赋值
			//												Name:applyer,
			//												Conn: OnlineUsers[applyer].Conn,//这里还是需要在OnlineUsers中找,导致速度还是较慢
			//											})
			//OnlineUsers[applyer].Conn.Write([]byte("加群申请通过，你已是群名为："+appliedGroup+"：的群成员"))

			//更新数据库grouptable中该群的群成员信息,同时更新该用户数据表中所加全部群的信息
			result := dao.AddGroupMembers(appliedGroup, applyer)
			if result {
				fmt.Println("新群成员存入群数据库中")
			} else {
				fmt.Println("新群成员存入群数据库中失败")
			}
			origin.OnlineUsers[applyer].Conn.Write([]byte(applyer + "@你申请加群被允许@" + appliedGroup))
		case "refuse":
			var refuseWhat string //strSlice[0]+"@"+strSlice[2] 即 申请加群人@群名
			json.Unmarshal([]byte(msg.Body), &refuseWhat)
			splitResult := strings.Split(refuseWhat, "@")
			applyer := splitResult[0]
			appliedGroup := splitResult[1]
			origin.OnlineUsers[applyer].Conn.Write([]byte("申请加入群：" + appliedGroup + ":未通过"))
		case "friendsList":
			fmt.Println(thisUser.Name)
			result := dao.ShowFirendsList(thisUser.Name)
			conn.Write(result)
		case "groupList":
			result := dao.ShowGroupsList(thisUser.Name) //返回该用户所在全部群的列表
			conn.Write(result)
		case "intoThisGroup":
			var intoThisGroup string
			json.Unmarshal([]byte(msg.Body), &intoThisGroup)
			_, ok := origin.OnlineGroup[intoThisGroup]
			if !ok { //这个群不存在
				conn.Write([]byte("ha@该群不存在，请正确输入群列表中的群名"))
				result := dao.ShowGroupsList(thisUser.Name) //返回该用户所在全部群的列表,让客户端再次选择进哪个群
				conn.Write(result)
			} else { //客户端想进入该群进行群聊
				//todo:标记一下此时用户已进入特定群，实现不接收其他群的消息(但接收私聊消息)的功能
				thisUser.Location = intoThisGroup //此时该用户所在的所有群里的该用户的该字段也更改了，因为thisUser是指针类型
				conn.Write([]byte(intoThisGroup + "@请开始群聊"))
			}
		case "groupmessage":
			var groupSentence_groupname string
			json.Unmarshal([]byte(msg.Body), &groupSentence_groupname)
			gg := strings.Split(groupSentence_groupname, "@")
			groupSentence := gg[0]
			groupname := gg[1]
//todo:将客户端发送来的群聊消息广播转发送给群成员
			serve.TransferGroupMessage(thisUser, groupSentence, groupname)
		case "withdrawGroup":
			//var groupName string
			//json.Unmarshal([]byte(msg.Body),&groupName)
			thisUser.Location = "" //将之前所在群的标记去掉，这样就不会再接收到该群的消息
		//todo:客户端发来聊天消息
		case "chatmessage":
			var sentence string
			json.Unmarshal([]byte(msg.Body), &sentence)
			//todo:将客户端发送来的聊天消息转发送给特定的聊天对象
			serve.TransferMessage(sentence, thisUser)
		case "setupGroup":
			var groupName string
			json.Unmarshal([]byte(msg.Body), &groupName)
			newGroup := common.NewGroup(groupName, thisUser.Name, conn)
			newGroup.Members = append(newGroup.Members, thisUser) //这样就可统一接受thisUser.Location = intoThisGroup的更改了
			result := dao.SetUpGroup(newGroup)
			if result == "ha@新建群成功" {
				//todo:将新群添加进服务器维护的所有在线群列表中
				origin.OnlineGroup[newGroup.GroupName] = newGroup
				fmt.Println("向OnlineGroup上架新群成功")
			}
			conn.Write([]byte(result))
		case "addGroup":
			var groupName string
			json.Unmarshal([]byte(msg.Body), &groupName)
			goalGroup, ok := origin.OnlineGroup[groupName]
			if ok { //这个群是存在的
				conn.Write([]byte("加群申请等待通过,你可先进行其他操作"))
				//todo:开启一个背后的小协程，默默地等待被申请加好友的用户上线，告知有人想加群
				go serve.Wait2conn(goalGroup.Owner.UserName, 2*time.Second, thisUser.Name+"@申请加群@"+groupName)
			} else {
				conn.Write([]byte("你要找的群不存在"))
			}
		case "allOnlineUsers":
			var ss string
			for k, _ := range origin.OnlineUsers {
				ss += k + "\t"
			}
			thisUser.Conn.Write([]byte("服务器全部在线用户有" + ss))
		}
	}
}
