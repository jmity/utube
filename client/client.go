package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

// TODO: scrape videoDetails json instead

const (
	videoURL     = "https://www.youtube.com/watch?v=%s"
	titleRegex   = `\<title\>(.+)\<\/title\>`
	secondsRegex = `"lengthSeconds":"([0-9]+)"`
)

var (
	titleRegexp   *regexp.Regexp
	secondsRegexp *regexp.Regexp
)

func init() {
	var err error
	titleRegexp, err = regexp.Compile(titleRegex)
	if err != nil {
		panic("failed to compile title regex")
	}

	secondsRegexp, err = regexp.Compile(secondsRegex)
	if err != nil {
		panic("failed to compile seconds regex")
	}
}

// UtubeClient Client for aquiring metada about youtube videos
type UtubeClient interface {
	GetMetadta(id string) (*VideoMetadata, error)
}

func New() UtubeClient {
	return &utubeClient{}
}

type VideoMetadata struct {
	VideoLength time.Duration
	VideoName   string
}

type utubeClient struct {
}

func (u *utubeClient) GetMetadta(id string) (*VideoMetadata, error) {
	url := fmt.Sprintf(videoURL, id)
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get video metadata for video with id %s", id)
	}

	rawPage, err := ioutil.ReadAll(resp.Body)

	submatches := titleRegexp.FindAllStringSubmatch(string(rawPage), -1)

	if len(submatches) == 0 || len(submatches[0]) < 2 {
		return nil, fmt.Errorf("no title found for video %s", id)
	}

	title := submatches[0][1]

	submatches = secondsRegexp.FindAllStringSubmatch(string(rawPage), -1)

	if len(submatches) == 0 || len(submatches[0]) < 2 {
		return nil, fmt.Errorf("no time found for video %s", id)
	}

	seconds, err := strconv.Atoi(submatches[0][1])
	if err != nil {
		return nil, fmt.Errorf("failed to get time for video with duration %d", seconds)
	}

	return &VideoMetadata{VideoName: title, VideoLength: time.Second * time.Duration(seconds)}, nil
}
