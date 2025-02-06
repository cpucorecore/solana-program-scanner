package http_client

import (
	"net/http"
)

var GetBlockHttpClient = &http.Client{
	Transport: getBlockTransport,
}

var CommonHttpClient = &http.Client{
	Transport: commonTransport,
}
