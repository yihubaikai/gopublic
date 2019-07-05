package redisx

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
)

var g_Conn redis.Conn = nil

var Pool redis.Pool

func Redis_init(ip, port string) { //init 用于初始化一些参数，先于main执行
	Pool = redis.Pool{
		MaxIdle:     16,
		MaxActive:   32,
		IdleTimeout: 120,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", ip+":"+port)
		},
	}

	fmt.Println("---Redis_init---")
}

func GetDB() redis.Conn {
	if g_Conn == nil {
		Redis_init("127.0.0.1", "6379")
		g_Conn = Pool.Get()
		//defer g_Conn.Close()
	}
	return g_Conn
}

func Redis_AddNick(key, val string, iDB int) {
	//获取数据库
	conn := Pool.Get()

	//选择数据库
	conn.Do("SELECT", iDB)

	//执行设置指令
	conn.Do("SETNX", key, val)

}

func Redis_GetNick() {

}

/*
var g_Conn *Redis.Conn = nil

func getdb() Redis.Conn {
	if g_Conn == nil {
		redis.Redis_init("127.0.0.1", "6379")
		g_Conn = Pool.Get()
	}
	return g_Conn
}

func main1() {

	//conn.Do("SELECT", 1)
	conn := getdb()

	for i := 0; i < 20000000; i++ {
		k := fmt.Sprintf("name-68287445474-%d", i)
		v := fmt.Sprintf("|%d|听涛沁心|镜头和现实的差距等于梦和现实的距离||70后|6697812004029222155|68287445474|", i)
		res, _ := conn.Do("SETNX", k, v)
		//conn.Do("expire", k, 10)
		//res, err := conn.Do("HSET", "huoshan", k, v)
		if (i % 10000) == 0 {
			fmt.Println(i, res)
		}

	}

	res1, err := redis.String(conn.Do("HGET", "student", "name"))
	fmt.Printf("res:%s,error:%v", res1, err)
}*/
