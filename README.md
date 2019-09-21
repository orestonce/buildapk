# buildapk
 golang 版本的buildapk : https://github.com/yipianfengye/buildapk
# 参考
* [x] 合并zip: https://github.com/rsc/zipmerge
# 用法
````
go get github.com/orestonce/buildapk/cmd
go build -o buildapk github.com/orestonce/buildapk/cmd
buildapk -apk test_v1.0.0.apk -channelList channel.txt
````
# 简单点, 也可以...
````
zip -u test_v1.0.0.apk META-INF/channel_uc.txt
````
