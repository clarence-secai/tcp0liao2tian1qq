package serve

import (
	"fmt"
	"go7tcp0liao2tian1app/20200619tcp0qq0ruan3jian4/common"
	"go7tcp0liao2tian1app/20200619tcp0qq0ruan3jian4/server/dao"
	"go7tcp0liao2tian1app/20200619tcp0qq0ruan3jian4/server/origin"
	"strings"
	"time"
)

func TransferGroupMessage(thisUser *common.OnlineUser, sentence string, groupname string) {
	for _, member := range origin.OnlineGroup[groupname].Members {
		//单个用户一上线就需在他所加的所有群里都去上线，否则后面收不到各个群的消息
		//还有要避免在一个群聊中收到其他群的群聊消息，因为无法区分到底是哪个群的群聊消息
		fmt.Println("member.Location=", member.Location)
		if member.Location != thisUser.Location {
			continue //没进入该群聊天的在线用户也接收不到群消息
		}
		member.Conn.Write([]byte(thisUser.Name + ":" + sentence + "#all" + "||群聊来自群：" + groupname))
	}
}

func TransferMessage(sentence string, thisuser *common.OnlineUser) {
	sentenceSlice := strings.Split(sentence, "#")
	//查看聊天消息的目的用户是否在线
	goalUser, ok := origin.OnlineUsers[sentenceSlice[1]] //todo:私聊每次都需要去所有在线用户OnlineUsers中去找，实际生产中应该效率较低，须有其他解决方案
	if !ok {                                      //该好友不在线，接着判断是该好友没上线还是根本就不存在这个用户
		downLineGoalUser := common.User{UserName: sentenceSlice[1]}
		//从数据库全部QQ用户数据表中查询是否有该用户，其实一般不需要，因实际qq软件界面不会让客户有这种操作可能
		flag := dao.Check(&downLineGoalUser)
		if flag { //有这个QQ用户但当前不在线
			thisuser.Conn.Write([]byte("该好友不在线，请勿超过3条留言！！！")) //todo:也可通过写入历史消息文件，在目标用户上线后再发送
			origin.LeaveWordTime += 1
			if origin.LeaveWordTime == 4 {
				thisuser.Conn.Write([]byte("留言超3条，该条留言不会留言给好友"))
				return
			}
			go func() { //存在留言消息不一定按先后顺序在用户上线后予以显示的问题
				timeStap := time.Now().Format("01-02 15:04:05")
				for {
					goalUser, ok := origin.OnlineUsers[sentenceSlice[1]]
					if !ok {
						time.Sleep(time.Minute * 1) //当然也可改成其他时间
					} else {
						goalUser.Conn.Write([]byte(thisuser.Name + ":" + sentence + "(--留言消息--)" + timeStap))
						break
					}
				}
			}()
		} else { //不存在这个QQ用户
			thisuser.Conn.Write([]byte("不存在该用户，消息未发送"))
		}
	} else {
		//好友在线，则直接将聊天消息转发给该好友
		goalUser.Conn.Write([]byte(thisuser.Name + ":" + sentence))
	}
}

func Wait2conn(key string, waitDuration time.Duration, content string) {
	for { //todo:间歇循环查看,等待该用户上线，告知有人想加他好友或加群
		user, ok := origin.OnlineUsers[key]
		if !ok {
			time.Sleep(waitDuration)
			continue
		} else {
			user.Conn.Write([]byte(content))
			break
		}
	}

}

func AddFriend(newFriend string, thisUser *common.OnlineUser) string {
	return dao.AddFriend(newFriend, thisUser.Name)
}

func DealExploreFriend(body string) bool {
	var user common.User
	user.UserName = body
	return dao.Check(&user)
}



