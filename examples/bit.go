package examples

import "fmt"

//位操作测试
func bitTest() {
	// i := int64(13)

	// val := getBinaryValueAt(i, 0)
	// fmt.Println("i=", i, ",val=", val)

	// svrNo := uint16(10)
	// stype := uint8(2)
	// areaID := uint8(1011)

	// buf := new(bytes.Buffer)
	// binary.Write(buf, binary.BigEndian, svrNo)
	// binary.Write(buf, binary.BigEndian, stype)
	// binary.Write(buf, binary.BigEndian, clusterID)
	// svcSvr := binary.BigEndian.Uint32(buf.Bytes())

	// fmt.Println("svcSvr=", svcSvr)

	// no := uint16((svcSvr >> 16) & 0xFFFF) //取高16位

	// fmt.Println("no=", no)

	// //取低8位
	// cID := svcSvr & 0xFF
	// fmt.Println("cID=", cID)

	//------------------位操作：高10位是svnNo,低22位是connID--------------------------

	// svrNo := uint32(10)
	// connID := uint32(108)

	// clientID := (svrNo << 22) | connID
	// fmt.Println("clientID=", clientID)
	// clientID := 12582942
	// // //取低22位
	// tID := clientID & 0x3fffff
	// fmt.Println("tID=", tID)

	// // // //取高10位
	// gatewayNo := (clientID >> 22) & 0x3ff
	// fmt.Println("gatewayNo=", gatewayNo)

	// //-----------------------位操作：高4位是虚拟机index， 低60位是序号---------------------

	// index := uint64(3)
	// num := uint64(127)
	// key := (index << 60) | num

	// //取低60位
	// tnum := key & 0xfffffffffffffff

	// //取高4位
	// indexNu := (key >> 60) & 0xf

	// fmt.Println("key=", key, ",indexNu=", indexNu, ",tnum=", tnum)

	//------------------------位操作：高12位是服务编号，中间4位是服务类型，低16位是区服ID-----------------------------
	// world := uint32(1008)                          //区服ID
	// srvType := uint32(6)                           //登录服
	// svrNo := uint32(10)                            //编号
	// worldID := (svrNo << 20) | srvType<<16 | world //中间4位，28 + 4 => 12 + 4 + 16,28-12=16

	// worldID = 1572864
	// fmt.Println("worlID=", worldID)

	// //取低16位
	// tWorld := worldID & 0xffff
	// //取中间4位
	// tSrvType := (worldID >> 16) & 0xf
	// //取高12位
	// tSvrNo := (worldID >> 20) & 0xfff

	// fmt.Println("tSvrNo=", tSvrNo, ",tSrvType=", tSrvType, ",tWorld=", tWorld)

	//------------------------位操作：高10位是聊天服编号 + 18位序号+ 4位频道类型 + 32位时间戳-----------------------------

	// chatSvrNum := uint64(1)
	// orderId := uint64(208)
	// channelType := uint64(4) //4是单聊
	// timestamp := uint64(time.Now().Unix())

	// fmt.Println("chatSvrNum=", chatSvrNum, ",orderId=", orderId, ",channelType=", channelType, ",timestamp=", timestamp)

	// channelId := (chatSvrNum << 54) | (orderId << 36) | (channelType << 32) | timestamp

	channelId := 18014554827320190
	fmt.Println("channelId=", channelId)

	//取高10位
	tChatSvrNum := channelId >> 54

	//取18位序号
	tOrderId := (channelId >> 36) & 0x3ffff

	//取4位频道
	tChannelType := (channelId >> 32) & 0xf

	//取低32位
	tTimeStame := channelId & 0xffffffff

	fmt.Println("tChatSvrNum=", tChatSvrNum, ", tOrderId=", tOrderId, ",tChannelType=", tChannelType, ", tTimeStame=", tTimeStame)
}
