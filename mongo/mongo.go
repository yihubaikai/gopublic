package mongo

import (
	"fmt"
	"github.com/yihubaikai/gopublic"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	//"log"
)

//单条数据节点
type DataNode struct {
	Nick      string `json:"nick"`
	Class     string `json:"class"`
	Flag      int    `json:"flag"`
	Starttime string `json:"starttime"`
}

type Person struct {
	NAME  string
	PHONE string
}
type Men struct {
	Persons []Person
}

var session *mgo.Session
var database *mgo.Database = nil

// get mongodb
//db: data
//host： 127.0.0.1
//
func GetDB() *mgo.Database {
	if database == nil {
		session, err := mgo.Dial("127.0.0.1:27017")
		if err != nil {
			panic(err)
		}
		//defer session.close()
		session.SetMode(mgo.Monotonic, true)
		database = session.DB("data")
		return database
	} else {
		return database
	}

	/*
		var err error
			dialInfo := &mgo.DialInfo{
		        Addrs:     []string{config.Hosts},
		        Direct:    false,
		        Timeout:   time.Second * 1,
		        PoolLimit: 4096, // Session.SetPoolLimit    }
		    //创建一个维护套接字池的session
		    session, err = mgo.DialWithInfo(dialInfo)

		    if err != nil {
		        log.Println(err.Error())
		    }
		    session.SetMode(mgo.Monotonic, true)
		    //使用指定数据库
		    database = session.DB(config.Database)
	*/
}

//插入记录
func Insert(db *mgo.Database, item, nick, class string) bool {
	c := db.C(item)
	//defer c.Close()
	//result := DataNode{}
	err, _ := c.Upsert(bson.M{"nick": nick}, &DataNode{Nick: nick, Class: class, Flag: 0, Starttime: hPub.Gettime()})
	if err != nil {
		fmt.Println("添加记录成功:", item, nick, class, hPub.Gettime())
		return true
	} else {
		fmt.Println("更新记录成功:", item, nick, class, hPub.Gettime())
		return true
	}
	//fmt.Println(rs)
	return true
}

//获取记录
func GetOne(db *mgo.Database, item string, find, update int) (nick, class string) {
	nick = ""
	class = ""
	c := db.C(item)
	result := DataNode{}
	err := c.Find(bson.M{"flag": find}).One(&result)
	if err == nil {
		fmt.Println("nick", result.Nick, "class", result.Class)
		nick = result.Nick
		class = result.Class

		//获取记录了之后更新一下记录
		//err = c.Update(bson.M{"nick": nick}, bson.M{"$set": bson.M{"flag": 2, "starttime": hPub.Gettime()}})
		changeInfo, err := c.UpdateAll(bson.M{"nick": nick}, bson.M{"$set": bson.M{"flag": update}})
		if err == nil {
			fmt.Println("获取记录：更新记录成功:", changeInfo)
		} else {
			fmt.Println("获取记录：更新记录失败:", changeInfo)
		}
	} else {
		fmt.Println("获取记录失败")
	}
	return nick, class
}

//--------------------------------------------------------------------
//插入单条数据
func Insert2(db *mgo.Database) {
	//db := GetDB()

	c := db.C("user")
	type User struct {
		Name string "bson:`name`"
		Age  int    "bson:`age`"
	}

	err := c.Insert(&User{Name: "Tom", Age: 20})
	if err != nil {
		panic(err)
	}
	fmt.Println(err)
}
func Insert3(db *mgo.Database, item, nick, class string) bool {
	c := db.C(item)
	result := DataNode{}
	err := c.Find(bson.M{"nick": nick}).One(&result)
	if err != nil {
		//err := c.Upsert(&DataNode{Nick: nick, Class: class, Flag: 0, Starttime: hPub.Gettime()})
		err, _ := c.Upsert(bson.M{"nick": nick}, &DataNode{Nick: nick, Class: class, Flag: 0, Starttime: hPub.Gettime()})
		if err != nil {
			fmt.Println("添加记录成功:", item, nick, class, hPub.Gettime())
			return true
		} else {
			fmt.Println("更新记录成功:", item, nick, class, hPub.Gettime())
			return true
		}
	} else {
		fmt.Println("可能存在")
		return true
	}
}

//一次插入多条记录
func insertMuti() {
	db := GetDB()

	c := db.C("user")
	type User struct {
		Name string "bson:`name`"
		Age  int    "bson:`age`"
	}

	err := c.Insert(&User{Name: "Tom", Age: 20}, &User{Name: "Anny", Age: 28})
	if err != nil {
		panic(err)
	}
	fmt.Println(err)
}

//插入数组格式数据
func insertArray() {
	db := GetDB()
	c := db.C("user")

	type User struct {
		Name   string   "bson:`name`"
		Age    int      "bson:`age`"
		Groups []string "bson:`groups`"
	}

	err := c.Insert(&User{
		Name:   "Tom",
		Age:    20,
		Groups: []string{"news", "sports"},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(err)
}

//插入嵌套数据
func insertNesting() {
	db := GetDB()

	c := db.C("user")

	type Toy struct {
		Name string "bson:`name`"
	}
	type User struct {
		Name string "bson:`name`"
		Age  int    "bson:`age`"
		Toys []Toy
	}

	err := c.Insert(&User{
		Name: "Tom",
		Age:  20,
		Toys: []Toy{{Name: "dog"}},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(err)
}

//插入map格式的数据
func insertMap() {
	db := GetDB()
	c := db.C("user")

	user := map[string]interface{}{
		"name":   "Tom",
		"age":    20,
		"groups": []string{"news", "sports"},
		"toys": []map[string]interface{}{
			{
				"name": "dog",
			},
		},
	}

	err := c.Insert(&user)
	if err != nil {
		panic(err)
	}
	fmt.Println(err)
}

//插入关联其它集合ObjectId的数据
//要使用bson.ObjectIdHex函数对字符串进行转化，bson.ObjectIdHex函数原型
func insertObjectId() {
	db := GetDB()
	c := db.C("user")

	user := map[string]interface{}{
		"name":     "Tom",
		"age":      20,
		"group_id": bson.ObjectIdHex("540046baae59489413bd7759"),
	}

	err := c.Insert(&user)
	if err != nil {
		panic(err)
	}
	fmt.Println(err)
}

func QueryOne() {

	db := GetDB()
	c := db.C("user")
	//*****查询单条数据*******
	result := Person{}
	err := c.Find(bson.M{"NAME": "456"}).One(&result)
	if err != nil {
		panic(err)
	}
	fmt.Println("Phone:", result.NAME, result.PHONE)

}

func QueryMutil() {
	db := GetDB()
	c := db.C("user")
	//*****查询多条数据*******
	result := Person{}
	var personAll Men //存放结果
	iter := c.Find(nil).Iter()
	for iter.Next(&result) {
		fmt.Printf("Result: %v\n", result.NAME)
		personAll.Persons = append(personAll.Persons, result)
	}
}
