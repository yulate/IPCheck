# IPCheck



## 0x01 基础功能



- [x] 检测是否为单位IP
- [x] 检测端口是否开放
- [ ] IP反查域名



## 0x02 基础使用

### 默认页面

![image-20230622180255980](https://c.img.dasctf.com/images/2023622/1687428177329-3ee6908b-8923-4034-842f-9b69b148eece.png)

### run指令

启动分析

```
PS D:\ack\codes\go\IPCheck> .\IPCheck.exe run -h
启动分析
        分析等级：
                low     仅扫描80 443端口
                midden  扫描80, 443, 7000, 8080, 8081, 8443 端口
                high    扫描21, 22, 23, 80, 81, 82, 88, 8000, 8888, 888, 443, 8443, 5000, 7000 端口

Usage:
  IPCheck run [flags]

Flags:
  -h, --help            help for run
  -l, --level string    检测等级 (default "low")
  -p, --path string     需检测ip文件位置 (default "./ip.txt")
  -t, --thread string   扫描线程数 (default "10")
```

根据需要指定参数

```
.\IPCheck.exe run -t 300
```

开启300个协程 扫描5096个独立IP 耗时8s

![image-20230622181950619](C:\Users\yulate\AppData\Roaming\Typora\typora-user-images\image-20230622181950619.png)