package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

type config struct {
	Api string `json:"api"`
}

func loadConfig(dev bool) (*config, error) {
	var configfp string
	if !dev {
		runp, err := os.Executable()
		if err != nil {
			return nil, err
		}
		rund := filepath.Dir(runp)
		configfp = rund + "/config.json"
	} else {
		configfp = "/config.json"
	}
	f, err := os.Open(configfp)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg config
	err = json.NewDecoder(f).Decode(&cfg)
	return &cfg, err
}

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Text string `json:"text"`
}

func HttpRequest(api_url string) (*Response, error) {
	var response Response
	resp, err := http.Get(api_url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func urlGen(to string, from string, text string, base string) string {
	param := url.Values{}
	param.Set("to", to)
	if from != "" {
		param.Set("from", from)
	}
	param.Set("text", url.QueryEscape(text))
	return base + "?" + param.Encode()
}
