package main

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

var tr = &http.Transport{
	MaxIdleConnsPerHost: 12,
}

func One() (Code, error) {
	const url = "http://bgmstore.net/random"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	resp, err := tr.RoundTrip(req)
	if err != nil {
		return "", err
	}

	io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()

	loc := resp.Header.Get("Location")
	code, err := unwrapLoc(loc)

	return Code(code), err
}

// e.g /view/rxBQI/random
// 위 Location 값에서 가운데에 위치한 rxBQI 값을 추출한다.
func unwrapLoc(loc string) (string, error) {
	i := len("/view/")
	if len(loc) < i {
		return "", errors.New("too busy")
	}

	loc = loc[i:]

	// rxBQI 에 해당하는 문자의 길이가 항상 5 글자인 것 같으나 보장되지 않았으므로
	// 안전하게 "/"" 이전 까지를 반환한다.
	i = strings.IndexRune(loc, '/')

	return loc[:i], nil
}
