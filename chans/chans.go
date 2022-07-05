package chans

import (
	"fmt"
	"strings"
	"github.com/yihubaikai/gopublic/net"
	"github.com/jinzhu/configor"      //配置文件
	
)


//------------全局变量---------------------------- 
var jobs chan string;               //数据通道
var _URL      string                //请求URL
var _CHATID   string                //聊天ID
var BotKWords  map[string]string    //关键字字符串
var BotFWords map[string]string     //过滤字符串
var iStart    int = 0;              //初始化标识


//------------定义配置文件存放的结构体---------------
//954559766
var Config = struct {
    AppName string `default:"QQBot"`

    Telegram struct {
        Token           string `required:"true" default:"5435489225:AAHa1ch62IOihWUKi6Qir3WiGd3End6RU9E"`
        Chat_id         string `required:"true" default:"0"`
        Url_getupdates  string `required:"true" default:"https://api.telegram.org/bot5435489225:AAHa1ch62IOihWUKi6Qir3WiGd3End6RU9E/getUpdates"`
        Url_sendmessage string `required:"true" default:"https://api.telegram.org/bot5435489225:AAHa1ch62IOihWUKi6Qir3WiGd3End6RU9E/sendMessage"`
    }

   Filterwords  string `required:"true" default:"+群"`
   Keywords     string `required:"true" default:"谁|有没有|价格|多少钱|来一|带价|接单|哪里|全体成员|能买|报价|优先"`
}{}

 

//-------------查找关键字函数 ------------------------
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

//------------分割关键字函数----------------------------
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


//------------写一个测试函数----------------------
func Test(text string)bool{
  iStart = 0
  chans_init()
  fRet := Find_Order_Message(text, BotFWords)
  bRet := Find_Order_Message(text, BotKWords)
  fmt.Println("fRet:", fRet, "bRet:", bRet)
  fbRet := false
  if(!fRet && bRet){
	fbRet = true	
  }
  return fbRet
}



//--------------------初始化:当iStart==0的时候调用-----------------------
func chans_init(){
	if(iStart == 0){
		iStart = 1
		fmt.Println("****************chans_init**********************")
		configor.Load(&Config, "qqbot.yml")
		fmt.Println("Read Config:\n%v", Config)
		
		
		
		//BotKWords    = Split_Init( "谁|有没有|价格|多少钱|来一|带价|接单|哪里|全体成员|能买|报价|优先", "|")
		fmt.Println("Split_KeyWords:", Config.Keywords)
		BotKWords    = Split_Init( Config.Keywords, "|")
		
		//BotFWords = Split_Init( "+群|多赚钱|不禁言|价格优惠|价格实惠|收录|关键词|代写|换群|宠物|企业签|你喜欢的这都有|域名|欢迎", "|")
		fmt.Println("Split_Filterwords:", Config.Filterwords)
		BotFWords   = Split_Init( Config.Filterwords, "|")
		
		
		t := Config.Telegram
		//go RunWork("https://api.telegram.org/bot5435489225:AAHa1ch62IOihWUKi6Qir3WiGd3End6RU9E/sendMessage","954559766")
		fmt.Println("RunWork: Url:", t.Url_sendmessage, "Token:",  t.Token)
		go RunWork(t.Url_sendmessage, t.Token )
		
		
		fmt.Println("****************chans_init end******************")
	}
}

//传入请求
func PutString(text string) string {
	if(strings.Contains(text, "refresh_config")){ //当收到指令:refresh_config
	   fmt.Println("****************refresh_config******************")
	    iStart = 0
	}
	if(iStart == 0){
		fmt.Println("****************PutString chans_init****************")
		chans_init()
		fmt.Println("****************PutString chans_init end******************")
	}
	fRet := Find_Order_Message(text, BotFWords)
	if(fRet){
	        return text
	}
	bRet := Find_Order_Message(text, BotKWords)
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
FilterWords = Split_Init( "+群|多赚钱|不禁言|价格优惠|价格实惠|收录|关键词|代写|换群|宠物|企业签|你喜欢的这都有|域名|欢迎", "|")
			
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

func main() { //配置文件调用DEMO
    configor.Load(&Config, "qqbot.yml") //加载配置文件读取
    fmt.Println("config: %#v", Config)
    appname := Config.AppName
    filter  := Config.Filterwords
    Keyword := Config.Keywords
    fmt.Println("读取appname:", appname)
    fmt.Println("读取filter :", filter)
    fmt.Println("读取Keyword:", Keyword)
    t       := Config.Telegram
	fmt.Println("读取token 1:", t.Token)
	fmt.Println("读取chat_id:", t.Chat_id)
	fmt.Println("读取updates:", t.Url_getupdates)
	fmt.Println("读取sendmsg:", t.Url_sendmessage)
}

//测试方法: 加QQ好友,然后发送: refresh_config
*/
