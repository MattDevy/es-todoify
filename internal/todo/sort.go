package todo

// SortOrder represents the direction of sorting.
// This is a value object in DDD terminology.
type SortOrder string

const (
	SortOrderAsc  SortOrder = "asc"
	SortOrderDesc SortOrder = "desc"
)

// IsValid checks if the sort order is one of the defined values.
func (s SortOrder) IsValid() bool {
	switch s {
	case SortOrderAsc, SortOrderDesc:
		return true
	}
	return false
}

// String returns the string representation of the SortOrder.
func (s SortOrder) String() string {
	return string(s)
}

// SortField represents valid fields that can be used for sorting.
// This is a value object in DDD terminology.
type SortField string

const (
	SortFieldCreateTime SortField = "createTime"
	SortFieldUpdateTime SortField = "updateTime"
	SortFieldTitle      SortField = "title"
	SortFieldStatus     SortField = "status"
)

// IsValid checks if the sort field is one of the defined values.
func (s SortField) IsValid() bool {
	switch s {
	case SortFieldCreateTime, SortFieldUpdateTime, SortFieldTitle, SortFieldStatus:
		return true
	}
	return false
}

// String returns the string representation of the SortField.
func (s SortField) String() string {
	return string(s)
}

// AllSortFields returns all valid sort field values.
func AllSortFields() []SortField {
	return []SortField{
		SortFieldCreateTime,
		SortFieldUpdateTime,
		SortFieldTitle,
		SortFieldStatus,
	}
}
