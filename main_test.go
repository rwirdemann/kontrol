package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGetEnvironment (t *testing.T) {
	environment := getEnvironment()
	assert.NotEmpty(t, environment)
	assert.NotEmpty(t, environment.CertFile)
	assert.NotEmpty(t, environment.KeyFile)
	assert.NotEmpty(t, environment.Hostname)
}
