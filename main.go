package main

import (
	"archive/zip"
	"bytes"
	"compress/flate"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

func main() {
	var apkInFile string
	var channelName string
	var channelNameListFile string

	flag.StringVar(&apkInFile, `apkInFile`, ``, `input apk filename`)
	flag.StringVar(&channelName, `channelName`, ``, `single channel name`)
	flag.StringVar(&channelNameListFile, `channelNameListFile`, ``, `channel name file, every channel per line`)
	flag.Parse()
	if apkInFile == `` || (channelNameListFile == `` && channelName == ``) {
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

	log.Println(`begin read channel list`)
	channelNameCh := getChannelNameListCh(channelName, channelNameListFile)
	log.Println(`begin unarchive input apk file`)
	m := zipUnArchive(apkInFile)
	threadCount := runtime.NumCPU() * 2
	log.Println(`begin build output apk, thread count`, threadCount)

	var wg sync.WaitGroup
	for i := 0; i < threadCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for channelName := range channelNameCh {
				after := zipArchiveWithExtraEmptyFile(m, `META-INF/channel_`+channelName)
				outFileName := filepath.Join(apkOutDir, apkInNamebase+"_channel_"+channelName+".apk")
				err = ioutil.WriteFile(outFileName, after, 0666)
				handleError(err)
				log.Println(`success for channel `, channelName)
			}
		}()
	}
	wg.Wait()
}

func getChannelNameListCh(channelName string, channelNameListFile string) (ch <-chan string) {
	chIn := make(chan string, 100)
	go func() {
		putChannelName := func(name string) {
			name = strings.TrimSpace(name)
			if name == `` {
				return
			}
			chIn <- name
		}
		putChannelName(channelName)
		if channelNameListFile != `` {
			content, err := ioutil.ReadFile(channelNameListFile)
			handleError(err)
			for _, line := range bytes.Split(content, []byte{'\n'}) {
				putChannelName(string(line))
			}
		}
		close(chIn)
	}()
	return chIn
}

func zipUnArchive(inFileName string) (out map[string][]byte) {
	content, err := ioutil.ReadFile(inFileName)
	handleError(err)
	r, err := zip.NewReader(bytes.NewReader(content), int64(len(content)))
	handleError(err)
	out = map[string][]byte{}
	for _, f := range r.File {
		rc, err := f.Open()
		handleError(err)
		b, err := ioutil.ReadAll(rc)
		err = rc.Close()
		handleError(err)
		out[f.Name] = b
	}
	return
}

func zipArchiveWithExtraEmptyFile(m map[string][]byte, emptyFileName string) (content []byte) {
	_, ok := m[emptyFileName]
	if ok {
		panic(`Unexpected file in archive: ` + emptyFileName)
	}
	buf := &bytes.Buffer{}
	w := zip.NewWriter(buf)
	w.RegisterCompressor(zip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
		return flate.NewWriter(out, flate.BestCompression)
	})
	appendOneFile := func(filename string, content []byte) {
		f, err := w.Create(filename)
		handleError(err)
		_, err = f.Write(content)
		handleError(err)
	}
	for filename, content := range m {
		appendOneFile(filename, content)
	}
	appendOneFile(emptyFileName, nil)
	err := w.Close()
	handleError(err)
	return buf.Bytes()
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
