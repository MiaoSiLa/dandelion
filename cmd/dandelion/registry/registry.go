package registry

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/tengattack/dandelion/client"
	"github.com/tengattack/dandelion/cmd/dandelion/config"
)

// Client for registry
type Client struct {
	service    string
	endpoint   string
	username   string
	password   string
	verify     bool
	httpClient *http.Client
}

// ListTagsResponse {catalog}/tags/list response
type ListTagsResponse struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

type harborTag struct {
	Name    string `json:"name"`
	Created string `json:"created"`
}

type riderTag struct {
	Branch          string `json:"branch"`
	BuildTaskID     int64  `json:"build_task_id"`
	BuildTime       string `json:"build_time"`
	CommitID        string `json:"commit_id"`
	CommitMessage   string `json:"commit_message"`
	DockerImageName string `json:"docker_image_name"`

	RetagDockerImageName *string `json:"retag_docker_image_name"`
}

type riderTagListResponse struct {
	Data        []riderTag `json:"data"`
	Result      bool       `json:"result"`
	TotalPages  int64      `json:"totalPages"`
	TotalCounts int64      `json:"total_counts"`
}

// NewClient creates a new registry client
func NewClient(conf *config.SectionRegistry) *Client {
	c := new(Client)
	c.service = conf.Service
	c.endpoint = conf.Endpoint
	c.username = conf.Username
	c.password = conf.Password
	c.verify = conf.Verify
	if !c.verify {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		c.httpClient = &http.Client{Transport: tr}
	} else {
		c.httpClient = &http.Client{}
	}

	return c
}

// sort version like X.Y.Z-BuildNum
func lessDashVersion(a, b string) bool {
	version1 := strings.Split(a, "-")
	version2 := strings.Split(b, "-")
	for i := range version1 {
		if len(version2) <= i {
			// all previous parts equal but version2 has more parts than version1
			return false
		}
		v1, err1 := strconv.Atoi(version1[i])
		v2, err2 := strconv.Atoi(version2[i])
		if err1 == nil && err2 == nil {
			// num compare
			if v1 > v2 {
				return false
			} else if v1 < v2 {
				return true
			}
		} else {
			// non num less then num parts
			if err1 != nil && err2 == nil {
				return true
			} else if err1 == nil && err2 != nil {
				return false
			}
			if version1[i] > version2[i] {
				return false
			} else if version1[i] < version2[i] {
				return true
			}
		}
	}
	if len(version2) > len(version1) {
		return true
	}
	return false
}

func (c *Client) registryListTags(catalog string) (*ListTagsResponse, error) {
	url := fmt.Sprintf("%s/v2/%s/tags/list", c.endpoint, catalog)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	client.InitHTTPRequest(req, false)
	if c.username != "" && c.password != "" {
		req.SetBasicAuth(c.username, c.password)
	}
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var v ListTagsResponse
	err = json.Unmarshal(body, &v)
	if err != nil {
		return nil, err
	}
	if len(v.Tags) > 0 {
		sort.SliceStable(v.Tags, func(i, j int) bool {
			// reverse sort
			return !lessDashVersion(v.Tags[i], v.Tags[j])
		})
	}
	return &v, nil
}

func (c *Client) harborListTags(catalog string) (*ListTagsResponse, error) {
	url := fmt.Sprintf("%s/api/repositories/%s/tags", c.endpoint, catalog)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	client.InitHTTPRequest(req, false)
	if c.username != "" && c.password != "" {
		req.SetBasicAuth(c.username, c.password)
	}
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var tagsRes []harborTag
	err = json.Unmarshal(body, &tagsRes)
	if err != nil {
		return nil, err
	}
	v := &ListTagsResponse{
		Name: catalog,
		Tags: make([]string, 0, len(tagsRes)),
	}
	if len(tagsRes) > 0 {
		sort.SliceStable(tagsRes, func(i, j int) bool {
			return tagsRes[i].Created > tagsRes[j].Created
		})
		for _, tag := range tagsRes {
			v.Tags = append(v.Tags, tag.Name)
		}
	}
	return v, nil
}

func (c *Client) riderListTags(catalog string) (*ListTagsResponse, error) {
	url := fmt.Sprintf("%s/api/pkg/search_app_release", c.endpoint)
	query := make(map[string]interface{})
	query["tree_path"] = "bilibili." + catalog
	queryBuf, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(queryBuf))
	if err != nil {
		return nil, err
	}
	client.InitHTTPRequest(req, false)
	req.Header.Set("Content-Type", "application/json")
	if c.password != "" {
		req.Header.Set("rider-token", c.password)
	}
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var resp riderTagListResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}
	v := &ListTagsResponse{
		Name: catalog,
		Tags: make([]string, 0, len(resp.Data)),
	}
	if len(resp.Data) > 0 {
		for _, tag := range resp.Data {
			imageName := tag.DockerImageName
			if tag.RetagDockerImageName != nil && *tag.RetagDockerImageName != "" {
				imageName = *tag.RetagDockerImageName
			}
			pos := strings.LastIndex(imageName, ":")
			if pos < 0 {
				continue
			}
			v.Tags = append(v.Tags, imageName[pos+1:])
		}
	}
	return v, nil
}

// ListTags list tags for catalog
func (c *Client) ListTags(catalog string) (*ListTagsResponse, error) {
	if c.service == "rider" {
		return c.riderListTags(catalog)
	} else if c.service == "harbor" {
		return c.harborListTags(catalog)
	}
	// defaults to registry
	return c.registryListTags(catalog)
}
