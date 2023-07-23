package models

//用户A发消息给B
//A先通过websocket 给服务器  服务器对消息分类处理 然后通过UDP广播 给B  B通过UDP监听拿到消息后 分类 然后给websocket  websocket给前端页面
import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
	"net"
	"net/http"
	"strconv"
	"sync"
)

// 用来描述消息的结构体
type Message struct {
	gorm.Model
	//发送者
	FromId int64
	//接收者
	TargetId int64
	//发送类型 1私聊 2私聊 3广播
	Type int
	//消息类型 1文字 2表情包 3图片 4音频
	Media int
	//消息内容
	Content string
	Pic     string
	Url     string
	Desc    string
	//其他数字统计
	Amount int
}

func (table *Message) TableName() string {
	return "message"
}

// 注:websocket是全双工通讯
// 这个结构体是用来描述websocket连接的节点的 每个客户端连接到服务器时，会创建一个对应的Node示例来表示该连接
type Node struct {

	//表示一个websocket连接 可用于发送和接受消息
	Conn *websocket.Conn
	//表示一个数据队列  这里是一个管道
	DataQueue chan []byte
	//表示一个组的集合 它是一个集合数据结构 用于管理节点所属的组 可以通过该集合进行组的添加 删除和查找操作
	//以实现对节点的组管理功能
	GroupSets set.Interface
}

// 映射关系
// 这个clientMap变量可以用来存储和管理客户端的信息 每个客户端都可以通过唯一的键值来访问和操作对应的*Node
var clientMap map[int64]*Node = make(map[int64]*Node, 0)

// 读写锁
var rwLocker sync.RWMutex

// Chat函数用于处理websocket连接和消息发送
func Chat(writer http.ResponseWriter, request *http.Request) {
	//1.获取参数并校验token等合法性
	//获取请求中的URL
	query := request.URL.Query()
	Id := query.Get("userId")
	//将Id转化为int64 第一个参数是待转换的字符串 第二个参数表示进制  第三个参数表示结果的位数
	userId, _ := strconv.ParseInt(Id, 10, 64)
	//msgType:=query.Get("type")
	//targetId:=query.Get("targetId")
	//context:=query.Get("context")

	//这里先默认为true
	isvalida := true //checkToke()
	//将http连接升级为websocket连接
	conn, err := (&websocket.Upgrader{
		//token校验
		CheckOrigin: func(r *http.Request) bool {
			return isvalida
		},
	}).Upgrade(writer, request, nil)
	if err != nil {
		fmt.Println(err)
	}
	//声明一个连接示例
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}
	//上锁  这里是写入锁
	rwLocker.Lock()
	//通过发送者Id 建立映射关系
	clientMap[userId] = node
	//解锁
	rwLocker.Unlock()
	//完成发送逻辑
	go sendProc(node)
	//完成接受逻辑
	go recvProc(node)
	sendMsg(userId, []byte("欢迎进入聊天系统"))
}

// 处理消息的发送逻辑
func sendProc(node *Node) {
	for {
		//select是用于多个通信操作中选择一个可执行的操作  下面有多个case 它会选择能执行的执行
		//如果没有能执行的且没有default语句 则阻塞
		select {
		//它从dataQueue通道中接受待发送的消息 并将其通过websocket连接发送出去
		case data := <-node.DataQueue:
			//这句代码用于向websocket连接写入消息
			//第一个参数表示消息类型 这里是文本类型 第二个参数是消息内容
			//注意这里是给websocket那一端发送消息
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

// 处理消息的接受逻辑
func recvProc(node *Node) {
	for {
		//这里从websocket连接中读取消息
		//这里是从websocket测试那一端接受消息
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		//将接受到的消息传给broadMsg函数处理 这个函数的作用是将消息广播给它相关客户端
		broadMsg(data)
		//打印消息
		fmt.Println("[ws]<<<<<", data)
	}
}

// 定义一个管道  用于发送消息给UDP服务器
var udpsendChan chan []byte = make(chan []byte, 1024)

// 将信息写入UDP管道
func broadMsg(data []byte) {
	fmt.Println("2")
	udpsendChan <- data
}

// 启动UDP的功能
func init() {
	go udpSendProc()
	go udpRecvProc()
}

// 向指定的IP地址和端口发送消息
// 实现消息的广播和接收功能
func udpSendProc() {
	//这里创建了一个UDP连接
	//下面参数分别是 协议类型  本地地址 和目标地址
	con, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(192, 168, 0, 255),
		Port: 3000,
	})
	//延时关闭连接
	defer con.Close()
	if err != nil {
		fmt.Println(err)
	}
	for {
		//监听udpsendChan通道的数据 当通道中有数据时，将其取出并通过连接对象con发送出去
		select {
		case data := <-udpsendChan:
			//通过con.Write将消息数据发送到目标地址
			_, err := con.Write(data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

// 负责从该IP地址和端口接受消息
// 监听UPD连接并接受来自其他客户端的消息
func udpRecvProc() {
	//创建一个UDP连接
	con, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 3000,
	})
	//延迟关闭连接
	defer con.Close()
	if err != nil {
		fmt.Println(err)
	}
	for {
		//存储接收到的UDP数据
		var buf [512]byte
		//从UDP连接中读取数据  n是字节数
		n, err := con.Read(buf[0:])
		if err != nil {
			fmt.Println(err)
			return
		}
		//调用该函数 将接受到的数据传递给他进行处理
		dispatch(buf[0:n])
	}
}

// 将接收到的消息内容进行分发处理
func dispatch(data []byte) {
	//用于存储解析后的消息数据
	msg := Message{}
	//将收到的消息 反序列化为对应的结构体
	err := json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	//根据消息类型分类处理
	switch msg.Type {
	case 1: //私信
		//第一个参数是目标用户和消息类型
		sendMsg(msg.TargetId, data)
		//case 2:
		//	sendGroupMsg()
		//case 3:
		//	sendAllMsg()
	}
}

// 给指定用户发送消息
func sendMsg(userId int64, msg []byte) {
	//上锁 锁定读取map的操作 保证并发安全
	//注意这里是读锁
	rwLocker.RLock()
	//通过userId获取对应节点
	node, ok := clientMap[userId]
	//解锁读取操作
	rwLocker.RUnlock()
	if ok {
		//向管道传递消息
		node.DataQueue <- msg
	}
}
