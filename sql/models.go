package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/orm"
	"github.com/yihubaikai/gopublic"
	"log"
)

//
//测试列表
func U_Query() {
	/*for i := 0; i < 1; i++ {
		o := orm.NewOrm()
		fmt.Println(o)
		value := reflect.ValueOf(o)
		typ := value.Type()
		for i := 0; i < value.NumMethod(); i++ {
			fmt.Println(fmt.Sprintf("method[%d]%s and type is %v", i, typ.Method(i).Name, typ.Method(i).Type))
		}
		fmt.Println(i)
		time.Sleep(100)

		var maps []orm.Params
		num, err := o.Raw("SELECT * FROM card").Values(&maps)
		beego.Debug(num, err)
		for _, row := range maps {
			beego.Debug(row["id"], ":", row["account"])
		}

		o.Raw("SELECT * FROM bill ").Values(&maps)
		for _, row := range maps {
			beego.Debug(row["id"], ":", row["order_no"])
		}
	}
	*/
}

//收款账号
func U_Card_GetInfo(username string) (err error, rs []orm.Params) {
	sqlcmd := " select card.id,card.account,card.yue,card.yuetime,card.warningmoney,card.status,card.time,card.url,card_online.online_time  from card left join card_online  on card_online.card_id=card.id"
	if len(username) > 0 {
		sqlcmd = "select id,account,yue,yuetime,warningmoney,status,time,url from card where account='" + username + "'"
	}
	return Getrs(sqlcmd)
}

//增加在线收款账号
func U_Card_Online_add(card_id, account string) {
	sqlcmd := "select id from card_online where card_name='" + account + "';"
	//fmt.Println(sqlcmd)
	err, rs := Getrs(sqlcmd)
	if err == nil && len(rs) == 0 {
		sqlcmd = "insert into card_online (card_id,card_name,state) values('" + card_id + "','" + account + "','0');"
		//fmt.Println(sqlcmd)
		Dosql(sqlcmd)
	}
}

//增加收款账号
func U_Card_Add(account string) string {
	sqlcmd := "select id,account from card where account='" + account + "';"
	//fmt.Println(sqlcmd)
	err, rs := Getrs(sqlcmd)
	if err != nil {
		log.Fatal(err)
		return "查询数据库失败"
	}

	if len(rs) > 0 { //有记录
		id := fmt.Sprintf("%s", rs[0]["id"])
		U_Card_Online_add(id, account)
		return "exists"
	}
	sqlcmd = "insert into card (account,yue,yuetime,warningmoney,warningtime,time,status,user_name) values('" + account + "','0','" + hPub.Gettime() + "','0','" + hPub.Gettime() + "','" + hPub.Gettime() + "','0','" + account + "');"
	//fmt.Println(sqlcmd)
	Dosql(sqlcmd)
	ret := U_Card_Add(account)
	if ret == "exists" {
		return "succ"
	}
	return "error"
}

//收款账号在线
func U_Card_Online(account string) string {
	sqlcmd := "update card set time='" + hPub.Gettime() + "' where account='" + account + "';"
	//fmt.Println(sqlcmd)
	Dosql(sqlcmd)
	return sqlcmd
}

//增加客户
func U_User_Add(Username, Password, UserIdentity, UserSecureKey, Bili string) string {
	if len(Username) == 0 || len(Password) == 0 || len(UserIdentity) == 0 || len(UserSecureKey) == 0 || len(Bili) == 0 {
		return "参数不正确,有参数未传递值"
	}
	sqlcmd := "select id from user where username='" + Username + "';"
	err, rs := Getrs(sqlcmd)
	if err != nil {
		log.Fatal(err)
		return "查询数据库失败：" + Username
	}

	if len(rs) > 0 {
		return Username + ": 用户已经存在"
	}

	sqlcmd1 := "insert into user (username, userpass, useridentity,md5_key,nick,para_id,lilv ) values('" + Username + "','" + Password + "','" + UserIdentity + "','" + UserSecureKey + "','" + Username + "','" + Username + "','" + Bili + "');"
	Dosql(sqlcmd1)
	err, rs = Getrs(sqlcmd)
	if err == nil && len(rs) > 0 {
		return "添加用户:" + Username + "成功"
	}
	return "添加用户失败"
}

//获取客户
func U_User_GetInfo(para_id string) (err error, rs []orm.Params) {
	sqlcmd := "select * from user;"
	if len(para_id) > 0 {
		sqlcmd = "select * from user where username='" + para_id + "';"
	}
	return Getrs(sqlcmd)
}

//删除user
func U_User_DelUser(id string) string {
	sqlcmd := "delete from user where id='" + id + "';"
	//fmt.Println(sqlcmd)
	Dosql(sqlcmd)
	return "succ"
}

//

//创建支付表
func CreteTab_payed(tabname string) {
	sqlcmd := "CREATE TABLE `" + tabname + "` (`Id` int(11) NOT NULL AUTO_INCREMENT,   `order_no` varchar(50) DEFAULT NULL COMMENT '订单号', `money` varchar(20) DEFAULT NULL COMMENT '金钱',  `tradeno` varchar(50) DEFAULT NULL COMMENT '流水号', `account_id` varchar(20) DEFAULT NULL COMMENT '收款人ID',  `para_id` varchar(20) DEFAULT NULL COMMENT '客户ID', `start_time` datetime DEFAULT NULL COMMENT '开始时间', `end_time` datetime DEFAULT NULL COMMENT '结束时间', PRIMARY KEY (`Id`),UNIQUE KEY `order_no` (`order_no`), KEY `start_time` (`start_time`) ) ENGINE=MyISAM DEFAULT CHARSET=utf8;"
	Dosql(sqlcmd)
}

//创建申请表
func CreateTab_succ(tabname string) {
	sqlcmd := "CREATE TABLE `" + tabname + "` (`Id` int(11) NOT NULL AUTO_INCREMENT,`order_no` varchar(50) DEFAULT NULL COMMENT '订单号',`start_time` datetime DEFAULT NULL COMMENT '开始时间',`para_id` varchar(20) DEFAULT NULL COMMENT '客户ID',`channel_id` int(11) DEFAULT NULL COMMENT '频道ID',`account_id` int(11) DEFAULT NULL COMMENT '频道ID',PRIMARY KEY (`Id`),UNIQUE KEY `order_no` (`order_no`),KEY `start_time` (`start_time`)) ENGINE=MyISAM DEFAULT CHARSET=utf8;"
	Dosql(sqlcmd)
}

//创建统计表
func CreateTab_tongji(tabname string) {
	sqlcmd := "CREATE TABLE `" + tabname + "` (`Id` int(11) NOT NULL AUTO_INCREMENT,`day` date DEFAULT '0000-00-00' COMMENT '日期',`alloc` int(11) DEFAULT '0' COMMENT '所有申请的数量',`succ` int(11) DEFAULT '0' COMMENT '成功申请的数量',`para_id` varchar(10) DEFAULT '0' COMMENT '客户ID',`payed` int(11) DEFAULT '0' COMMENT '已经支付的数量',`account` varchar(50) DEFAULT '0' COMMENT '通道',`account_id` varchar(10) DEFAULT NULL,`money` bigint(11) DEFAULT '0' COMMENT '金钱',PRIMARY KEY (`Id`),KEY `Id` (`Id`,`day`,`para_id`,`account`)) ENGINE=MyISAM AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;"
	Dosql(sqlcmd)
}

//创建账户表
func CreateTab_card(tabname string) {
	sqlcmd := "CREATE TABLE `" + tabname + "` (`Id` int(11) NOT NULL AUTO_INCREMENT,`card_id` int(11) DEFAULT NULL COMMENT '卡的ID',`account` varchar(50) DEFAULT NULL COMMENT '账号',PRIMARY KEY (`Id`)) ENGINE=MyISAM DEFAULT CHARSET=utf8;"
	Dosql(sqlcmd)
}

//创建提现表
func CreateTab_tixian(tabname string) {
	sqlcmd := "CREATE TABLE `" + tabname + "` (`Id` int(11) NOT NULL AUTO_INCREMENT,`day` varchar(20) DEFAULT NULL COMMENT '日期——天',`card_id` int(11) DEFAULT NULL COMMENT '账号id',`money` varchar(50) DEFAULT NULL COMMENT '金钱',`tixian` varchar(50) DEFAULT NULL,PRIMARY KEY (`Id`)) ENGINE=MyISAM AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;"
	Dosql(sqlcmd)
}

//检查表是否存在
func CheckTab() {
	if !Check_Tab_Exists("_payed") {
		CreteTab_payed("_payed")
	}
	if !Check_Tab_Exists("_succ") {
		CreateTab_succ("_succ")
	}
	if !Check_Tab_Exists("_tongji") {
		CreateTab_tongji("_tongji")
	}
	if !Check_Tab_Exists("_card") {
		CreateTab_card("_card")
	}
	if !Check_Tab_Exists("_tixian") {
		CreateTab_tixian("_tixian")
	}
}

//同步card表
func sync_card() {
	sqlcmd := "select id, account from `card`"
	err, rs := Getrs(sqlcmd)
	if err == nil {
		for _, row := range rs {
			id := fmt.Sprintf("%s", row["id"])
			ac := fmt.Sprintf("%s", row["account"])
			beego.Debug(id, ac)
			if len(id) > 0 && len(ac) > 0 {
				sqlcmd1 := "select * from _card where account='" + ac + "' and card_id='" + id + "';"
				err1, rs1 := Getrs(sqlcmd1)
				if err1 == nil && len(rs1) == 0 {
					sqlcmd2 := "insert into `_card` (card_id, account) values('" + id + "','" + ac + "');"
					beego.Debug(sqlcmd1, sqlcmd2)
					Dosql(sqlcmd2)
				}
			}
		}
	}
}

//获取上次支付时间
func Get_payed_lastTime() string {
	lasttime := Get_config("lasttime") //读取配置
	if len(lasttime) == 0 {
		lasttime = "2019-01-01 00:00:00"
	}
	lasttime = "2019-01-01 00:00:00"
	/*sqlcmd := "select * from `_payed` order by end_time desc limit 0, 1"
	err, rs := Getrs(sqlcmd)
	if err == nil && len(rs) > 0 {
		for _, row := range rs { //读取最近时间
			lasttime = fmt.Sprintf("%s", row["end_time"])
			break
		}
	}*/
	return lasttime
}

//获取上次申请时间
func Get_succ_lastTime() string {
	lasttime := Get_config("lasttime") //读取配置
	if len(lasttime) == 0 {
		lasttime = "2019-01-01 00:00:00"
	}
	sqlcmd := "select * from `_succ` order by start_time desc limit 0, 1"
	err, rs := Getrs(sqlcmd)
	if err == nil && len(rs) > 0 {
		for _, row := range rs { //读取最近时间
			lasttime = fmt.Sprintf("%s", row["start_time"])
			break
		}
	}
	return lasttime
}

//获取

//同步支付订单
func sync_payed(daystr string) {
	tmp := make(map[string]string)
	sqlcmd := "select bill.order_no,round(bill.money,2) as money,bill.tradeno,bill.start_time,bill.end_time,card.id as account,bill.type as para_id from `bill` left join card on bill.qrcode_id=card.id where unix_timestamp(bill.end_time)>=unix_timestamp('" + daystr + "') and end_time is not null order by bill.end_time asc"
	err, rs := Getrs(sqlcmd)
	if err == nil && len(rs) > 0 {
		wheretj := ""
		for _, row := range rs {
			order_no := fmt.Sprintf("%s", row["order_no"])
			money := fmt.Sprintf("%s", row["money"])
			tradeno := fmt.Sprintf("%s", row["tradeno"])
			account_id := fmt.Sprintf("%s", row["account"])
			para_id := fmt.Sprintf("%s", row["para_id"])
			starttime := fmt.Sprintf("%s", row["start_time"])
			endtime := fmt.Sprintf("%s", row["end_time"])
			sqlcmd1 := "insert into `_payed` (order_no,money,tradeno,account_id,para_id,start_time,end_time) values('" + order_no + "','" + money + "','" + tradeno + "','" + account_id + "','" + para_id + "','" + starttime + "','" + endtime + "');"
			Dosql(sqlcmd1)
			if len(tmp[starttime[0:10]]) > 0 {
				continue
			}
			tmp[starttime[0:10]] = "1"
			if len(wheretj) == 0 {
				wheretj = " start_time like '" + starttime[0:10] + "%'"
			} else {
				wheretj = wheretj + " or start_time like '" + starttime[0:10] + "%'"
			}
		}

		//更新统计表
		sqlcmd2 := "select DATE_FORMAT(start_time, '%Y-%m-%d' ) as day,count(id) as num, 'payed' as item, sum(money) as money, account_id,para_id from `_payed` where " + wheretj + " group by DATE_FORMAT(start_time, '%Y-%m-%d' ),account_id,para_id"
		beego.Debug(sqlcmd2)
		err, rs = Getrs(sqlcmd2)
		if err == nil && len(rs) > 0 {
			//查询一记录是否存在
			for _, rows := range rs {
				money := fmt.Sprintf("%s", rows["money"])
				day := fmt.Sprintf("%s", rows["day"])
				payed := fmt.Sprintf("%s", rows["num"])
				item := fmt.Sprintf("%s", rows["item"])
				para_id := fmt.Sprintf("%s", rows["para_id"])
				account_id := fmt.Sprintf("%s", rows["account_id"])
				beego.Debug(money, day, payed, item, para_id, account_id)
				sqlcmd3 := "select id,payed,money from `_tongji` where day='" + day + "' and para_id='" + para_id + "' and account_id='" + account_id + "';"
				err3, rs3 := Getrs(sqlcmd3)
				if err3 == nil {
					if len(rs3) > 0 {
						for _, row3 := range rs3 {
							_id := fmt.Sprintf("%s", row3["id"])
							_payed := fmt.Sprintf("%s", row3["payed"])
							_money := fmt.Sprintf("%s", row3["money"])
							if _payed != payed || _money != money {
								sqlx := "update `_tongji` set payed='" + payed + "', money='" + money + "' where id ='" + _id + "';"
								beego.Debug(sqlx)
								Dosql(sqlx)
							}
						}
					} else {
						sqlx := "insert into `_tongji` (day,payed,money,para_id,account_id) values('" + day + "','" + payed + "','" + money + "','" + para_id + "','" + account_id + "');"
						beego.Debug(sqlx)
						Dosql(sqlx)
					}
				}

			}
		}

	}
}

//同步申请订单
func sync_succ(daystr string) {
	tmp := make(map[string]string) //去重复
	sqlcmd := "select bill.order_no,bill.start_time,card.id as account_id,bill.type as para_id from `bill` left join card on bill.qrcode_id=card.id where  unix_timestamp(bill.start_time)>=unix_timestamp('" + daystr + "') and bill.qrcode_id is not null order by bill.end_time asc"
	beego.Debug(sqlcmd)
	err, rs := Getrs(sqlcmd)
	if err == nil && len(rs) > 0 {
		wheretj := ""
		for _, row := range rs {
			order_no := fmt.Sprintf("%s", row["order_no"])
			para_id := fmt.Sprintf("%s", row["para_id"])
			starttime := fmt.Sprintf("%s", row["start_time"])
			account_id := fmt.Sprintf("%s", row["account_id"])
			sqlcmd1 := "insert into `_succ` (order_no,para_id,start_time,account_id) values('" + order_no + "','" + para_id + "','" + starttime + "','" + account_id + "');"
			//beego.Debug(sqlcmd1)
			Dosql(sqlcmd1)
			if len(tmp[starttime[0:10]]) > 0 {
				continue
			}
			tmp[starttime[0:10]] = "1"
			if len(wheretj) == 0 {
				wheretj = " start_time like '" + starttime[0:10] + "%'"
			} else {
				wheretj = wheretj + " or start_time like '" + starttime[0:10] + "%'"
			}
		}
		beego.Debug(wheretj)
		sqlcmd2 := "select DATE_FORMAT(start_time, '%Y-%m-%d' ) as day,count(id) as num, 'succ' as item, para_id, account_id from `_succ`  where " + wheretj + "  group by DATE_FORMAT(start_time, '%Y-%m-%d' ),para_id, account_id"
		beego.Debug(sqlcmd2)
		err, rs = Getrs(sqlcmd2)
		if err == nil && len(rs) > 0 {
			//查询一记录是否存在
			for _, rows := range rs {
				day := fmt.Sprintf("%s", rows["day"])
				succ := fmt.Sprintf("%s", rows["num"])
				para_id := fmt.Sprintf("%s", rows["para_id"])
				account_id := fmt.Sprintf("%s", rows["account_id"])
				sqlcmd3 := "select id,succ from `_tongji` where day='" + day + "' and para_id='" + para_id + "' and account_id='" + account_id + "'"
				err3, rs3 := Getrs(sqlcmd3)
				if err3 == nil {
					if len(rs3) > 0 {
						for _, row3 := range rs3 {
							_id := fmt.Sprintf("%s", row3["id"])
							_succ := fmt.Sprintf("%s", row3["succ"])
							if _succ != succ {
								sqlx := "update `_tongji` set succ='" + succ + "'  where id ='" + _id + "';"
								beego.Debug(sqlx)
								Dosql(sqlx)
							}
						}
					} else {
						sqlx := "insert into `_tongji` (day,payed,para_id,account_id) values('" + day + "','" + succ + "','" + para_id + "','" + account_id + "');"
						beego.Debug(sqlx)
						Dosql(sqlx)
					}
				}
			}
		}
	}
}

//同步数据库
func Syncbase() {
	//同步用户
	sync_card()

	//同步申请订单
	lasttime := Get_payed_lastTime()
	sync_payed(lasttime)

	//同步支付订单
	lasttime = Get_succ_lastTime()
	sync_succ(lasttime)
}
