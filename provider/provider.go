package provider

import "fmt"

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
}