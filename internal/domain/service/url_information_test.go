package service

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetWebsiteInformation(t *testing.T) {
	r := ioutil.NopCloser(strings.NewReader("<!DOCTYPE html><head></head><body><a href=\"http://example.com\"></a><a href=\"#internal\"></a><h1>test</h1><h2>test</h2></body>"))
	result, err := GetWebsiteInformation(r)
	assert.NoError(t, err)

	assert.NotNil(t, result)
	assert.Equal(t, 1, result.Headings["h1"])
	assert.Equal(t, 0, result.Headings["h3"])
	assert.Equal(t, 1, len(result.InternalLinks))
	assert.Equal(t, 1, len(result.ExternalLinks))
	assert.Equal(t, 1, len(result.ExternalLinks))
}
