package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/orestonce/buildapk/zipmerge"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var apkInFile string
	var channelNameListFile string

	flag.StringVar(&apkInFile, `apk`, ``, `input apk filename`)
	flag.StringVar(&channelNameListFile, `channelList`, ``, `channel name file, every channel per line`)
	flag.Parse()
	if apkInFile == `` || channelNameListFile == `` {
		flag.Usage()
		os.Exit(-1)
	}
	apkOutDir := filepath.Join(filepath.Dir(apkInFile), `buildapk_output`)
	err := os.MkdirAll(apkOutDir, 0777)
	handleError(err)

	apkInNamebase := ``
	if strings.HasSuffix(strings.ToLower(apkInFile), `.apk`) {
		apkInNamebase = filepath.Base(strings.TrimSuffix(strings.ToLower(apkInFile), `.apk`))
	} else {
		fmt.Println(`It seems not an apk file: `, apkInFile)
		os.Exit(-1)
	}

	channelNameList := getChannelNameList(channelNameListFile)
	origin, err := ioutil.ReadFile(apkInFile)
	handleError(err)

	for _, channelName := range channelNameList {
		after, err := zipmerge.ZipMerge(origin, map[string][]byte{
			`META-INF/channel_` + channelName: nil,
		})
		handleError(err)
		outFileName := filepath.Join(apkOutDir, apkInNamebase+"_channel_"+channelName+".apk")
		err = ioutil.WriteFile(outFileName, after, 0666)
		handleError(err)
		log.Println(`success build channel`, channelName)
	}
}

func getChannelNameList(channelNameListFile string) (channelList []string) {
	if channelNameListFile != `` {
		content, err := ioutil.ReadFile(channelNameListFile)
		handleError(err)
		for _, line := range bytes.Split(content, []byte{'\n'}) {
			name := strings.TrimSpace(string(line))
			if name == `` {
				continue
			}
			channelList = append(channelList, name)
		}
	}
	return channelList
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
