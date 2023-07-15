package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

//Md5工具类及注册密码加密

// 对输入的字符串进行MD5加密 返回加密后的结果
func Md5Encode(data string) string {
	//创建一个MD5加密对象
	h := md5.New()
	//将输入的字符串转化为字节数组并写入加密对象
	h.Write([]byte(data))
	//返回加密后的字节数组  参数nil表示我们不需要存储在指定的字节数组 让他自己返回一个即可
	tempStr := h.Sum(nil)
	//将加密后的字节数组转化为十六进制字符串
	return hex.EncodeToString(tempStr)
}

// 返回加密后的大写结果
func MD5Encode(data string) string {
	//调上面的函数 转化成大写
	return strings.ToUpper(Md5Encode(data))
}

// 生成密码的哈希值
// 将明文密码和盐值拼接后进行MD5加密
// 明文密码就是用户输入的原始密码 盐值是一个随机生成的字符串 用于增加密码的复杂度和安全性
func MakePassword(plainpwd, salt string) string {
	return Md5Encode(plainpwd + salt)
}

// 验证密码正确性 password是数据库中的密码
func ValidPassword(plainpwd, salt string, password string) bool {
	return Md5Encode(plainpwd+salt) == password
}
