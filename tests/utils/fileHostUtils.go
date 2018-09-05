package utils

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/pkg/errors"
	"time"
)

type FileCatalog []ImageMeta

type ImageMeta struct {
	Name     string `json:"name"`
	ItemType string `json:"type"`
	Mtime    string `json:"mtime"`
	Size     int32  `json:"size"`
}

func GetFileHostCatalog(ep, accessKey, secretKey string) (*FileCatalog, error) {
	req, err := http.NewRequest("GET", ep, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "could not create request for ep %q", ep)
	}
	if accessKey != "" || secretKey != "" {
		req.SetBasicAuth(accessKey, secretKey)
	}

	client := http.Client{
		Timeout: 5 * time.Second,
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			if accessKey != "" || secretKey != "" {
				r.SetBasicAuth(accessKey, secretKey) // Redirects will lose basic auth, so reset them manually
			}
			return nil
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "could not send request")
	}
	if resp.StatusCode >= 400 {
		buf, _ := ioutil.ReadAll(resp.Body)
		return nil, errors.Errorf("Got response code %q, message:\n%s", resp.StatusCode,string(buf))
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "could not read response body")
	}

	fc := new(FileCatalog)
	err = json.Unmarshal(buf, fc)
	if err != nil {
		return nil, errors.Wrap(err, "could not unmarshall response body")
	}

	return fc, nil
}
