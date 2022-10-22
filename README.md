
![image](https://github.com/ypcd/gstunnel/blob/master/img/gstunnel%20flow%20Diagram.png)

![image](https://github.com/ypcd/gstunnel/blob/master/img/gstunnel%20Class%20Diagram.png)


This is a secure network tunnel.

注意：
推荐使用2.7版本或高于2.7版本的项目源代码。
低于2.7版本的源代码存在较为严重的安全漏洞。

【gstunnel介绍】

gstunnel 是基于go 语言开发的高性能、高并发的跨平台轻量级安全网络加密管道，支持tcp协议。

项目采用多协程和无锁模式，保证了gstunnel的高性能和高并发。

多协程，可以充分使用多个cpu核心进行并行计算，发挥多核cpu的最大性能。

无锁模式的采用最大限度的减少了数据竞争的发生，和因为数据竞争带来的安全问题。

golang并没有在语言层面提供完整的内存安全性保证。golang使用gc管理内存，提供了部分内存安全性，但是依然可能因为数据竞争出现内存安全问题。

基于go语言开发，使用go默认net库。Go 语言层面采用阻塞+多协程模式，进行网络通信。

网络模型： 因为go的默认net库，底层基于非阻塞+多路复用模型（windows iocp、linux epoll）， 所以gstunnel实质模型为非阻塞+多路复用模型。保证了gstunnel网络通信的高性能。

项目采用AES加密算法，使用动态的对称密钥进行加密。默认情况下，每隔一分钟密钥就会进行更新，使用新密钥替换旧密钥。动态密钥机制，将会大幅增加攻击者破解加密数据的难度，提供更好的安全性。

动态密钥机制，会带来轻微的性能损耗（小于5%)，为了更高的安全性，这样的成本支出是合理并且必要的。

为了保证更高的安全性，项目采用基于硬件的强随机数生成器。

为了更好的性能，使用protobuf作为序列化格式。（3.8.1及之前的版本使用json作为序列化格式。使用protobuf格式后，性能为json版的4倍。）


不建议将gstunnel作为vpn（openvpn、ipsec等）的完全替代产品使用。

Gstunnel只是一个轻量级的网络加密管道，只提供有限的安全性，安全性低于openvpn和ipsec等主流的vpn产品，无法替代主流的vpn。

gstunnel相对于传统vpn，也有一些安全优势。
gstunnel只是普通的用户态程序，不需要专属的vpn驱动程序（一般为内核态驱动，具有最高权限）。所以即使gstunnel出现严重的安全漏洞，一般情况下，也不会危害到整个操作系统的安全。
gstunnel的动态密钥，比传统vpn使用的静态密钥具有更好的安全性。

gstunnel加密隧道可以是长连接，也可以是短连接，具体表现取决于所承载业务连接特性。如果业务连接是长连接，加密隧道也会保持长连接，如果业务连接是短连接，加密隧道就会是短连接。

gstunnel client的一个加密隧道连接，在mt模式下，将产生4个处理协程；在非mt模式下将产生2个处理协程。

gstunnel client的一个加密隧道连接，对应一个单独独享的加密密钥，每个加密隧道连接使用不同的加密密钥。
如果gstunnel client有10个加密隧道连接，就会有十个不的加密密钥，每个加密密钥单独负责一个加密隧道连接的加密解密。

支持的平台：
windows、linux、mac os

支持的应用：

http proxy（squid3等）、email、socks 5 proxy等基于tcp开发的应用。

注意：项目存在一些bug，暂时没有修复。这些bug并不影响正常使用。

gstunnel分为client和server两部分。

gstunnel 基于aes进行数据加密。

流程示意：

网络中，a到b的网络通信。

a-->b

使用gstunnel 后， a到b的网络通信。

a-->gstunnel client -->gstunnel server -->b

gstunnel 为a、b之间的网络通信提供了一个加密层。

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

可执行文件，接受基于命令行的参数输入（不推荐使用）和基于配置文件（json）的参数设置。
推荐使用配置文件（json）配置参数。
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
命令行格式:

可执行文件名 监听地址 目标地址 aes密码

注意：aes密码只能是16、24、32字节大小的字符串。

举例说明：

Linux bash：

user@ubuntu:~$ ./gstunnel_client 127.0.0.1:3128 1.2.3.4:43210 “1234567890123456“

user@ubuntu:~$ ./gstunnel_server 1.2.3.4:43210 1.2.3.4:3128 “1234567890123456“

注意：请保证client在linux系统中为可执行文件。是否是可执行文件，请查看client文件的文件属性。

Windows cmd：

C:> ./gstunnel_client.exe 127.0.0.1:3128 1.2.3.4:43210 “1234567890123456“

C:> ./gstunnel_server.exe 1.2.3.4:43210 1.2.3.4:3128 “1234567890123456“

日志

gstunnel在工作目录下自动生成日志文件，日志文件记录gstunnel运行时产生的错误信息。
日志文件名为gstunnel_client.err.log、gstunnel_server.err.log。

项目地址：https://github.com/ypcd/gstunnel

项目基于GPLv3协议开源。

---------------------------------------------
This is a secure network tunnel.

Note: It is recommended to use the project source code of version 2.7 or higher. The source code lower than 2.7 has serious security vulnerabilities.

[Introduction of gstunnel]

gstunnel is a high-performance, high-concurrency cross-platform lightweight security network encryption pipeline developed based on the go language, and it supports the tcp protocol.

The project adopts multi-coroutine and lock-free mode to ensure the high performance and high concurrency of gstunnel.

Multi-coroutine, you can make full use of multiple CPU cores for parallel computing, and give full play to the maximum performance of multi-core CPUs.

The adoption of the lock-free mode minimizes the occurrence of data competition and the security problems caused by data competition.

Golang does not provide complete memory safety guarantees at the language level. Golang uses gc to manage memory and provides some memory security, but memory security problems may still occur due to data competition.

Based on go language development, using go default net library. The Go language layer uses the blocking + multi-coroutine mode for network communication.

Network model: Because the default net library of go, the bottom layer is based on non-blocking + multiplexing model (windows iocp, linux epoll), so the essential model of gstunnel is non-blocking + multiplexing model. Ensure the high performance of gstunnel network communication.

The project uses the AES encryption algorithm and uses a dynamic symmetric key for encryption. By default, the key is updated every minute, replacing the old key with a new key. The dynamic key mechanism will greatly increase the difficulty for attackers to crack encrypted data and provide better security.

The dynamic key mechanism will bring a slight performance loss (less than 5%). For higher security, such a cost is reasonable and necessary.

In order to ensure higher security, the project uses a strong random number generator based on hardware.

For better performance, use protobuf as the serialization format. (Version 3.8.1 and earlier use json as the serialization format. After using the protobuf format, the performance is 4 times that of the json version.)

It is not recommended to use gstunnel as a complete substitute for VPN (openvpn, ipsec, etc.).

Gstunnel is just a lightweight network encryption pipeline, which only provides limited security. The security is lower than that of mainstream VPN products such as openvpn and ipsec, and cannot replace mainstream VPNs.

Compared with traditional VPN, gstunnel also has some security advantages. gstunnel is just an ordinary user-mode program and does not require a dedicated vpn driver (usually a kernel-mode driver with the highest authority). Therefore, even if gstunnel has serious security vulnerabilities, under normal circumstances, it will not endanger the security of the entire operating system. The dynamic key of gstunnel has better security than the static key used by traditional VPN.

The gstunnel encrypted tunnel can be a long connection or a short connection, and the specific performance depends on the connection characteristics of the carried service. If the business connection is a long connection, the encrypted tunnel will also remain a long connection. If the business connection is a short connection, the encrypted tunnel will be a short connection.

An encrypted tunnel connection of the gstunnel client, in the mt mode, will generate 4 processing goroutines; in the non-mt mode, it will generate 2 processing goroutines.

An encrypted tunnel connection of the gstunnel client corresponds to a separate and exclusive encryption key, and each encrypted tunnel connection uses a different encryption key. If the gstunnel client has 10 encrypted tunnel connections, there will be ten different encryption keys, and each encryption key is solely responsible for the encryption and decryption of an encrypted tunnel connection.

Supported platforms: windows, linux, mac os

Supported applications:

HTTP proxy (squid3, etc.), email, socks 5 proxy and other applications developed based on tcp.

Note: There are some bugs in the project, which have not been fixed for the time being. These bugs do not affect normal use.

gstunnel is divided into two parts: client and server.

gstunnel encrypts data based on aes.

Process schematic:

In the network, the network communication from a to b.

a-->b

After using gstunnel, the network communication from a to b.

a-->gstunnel client -->gstunnel server -->b

gstunnel provides an encryption layer for the network communication between a and b.

Make the communication data of a and b become encrypted data, so that the third party cannot know the communication content of a and b. So as to ensure the safety of a and b network communication.

Instructions:

It can be installed through the "go get" tool.

Or after downloading the project source code, copy it to the "$GOPATH\src" directory.

Use the command line tool to compile or install the project source code.

"Go build gstunnel_client" "go build gstunnel_server"

"Go install gstunnel_client" "go install gstunnel_server"

At this time you get two executable files gstunnel_client and gstunnel_server.

3.8.1 and earlier versions, if there is a problem in compiling the source code, please try to enter the command "set GO111MODULE=off" to turn off the go module function.

Configuration parameter

Executable file, accepts command line-based parameter input (not recommended) and configuration file (json)-based parameter settings. It is recommended to use a configuration file (json) to configure the parameters. The configuration file name of the client: config.client.json The configuration file name of the server: config.server.json Configuration file parameters:
```
type GsConfig struct {
Listen string
Servers []string
Key string
Debug bool
Tmr_display_time int
Tmr_changekey_time int
Mt_model bool
}
```
```
Required parameters
listen: listening address (string)
servers: target address (string array)
key: aes encryption key (string)

Optional parameters
debug: whether to enable debug mode (true or false)
Tmr_display_time Set the interval time (in seconds) of information output to the standard output stream
Tmr_changekey_time Set how long it takes for the dynamic key to be changed (unit: second)
Mt_model Whether to enable multi-coroutine mode in the main logic module (true or false)
```
```
Example configuration file:

{"listen": "127.0.0.1:33128", "servers": ["127.0.0.1:10036"], "key": "12345678901234567890123456789012"}

listen: listening address
servers: destination address
key: aes encryption key
```
Command line format:

Executable file name Listening address Target address aes password

Note: the aes password can only be a string of 16, 24, 32 bytes.

for example:

Linux bash:

user@ubuntu:~$ ./gstunnel_client 127.0.0.1:3128 1.2.3.4:43210 "12345678901234567890123456789012"

user@ubuntu:~$ ./gstunnel_server 1.2.3.4:43210 1.2.3.4:3128 "12345678901234567890123456789012"

Note: Please ensure that the client is an executable file in the Linux system. Whether it is an executable file, please check the file attributes of the client file.

Windows cmd:

C:> ./gstunnel_client.exe 127.0.0.1:3128 1.2.3.4:43210 "12345678901234567890123456789012"

C:> ./gstunnel_server.exe 1.2.3.4:43210 1.2.3.4:3128 "12345678901234567890123456789012"

Log

gstunnel automatically generates a log file in the working directory, and the log file records error messages generated when gstunnel is running. The log file names are gstunnel_client.err.log, gstunnel_server.err.log.

Project address: https://github.com/ypcd/gstunnel

The project is open source based on the GPLv3 agreement.
