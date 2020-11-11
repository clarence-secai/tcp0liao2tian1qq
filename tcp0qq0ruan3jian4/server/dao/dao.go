package dao

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"go7tcp0liao2tian1app/20200619tcp0qq0ruan3jian4/common"
	"go7tcp0liao2tian1app/20200619tcp0qq0ruan3jian4/server/utils"
)

var mydb *sql.DB
var err1 error

//运行程序前初始化得到打开数据库的mydb
func init() {
	mydb, err1 = sql.Open("mysql", "root:413188ok@tcp(localhost:3306)/qq")
	utils.ErrContinue("11行", err1)
}

func Regist(username string, password int64, friends []string, groups []string) string {
	var user = common.User{UserName: username, PassWord: password}
	flag := Check(&user)
	if flag {
		return "ha@用户已存在，无需注册，请登录" //已存在，无需注册
	} else {
		friendsByte, err11 := json.Marshal(friends)
		groupsByte, err11 := json.Marshal(groups)
		if err11 != nil {
			fmt.Println("25行序列化好友切片出错：", err11)
			return "ha@注册失败"
		}
		sql := "insert into usertable (username,password,friends,groups) values(?,?,?,?)"
		_, err := mydb.Exec(sql, username, password, string(friendsByte), string(groupsByte))
		if err != nil {
			fmt.Println("dao包23行插入数据出错")
			return "ha@注册失败"
		} else {
			return "ha@注册成功，请登录" //成功插入一条注册数据
		}
	}
}

func Login(username string, password int64) string {
	sql := "select id from usertable where username=? and password=?"
	row := mydb.QueryRow(sql, username, password)
	var id int
	row.Scan(&id)
	if id > 0 {
		return username + "@" + "登录成功"
	} else {
		return "ha@用户名或密码错误"
	}

}

//检查数据库中是否有该用户
func Check(user *common.User) bool {
	sql := "select id from usertable where username=?"
	row := mydb.QueryRow(sql, user.UserName)
	var id int
	row.Scan(&id)
	if id > 0 {
		return true
	} else {
		return false
	}

}


func AddFriend(newFriend string, thisUserName string) string {
	sql := "select friends from usertable where username= ?"
	row := mydb.QueryRow(sql, thisUserName)
	var friendsStr string
	row.Scan(&friendsStr)
	var friendsSlice []string
	json.Unmarshal([]byte(friendsStr), &friendsSlice)
	friendsSlice = append(friendsSlice, newFriend)
	newFriends, err := json.Marshal(friendsSlice)
	if err != nil {
		fmt.Println("78行序列化更新好友切片出错", err)
		return "ha@添加好友失败"
	}
	sql = "update usertable set friends=? where username=?"
	mydb.Exec(sql, string(newFriends), thisUserName)
	return "ha@添加好友成功"
}

func ShowFirendsList(thisUserName string) []byte {
	sql := "select friends from usertable where username= ? "
	//fmt.Println("thisUserName=",thisUserName)
	row := mydb.QueryRow(sql, thisUserName)
	var str string
	err := row.Scan(&str)
	utils.ErrExit("89行", err)
	//fmt.Println("str=",str)
	return []byte(str + "@" + "好友列表")
	//var friendsSlice []string
	//json.Unmarshal([]byte(str),&friendsSlice)
	//return friendsSlice
}

func ShowGroupsList(thisUserName string) []byte {
	sql := "select groups from usertable where username=?"
	row := mydb.QueryRow(sql, thisUserName)
	var groupsStr string
	row.Scan(&groupsStr)
	return []byte(groupsStr + "@所在全部群的列表")
}

func SetUpGroup(newGroup *common.Group) string {
	var memberSlice []string
	for _, v := range newGroup.Members {
		memberSlice = append(memberSlice, v.Name) //只存群成员的名字，成员的conn、地址等因易变故不存
	}
	membersByte, err := json.Marshal(memberSlice)
	if err != nil {
		fmt.Println("102行序列化出错", err)
		return "ha@新建群失败"
	}
	//将新建的群存进数据库中
	sql := "insert into grouptable(groupname,owner,members)values(?,?,?)"
	_, err = mydb.Exec(sql, newGroup.GroupName, newGroup.Owner.UserName, string(membersByte))
	if err != nil {
		fmt.Println("108行插入新建群记录出错", err)
		return "ha@新建群失败"
	}
	//自己新建了群，在usertable中自己的所加群的groups字段也应更新加上自己为群主的这个群
	sql = "select groups from usertable where username=?"
	row := mydb.QueryRow(sql, newGroup.Owner.UserName)
	var groupsStr string
	row.Scan(&groupsStr)
	var groupsSlice []string
	json.Unmarshal([]byte(groupsStr), &groupsSlice)
	groupsSlice = append(groupsSlice, newGroup.GroupName)
	sql = "update usertable set groups=? where username=?"
	groupsByte, _ := json.Marshal(groupsSlice)
	mydb.Exec(sql, string(groupsByte), newGroup.Owner.UserName)
	return "ha@新建群成功"
}

func AddGroupMembers(appliedGroup string, applyer string) bool {
	sql := "select members from grouptable where groupname=?"
	row := mydb.QueryRow(sql, appliedGroup)
	var membersStr string
	row.Scan(&membersStr)
	var memberSlice []string
	json.Unmarshal([]byte(membersStr), &memberSlice)
	memberSlice = append(memberSlice, applyer)
	sql = "update grouptable set members=? where groupname=?"
	memberByte, err := json.Marshal(memberSlice)
	if err != nil {
		return false
	}
	mydb.Exec(sql, string(memberByte), appliedGroup)

	//todo:更新申请加群人在usertable数据表中的groups一栏
	sql = "select groups from usertable where username=?"
	row = mydb.QueryRow(sql, applyer)
	var groupsStr string
	row.Scan(&groupsStr)
	var groupsSlice []string
	json.Unmarshal([]byte(groupsStr), &groupsSlice)
	groupsSlice = append(groupsSlice, appliedGroup)
	sql = "update usertable set groups=? where username=?"
	groupsByte, err := json.Marshal(groupsSlice)
	if err != nil {
		return false
	}
	mydb.Exec(sql, string(groupsByte), applyer)
	return true
}
