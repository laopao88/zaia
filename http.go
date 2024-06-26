package zaia

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"
)

func HttpPost(url string, content string, timeout time.Duration) error {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(content)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{
		Timeout: timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	body, _ := io.ReadAll(resp.Body)
	return errors.New(resp.Status + " " + string(body))
}

func HttpRequest(url string, post bool, content []byte, contentType string, timeout time.Duration, headerCallback func(header http.Header)) ([]byte, error) {
	var req *http.Request
	var err error
	if post {
		if content != nil {
			req, err = http.NewRequest("POST", url, bytes.NewBuffer(content))
		} else {
			req, err = http.NewRequest("POST", url, nil)
		}
	} else {
		req, err = http.NewRequest("GET", url, nil)
	}
	if err != nil {
		return nil, err
	}
	if post {
		req.Header.Set("Content-Type", contentType)
	}
	if headerCallback != nil {
		headerCallback(req.Header)
	}
	client := &http.Client{
		Timeout: timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusOK {
		return body, err
	}
	return nil, errors.New(resp.Status + " " + string(body))
}

func HttpAddParam(address string, k string, v string) (string, error) {
	u, err := url.Parse(address)
	if err != nil {
		return address, err
	}
	q, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return address, err
	}
	for key, _ := range q {
		if key == k {
			return address, nil
		}
	}

	q.Add(k, v)
	u.RawQuery = q.Encode()
	return u.String(), nil
}
