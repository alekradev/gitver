package cmd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateFlags(t *testing.T) {
	tagFlag = true
	pushFlag = true
	autoFlag = true

	result, err := validateFlags()
	assert.True(t, result)
	assert.NoError(t, err)

}

func TestValidateFlags_ErrorTagAndPush(t *testing.T) {
	tagFlag = false
	pushFlag = true

	result, err := validateFlags()
	assert.False(t, result)
	assert.Error(t, err)

}

func TestValidateFlags_ErrorTagAndAmend(t *testing.T) {
	tagFlag = false
	amendFlag = true

	result, err := validateFlags()
	assert.False(t, result)
	assert.Error(t, err)
}

func TestValidateFlags_ErrorMode(t *testing.T) {
	flags := []struct {
		name  string
		value *bool
	}{
		{"autoFlag", &autoFlag},
		{"commitFlag", &commitFlag},
		{"majorFlag", &majorFlag},
		{"minorFlag", &minorFlag},
		{"patchFlag", &patchFlag},
	}

	// Erstelle alle Kombinationen, bei denen mindestens zwei Flags auf true gesetzt sind
	for i := 0; i < len(flags); i++ {
		for j := i + 1; j < len(flags); j++ {
			// Setze alle Flags zurück
			autoFlag = false
			commitFlag = false
			majorFlag = false
			minorFlag = false
			patchFlag = false

			// Setze zwei Flags auf true
			*flags[i].value = true
			*flags[j].value = true

			// Führe den Test aus
			t.Run(flags[i].name+"And"+flags[j].name, func(t *testing.T) {
				result, err := validateFlags()
				assert.False(t, result, "Sollte false zurückgeben, wenn mehr als ein Flag gesetzt ist")
				assert.Error(t, err, "Sollte einen Fehler zurückgeben, wenn mehr als ein Flag gesetzt ist")
			})
		}
	}
}
