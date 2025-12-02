# 备案通用模板

## 一、修改 conf/hosts.txt 可以定义对应域名显示的数据

### 1、每一行描述标题、域名列表，版权信息和备案号四类信息。每类信息以竖线(|) 分隔

### 2、域名列表以逗号(,)分隔

### 3、每个类别信息字符串前后两端可以加任意空格，空格会被自动清理。每类信息如果想不显示，使用空格代替即可，对应的竖线不能少。信息为空，对应页面位置的信息显示为空白

### 4、实现了自动加载机制。如果配置文件或者页面模板被修改，会自动加载，所以修改信息不用重启应用

### 5、以井号(#) 开始的行为注释，被忽略

## 二、修改 template/template.tmpl 可以设置页面

## 三、访问日志文件保存在 log/acess.txt 中

## 四、性能优化

### 1、取消了文件修改后的自动加载

### 2、使用 /reaload 接口提供热加载数据

### 3、使用scratch构建镜像，将镜像从 19.6MB 降到了 11.3MB

### 4、运行容器时挂载主机的时区文件:经过测试，可用

```docker
docker run -p 8801:8901 -v /etc/localtime:/etc/localtime:ro -v /etc/timezone:/etc/timezone:ro --label remarks="通用备案模板" --label 创建日期="2025-10-31" --label version="1.1" -d --name beian beian-universal-template
```

## 五、定制镜像

可以利用备案通用模板的镜像，制定自己的个性化镜像。
创建一个dockerfile，将自己的域名备案信息和配置文件放到dockerfile同一个路径的conf目录下。

dockerfile如下：

```dockerfile
FROM beian-universal-template:latest
COPY conf/ /conf
```

然后，构建

```sh
docker build -t custom-beian-instance .
```

```sh
docker run -p 8901:8901 -v /etc/localtime:/etc/localtime:ro -v /etc/timezone:/etc/timezone:ro --label description="beian备案镜像实例" --label 启动日期="2025-11-29" -d --name beian custom-beian-instance
```

## 六、开发计划

### 1、添加日志文件 - 完成 DOEN

添加访问日志文件，可以记录访问 `www.yourdomain.com` 下的嗅探客户端 IP 。

### 2、设置刷新配置文件 - 完成 DOEN

提供一个时间设置，通过这个时间控制 /reload 是否允许执行。
防止有人恶意请求 /reload 接口。
在设置网站初期，可以将这个时间值设置的短一点，比如允许 2 秒刷新一次。
一旦设置完成，不再需要热加载配置数据的时候，将此数值设置大一些即可。

### 3、增加刷新Ip的限制 - 完成 DOEN

防止恶意热加载，增加IpFilter。记录第一次访问应用的Ip地址，只允许第一次访问的Ip访问此接口。
