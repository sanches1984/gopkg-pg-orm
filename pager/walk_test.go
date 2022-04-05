package pager

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWalk(t *testing.T) {
	t.Run("One page", func(t *testing.T) {
		walkTimes := 0

		err := Walk(10, func(page int32) (int64, error) {
			walkTimes++
			return 7, nil
		})

		assert.Nil(t, err)
		assert.Equal(t, 1, walkTimes)
	})

	t.Run("Many pages", func(t *testing.T) {
		walkTimes := 0

		err := Walk(10, func(page int32) (int64, error) {
			walkTimes++
			return 55, nil
		})

		assert.Nil(t, err)
		assert.Equal(t, 6, walkTimes)
	})

	t.Run("Full pages", func(t *testing.T) {
		walkTimes := 0

		err := Walk(10, func(page int32) (int64, error) {
			walkTimes++
			return 60, nil
		})

		assert.Nil(t, err)
		assert.Equal(t, 6, walkTimes)
	})
}
