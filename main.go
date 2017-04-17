package main

import (
	"bufio"
	"encoding/json"
  "flag"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"io/ioutil"
	"net/url"
	"os"
)

type ApiConf struct {
	ConsumerKey       string `json:"consumer_key"`
	ConsumerSecret    string `json:"consumer_secret"`
	AccessToken       string `json:"access_token"`
	AccessTokenSecret string `json:"access_token_secret"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func TimeLine(api anaconda.TwitterApi, v url.Values, tlchan chan anaconda.Tweet){
	twitterStream := api.UserStream(v)
	for {
		x := <-twitterStream.C
		switch tweet := x.(type) {
		case anaconda.Tweet:
			tlchan <- tweet
		default:
		}
	}
}

func Post(api anaconda.TwitterApi, v url.Values, poschan chan anaconda.Tweet){
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		tweet, _ := api.PostTweet(scanner.Text(),v)
		poschan <- tweet
	}
}

func Controller(tlchan chan anaconda.Tweet, postchan chan anaconda.Tweet){
	for{
		select {
		case tl := <- tlchan:
			if tl.Retweeted == true {
	      fmt.Println(">> "+tl.User.ScreenName)
				tl = *tl.RetweetedStatus
				fmt.Println("from "+tl.User.ScreenName)
				fmt.Println("RT:"+tl.Text)
			}else{
				fmt.Println(">> "+tl.User.ScreenName)
				fmt.Println(tl.Text)
			}
			fmt.Println("--------------")
		case post := <- postchan:
			fmt.Println("post!:",post.Text)
			fmt.Println("--------------")
		default:
		}
	}
}

func main() {
	tlchan := make(chan anaconda.Tweet)
	postchan := make(chan anaconda.Tweet)
	var apiConf ApiConf
	{
		apiConfPath := flag.String("conf", "config.json", "API Config File")
		flag.Parse()
		data, err_file := ioutil.ReadFile(*apiConfPath)
		check(err_file)
		err_json := json.Unmarshal(data, &apiConf)
		check(err_json)
	}
	anaconda.SetConsumerKey(apiConf.ConsumerKey)
	anaconda.SetConsumerSecret(apiConf.ConsumerSecret)
	api := anaconda.NewTwitterApi(apiConf.AccessToken, apiConf.AccessTokenSecret)
	v := url.Values{}
  go TimeLine(*api,v,tlchan)
	go Post(*api,v,postchan)
	Controller(tlchan,postchan)
}
