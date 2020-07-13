package clientProcess

import (
	"IPFS/common/message"
	"IPFS/common/utils"
	"encoding/json"
	"fmt"
	"net"
)


func Login(userId int) (conn net.Conn, code int) {
	conn, err := net.Dial("tcp","localhost:8889")  //连接到服务器
	if err != nil{
		fmt.Println("net.Dail err=", err)
		return
	}

	//2. 通过conn发送消息给服务
	var mes message.Message  //先发送一个message包
	mes.Type = message.LoginMesType    //分别确定message包的类型和数据
	//3.创建一个LoginMes 结构体
	var loginMes message.LoginMes
	loginMes.UserId = userId
	//4. 将loginMes 序列化
	data, err := json.Marshal(loginMes)
	if err != nil{
		fmt.Println("loginMes json.Marshal err=", err)
		return
	}
	//5. 把data赋值给message 的data字段
	mes.Data = string(data)

	tf := &utils.Transfer{
		Conn: conn,
	}

	err = tf.WritePkg(mes)//  发送mes包
	//接收服务器反馈
	resMes, err := tf.ReadPkg()

	if err!= nil{
		fmt.Println("readPkg(resMes) err=", err)
	}
	var loginResMes message.ResMes
	err = json.Unmarshal([]byte(resMes.Data), &loginResMes)
	if err != nil{
		fmt.Println("json.Unmarshal(resMes.Data) err=", err)
	}
	code = loginResMes.Code
	return
}
func UpLoad(conn net.Conn)  (err error){
	fmt.Println("请输入待上传的消息:")
	var UpLoadMes string
	fmt.Scanf("%s\n", &UpLoadMes)
	fmt.Println(UpLoadMes)
	//对输入的字符串进行加密

	var mes message.Message
	mes.Type = message.UpLoadMesType

	var upMes message.UpLoadMes
	upMes.Cipher = UpLoadMes
	data, err := json.Marshal(upMes)
	if err != nil{
		fmt.Println("upMes json.Marshal err=", err)
		return
	}
	mes.Data = string(data)
	tf := &utils.Transfer{
		Conn: conn,
	}
	err = tf.WritePkg(mes)//  发送mes包
	//接收服务器反馈
	resMes, err := tf.ReadPkg()
	if err!= nil{
		fmt.Println("readPkg(resMes) err=", err)
	}
	var UpResMes message.ResMes
	err = json.Unmarshal([]byte(resMes.Data), &UpResMes)
	if err != nil{
		fmt.Println("json.Unmarshal(UpResMes.Data) err=", err)
	}
	code := UpResMes.Code
	if code==300 {
		fmt.Println("上传成功")
	}else {
		fmt.Println("上传失败")
	}
	return
}

func DownLoad(conn net.Conn) (err error){

	var mes message.Message
	mes.Type = message.DownloadReqType
	tf := &utils.Transfer{
		Conn: conn,
	}
	err = tf.WritePkg(mes)//  发送mes包
	if err != nil{
		fmt.Println("tf.WritePkg err=", err)
	}

	//接收服务器反馈
	resMes, err := tf.ReadPkg()
	if err!= nil{
		fmt.Println("readPkg(resMes) err=", err)
	}
	var DlResMes message.DownloadRes
	err = json.Unmarshal([]byte(resMes.Data), &DlResMes)
	if err != nil{
		fmt.Println("json.Unmarshal(DlResMes.Data) err=", err)
	}
	MesNum := DlResMes.MesNum
	MesAdd := DlResMes.ResMes
	if MesNum==0 {
		fmt.Println("当前无未读消息")
		return
	}else {
		fmt.Printf("当前未读消息的条数为%d ,消息地址如下所示：\n", MesNum)
		for _,v := range MesAdd{
			fmt.Println(v)
		}
		fmt.Println("请输入你想取出的消息地址")
		fmt.Println("----------------------------------------------------")
		var add int
		fmt.Scanf("%d\n", &add)
		//将地址回送给服务器
		var mes message.Message
		mes.Type = message.DownloadAddrType

		var DlAddr message.DownloadAddr
		DlAddr.Addr = add
		data, err := json.Marshal(DlAddr)
		if err != nil{
			fmt.Println("DlAddr json.Marshal err=", err)
		}
		mes.Data = string(data)
		tf := &utils.Transfer{
			Conn: conn,
		}
		err = tf.WritePkg(mes)//  发送mes包

		//接收服务器反馈
		resMes, err := tf.ReadPkg()
		if err!= nil{
			fmt.Println("readPkg(resMes) err=", err)
		}
		var DlContMes message.DownloadCont
		err = json.Unmarshal([]byte(resMes.Data), &DlContMes)
		if err != nil{
			fmt.Println("json.Unmarshal(DlResMes.Data) err=", err)
		}
		if DlContMes.Code == 404{
			fmt.Println("下载出错，请重试")
		}else {
			fmt.Println("你读取的消息内容如下:")
			fmt.Println(DlContMes.Cipher)
		}

	}
    return
}

