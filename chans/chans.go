package chans

import (
	"fmt"
	"time"
	"github.com/yihubaikai/gopublic/net"
)


//------------群处理模块----------------------------start  --------------------
//全局变量

var jobs chan string;
var _URL      string
var _CHATID   string
var iStart    int = 0;

//--------------------------------------------------
//定义两个结构体
type Server struct {	
	Name string `json:"name"`
	Link string `json:"link"`
}

func Sleep(){
	time.Sleep( 1 )
}


func chans_init(){
	if(iStart == 0){
		iStart = 1
		fmt.Println("****************RunWork start****************")
		fmt.Println("****************RunWork end******************")
		go RunWork("https://api.telegram.org/bot5435489225:AAHa1ch62IOihWUKi6Qir3WiGd3End6RU9E/sendMessage","954559766")
	}
}

//传入请求
func PutString(text string) string {
	if(iStart == 0){
		fmt.Println("****************chans_init start****************")
		chans_init()
		fmt.Println("****************chans_init end******************")
	}
	jobs <- string(text)
	fmt.Println("PutString:", text)
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


func RunWork(url, _chat_id string) {
	//两个channel，一个用来放置工作项，一个用来存放处理结果。
	jobs = make(chan string, 1000)
	_URL    = url
	_CHATID = _chat_id
	// 开启三个线程，也就是说线程池中只有3个线程，实际情况下，我们可以根据需要动态增加或减少线程。
	for w := 0; w < 10; w++ {
		go worker(w, jobs)
	}
}

/*00000
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
