package provider

import (
	"fmt"
	"io"
	"log"
	"os"
)

func NewYoutube() *Youtube {
	return new(Youtube)
}

type stream map[string]string

type Youtube struct {
	StreamList []stream
	VideoID string
	VideoInfo string
	DownloadPercentage chan int64
	contentLength float64
	totalWrittenBytes float64
	downloadLevel float64
}

func (y *Youtube) DecodeURL(url string) error {
	err := y.findVideoId(url)
	if err != nil {
		return fmt.Errorf("findvideoID error=%s", err)
	}

	err = y.getVideoInfo()
	if err != nil {
		return fmt.Errorf("getVideoInfo error=%s", err)
	}

	err = y.parseVideoInfo()
	if err != nil {
		return fmt.Errorf("parse video info failed, err=%s", err)
	}

	return nil
}

func (y *Youtube) StartDownload(dstDir string) error {
	y.DownloadPercentage = make(chan int64, 100)
	//download highest resolution on [0]
	targetStream := y.StreamList[0]
	url := targetStream["url"] + "&signature=" + targetStream["sig"]
	log.Println("Download url=", url)

	targetFile := fmt.Sprintf("%s%c%s.%s", dstDir, os.PathSeparator, targetStream["title"], "mp4")
	//targetStream["title"], targetStream["author"])
	log.Println("Download to file=", targetFile)
	err := y.videoDLWorker(targetFile, url)
	return err
}

func SetLogOutput(w io.Writer) {
	log.SetOutput(w)
}

