build	:
	go build main.go
	mkdir -p bin
	mv main bin/buildapk
run	:
	./bin/buildapk -apkInFile test_v1.0.0.apk -channelNameListFile channel.txt
