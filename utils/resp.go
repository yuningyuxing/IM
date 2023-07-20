package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// resp表示程序或函数返回的响应结果
type H struct {
	//响应状态码
	Code int
	//响应消息 用于描述响应状态的信息
	Msg string
	//响应数据 用于存储返回具体的业务数据
	Data interface{}
	//响应列表数据 用于存储返回的数据列表
	Rows interface{}
	//响应列表数据的总数  通常用于分页场景 表示总共有多少条数据
	Total interface{}
}

// Resp函数用于构造并发送JSON格式的HTTP响应
// w表示HTTP响应的写入器 用于将响应发送给客户端  data响应数据 msg响应消息
func Resp(w http.ResponseWriter, code int, data interface{}, msg string) {
	//设置响应头Content-Type为application/json表示响应内容为JSON格式
	w.Header().Set("Content-Type", "application/json")
	//设置HTTP响应状态码为200 表示请求成功
	w.WriteHeader(http.StatusOK)
	//存储到H示例里
	h := H{
		Code: code,
		Data: data,
		Msg:  msg,
	}
	//将结构体序列化为JSON格式的字节切片 用于传输
	ret, err := json.Marshal(h)
	if err != nil {
		fmt.Println(err)
	}
	//将序列化后的JSON格式字节切片写入HTTP响应的Body中 完成响应发送
	w.Write(ret)
}

// 跟上面一样 只不过这里返回的数据内容是列表和数据总数
func RespList(w http.ResponseWriter, code int, data interface{}, total interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	h := H{
		Code:  code,
		Rows:  data,
		Total: total,
	}
	ret, err := json.Marshal(h)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(ret)
}

// 构造表示发送失败的JSON格式的HTTP响应
func RespFail(w http.ResponseWriter, msg string) {
	Resp(w, -1, nil, msg)
}

// 表示发送成功 注意这里是调用Resp函数 下面同理
func RespOK(w http.ResponseWriter, data interface{}, msg string) {
	Resp(w, 0, data, msg)
}

// 这里也表示发送成功 但包含列表数据和总数
func RespOKList(w http.ResponseWriter, data interface{}, total interface{}) {
	RespList(w, 0, data, total)
}
