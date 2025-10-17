package todo

import "testing"

func TestSortOrder_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		order SortOrder
		want  bool
	}{
		{"valid asc", SortOrderAsc, true},
		{"valid desc", SortOrderDesc, true},
		{"invalid empty", SortOrder(""), false},
		{"invalid random", SortOrder("random"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.order.IsValid(); got != tt.want {
				t.Errorf("SortOrder.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSortField_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		field SortField
		want  bool
	}{
		{"valid createTime", SortFieldCreateTime, true},
		{"valid updateTime", SortFieldUpdateTime, true},
		{"valid title", SortFieldTitle, true},
		{"valid status", SortFieldStatus, true},
		{"invalid empty", SortField(""), false},
		{"invalid random", SortField("invalidField"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.field.IsValid(); got != tt.want {
				t.Errorf("SortField.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAllSortFields(t *testing.T) {
	fields := AllSortFields()

	if len(fields) != 4 {
		t.Errorf("AllSortFields() returned %d fields, want 4", len(fields))
	}

	// Verify all returned fields are valid
	for _, field := range fields {
		if !field.IsValid() {
			t.Errorf("AllSortFields() returned invalid field: %v", field)
		}
	}
}
