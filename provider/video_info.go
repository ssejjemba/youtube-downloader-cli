package provider

import (
	"io/ioutil"
	"log"
	"net/http"
)

func (y *Youtube) getVideoInfo() error {
	url := "http://youtube.com/get_video_info?video_id=" + y.videoID
	log.Printf("url: %s", url)

	response, err := http.Get(url)
	if err != nil {
		return err
	}

	defer response.Body.Close()
	if response.StatusCode != 200 {
		return err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	y.videoInfo = string(body)
	return nil
}