# Wall-E 钉钉资产监控


                               _-
                _w*P~M        y~~F*m__
             p*T^  _, l  ___ _' p    ~MNq_
         ,w*~  _   ~!  #RMA0#f  ~   ,_,  ^7#g          ██╗    ██╗ █████╗ ██╗     ██╗                   ███████╗
      _*~   _pF~~4q     M0#0Q    _P^ ]00q    M         ██║    ██║██╔══██╗██║     ██║                   ██╔════╝
      "     NNp__#0&    #0MB0    R&gyM000S    4        ██║ █╗ ██║███████║██║     ██║         █████╗    █████╗
     f qf  ]0N0000NN    l00Bf    N0MM00M0#  #6 I       ██║███╗██║██╔══██║██║     ██║         ╚════╝    ██╔══╝
           ]0MM0N0MMn   jb###    MN0#M0MM#  ^  ]       ╚███╔███╔╝██║  ██║███████╗███████╗              ███████╗
           "00M000M0   j00M#0q   RN0#0#0#'     p        ╚══╝╚══╝ ╚═╝  ╚═╝╚══════╝╚══════╝              ╚══════╝
     t      ^M00NMM   g0N00#N0k,  "#0#M~       !
     ^,   q   ^^^  _pM         ^M&q_     4r  _P                                             
      "q  ~    __aM^               ^~Mmqq_,g#~                                              
        Mqq,gv@^                         ^^


## 功能
- 权限管理
- 资产管理
- 端口监控
- 新增端口推送

## 特点
这个系统采用了客户端-服务器结构，并利用RabbitMQ进行消息队列处理，以确保任务的可靠性运行和实现分布式部署，同时根据监控结果进行消息推送。

## 产品逻辑
![Wall.png](Image%2FArchitecture.png)

## 部署方式
### 数据库
```
sudo docker run -itd \
  --name mysql-wall \
  -p 3306:3306 \
  -e MYSQL_ROOT_PASSWORD=changeMe \
  mysql
```
### MQ
建议使用Docker部署
```
sudo docker run -d \
  --name rabbitmq-wall \
  -p 5672:5672 \
  -p 15672:15672 \
  -e RABBITMQ_DEFAULT_USER=admin \
  -e RABBITMQ_DEFAULT_PASS=changeMe \
  -e RABBITMQ_DEFAULT_VHOST=wall \
  --restart=always \
  rabbitmq:management
```
建议使用Docker部署
## 钉钉
参考钉钉开发文档，申请应用能力并添加机器人。  
[https://open.dingtalk.com/document/orgapp/the-creation-and-installation-of-the-application-robot-in-the](https://open.dingtalk.com/document/orgapp/the-creation-and-installation-of-the-application-robot-in-the)

记录应用凭证:
- Client ID
- Client Secret

添加机器人后，拉该机器人进入企业内部群，查看机器人Webhook地址。 （因Stream流的方式无法主动推送告警）
[https://open.dingtalk.com/document/orgapp/webhook-robot](https://open.dingtalk.com/document/orgapp/webhook-robot)

记录Webhook地址中的access_token：
- Webhook

## 监控程序部署

### 安装基础功能

- Centos
```
yum install nmap masscan -y
yum install libpcap-devel
```
- Ubuntu
```
sudo apt-get install nmap masscan -y
sudo apt-get install libpcap-dev
sudo ln /usr/libx86_64-linux-gnu/libpcap.so.0.8 /usr/libx86_64-linux-gnu/libpcap.so.1
```
- Windows
    - [Nmap下载地址](https://nmap.org/download.html#windows)
    - [Masscan下载地址](https://github.com/zan8in/masscan/blob/main/bin/masscan-win/masscan.exe)

### 下载部署
[Wall_E项目地址](https://github.com/Mustard404/Wall_E/releases)     
分别修改Server端和Agent端内的config.ini, 执行即可。


