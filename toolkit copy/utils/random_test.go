// nolint
package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandomString(t *testing.T) {
	tests := []struct {
		length int
	}{
		{10},
		{20},
		{30},
		{0},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Length: %d\n", test.length), func(t *testing.T) {
			result, err := RandomString(test.length)
			require.NoError(t, err)
			require.Len(t, result, test.length)

		})
	}
}
