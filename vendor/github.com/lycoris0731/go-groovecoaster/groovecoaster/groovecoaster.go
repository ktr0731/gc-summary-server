package groovecoaster

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// APIClient is the structure that represents Groove Coaster client
type APIClient struct {
	client *http.Client
}

// New creates a groovecoaster client
func New() *APIClient {
	client, err := tryLogin()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	return &APIClient{
		client: client,
	}
}

func (c *APIClient) get(uri string) ([]byte, error) {
	res, err := c.client.Get("https://" + uri)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch %s failed (Status: %s)", uri, res.Status)
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("response body is empty")
	}

	return data, nil
}

func (c *APIClient) unmarshal(data []byte, i interface{}) error {
	if err := json.Unmarshal(data, &i); err != nil {
		return fmt.Errorf("unmarshal failed: %s", err)
	}

	return nil
}
