package processing

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestBerechne_Gewerbesteuer(t *testing.T) {
	var gs = 0.0

	gs =  berechne_Gewerbesteuer(10.000)
	assert.Equal(t, 0.0, gs)

	gs = berechne_Gewerbesteuer(Gewerbesteuer_Freibetrag)
	assert.Equal(t, 0.0, gs)

	gs = math.Round(100*berechne_Gewerbesteuer((Gewerbesteuer_Freibetrag+100.0)))/100.0
	assert.Equal(t, 16.45, gs)
	}