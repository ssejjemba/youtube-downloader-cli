package provider

import (
	"fmt"
	"log"
	"os"
)

func NewYoutube() *Youtube {
	return new(Youtube)
}

type stream map[string]string

type Youtube struct {
	StreamList []stream
	videoID string
	videoInfo string
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
	//download highest resolution on [0]
	targetStream := y.StreamList[0]
	url := targetStream["url"] + "&signature=" + targetStream["sig"]
	log.Println("Download url=", url)

	targetFile := fmt.Sprintf("%s%v%s.%s", dstDir, os.PathSeparator, targetStream["title"], "mp4")
	//targetStream["title"], targetStream["author"])
	log.Println("Download to file=", targetFile)
	err := videoDLWorker(targetFile, url)
	return err
}

