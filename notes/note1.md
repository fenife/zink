## 视频笔记（typora）

### 03-zinxV0.1-基础server模块定义

<img src="assets/zinx-architecture.png" alt="zinx-architecture" style="zoom: 33%;" />

<img src="assets/image-20211114120809477.png" alt="image-20211114120809477" style="zoom: 33%;" />

<img src="assets/image-20211113115448658.png" alt="image-20211113115448658" style="zoom:50%;" />

### 04-zinxV0.1-基础server模块启动实现

<img src="assets/image-20211114120555112.png" alt="image-20211114120555112" style="zoom:50%;" />

### 05-zinxV0.1-开发服务器应用

同上

### 06-zinxV0.2-链接模块的封装(方法与属性) 

<img src="assets/image-20211114174302325.png" alt="image-20211114174302325" style="zoom:50%;" />

```go
func (c *Connection) Start() {
   panic("implement me")
}

func (c *Connection) Stop() {
   panic("implement me")
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
   panic("implement me")
}

func (c *Connection) GetConnID() uint32 {
   panic("implement me")
}

func (c *Connection) RemoteAddr() net.Addr {
   panic("implement me")
}

func (c *Connection) Send(data []byte) error {
   panic("implement me")
}
```