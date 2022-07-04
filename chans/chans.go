package chans

import (
	"fmt"
	"time"
	"strings"
	"github.com/yihubaikai/gopublic/net"
)


//------------群处理模块----------------------------start  --------------------
//全局变量
var jobs chan string;              //数据通道
var _URL      string               //请求URL
var _CHATID   string               //聊天ID
var KeyWords  map[string]string     //关键字字符串
var FilterWords map[string]string   //过滤字符串
var iStart    int = 0;             //初始化标识

//--------------------------------------------------
//定义两个结构体
type Server struct {	
	Name string `json:"name"`
	Link string `json:"link"`
}


func Find_Order_Message(inText string, Filter map[string]string)bool{
	bRet := false
	for key,_ := range Filter{
		if(strings.Contains(inText, key)){
			bRet = true
			break
		}
	}
	return bRet
}


func Split_Init(text, Filt string) (map[string]string){
	Ret := make(map[string]string)
	arr := strings.Split(text, Filt)
	for _, _val := range arr {
		if(len(_val)>0){
			Ret[_val] = _val
		}
	}
	return Ret
}




func Sleep(timeN time.Duration){
	time.Sleep( timeN )
}


func chans_init(){
	if(iStart == 0){
		iStart = 1
		fmt.Println("****************chans_init Filer****************")
		KeyWords    = Split_Init( "谁|有没有|价格|多少钱|来一|带价|接单|哪里|全体成员|能买|报价|优先", "|")
		FilterWords = Split_Init( "+群|多赚钱|不禁言|价格优惠|价格实惠|收录|关键词|代写|换群|宠物|企业签|你喜欢的这都有|腾讯云|阿里云|", "|")
		fmt.Println("****************chans_init start****************")
		go RunWork("https://api.telegram.org/bot5435489225:AAHa1ch62IOihWUKi6Qir3WiGd3End6RU9E/sendMessage","954559766")
		fmt.Println("****************chans_init end******************")
	}
}

//传入请求
func PutString(text string) string {
	if(iStart == 0){
		fmt.Println("****************chans_init start****************")
		chans_init()
		fmt.Println("****************chans_init end******************")
	}
	fRet := Find_Order_Message(text, FilterWords)
	if(fRet){
	        return text
	}
	bRet := Find_Order_Message(text, KeyWords)
	if(bRet){
	    jobs <- string(text)
	}
	return text
}

//解析传入字符
func GetString(szText string) string {
  	s := make(map[string]string)
	s["chat_id"] = _CHATID
	s["text"]    = szText
	ret := hNet.Httppost(_URL, s)
    return ret
}

//创建一个解析线程
//这个是工作线程，处理具体的业务逻辑，将jobs中的任务取出，处理后将处理结果放置在results中。
func worker(id int, jobs <-chan string) {
	for j := range jobs {
	     GetString(j)
	}
}


func RunWork(_url, _chat_id string) {
	//两个channel，一个用来放置工作项，一个用来存放处理结果。
	jobs = make(chan string, 1000)
	_URL    = _url
	_CHATID = _chat_id
	// 开启三个线程，也就是说线程池中只有3个线程，实际情况下，我们可以根据需要动态增加或减少线程。
	for w := 0; w < 10; w++ {
		go worker(w, jobs)
	}
}

/*
KeyWords    = Split_Init( "谁|有没有|价格|多少钱|来一|带价|接单|哪里|全体成员|能买|报价|优先", "|")
FilterWords = Split_Init( "+群|多赚钱|不禁言|价格优惠|价格实惠|收录|关键词|代写|换群|宠物", "|")
			
#telegram-send --configure
#token =  "5435489225:AAHa1ch62IOihWUKi6Qir3WiGd3End6RU9E"

获取ID
https://api.telegram.org/bot5435489225:AAHa1ch62IOihWUKi6Qir3WiGd3End6RU9E/getUpdates

发送ID:
https://api.telegram.org/bot5435489225:AAHa1ch62IOihWUKi6Qir3WiGd3End6RU9E/sendMessage
curl -X POST "https://api.telegram.org/bot5435489225:AAHa1ch62IOihWUKi6Qir3WiGd3End6RU9E/sendMessage" -d "chat_id=954559766&text=send msg"


func main(){
	var name string
	go chans.RunWork("http://156.251.30.67:8000/telegram/")
	time.Sleep( 1 )

	for i := 0; i < 10; i++ {
		s := fmt.Sprintf("%v", i)
		chans.PutString(s)
	}
	fmt.Print  ("按任意键退出...")
  	fmt.Scanln(&name)
}
*/
