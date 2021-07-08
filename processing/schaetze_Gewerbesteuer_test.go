package processing

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestBerechne_Gewerbesteuer(t *testing.T) {
	var gs = 0.0

	gs =  berechneGewerbesteuer(0)
	assert.Equal(t, 0.0, gs)

	gs = berechneGewerbesteuer(Gewerbesteuer_Freibetrag)
	assert.Equal(t, 0.0, gs)

	gs  = math.Round( berechneGewerbesteuer(652733.50) )
	assert.Equal(t, 103339.0, gs)

	gs = math.Round( berechneGewerbesteuer(596219.72) )
	assert.Equal(t, 94045.0, gs)
	}