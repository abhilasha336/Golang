package utils

import (
	"gitlab.com/tuneverse/toolkit/consts"
	"gitlab.com/tuneverse/toolkit/models"
)

// Paginate sets default values for page and limit if not exists.
func Paginate(page, limit, defaultLimit int32) (int32, int32) {
	if page <= 0 {
		page = consts.DefaultPage
	}
	if (limit <= 0 || limit > consts.MaxLimit) && defaultLimit > 0 {
		limit = defaultLimit
	} else if limit <= 0 {
		limit = consts.DefaultLimit
	}
	return page, limit
}

// MetaDataInfo calculates values for pagination.
func MetaDataInfo(metaData *models.MetaData) *models.MetaData {
	if metaData.Total < 1 {
		return nil
	}
	if int64(metaData.CurrentPage)*int64(metaData.PerPage) < metaData.Total {
		metaData.Next = metaData.CurrentPage + 1
	}
	if metaData.CurrentPage > 1 {
		metaData.Prev = metaData.CurrentPage - 1
	}
	return metaData
}
