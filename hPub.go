package hPub

import (
	"crypto/md5"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"io/ioutil"
	"log"

	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func Test() {
	fmt.Println("我是hPub中的test")
}

/*获取当前路径
"path/filepath"
"strings" //需要引入2个库
*/
func GetCurrentDir(file string) string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	ret := strings.Replace(dir, "\\", "/", -1)
	ret += "/" + file
	return ret
}

/*保存文件（优化版）*/
func SaveLog(m_FilePath string, val string) {
	var dir, filename string
	filename = filepath.Base(m_FilePath)
	if len(m_FilePath) > 1 && string([]byte(m_FilePath)[1:2]) == ":" {
		filename = filepath.Base(m_FilePath)
		dir = strings.TrimSuffix(m_FilePath, filename)
		//print("abspath:filename:" + filename + "\n" + "dir:" + dir + "\n")
	} else {
		dir, _ = filepath.Abs(filepath.Dir(os.Args[0]))
		dir = dir + "/" + m_FilePath
		filename = filepath.Base(m_FilePath)
		dir = strings.TrimSuffix(dir, filename)
		//print("noptabspath:filename:" + filename + "\n" + "dir:" + dir + "\n")
	}

	p := dir + "/" + filename
	p = strings.Replace(p, "\\", "/", -1)
	p = strings.Replace(p, "//", "/", -1)
	//print("fullpath" + p + "\n")
	_, err := os.Stat(dir)
	if err != nil {
		if !os.IsExist(err) {
			os.MkdirAll(dir, os.ModePerm)
		}
	}
	fl, err := os.OpenFile(p, os.O_APPEND|os.O_CREATE, 0644)
	defer fl.Close()

	if err != nil {
		fmt.Println("SaveLog:error")
	} else {
		io.WriteString(fl, val)
	}
}

//SaveLog扩展函数，可以输出当前调用函数
func SaveLogEx(val string) {
	funcname := Get_FuncName(2)
	var m_FilePath string
	m_FilePath = "log/" + Getday() + ".log"
	SaveLog(m_FilePath, "["+funcname+"]"+val)
	fmt.Print("[" + funcname + "]" + val)
}

/*读取文件*/
func ReadLog(m_FilePath string) string {
	var dir, filename string
	filename = filepath.Base(m_FilePath)
	if len(m_FilePath) > 1 && string([]byte(m_FilePath)[1:2]) == ":" {
		filename = filepath.Base(m_FilePath)
		dir = strings.TrimSuffix(m_FilePath, filename)
		//print("abspath:filename:" + filename + "\n" + "dir:" + dir + "\n")
	} else {
		dir, _ = filepath.Abs(filepath.Dir(os.Args[0]))
		dir = dir + "/" + m_FilePath
		filename = filepath.Base(m_FilePath)
		dir = strings.TrimSuffix(dir, filename)
		//print("noptabspath:filename:" + filename + "\n" + "dir:" + dir + "\n")
	}
	p := dir + "/" + filename
	p = strings.Replace(p, "\\", "/", -1)
	p = strings.Replace(p, "//", "/", -1)
	if !File_Exists(p) {
		return ""
	}
	fi, err := os.Open(p)
	defer fi.Close()
	if err != nil {
		panic(err)
	}
	fd, err := ioutil.ReadAll(fi)
	return string(fd)
}

/* 判断文件是否存在  存在返回 true 不存在返回false*/
func File_Exists(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}
func File_exename(pathx string) string {
	_, fileName := filepath.Split(pathx)
	return fileName
}
func file_path(pathx string) string {
	paths, _ := filepath.Split(pathx)
	return paths
}

/*获取当前时间*/
func Gettime() string {
	Year := time.Now().Year()     //年[:3]
	Month := time.Now().Month()   //月
	Day := time.Now().Day()       //日
	Hour := time.Now().Hour()     //小时
	Minute := time.Now().Minute() //分钟
	Second := time.Now().Second() //秒
	//Nanosecond:=time.Now().Nanosecond()//纳秒
	var timestr string
	timestr = fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", Year, Month, Day, Hour, Minute, Second)
	return timestr
}

/*获取日期*/
func Getday() string {
	var timestr string
	timestr = Gettime()
	timestr = string([]byte(timestr)[:10])
	return timestr
}

/*获取系统当前时间戳*/
func Gettimecuo() string {
	t := time.Now()
	timestamp := strconv.FormatInt(t.UnixNano(), 10)
	timestamp = timestamp[0:13]
	//fmt.Println(timestamp)
	//fmt.Println(t.Unix())
	return timestamp
}

//时间转到到时间戳
func Gettime_t2c(time_string string) string {
	//时间 to 时间戳
	var retstr string
	//loc, _ := time.LoadLocation("Asia/Shanghai")                                  //设置时区
	tt, _ := time.ParseInLocation("2006-01-02 15:04:05", time_string, time.Local) //2006-01-02 15:04:05是转换的格式如php的"Y-m-d H:i:s"
	retstr = fmt.Sprintf("%d", tt.Unix())
	return retstr
}

//时间戳转换成时间
func Gettime_c2t(timeUnix string) string {
	var s int64
	s, err := strconv.ParseInt(timeUnix, 10, 64)
	if err != nil {
		return ""
	}
	formatTimeStr := time.Unix(s, 0).Format("2006-01-02 15:04:05")
	return formatTimeStr
}

//获取函数名 需要导入包 "strings" "runntime"
//直接调用显示函数名 Get_FuncName（1）
//SaveLog调用Get_FuncName(2)
func Get_FuncName(iDeep int) string {
	funcAddr, _, _, _ := runtime.Caller(iDeep)
	funcName := runtime.FuncForPC(funcAddr).Name()
	ret := strings.Split(funcName, ".")
	return ret[1]
}

//加个MD5字符串
func Md5str(mingwen_text string) string {
	data := []byte(mingwen_text)
	has := md5.Sum(data)
	md5str1 := fmt.Sprintf("%x", has)
	return md5str1
}
