
![image](https://github.com/ypcd/gstunnel/blob/master/img/gstunnel%20flow%20Diagram.png)

![image](https://github.com/ypcd/gstunnel/blob/master/img/gstunnel%20Class%20Diagram.png)

gstunnel（go security tunnel，go语言开发的安全网络隧道）

This is a secure network tunnel.
```
注意：
推荐使用6.2.30.2版本或更新版本。
低于6.2.30.2版本的源代码存在较为严重的安全漏洞。
```

【项目介绍】


gstunnel 是基于go 语言开发的高性能、高并发的跨平台轻量级高安全网络加密隧道，支持tcp协议。

```
gstunnel采用安全优先原则，不追求最大化性能，所以没有使用一些会降低安全性的设计。
比如，协程池、内存对象池、socket连接复用（在复杂的网络环境中，反而会降低网络通信性能）等设计。
```
高性能：
```
项目采用多go协程和无锁模式，保证了gstunnel的高性能和高并发。

多go协程，可以充分使用多个cpu核心进行并行计算，发挥多核cpu的最大性能。

无锁模式的采用最大限度的减少了数据竞争的发生，和因为数据竞争带来的安全问题。（通过良好的项目设计，实现项目核心代码无锁）

使用go标准库的net库。Go 语言层面采用阻塞+多协程模式。

网络模型： 

因为go的默认net库，底层基于非阻塞+多路复用模型（windows iocp、linux epoll）， 所以gstunnel实质模型为非阻塞+多路复用模型。
保证了gstunnel网络通信的高性能。

为了更好的性能，使用protobuf作为序列化格式。（3.8.1及之前的版本使用json作为序列化格式。使用protobuf格式后，性能为json版的4倍。）

```
安全性：
```
golang并没有在语言层面提供完整的内存安全性（比如，rust）保证。
golang使用gc管理内存，提供了部分内存安全性，但是依然可能因为数据竞争出现内存安全问题。


项目采用AES-GCM-256加密算法，使用动态的对称密钥进行加密，通过安全的gstunnel协议进行网络通信。
动态密钥机制，默认情况下，每隔一分钟密钥就会进行更新，使用新密钥替换旧密钥。
动态密钥机制，将会大幅增加攻击者破解加密数据的难度，提供更好的安全性。

动态密钥机制，会带来轻微的性能损耗（小于5%)，为了更高的安全性，采用了这样的设计。

为了保证更高的安全性，项目采用基于硬件的真随机数生成器（Intel Digital Random Number Generator (DRNG)等）。
```

```
不建议将gstunnel作为vpn（openvpn、ipsec等）的完全替代产品使用。
```
```
Gstunnel只是一个轻量级的网络加密管道，只提供有限的安全性，安全性低于openvpn和ipsec等主流的vpn产品，无法替代主流的vpn。

gstunnel相对于传统vpn，也有一些安全优势。

gstunnel只是普通的用户态程序，不需要专属的vpn驱动程序（一般为内核态驱动，具有最高权限）。
所以即使gstunnel出现严重的安全漏洞，一般情况下，也不会危害到整个操作系统的安全。

gstunnel的动态密钥，比传统vpn使用的静态密钥具有更好的安全性。
```
gstunnel加密隧道连接是有状态连接。
它可以是长连接，也可以是短连接，具体表现取决于所承载业务连接特性。
如果业务连接是长连接，加密隧道也会保持长连接，如果业务连接是短连接，加密隧道就会是短连接。

gstunnel client的一个加密隧道连接，在mt模式下，将产生4个处理协程；在非mt模式下将产生2个处理协程。

gstunnel client的一个加密隧道连接，对应一个单独独享的加密密钥，每个加密隧道连接使用不同的加密密钥。
如果gstunnel client有10个加密隧道连接，就会有十个不的加密密钥，每个加密密钥单独负责一个加密隧道连接的加密解密。
```
支持的平台：
windows、linux、mac os

支持的应用：

http proxy（squid3等）、email、socks 5 proxy等基于tcp开发的应用。

注意：项目存在一些bug，暂时没有修复。这些bug并不影响正常使用。
```
性能表现：
```
在intel i5-1135G7 2.40GHz cpu(低压cpu，笔记本电脑cpu）配置:

gstunnel（v6.2.30.2)可以实现单端（gstunnel server或gstunnel client）超过120MBytes/s的性能。
```

gstunnel分为client和server两部分。

gstunnel 基于aes进行数据加密。

流程示意：

网络中，a到b的网络通信。

a-->b

使用gstunnel 后， a到b的网络通信。

a-->gstunnel client -->gstunnel server -->b

gstunnel 为a、b之间的网络通信添加了一个加密层。

使得a、b的通信数据，变为了加密数据，这样第三方就不能获知a、b的通信内容。从而保证了a、b网络通信的安全。

使用方法:

可以通过“go get”工具安装。

或者将项目源代码下载后，拷贝到”$GOPATH\src”目录下。

使用命令行工具编译或者安装项目源代码。

“go build gstunnel_client” “go build gstunnel_server”

“go install gstunnel_client” “go install gstunnel_server”

这时你得到了两个可执行文件gstunnel_client、gstunnel_server。

3.8.1及之前的版本，编译源代码如果出现问题，请尝试输入命令“set GO111MODULE=off”，关闭go模块功能。

配置参数

可执行文件，接受基于配置文件（json）的参数设置。
client的配置文件名：config.client.json
server的配置文件名：config.server.json
配置文件参数：
```
type GsConfig struct {
	Listen             string
	Servers            []string
	Key                string
	Debug              bool
	Tmr_display_time   int
	Tmr_changekey_time int
	Mt_model           bool
}
```
```
必选参数
listen:	监听地址（字符串）
servers:目标地址（字符串数组）
key:	aes加密密钥（字符串）

注意：aes加密密钥只能是32字节大小的字符串。

可选参数
debug:			是否开启调试模式（true或false）
Tmr_display_time	设置输出到标准输出流的信息的间隔时间（单位为秒）
Tmr_changekey_time 	设置动态密钥经过多长时间进行更换（单位为秒）
Mt_model           	是否在主逻辑模块开启多协程模式（true或false）
```
```
配置文件示例：

{"listen": "127.0.0.1:33128", "servers": ["127.0.0.1:10036"], "key": "12345678901234567890123456789012"}

listen:		监听地址
servers:	目标地址
key:		aes加密密钥
```

举例说明：

Linux bash：

user@ubuntu:~$ ./gstunnel_client

user@ubuntu:~$ ./gstunnel_server

注意：请保证client在linux系统中为可执行文件。是否是可执行文件，请查看client文件的文件属性。

Windows cmd：

C:> ./gstunnel_client.exe

C:> ./gstunnel_server.exe

日志

gstunnel在工作目录下自动生成日志文件，日志文件记录gstunnel运行时产生的错误信息。
日志文件名为gstunnel_client.err.log、gstunnel_server.err.log。

项目地址：https://github.com/ypcd/gstunnel

项目基于GPLv3协议开源。

---------------------------------------------

gstunnel (go security tunnel, secure network tunnel developed in go language)

This is a secure network tunnel.
```
Note:
6.2.30.2 or later is recommended.
Source code versions earlier than 6.2.30.2 have serious security vulnerabilities.
```

【 Project Introduction 】


gstunnel is a high performance, high concurrency cross-platform lightweight high security network encryption tunnel developed based on go language, which supports tcp protocol.

```
gstunnel uses a security first principle and does not seek to maximize performance, so it does not use designs that would reduce security.
For example, coroutine pooling, memory object pooling, socket connection reuse (in complex network environment, but will reduce the performance of network communication) and other designs.
```
High performance:
```
The project uses multi-GO coroutines and lock-free mode to ensure the high performance and high concurrency of gstunnel.

Multi-go coroutines can make full use of multiple cpu cores for parallel computing to maximize the performance of multi-core cpus.

The lock-free mode minimizes the occurrence of data contention and the security problems caused by data contention. (Through good project design, the core code of the project is lock-free)

net library using the go standard library. At the Go language level, blocking + multi-coroutine mode is adopted.

Network model:

Because go's default net library is based on non-blocking + multiplexing models (windows iocp, linux epoll), the gstunnel substantial model is a non-blocking + multiplexing model.
The high performance of gstunnel network communication is guaranteed.

For better performance, use protobuf as the serialization format. (Versions 3.8.1 and earlier use json as the serialization format. In protobuf format, the performance is four times that of the json version.)

```
Security:
```
golang does not provide full memory security (e.g., rust) guarantees at the language level.
golang uses gc to manage memory, which provides some memory security, but memory security issues can still occur due to data contention.


The project uses AES-GCM-256 encryption algorithm, using dynamic symmetric key encryption, network communication through the secure gstunnel protocol.
Dynamic key mechanism. By default, the key is updated every minute, replacing the old key with the new one.
Dynamic key mechanism will greatly increase the difficulty of attackers to crack encrypted data and provide better security.

The dynamic key mechanism, which brings a slight performance loss (less than 5%), is designed for higher security.

In order to ensure higher security, the project uses Intel Digital Random Number Generator (DRNG) based on hardware.
```

```
gstunnel is not recommended as a complete alternative to vpn (openvpn, ipsec, etc.).
```
```
Gstunnel is a lightweight network encryption pipeline that provides limited security compared to mainstream vpn products such as openvpn and ipsec, and cannot replace mainstream VPNS.

gstunnel also offers several security advantages over traditional VPNS.

gstunnel is a common user-mode program that does not require a dedicated vpn driver (usually a kernel driver with the highest permissions).
Therefore, even if gstunnel has serious security vulnerabilities, it will not endanger the security of the entire operating system under normal circumstances.

The dynamic key of gstunnel provides better security than the static key used in traditional VPNS.
```
A gstunnel encrypted tunnel connection is a stateful connection.
It can be either a long or a short connection, depending on the nature of the hosted service connection.
If the service connection is a long connection, the encryption tunnel also remains a long connection. If the service connection is a short connection, the encryption tunnel is a short connection.

An encrypted tunnel connection of the gstunnel client, in mt mode, produces four processing coroutines; In non-MT mode, two processing coroutines are generated.

An encryption tunnel connection of the gstunnel client, corresponding to a separate and exclusive encryption key. Each encryption tunnel connection uses a different encryption key.
If the gstunnel client has 10 cryptographic tunnel connections, there will be 10 cryptographic keys, each of which is solely responsible for encrypting and decrypting one cryptographic tunnel connection.
```
Supported platforms:
windows, linux, mac os

Supported applications:

Applications such as http proxy (squid3), email, and socks 5 proxy are developed based on tcp.

Note: There are some bugs in the project that have not been fixed yet. These bugs do not affect normal use.
```
Performance Performance:
```
intel i5-1135G7 2.40GHz cpu(low-voltage cpu, laptop cpu) configuration:

gstunnel (v6.2.30.2) can achieve single-ended (gstunnel server or gstunnel client) performance of more than 120MBytes/s.
```

gstunnel is divided into client and server parts.

gstunnel encrypts data based on aes.

Schematic diagram of process:

In network, the network communication between a and b.

a-->b

After using gstunnel, a to b network communication.

a-->gstunnel client -->gstunnel server -->b

gstunnel adds an encryption layer for network communication between a and b.

As a result, the communication data of a and b becomes encrypted data, so that a third party cannot know the communication content of A and b. This ensures the security of network communication between a and b.

Method of use:

You can install it using the go get tool.

Or download the project source code and copy it to the $GOPATH\src directory.

Compile or install the project source code using a command-line tool.

"go build gstunnel_client" "go build gstunnel_server"

go install gstunnel_client go install gstunnel_server

You have two executables, gstunnel_client and gstunnel_server.

3.8.1 and earlier versions, if there is a problem compiling the source code, please try to enter the command "set GO111MODULE=off" to turn off the go module function.

Configuration parameters

An executable that accepts configuration file (json) based parameter Settings.
The client configuration file name is config.client.json
server configuration file name: config.server.json
Configuration file parameters:
```
type GsConfig struct {
Listen             string
Servers            []string
Key                string
Debug              bool
Tmr_display_time   int
Tmr_changekey_time int
Mt_model           bool
}
```
```
Will choose parameters
listen: listening address (string)
servers: destination address (array of strings)
key: aes encryption key (string)

Note: The aes encryption key can only be a string of 32 bytes.

Optional parameters
debug: Whether to enable the debug mode (true or false)
Tmr_display_time Sets the interval (in seconds) at which information is output to the standard output stream
Tmr_changekey_time Specifies the interval in seconds for changing dynamic keys.
Mt_model Whether to enable multiple coroutines in the main logic module (true or false)
```
```
Example configuration file:

{" listen ", "127.0.0.1:33128", "the servers:" [] "127.0.0.1:10036", "key" : "12345678901234567890123456789012"}

listen: indicates a listening address
servers: target address
key: aes encryption key
```

For example:

Linux bash:

user@ubuntu:~$ ./gstunnel_client

user@ubuntu:~$ ./gstunnel_server

Note: Ensure that client is an executable file in linux. If it is an executable file, view the file properties of the client file.

Windows CMD:

C:> ./gstunnel_client.exe

C:> ./gstunnel_server.exe

The log

gstunnel automatically generates a log file in the working directory. The log file records the error information generated when gstunnel is running.
The log files are named gstunnel_client.err.log and gstunnel_server.err.log.

The address of the project: https://github.com/ypcd/gstunnel

The project is open source based on GPLv3 protocol.
