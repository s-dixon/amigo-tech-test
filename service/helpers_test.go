package service

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"net/url"
	"net/http/httptest"
)

func Test_ShouldReturnDefaultValueDueToParamNotFound(t *testing.T) {
	values := url.Values{"offset": []string{"0"}}

	result := getQueryParamOrDefault(values, "limit", "10")

	assert.Equal(t, "10", result, "Default value was not returned")
}

func Test_ShouldReturnDefaultValueDueToMultipleValuesFoundForKey(t *testing.T) {
	values := url.Values{"offset": []string{"0","1"}}

	result := getQueryParamOrDefault(values, "offset", "10")

	assert.Equal(t, "10", result, "Default value was not returned")
}

func Test_ShouldReturnQueryParamValue(t *testing.T) {
	values := url.Values{"offset": []string{"0"}}

	result := getQueryParamOrDefault(values, "offset", "10")

	assert.Equal(t, "0", result, "Query param value was not returned")
}

func Test_ShouldReturnIpAddress(t *testing.T) {
	request := httptest.NewRequest("POST", "/test", nil)
	request.RemoteAddr = "192.168.200.201:1234"

	ip, err := getClientIp(request)

	assert.Nil(t, err)
	assert.Equal(t, "192.168.200.201", ip.String(), "Ip address was not extracted correct")
}

func Test_ShouldReturnErrorWhenRetrievingIpAddressDueToMissingPost(t *testing.T) {
	request := httptest.NewRequest("POST", "/test", nil)
	request.RemoteAddr = "192.168.200.201"

	ip, err := getClientIp(request)

	assert.NotNil(t, err)
	assert.Nil(t, ip)
}

func Test_ShouldReturnErrorWhenRetrievingIpAddressDueInvalidFormat(t *testing.T) {
	request := httptest.NewRequest("POST", "/test", nil)
	request.RemoteAddr = "not_an_ip:1234"

	ip, err := getClientIp(request)

	assert.NotNil(t, err)
	assert.Nil(t, ip)
}

func Test_ShouldRespondWithJson(t *testing.T) {
	responseRecorder := httptest.NewRecorder()

	respondWithJSON(responseRecorder, 200, map[string]string{"key":"value"})

	assert.Equal(t, "application/json", responseRecorder.Header().Get("Content-Type"), "Response content type not set as expected")
	assert.Equal(t, 200, responseRecorder.Code, "Response status code not as expected")
	assert.Equal(t, "{\"key\":\"value\"}", responseRecorder.Body.String(), "Response body not as expected")
}

func Test_ShouldRespondWithError(t *testing.T) {
	responseRecorder := httptest.NewRecorder()

	respondWithError(responseRecorder, 400, "Bad request")

	assert.Equal(t, "application/json", responseRecorder.Header().Get("Content-Type"), "Response content type not set as expected")
	assert.Equal(t, 400, responseRecorder.Code, "Response status code not as expected")
	assert.Equal(t, "{\"error\":\"Bad request\"}", responseRecorder.Body.String(), "Response body not as expected")
}

func Test_ShouldRespondWithString(t *testing.T) {
	responseRecorder := httptest.NewRecorder()

	respondWithString(responseRecorder, 200, "Test message value")

	assert.Equal(t, "", responseRecorder.Header().Get("Content-Type"), "Response content type should not be set")
	assert.Equal(t, 200, responseRecorder.Code, "Response status code not as expected")
	assert.Equal(t, "Test message value", responseRecorder.Body.String(), "Response body not as expected")
}