package utils_test

import (
	"testing"

	"github.com/alvesrafa/video-encoder/framework/utils"
	"github.com/stretchr/testify/require"
)

func TestUtils(t *testing.T) {
	t.Run("should be a json", func(t *testing.T) {

		json := `{
								"id": 1,
								"file_path": "path",
								"status": "PENDING"
						}`

		err := utils.IsJson(json)

		require.Nil(t, err)
	})

	t.Run("should'nt be a json", func(t *testing.T) {

		json := `{fakeJSON
								id:1,
								file_path: "path",
								status: "PENDING"
								
						}`

		err := utils.IsJson(json)

		require.Error(t, err)
	})
}
