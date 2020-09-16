# kk_core
框架的核心类


## 需要实现自己的Conn接口和IPacket接口
因每个项目协议结构和包体的数据交换格式都不相同，需要自己实现Conn接口和IPacket接口，项目中分别实现了一个WSConn(websocket协议)， TcpConn(tcp协议)，包头经过加密


## 编程模型及业务层调度原理


##  性能测试


