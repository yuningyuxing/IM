package models

//用户A发消息给B
//A先通过websocket 给服务器  服务器对消息分类处理 然后通过UDP广播 给B  B通过UDP监听拿到消息后 分类 然后给websocket  websocket给前端页面
import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
	"main/utils"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// 用来描述消息的结构体
type Message struct {
	gorm.Model
	//发送者
	UserId int64
	//接收者
	TargetId int64
	//发送类型 1私聊 2私聊 3广播
	Type int
	//消息类型 1文字 2表情包 3音频 4图片
	Media int
	//消息内容
	Content string
	//创建时间
	CreateTime uint64
	//读取时间
	ReadTime uint64
	Pic      string
	Url      string
	Desc     string
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
	//客户端地址
	Addr string
	//首次连接时间
	FirstTime uint64
	//心跳时间
	HeartbeatTime uint64
	//登录时间
	LoginTime uint64
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
	//2.获取Conn

	//获取当前时间 用来初始化心跳时间和登录时间
	currentTime := uint64(time.Now().Unix())
	//声明一个连接示例 并初始化一些参数
	node := &Node{
		Conn: conn,
		Addr: conn.RemoteAddr().String(),
		//初始化心跳时间
		HeartbeatTime: currentTime,
		//初始化登录时间
		LoginTime: currentTime,
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
	//加入在线用户到缓存
	SetUserOnlineInfo("online_"+Id, []byte(node.Addr), time.Duration(viper.GetInt("timeout.RedisOnlineTime"))*time.Hour)
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
		//将接受来的消息进行反序列化
		msg := Message{}
		err = json.Unmarshal(data, &msg)
		if err != nil {
			fmt.Println("json unmarshal fail err = ", err)
			return
		}
		//进行心跳检测
		if msg.Type == 3 {
			//这里更新用户心跳
			currentTime := uint64(time.Now().Unix())
			//此函数用来更新用户心跳
			node.Heartbeat(currentTime)
		} else {
			dispatch(data)
			//将消息广播到局域网
			//broadMsg(data)
			fmt.Println("[ws] recvProc <<<<", string(data))
		}
	}
}

// 定义一个管道  用于发送消息给UDP服务器
var udpsendChan chan []byte = make(chan []byte, 1024)

// 将信息写入UDP管道
func broadMsg(data []byte) {
	udpsendChan <- data
}

// 启动UDP的功能
func init() {
	go udpSendProc()
	go udpRecvProc()
	fmt.Println("init goroutine")
}

// 向指定的IP地址和端口发送消息
// 实现消息的广播和接收功能
func udpSendProc() {
	//这里创建了一个UDP连接
	//下面参数分别是 协议类型  本地地址 和目标地址
	con, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(192, 168, 0, 255),
		Port: viper.GetInt("port.udp"),
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
		Port: viper.GetInt("port.udp"),
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
	msg.CreateTime = uint64(time.Now().Unix())
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
	case 2: //群发
		sendGroupMsg(msg.TargetId, data)
		//case 3: //更新心跳
	}
}

// 第一个参数就表示群ID
func sendGroupMsg(targetId int64, msg []byte) {
	fmt.Println("开始群发消息")
	//找到所有群成员
	userIds := SearchUserByGroupId(uint(targetId))
	for i := 0; i < len(userIds); i++ {
		//注意这里排除自己
		if targetId != int64(userIds[i]) {
			sendMsg(int64(userIds[i]), msg)
		}
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
	//将消息反序列化
	jsonMsg := Message{}
	json.Unmarshal(msg, &jsonMsg)
	//创建一个空的上下文
	ctx := context.Background()
	//获得接收者的ID 并转化为字符串
	targetIdStr := strconv.Itoa(int(userId))
	//将发送者的ID 转化为字符串
	userIdStr := strconv.Itoa(int(jsonMsg.UserId))
	//将消息的创建时间设置为当前
	jsonMsg.CreateTime = uint64(time.Now().Unix())
	//获取接受者用户的在线信息  根据我们刚获得的ID  .Result()的作用是获得redis命令的返回值和错误
	r, err := utils.Red.Get(ctx, "online_"+targetIdStr).Result()
	if err != nil {
		fmt.Println("redis get fail err=", err)
		return
	}
	//如果接收者在线 则发送消息给接收者
	if r != "" && ok {
		fmt.Println("sendMsg>>>userID: ", userId, " msg:", string(msg))
		node.DataQueue <- msg
	}
	//之后我们将消息缓存到redis中
	//我们将两个人的消息都缓存在一个有序集合中 然后根据消息的结构体来是谁发谁的即可  score可以代表顺序
	//确定存储消息的redis键名
	var key string
	//我们这里将ID小的人放前面
	if userId > jsonMsg.UserId {
		key = "msg_" + userIdStr + "_" + targetIdStr
	} else {
		key = "msg_" + targetIdStr + "_" + userIdStr
	}
	//在redis中这里获取指定键名的有序集合中的所有元素  这里start stop表示从第一个元素开始到最后一个元素
	//有序集合中依靠一个score分数排序
	res, err := utils.Red.ZRevRange(ctx, key, 0, -1).Result()
	if err != nil {
		fmt.Println("redis zrevrange fail err=", err)
		return
	}
	//指定一个分数 以便将新构造的消息加入有序集合
	//新分数为集合容量+1
	score := float64(cap(res)) + 1
	//在指定的有序集合加入新消息
	ress, err := utils.Red.ZAdd(ctx, key, &redis.Z{score, msg}).Result()
	if err != nil {
		println("redis zadd fail err=", err)
		return
	}
	fmt.Println(ress)
}

// 重写该方法 来使消息得以序列化
func (msg Message) MarshalBinary() ([]byte, error) {
	return json.Marshal(msg)
}

// 获取缓存里的消息
// 参数userIdA和userIdB表示获取谁和谁的信息
// start和end表示需要获取的消息记录的范围
// isRev是一个布尔值 表示是否按照时间倒序获取消息
// 返回一个字符串切片 包含的就是消息记录
func RedisMsg(userIdA int64, userIdB int64, start int64, end int64, isRev bool) []string {
	//创建一个上下文 用于在redis操作中传递上下文信息
	ctx := context.Background()
	//将两个人的Id转化为字符串
	userIdStr := strconv.Itoa(int(userIdA))
	targetIdStr := strconv.Itoa(int(userIdB))
	//拿到我们存储缓存信息的键值 就是两个人的ID一起 小的在前面
	var key string
	if userIdA > userIdB {
		key = "msg_" + targetIdStr + "_" + userIdStr
	} else {
		key = "msg_" + userIdStr + "_" + targetIdStr
	}
	//用于存储消息的字符串切片
	var rels []string
	var err error
	//看看要不要倒序  这里score就表示时间
	if isRev {
		//这里是按照score从小到大的顺序获取信息
		//也就是遍历字符串的时候 是老信息在前面
		rels, err = utils.Red.ZRange(ctx, key, start, end).Result()
	} else {
		//这里是按照从大到小的顺序
		//这里就是新信息在前面
		rels, err = utils.Red.ZRange(ctx, key, start, end).Result()
	}
	if err != nil {
		fmt.Println("redisMsg find fail err=", err)
	}
	return rels
}

// 更新心跳时间的函数
func (node *Node) Heartbeat(currentTime uint64) {
	node.HeartbeatTime = currentTime
}

// 清理超时连接
func CleanConnection(param interface{}) (result bool) {
	result = true
	//捕获异常
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("cleanConnection err=", err)
		}
	}()
	//获取当前时间
	currentTime := uint64(time.Now().Unix())
	//遍历所有的websocket连接
	for i := range clientMap {
		node := clientMap[i]
		//根据写的判断函数 判断是否超时
		if node.IsHeartbeatTimeOut(currentTime) {
			fmt.Println("心跳超时 连接关闭：", node)
			node.Conn.Close()
		}
	}
	return result
}

// 判断用户心跳是否超时
// 这里就是node的上次心跳时间+最大心跳时间 如果小于等于当前时间 就表示隔的太久了 让他下线
func (node *Node) IsHeartbeatTimeOut(currentTime uint64) (timeout bool) {
	if node.HeartbeatTime+viper.GetUint64("timeout.HeartbeatMaxTime") <= currentTime {
		fmt.Println("心跳超时，自动下线", node)
		timeout = true
	}
	return
}
