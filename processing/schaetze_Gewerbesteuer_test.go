package processing

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestBerechne_Gewerbesteuer(t *testing.T) {
	var gs = 0.0

	gs =  berechne_Gewerbesteuer(0)
	assert.Equal(t, 0.0, gs)

	gs = berechne_Gewerbesteuer(Gewerbesteuer_Freibetrag)
	assert.Equal(t, 0.0, gs)

	gs  = math.Round( berechne_Gewerbesteuer(652733.50) )
	assert.Equal(t, 103339.0, gs)

	gs = math.Round( berechne_Gewerbesteuer(596219.72) )
	assert.Equal(t, 94045.0, gs)
	}