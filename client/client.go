package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/labbcb/brave/search"
	"github.com/labbcb/brave/variant"
	"io/ioutil"
	"net/http"
)

type Client struct {
	Host     string
	Username string
	Password string
}

// SearchVariants requests variants by submitting queries.
func (c *Client) SearchVariants(input *search.Input) (*search.Response, error) {
	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(input); err != nil {
		return nil, err
	}

	resp, err := http.Post(c.Host+"/search", "application/json", &b)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("%d: %s", resp.StatusCode, string(body))
	}

	var response search.Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

// InsertVariant submits a variant to BraVE server.
func (c *Client) InsertVariant(v *variant.Variant) error {
	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(v); err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, c.Host+"/variants", &b)
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.Username, c.Password)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("%d: %s", resp.StatusCode, string(body))
	}

	var res map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return err
	}
	v.ID = res["id"]
	return nil
}

func (c *Client) RemoveVariants(datasetID, assemblyID string) error {
	url := fmt.Sprintf("%s/variants?dataset=%s&assembly=%s", c.Host, datasetID, assemblyID)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.Username, c.Password)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("%d: %s", resp.StatusCode, string(body))
	}

	return nil
}
