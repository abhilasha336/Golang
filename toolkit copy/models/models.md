# Package Models
This package defines the common models used in the application.

### MetaData

The MetaData struct is used to represent metadata associated with a collection of data,  used in APIs or data retrieval scenarios.

1. `Total`: 
An integer field (int64) that represents the total number of items in the collection. It signifies the overall count of items without pagination.

2. `PerPage`: 
An integer field (int32) that indicates the number of items displayed per page or per API request. This value represents the size of each page in a paginated dataset.

3. `CurrentPage`: 
An integer field (int32) that denotes the current page or section of the dataset being accessed. It provides context about where you are within the paginated data.

4. `Next`: 
An integer field (int32) that typically holds the number of the next page in a paginated dataset. It's used to navigate to the next page of data when fetching paginated results.

5. `Prev`: 
An integer field (int32) that usually contains the number of the previous page in a paginated dataset. It allows you to navigate to the previous page of data when needed.