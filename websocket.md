# websocket

## 什么是websocket

首先，我们需要弄明白，WebSocket本质上一种计算机网络应用层的协议，用来弥补http协议在持久通信能力上的不足。

我们知道http协议本身是无状态协议，每一个新的http请求，只能通过客户端主动发起，通过 建立连接-->传输数据-->断开连接 的方式来传输数据，传送完连接就断开了，也就是这次http请求已经完全结束了（虽然http1.1增加了keep-alive请求头可以通过一条通道请求多次，但本质上还是一样的）。

并且服务器是不能主动给客户端发送数据的（因为之前的请求得到响应后连接就断开了，之后服务器根本不知道谁请求过），客户端也不会知道之前请求的任何信息。

所以说，http协议本身是没有持久通信能力的，但是我们在实际的应用中，是很需要这种能力的，所以WebSocket协议由此而生，于2011年被IETF定为标准RFC6455，并被RFC7936所补充规范。

并且在HTML5标准中增加了有关WebSocket协议的相关api，所以只要实现了HTML5标准的客户端，就可以与支持WebSocket协议的服务器进行全双工的持久通信了。

ps：这里的持久通信能力指的是协议本身的能力，我们当然可以通过编程的方式实现这种功能，比如轮询的方式，但谁不喜欢原生就支持呢？

## websocket协议原理

与http协议一样，WebSocket协议也需要通过已建立的TCP连接来传输数据。具体实现上是通过http协议建立通道，然后在此基础上用真正的WebSocket协议进行通信，所以WebSocket协议和http协议是有一定的交叉关系的。

websocket请求头：



请求头中重要的字段：

```
Connection:Upgrade

Upgrade:websocket

Sec-WebSocket-Extensions:permessage-deflate; client_max_window_bits

Sec-WebSocket-Key:mg8LvEqrB2vLpyCNnCJV3Q==

Sec-WebSocket-Version:13
```

1. Connection和Upgrade字段告诉服务器，客户端发起的是WebSocket协议请求

2. Sec-WebSocket-Extensions表示客户端想要表达的协议级的扩展

3. Sec-WebSocket-Key是一个Base64编码值，由浏览器随机生成

4. Sec-WebSocket-Version表明客户端所使用的协议版本

而得到的响应头中重要的字段：

```
Connection:Upgrade

Upgrade:websocket

Sec-WebSocket-Accept:AYtwtwampsFjE0lu3kFQrmOCzLQ=
```

1. Connection和Upgrade字段与请求头中的作用相同

2. Sec-WebSocket-Accept表明服务器接受了客户端的请求

Status Code:101 Switching Protocols

并且http请求完成后响应的状态码为101，表示切换了协议，说明WebSocket协议通过http协议来建立运输层的TCP连接，之后便与http协议无关了。

## websocket的优缺点

优点：

· WebSocket协议一旦建议后，互相沟通所消耗的请求头是很小的

· 服务器可以向客户端推送消息了

缺点：

· 少部分浏览器不支持，浏览器支持的程度与方式有区别

