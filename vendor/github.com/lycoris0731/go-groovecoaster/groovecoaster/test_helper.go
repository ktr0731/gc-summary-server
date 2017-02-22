package groovecoaster

import "net/http"

const scheme = "https://"

var testClient APIClient

func init() {
	var client http.Client
	testClient.client = &client
}
