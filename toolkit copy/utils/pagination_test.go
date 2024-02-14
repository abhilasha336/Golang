// nolint
package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/tuneverse/toolkit/consts"
	"gitlab.com/tuneverse/toolkit/models"
)

func TestPaginate(t *testing.T) {

	testCases := []struct {
		name          string
		page, limit   int32
		defaultLimit  int32
		expectedPage  int32
		expectedLimit int32
	}{
		{
			name:          "Page and limit both zero, default limit provided",
			page:          0,
			limit:         0,
			defaultLimit:  50,
			expectedPage:  consts.DefaultPage,
			expectedLimit: 50,
		},
		{
			name:          "Page zero, limit negative, default limit provided",
			page:          0,
			limit:         -5,
			defaultLimit:  50,
			expectedPage:  consts.DefaultPage,
			expectedLimit: 50,
		},
		{
			name:          "Page negative, limit zero, default limit provided",
			page:          -2,
			limit:         0,
			defaultLimit:  50,
			expectedPage:  consts.DefaultPage,
			expectedLimit: 50,
		},
		{
			name:          "Page and limit both positive",
			page:          3,
			limit:         20,
			defaultLimit:  50,
			expectedPage:  3,
			expectedLimit: 20,
		},
		{
			name:          "Page zero, limit zero, default limit zero",
			page:          0,
			limit:         0,
			defaultLimit:  0,
			expectedPage:  consts.DefaultPage,
			expectedLimit: consts.DefaultLimit,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualPage, actualLimit := Paginate(tc.page, tc.limit, tc.defaultLimit)
			require.Equal(t, tc.expectedPage, actualPage)
			require.Equal(t, tc.expectedLimit, actualLimit)
		})
	}
}

func TestMetaDataInfo(t *testing.T) {
	testCases := []struct {
		name          string
		metaDataInput *models.MetaData
		expectedNext  int32
		expectedPrev  int32
	}{
		{
			name: "First page with more data available",
			metaDataInput: &models.MetaData{
				Total:       100,
				PerPage:     10,
				CurrentPage: 1,
			},
			expectedNext: 2,
			expectedPrev: 0, // No previous page on the first page
		},
		{
			name: "Second page",
			metaDataInput: &models.MetaData{
				Total:       100,
				PerPage:     10,
				CurrentPage: 2,
			},
			expectedNext: 3,
			expectedPrev: 1,
		},
		{
			name: "Last page with more data available",
			metaDataInput: &models.MetaData{
				Total:       100,
				PerPage:     10,
				CurrentPage: 10,
			},
			expectedNext: 0, // No next page on the last page
			expectedPrev: 9,
		},
		{
			name: "Last page with no more data",
			metaDataInput: &models.MetaData{
				Total:       95,
				PerPage:     10,
				CurrentPage: 10,
			},
			expectedNext: 0,
			expectedPrev: 9,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := MetaDataInfo(tc.metaDataInput)
			require.Equal(t, tc.expectedNext, result.Next)
			require.Equal(t, tc.expectedPrev, result.Prev)
		})
	}
}
