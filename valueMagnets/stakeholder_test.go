package valueMagnets

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetKommitmenschenRepository(t *testing.T) {

	// generate accounts for all stakeholders
	repo := KommimtmentYear{}
	assert.Equal(t, "AN", repo.All(2016)[0].Id)
	assert.Equal(t, "JM", repo.All(2016)[2].Id)
//	assert.Equal(t, "0%", repo.All(2018)[1].Arbeit)
	assert.True(t, len(repo.All(2016)) > 1 )

}

func TestKommimtmentYear_All(t *testing.T) {

	repo := Stakeholder{}
	assert.Equal(t, "JM", repo.All(2016)[3].Id)
}

func TestIsValidStakeholder (t *testing.T) {
	repo := Stakeholder{}
	assert.True(t, repo.IsValidStakeholder("K"))
	assert.True(t, repo.IsValidStakeholder("JM"))
	assert.True(t, repo.IsValidStakeholder("Rest"))
	assert.True(t, repo.IsValidStakeholder("Extern"))
	assert.False(t, repo.IsValidStakeholder("krümelhügliplis"))

}

func TestKommimtmentYear_Liqui(t *testing.T) {
	repo := KommimtmentYear{}
	assert.Equal(t, 42.23, repo.Liqui(2016))
}

func TestIsEmployee (t *testing.T) {
	repo := Stakeholder{}
	assert.False(t, repo.IsEmployee("JM"))
}

func TestIsPartner (t *testing.T) {
	repo := Stakeholder{}
	assert.True(t, repo.IsPartner("JM"))
}