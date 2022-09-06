package provider

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func (y *Youtube) getVideoInfo() error {
	url := "http://youtube.com/get_video_info?video_id=" + y.VideoID
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

	y.VideoInfo = string(body)
	return nil
}

func (y *Youtube) parseVideoInfo() error {
	answer, err := url.ParseQuery(y.VideoInfo)
	if err != nil {
		return err
	}

	status, ok := answer["status"]
	if !ok {
		err = fmt.Errorf("no response status found in the server's answer")
		return err
	}

	if status[0] == "fail" {
		reason, ok := answer["reason"]
		if ok {
			err = fmt.Errorf("'fail' response status found in the server's answer, reason: '%s'", reason[0])
		} else {
			err = errors.New(fmt.Sprint("'fail' response status found in the server's answer, no reason given"))
		}
		return err
	}

	if status[0] != "ok" {
		err = fmt.Errorf("non-success response status found in the server's answer (status: '%s')", status)
		return err
	}

	// read the streams map
	stream_map, ok := answer["url_encoded_fmt_stream_map"]
	if !ok {
		err = errors.New(fmt.Sprint("no stream map found in the server's answer"))
		return err
	}

	// read each stream
	streams_list := strings.Split(stream_map[0], ",")

	var streams []stream

	for stream_pos, stream_raw := range streams_list {
		stream_qry, err := url.ParseQuery(stream_raw)
		if err != nil {
			log.Println(fmt.Sprintf("An error occured while decoding one of the video's stream's information: stream %d: %s\n", stream_pos, err))
			continue
		}
		var sig string
		if _, exist := stream_qry["sig"]; exist {
			sig = stream_qry["sig"][0]
		}

		stream := stream{
			"quality": stream_qry["quality"][0],
			"type":    stream_qry["type"][0],
			"url":     stream_qry["url"][0],
			"sig":     sig,
			"title":   answer["title"][0],
			"author":  answer["author"][0],
		}
		streams = append(streams, stream)
		log.Printf("Stream found: quality '%s', format '%s'", stream_qry["quality"][0], stream_qry["type"][0])
	}

	y.StreamList = streams
	return nil
}