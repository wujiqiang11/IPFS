package utils

import (
	"IPFS/common/message"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)

//这里将这些方法关联到结构体
type Transfer struct{
	//
	Conn net.Conn
	Buf [8096]byte  //传输时使用的缓冲
}

//读包，解包函数
func (this *Transfer)ReadPkg() (mes message.Message, err error){
	buf := make([]byte, 8086)
	_, err = this.Conn.Read(buf[:4]) //先读出包的长度
	if err != nil{
		return
	}
	//根据buf[:4] 转成一个 unit32类型
	var pkgLen uint32
	pkgLen = binary.BigEndian.Uint32(buf[0:4])

	//根据 长度 pkgLen 读取下一个内容包
	n, err := this.Conn.Read(buf[:pkgLen])
	if n!=int(pkgLen) || err != nil{
		return
	}
	//将 buf[:pkgLen] 反序列化成 -> message.Message
	err = json.Unmarshal(buf[:pkgLen], &mes)
	if err != nil {
		fmt.Println("readPkg json.Unmarshal err=", err)
		return
	}
	return

}

func (this *Transfer) WritePkg(mes message.Message)(err error)  {
	data, err := json.Marshal(mes)
	if err != nil{
		fmt.Println("mes json.Marshal err=", err)
		return
	}
	//7. 此时data为待发送数据包
	//7.1 现将data的字节数发送给对方进行检错
	//先将 data长度->转成一个byte切片
	var pkgLen uint32
	pkgLen = uint32(len(data))
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[0:4], pkgLen)

	n, err := this.Conn.Write(buf[:4])  //发送数据字节长度
	if n!=4 || err != nil {
		fmt.Println("conn.Write(head) err=", err)
		return
	}
	_, err = this.Conn.Write(data)  //发送数据字节长度
	if err != nil {
		fmt.Println("conn.Write(body) err=", err)
		return
	}
	return
}