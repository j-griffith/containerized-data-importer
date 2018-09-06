package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

// Catalog holds all ImageMeta objects parsed from a response by the cdi file host.
type Catalog []ImageMeta

// ImageMeta is defined to match a single file's meta data element in a json list response
// from the cdi file host when the request is sent to the server root path
type ImageMeta struct {
	Name     string `json:"name"`
	ItemType string `json:"type"`
	Mtime    string `json:"mtime"`
	Size     int32  `json:"size"`
}

type catalog interface {
	List() []string
	Size(string) (int, error)
}

var _ catalog = Catalog{}

// GetCatalog sends a request to the endpoint `ep` (with credentials if defined),
// which is expected to be the cdi file host's server root `/`
// The file server is configured to return a json array of objects matching ImageMeta
// definition. The response is parsed into a Catalog object, which is returned.
func GetCatalog(ep, accessKey, secretKey string) (*Catalog, error) {
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
	buf, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return nil, errors.Errorf("Got response code %q, message:\n%s", resp.StatusCode, string(buf))
	}
	if err != nil {
		return nil, errors.Wrap(err, "could not read response body")
	}

	fc := new(Catalog)
	err = json.Unmarshal(buf, fc)
	if err != nil {
		return nil, errors.Wrap(err, "could not unmarshall response body")
	}

	return fc, nil
}

// List generates a slice of file names from the Catalog.
func (c Catalog) List() []string {
	var files []string
	for _, i := range c {
		files = append(files, i.Name)
	}
	return files
}

// Size searches the catalog for a file name matching the parameter and returns the associated size
// If no match is found, returns an error and size of -1
func (c Catalog) Size(file string) (int, error) {
	for _, i := range c {
		if i.Name == file {
			return int(i.Size), nil
		}
	}
	return -1, errors.Errorf("File %q not found in catalog", file)
}
