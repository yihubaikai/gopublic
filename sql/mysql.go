package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

//初始化连接
func M_init() {
	//beego.Debug("m_init")
	//orm.RegisterDataBase("default", "mysql", "server02:root@123...@tcp(127.0.0.1:3306)/square?charset=utf8", 50, 50)
	//orm.RegisterDataBase("default", "mysql", "root:Abcf8765D5@tcp(rm-bp14tynb2bq1gg55k2o.mysql.rds.aliyuncs.com:3306)/power-pay?charset=utf8", 200, 200)
	baseaddress := Get_config("baseaddress")
	baseport := Get_config("baseport")
	basename := Get_config("basename")
	sqluser := Get_config("sqluser")
	sqlpass := Get_config("sqlpass")

	//beego.Debug(sqluser + ":" + sqlpass + "@tcp(" + baseaddress + ":" + baseport + ")/" + basename + "?charset=utf8")
	orm.RegisterDataBase("default", "mysql", sqluser+":"+sqlpass+"@tcp("+baseaddress+":"+baseport+")/"+basename+"?charset=utf8", 200, 200)
	//orm.RegisterModel(new(User))
}

//切换数据库
func M_Using(dbname string) {
	o := orm.NewOrm()
	o.Using(dbname)
}

//获取有返回值的sql语句，比如select, shwo database
func Getrs(sqlcmd string) (err error, rs []orm.Params) {
	o := orm.NewOrm()
	_, err = o.Raw(sqlcmd).Values(&rs)
	return err, rs
}

//执行sql语句，比如： insert into， update ，deltete， create 等
func Dosql(sqlcmd string) {
	o := orm.NewOrm()
	o.Raw(sqlcmd).Exec()
}

//检查数据库中的表是否存在
//检查数据库中的表是否存在
func Check_Tab_Exists(tabname string) bool {

	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		err := recover()
		if err != nil {
			beego.Debug(tabname, "查询异常，表", tabname, "不存在")
		}
	}()

	sqlcmd := "select * from `" + tabname + "` limit 0, 1"
	err1, rs := Getrs(sqlcmd)
	if err1 != nil {
		beego.Debug("查询失败 视为表", tabname, "不存在")
		return false //查询失败 视为表不存在
	}
	if len(rs) > 0 {
		beego.Debug("表", tabname, "存在")
		return true
	}
	beego.Debug(tabname, "表补存在")
	return false
}

//显示当前连接中所有的库取值字段：Database
func ShowDatabases() (err error, rs []orm.Params) {
	//var rs []orm.Params
	err, rs = Getrs("show databases")
	beego.Debug("ShowDatabases 取值字段：Database", rs)
	/*if err != nil {
		beego.Debug(err)
	} else {
		for num, row := range rs {
			beego.Debug(num, row["Database"])
		}
	}*/
	return err, rs
}

//显示当前库下所有的数据表 取值字段：Tables_in_square
func ShowTables(dbname string) (err error, rs []orm.Params) {
	M_Using(dbname)
	sqlcmd := "show tables;"
	err, rs = Getrs(sqlcmd)
	beego.Debug("ShowTables 取值字段：Tables_in_square", rs)
	return err, rs
}

//显示当前表下所有的字段  取值字段：Field
//调用这个函数之前必须调用showtables防止未切换数据库
func ShowField(tabname string) (err error, rs []orm.Params) {
	sqlcmd := "desc " + tabname + ";"
	err, rs = Getrs(sqlcmd)
	beego.Debug("ShowField 取值字段：Field", rs)
	return err, rs
}

func Get_config(key string) string {
	conf, err := config.NewConfig("ini", "conf/config.ini")
	if err != nil {
		beego.Debug("new config failed, err:", err)
		return ""
	}
	ret := conf.String("config::" + key)
	return ret
}
