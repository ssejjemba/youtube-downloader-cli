package provider

import (
	"errors"
	"log"
	"regexp"
	"strings"
)

func (y *Youtube) findVideoId(url string) error{
	videoId := url
	
	if strings.Contains(videoId, "youtu") || strings.ContainsAny(videoId, "\"?&/<%=") {
		// This is a url and we need to strip off the id
		// These Regex have been inspired from https://github.com/kkdai/youtube
		regex_list := []*regexp.Regexp{
			regexp.MustCompile(`(?:v|embed|watch\?v)(?:=|/)([^"&?/=%]{11})`),
			regexp.MustCompile(`(?:=|/)([^"&?/=%]{11})`),
			regexp.MustCompile(`([^"&?/=%]{11})`),
		}

		for _, reg := range regex_list {
			if is_match := reg.MatchString(videoId); is_match {
				sub_strings := reg.FindStringSubmatch(videoId)
				videoId = sub_strings[1]
			}
		}
	}

	log.Printf("Found video id: '%s'", videoId)
	y.videoID = videoId

	return validateVideoId(videoId)
}

func validateVideoId(id string) error{
	
	if strings.ContainsAny(id, "?&/<%=") {
		return errors.New("invalid characters in video id")
	}
	if len(id) < 10 {
		return errors.New("the video id must be at least 10 characters long")
	}
	return nil
}