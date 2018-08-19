package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/ahojsenn/kontrol/util"
)

func TestGetEnvironment (t *testing.T) {
	environment := util.GetEnv()
	assert.NotEmpty(t, environment)
	assert.NotEmpty(t, environment.CertFile)
	assert.NotEmpty(t, environment.KeyFile)
	assert.NotEmpty(t, environment.Hostname)
}
