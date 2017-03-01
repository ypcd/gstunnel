# gstunnel
This is a secure network tunnel.

项目简介：

gstunnel 是 基于go 语言开发的一个安全网络管道，支持tcp协议。

gstunnel分为client和server两部分。

gstunnel 基于aes进行数据加密。

流程示意：

网络中，a到b的网络通信。

a-->b

使用gstunnel 后， a到b的网络通信。

a-->gstunnel client -->gstunnel server -->b

gstunnel 为a、b之间的网络通信提供了一个加密层。

使得a、b的通信数据，变为了加密数据，这样第三方就不能获知a、b的通信内容。从而保证了a、b网络通信的安全。

支持的应用：

http proxy（squid3等）、email、socks 5 proxy等基于tcp开发的应用。

使用方法:

进入"gstunnel"目录下，使用"go build server.go"和"go build client.go"分别编译.go文件。

这时你得到了两个可执行文件client、server。

可执行文件，接受基于命令行的参数输入。

格式:

可执行文件名 监听地址 目标地址 aes密码

注意：aes密码只能是16、24、32字节大小的字符串。

举例说明：

Linux bash：

root@ubuntu:~# ./client 127.0.0.1:3128 1.2.3.4:43210 “1234567890123456“

root@ubuntu:~# ./server 1.2.3.4:43210 1.2.3.4:3128 “1234567890123456“

注意：请保证client在linux系统中为可执行文件。是否是可执行文件，请查看client文件的文件属性。

Windows cmd：

C:\> ./client 127.0.0.1:3128 1.2.3.4:43210 “1234567890123456“

C:\> ./server 1.2.3.4:43210 1.2.3.4:3128 “1234567890123456“

项目地址：https://github.com/ypcd/gstunnel

项目基于GPLv3协议开源。




The following from google translation:


Project Description:

Gstunnel is based on the development of a language network gateway, support tcp protocol.

Gstunnel divided into client and server two parts.

Gstunnel is based on aes for data encryption.

Flow diagram:

Network, a to b network communication.

A -> b

After using gstunnel, a to b network traffic.

A -> gstunnel client -> gstunnel server -> b

Gstunnel provides an encryption layer for network traffic between a and b.

So that a, b of communication data, into encrypted data, so that third parties can not know a, b communication content. Thus ensuring a, b network communication security.

Supported applications:

Http proxy (squid3, etc.), email, socks 5 proxy and other applications based on tcp development.

Instructions:

Go to the "gstunnel" directory and compile the .go file using "go build server.go" and "go build client.go".

At this point you get two executable files client, server.

Executable file, accept command-based parameter input.

format:

Executable File Name List Addresses Destination Address aes Password

Note: The aes password can only be a string of 16, 24, 32 bytes.

for example:

Linux bash:

Root @ ubuntu: ~ # ./client 127.0.0.1:3128 1.2.3.4:43210 "1234567890123456"

Root @ ubuntu: ~ # ./server 1.2.3.4:43210 1.2.3.4:3128 "1234567890123456"

Note: Please ensure that the client in the linux system for the executable file. Whether it is an executable file, please see the client file file attributes.

Windows cmd:

C:> ./client 127.0.0.1:3128 1.2.3.4: 43210 "1234567890123456"

C:> ./server 1.2.3.4:43210 1.2.3.4:3128 "1234567890123456"

Project address: https: //github.com/ypcd/gstunnel

The project is based on the GPLv3 protocol.
