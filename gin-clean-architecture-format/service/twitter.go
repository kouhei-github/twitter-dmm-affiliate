package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Image struct {
	ImageType string `json:"image_type"`
	W         int    `json:"w"`
	H         int    `json:"h"`
}

type ImageIdRequest struct {
	MediaId          int    `json:"media_id"`
	MediaIdString    string `json:"media_id_string"`
	Size             int    `json:"size"`
	ExpiresAfterSecs int    `json:"expires_after_secs"`
	Image            Image
}

func (oauth *OAuth1) GetMediaId(imagePath string) (*ImageIdRequest, error) {
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	file, errFile2 := os.Open(imagePath)
	defer file.Close()
	part2, errFile2 := writer.CreateFormFile("media", filepath.Base(imagePath))
	_, errFile2 = io.Copy(part2, file)
	if errFile2 != nil {
		return &ImageIdRequest{}, errFile2
	}
	err := writer.Close()
	if err != nil {
		return &ImageIdRequest{}, err
	}

	method := http.MethodPost
	url := "https://upload.twitter.com/1.1/media/upload.json"

	authHeader := oauth.BuildOAuth1Header(method, url, map[string]string{})
	client := &http.Client{}
	req, _ := http.NewRequest(method, url, payload)
	req.Header.Add("Authorization", authHeader)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		return &ImageIdRequest{}, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return &ImageIdRequest{}, err
	}
	var image ImageIdRequest
	err = json.Unmarshal(body, &image)
	if err != nil {
		return &ImageIdRequest{}, err
	}
	return &image, nil
}

func SeparateArray(number int, images []*ImageIdRequest) [][]*ImageIdRequest {
	var ImageIdRequests [][]*ImageIdRequest
	var small []*ImageIdRequest
	for index, image := range images {
		small = append(small, image)
		idx := index + 1
		if idx%number == 0 {
			fmt.Println(index + 1)
			ImageIdRequests = append(ImageIdRequests, small)
			small = []*ImageIdRequest{}
		}
	}

	return ImageIdRequests
}

func (oauth *OAuth1) PostToTwitterWithAttachment(separateMedias [][]*ImageIdRequest) error {
	length := strconv.Itoa(len(separateMedias))
	tweetId := ""
	for i, medias := range separateMedias {
		var imageIds []string
		for _, media := range medias {
			imageIds = append(imageIds, media.MediaIdString)
		}
		idx := i + 1
		index := strconv.Itoa(idx)
		title := ""
		text := ""
		if idx == 1 {
			title = "ここがすごい" + "\n"
		}
		if title != "" {
			text += title
		}
		text += "(" + index + "/" + length + ")"

		tweet, err := oauth.PostTwitter(text, imageIds, tweetId)
		if err != nil {
			return err
		}
		tweetId = tweet.Data.Id
	}
	return nil
}

type MediaId struct {
	MediaIds []string `json:"media_ids"`
}

type Reply struct {
	InReplyToTweetId string `json:"in_reply_to_tweet_id"`
}

type AttachmentReply struct {
	Text  string  `json:"text"`
	Media MediaId `json:"media"`
	Reply Reply   `json:"reply"`
}

type Attachment struct {
	Text  string  `json:"text"`
	Media MediaId `json:"media"`
}

type Tweet struct {
	Id   string `json:"id"`
	Text string `json:"text"`
}

type PostTwitter struct {
	Data Tweet `json:"data"`
}

func (oauth OAuth1) PostTwitter(text string, attachments []string, tweetId string) (PostTwitter, error) {
	mediaId := MediaId{MediaIds: attachments}
	var jsonValue []byte
	var err error
	if tweetId != "" {
		reply := Reply{InReplyToTweetId: tweetId}
		attachment := AttachmentReply{
			Text:  text,
			Media: mediaId,
			Reply: reply,
		}
		jsonValue, err = json.Marshal(attachment)
		if err != nil {
			myError := MyError{Message: "Jsonに変換できませんでした"}
			return PostTwitter{}, myError
		}
	} else {
		attachment := Attachment{
			Text:  text,
			Media: mediaId,
		}
		jsonValue, err = json.Marshal(attachment)
		if err != nil {
			myError := MyError{Message: "Jsonに変換できませんでした"}
			return PostTwitter{}, myError
		}
	}

	payload := strings.NewReader(string(jsonValue))

	url := "https://api.twitter.com/2/tweets"
	method := "POST"

	authHeader := oauth.BuildOAuth1Header(method, url, map[string]string{})
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		fmt.Println(err)
		return PostTwitter{}, err
	}
	req.Header.Add("Authorization", authHeader)
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return PostTwitter{}, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return PostTwitter{}, err
	}
	var tweet PostTwitter
	err = json.Unmarshal(body, &tweet)
	if err != nil {
		return PostTwitter{}, err
	}
	return tweet, nil
}

type TweetIncludeAuthor struct {
	Id       string `json:"id"`
	AuthorId string `json:"author_id"`
	Text     string `json:"text"`
}

type TwitterSearch struct {
	Data []TweetIncludeAuthor `json:"data"`
}

func (oauth OAuth1) SearchHashTagOnTwitter(searchText string, max int) (TwitterSearch, error) {
	url := "https://api.twitter.com/2/tweets/search/recent?query=" + searchText + "&max_results=" + strconv.Itoa(max) + "&expansions=author_id"
	method := http.MethodGet
	oauth.BuildOAuth1Header(method, url, map[string]string{})
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return TwitterSearch{}, err
	}
	req.Header.Add("Authorization", "Bearer "+oauth.BearerToken)
	req.Header.Add("accept", "application/json")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return TwitterSearch{}, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return TwitterSearch{}, err
	}
	fmt.Println(string(body))
	var tweet TwitterSearch
	err = json.Unmarshal(body, &tweet)
	if err != nil {
		return TwitterSearch{}, err
	}
	return tweet, nil
}
