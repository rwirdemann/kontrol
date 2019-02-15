package valueMagnets

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetKommitmenschenRepository(t *testing.T) {

	// generate accounts for all stakeholders
	repo := KommitmenschenRepository{}
	assert.Equal(t, "AN", repo.All(2016)[0].Id)
	assert.Equal(t, "JM", repo.All(2016)[2].Id)
//	assert.Equal(t, "0%", repo.All(2018)[1].Arbeit)
	assert.True(t, len(repo.All(2016)) > 1 )

}

func TestStakeholderRepository_All(t *testing.T) {

	repo := Stakeholder{}
	assert.Equal(t, "JM", repo.All(2016)[2].Id)
}

func TestIsValidStakeholder (t *testing.T) {
	repo := Stakeholder{}
	assert.True(t, repo.IsValidStakeholder("K"))
	assert.True(t, repo.IsValidStakeholder("JM"))
	assert.True(t, repo.IsValidStakeholder("RR"))
	assert.True(t, repo.IsValidStakeholder("EX"))
	assert.False(t, repo.IsValidStakeholder("krümelhügliplis"))

}