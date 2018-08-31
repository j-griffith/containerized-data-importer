package utils

import (
	"net/http"
	"encoding/json"
	"time"
	"io/ioutil"
	"github.com/golang/glog"
)

type fileCatalog []imageMeta

type imageMeta struct {
	Name string `json:"name"`
	ItemType string `json:"type"`
	Mtime string `json:"mtime"`
	size int32 `json:"size"`
}

func FileHostCatalog(ep string) (error) {
	client := http.Client{
		Timeout: 20 * time.Second,
	}
	resp, err := client.Get(ep)
	if err != nil {
		return err
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fc := &fileCatalog{}
	err = json.Unmarshal(buf, fc)
	if err != nil {
		return err
	}
	glog.Infof("DEBUG ---- %#v", fc)
	return nil
}
