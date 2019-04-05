package config

import (
	"io/ioutil"
	"net/http"
	"net/url"
)

func fetch(location string) ([]byte, error) {
	u, err := url.Parse(location)
	if err != nil {
		return fetchFile(location)
	}

	switch u.Scheme {
	case "http", "https":
		return fetchHTTP(u)
	default:
		return fetchFile(u.Path)
	}
}

func fetchFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

func fetchHTTP(u *url.URL) ([]byte, error) {
	cl := http.Client{
		Timeout: timeout,
	}
	resp, err := cl.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
