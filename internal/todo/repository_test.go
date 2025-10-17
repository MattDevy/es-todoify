package todo

import (
	"testing"
	"time"
)

func TestListFilter_Validate(t *testing.T) {
	now := time.Now()
	tomorrow := now.Add(24 * time.Hour)
	yesterday := now.Add(-24 * time.Hour)

	tests := []struct {
		name    string
		filter  ListFilter
		wantErr bool
	}{
		{
			name:    "valid default filter",
			filter:  DefaultListFilter(),
			wantErr: false,
		},
		{
			name: "valid filter with all fields",
			filter: ListFilter{
				Status:      StatusPending,
				Labels:      []string{"urgent", "bug"},
				SearchQuery: "test",
				FromDate:    &yesterday,
				ToDate:      &tomorrow,
				Limit:       100,
				Offset:      0,
				SortBy:      SortFieldCreateTime,
				SortOrder:   SortOrderAsc,
			},
			wantErr: false,
		},
		{
			name: "invalid status",
			filter: ListFilter{
				Status: Status("invalid"),
			},
			wantErr: true,
		},
		{
			name: "invalid sort field",
			filter: ListFilter{
				SortBy: SortField("invalidField"),
			},
			wantErr: true,
		},
		{
			name: "invalid sort order",
			filter: ListFilter{
				SortOrder: SortOrder("invalidOrder"),
			},
			wantErr: true,
		},
		{
			name: "negative limit",
			filter: ListFilter{
				Limit: -1,
			},
			wantErr: true,
		},
		{
			name: "negative offset",
			filter: ListFilter{
				Offset: -1,
			},
			wantErr: true,
		},
		{
			name: "invalid date range (from after to)",
			filter: ListFilter{
				FromDate: &tomorrow,
				ToDate:   &yesterday,
			},
			wantErr: true,
		},
		{
			name: "valid date range",
			filter: ListFilter{
				FromDate: &yesterday,
				ToDate:   &tomorrow,
			},
			wantErr: false,
		},
		{
			name: "zero limit is valid",
			filter: ListFilter{
				Limit: 0,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.filter.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ListFilter.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDefaultListFilter(t *testing.T) {
	filter := DefaultListFilter()

	if filter.Limit != 50 {
		t.Errorf("DefaultListFilter().Limit = %d, want 50", filter.Limit)
	}

	if filter.Offset != 0 {
		t.Errorf("DefaultListFilter().Offset = %d, want 0", filter.Offset)
	}

	if filter.SortBy != SortFieldCreateTime {
		t.Errorf("DefaultListFilter().SortBy = %v, want %v", filter.SortBy, SortFieldCreateTime)
	}

	if filter.SortOrder != SortOrderDesc {
		t.Errorf("DefaultListFilter().SortOrder = %v, want %v", filter.SortOrder, SortOrderDesc)
	}

	// Default filter should be valid
	if err := filter.Validate(); err != nil {
		t.Errorf("DefaultListFilter() should be valid, got error: %v", err)
	}
}
