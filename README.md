# gstunnel
This is a secure network tunnel.

【gstunnel介绍】
 
gstunnel 是基于go 语言开发的高性能、高并发的轻量级安全网络加密管道，支持tcp协议。
 
项目采用多协程和无锁模式，保证了gstunnel的高性能和高并发。无锁模式的采用也最大限度的减少了数据竞争的发生，和因为数据竞争带来的安全问题。golang并没有在语言层面提供完整的内存安全性保证。golang使用gc管理内存，提供了部分内存安全性，但是依然可能因为数据竞争出现内存安全问题。
 
基于go语言开发，使用go默认net库。Go 语言层面采用阻塞+多协程模式，进行网络通信。
网络模型： 因为go的默认net库，底层基于非阻塞+多路复用模型（windows iocp、linux epoll）， 所以gstunnel实质模型为非阻塞+多路复用模型。保证了gstunnel网络通信的高性能。
 
项目采用AES加密算法，使用动态的对称密钥进行加密。默认情况下，每隔一分钟密钥就会进行更新，使用新密钥替换旧密钥。动态密钥机制，将会大幅增加攻击者破解加密数据的难度，提供更好的安全性。
动态密钥机制，会带来轻微的性能损耗（小于5%)，为了更高的安全性，这样的成本支出是合理并且必要的。
为了保证更高的安全性，项目采用基于硬件的强随机数生成器。
 
项目存在一些bug，暂时没有修复。这些bug并不影响正常使用。
 
不建议将gstunnel作为vpn（openvpn、ipsec等）的完全替代产品使用。
Gstunnel只是一个轻量级的网络加密管道，只提供有限的安全性，安全性低于openvpn和ipsec等主流的vpn产品，无法替代主流的vpn。
 
支持的应用：
http proxy（squid3等）、email、socks 5 proxy等基于tcp开发的应用。
 
-------------------------------------------------------------------------------
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
“go build gstunnel_client”		“go build gstunnel_server”
“go install gstunnel_client”		“go install gstunnel_server”
这时你得到了两个可执行文件gstunnel_client、gstunnel_server。
编译源代码如果出现问题，请尝试输入命令“set GO111MODULE=off”，关闭go模块功能。
 
可执行文件，接受基于命令行的参数输入。
格式:
可执行文件名 监听地址 目标地址 aes密码
注意：aes密码只能是16、24、32字节大小的字符串。
举例说明：
Linux bash：
root@ubuntu:~# ./gstunnel_client 127.0.0.1:3128 1.2.3.4:43210 “1234567890123456“
root@ubuntu:~# ./gstunnel_server 1.2.3.4:43210 1.2.3.4:3128 “1234567890123456“
注意：请保证client在linux系统中为可执行文件。是否是可执行文件，请查看client文件的文件属性。
Windows cmd：
C:> ./gstunnel_client 127.0.0.1:3128 1.2.3.4:43210 “1234567890123456“
C:> ./gstunnel_server 1.2.3.4:43210 1.2.3.4:3128 “1234567890123456“
项目地址：https://github.com/ypcd/gstunnel
项目基于GPLv3协议开源。



[Introduction of gstunnel]
 
gstunnel is a high-performance, high-concurrency, lightweight security network encryption pipeline developed based on the go language, and supports the tcp protocol.
 
The project adopts multi-coroutine and lock-free mode to ensure the high performance and high concurrency of gstunnel. The adoption of the lock-free mode also minimizes the occurrence of data competition and the security problems caused by data competition. Golang does not provide complete memory safety guarantees at the language level. Golang uses gc to manage memory and provides some memory security, but memory security problems may still occur due to data competition.
 
Based on go language development, using go default net library. The Go language layer uses the blocking + multi-coroutine mode for network communication.
Network model: Because the default net library of go, the bottom layer is based on non-blocking + multiplexing model (windows iocp, linux epoll), so the essential model of gstunnel is non-blocking + multiplexing model. Ensure the high performance of gstunnel network communication.
 
The project uses the AES encryption algorithm and uses a dynamic symmetric key for encryption. By default, the key is updated every minute, replacing the old key with a new key. The dynamic key mechanism will greatly increase the difficulty for attackers to crack encrypted data and provide better security.
The dynamic key mechanism will bring a slight performance loss (less than 5%). For higher security, such a cost is reasonable and necessary.
In order to ensure higher security, the project uses a strong random number generator based on hardware.
 
There are some bugs in the project, which have not been fixed for the time being. These bugs do not affect normal use.
 
It is not recommended to use gstunnel as a complete substitute for VPN (openvpn, ipsec, etc.).
Gstunnel is just a lightweight network encryption pipeline, which only provides limited security. The security is lower than mainstream VPN products such as openvpn and ipsec, and cannot replace mainstream VPNs.
 
Supported applications:
HTTP proxy (squid3, etc.), email, socks 5 proxy and other applications developed based on tcp.
 
-------------------------------------------------- -----------------------------
gstunnel is divided into two parts: client and server.
gstunnel encrypts data based on aes.
Process schematic:
In the network, the network communication from a to b.
a-->b
After using gstunnel, the network communication from a to b.
a-->gstunnel client -->gstunnel server -->b
gstunnel provides an encryption layer for the network communication between a and b.
Make the communication data of a and b become encrypted data, so that the third party cannot know the communication content of a and b. So as to ensure the security of a and b network communication.
 
Instructions:
It can be installed via the "go get" tool.
Or after downloading the project source code, copy it to the "$GOPATH\src" directory.
 
Use the command line tool to compile or install the project source code.
"Go build gstunnel_client" "go build gstunnel_server"
"Go install gstunnel_client" "go install gstunnel_server"
At this time you get two executable files gstunnel_client and gstunnel_server.
If there is a problem in compiling the source code, please try to enter the command "set GO111MODULE=off" to turn off the go module function.
 
Executable file, accepts command line-based parameter input.
format:
Executable file name Listening address Target address aes password
Note: the aes password can only be a string of 16, 24, 32 bytes.
for example:
Linux bash:
root@ubuntu:~# ./gstunnel_client 127.0.0.1:3128 1.2.3.4:43210 "1234567890123456"
root@ubuntu:~# ./gstunnel_server 1.2.3.4:43210 1.2.3.4:3128 "1234567890123456"
Note: Please ensure that the client is an executable file in the Linux system. Whether it is an executable file, please check the file attributes of the client file.
Windows cmd:
C:> ./gstunnel_client 127.0.0.1:3128 1.2.3.4:43210 "1234567890123456"
C:> ./gstunnel_server 1.2.3.4:43210 1.2.3.4:3128 "1234567890123456"
Project address: https://github.com/ypcd/gstunnel
The project is open source based on the GPLv3 agreement.
