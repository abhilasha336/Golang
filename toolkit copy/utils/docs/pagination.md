## Overview
This provides functions for working with pagination and calculating pagination metadata.

## Index
- [Paginate(page, limit, defaultLimit int32) (int32, int32)](#func-Paginate)
- [MetaDataInfo(metaData *models.MetaData) *models.MetaData](#func-MetaDataInfo)

### func Paginate

   Paginate(page, limit, defaultLimit int32) (int32, int32)

The Paginate function checks and updates the page and limit values based on specific criteria.

### func MetaDataInfo

    MetaDataInfo(metaData *models.MetaData) *models.MetaData

This function calculates pagination metadata based on the provided MetaData struct and updates it with information about the next and previous pages