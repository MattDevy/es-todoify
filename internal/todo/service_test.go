package todo

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockRepository is a mock implementation of the Repository interface.
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, todo *Todo) error {
	args := m.Called(ctx, todo)
	return args.Error(0)
}

func (m *MockRepository) Get(ctx context.Context, id string) (*Todo, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Todo), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, todo *Todo) error {
	args := m.Called(ctx, todo)
	return args.Error(0)
}

func (m *MockRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) List(ctx context.Context, filter ListFilter) ([]*Todo, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*Todo), args.Error(1)
}

func (m *MockRepository) Count(ctx context.Context, filter ListFilter) (int, error) {
	args := m.Called(ctx, filter)
	return args.Int(0), args.Error(1)
}

// Test helpers

func newTestService(t *testing.T) (*Service, *MockRepository) {
	t.Helper()
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)
	return service, mockRepo
}

func newValidTodo(t *testing.T) *Todo {
	t.Helper()
	todo, err := NewTodo("Test Todo", "Test Description", []string{"test"})
	require.NoError(t, err)
	return todo
}

func validUUID() string {
	return uuid.New().String()
}

func TestService_CreateTodo(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		title       string
		description string
		labels      []string
		setupMock   func(*MockRepository)
		wantErr     bool
		assertErr   func(*testing.T, error)
		assertTodo  func(*testing.T, *Todo)
	}{
		{
			name:        "successful creation",
			title:       "Test Todo",
			description: "Description",
			labels:      []string{"label1"},
			setupMock: func(m *MockRepository) {
				m.On("Create", ctx, mock.AnythingOfType("*todo.Todo")).Return(nil)
			},
			wantErr: false,
			assertTodo: func(t *testing.T, todo *Todo) {
				require.Equal(t, "Test Todo", todo.Title)
				require.Equal(t, "Description", todo.Description)
				require.Equal(t, []string{"label1"}, todo.Labels)
				require.Equal(t, StatusPending, todo.Status)
				require.NotEqual(t, uuid.Nil, todo.ID)
			},
		},
		{
			name:        "empty title returns error",
			title:       "",
			description: "Description",
			labels:      nil,
			setupMock:   func(m *MockRepository) {},
			wantErr:     true,
			assertErr: func(t *testing.T, err error) {
				require.ErrorIs(t, err, ErrInvalidInput)
			},
		},
		{
			name:        "repository error",
			title:       "Test",
			description: "Description",
			labels:      nil,
			setupMock: func(m *MockRepository) {
				m.On("Create", ctx, mock.AnythingOfType("*todo.Todo")).
					Return(errors.New("database error"))
			},
			wantErr: true,
			assertErr: func(t *testing.T, err error) {
				require.Contains(t, err.Error(), "failed to create todo")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := newTestService(t)
			tt.setupMock(mockRepo)

			todo, err := service.CreateTodo(ctx, tt.title, tt.description, tt.labels)

			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, todo)
				if tt.assertErr != nil {
					tt.assertErr(t, err)
				}
			} else {
				require.NoError(t, err)
				require.NotNil(t, todo)
				if tt.assertTodo != nil {
					tt.assertTodo(t, todo)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_GetTodo(t *testing.T) {
	ctx := context.Background()
	validID := validUUID()
	testTodo := newValidTodo(t)

	tests := []struct {
		name      string
		id        string
		setupMock func(*MockRepository)
		wantErr   bool
		assertErr func(*testing.T, error)
	}{
		{
			name: "successful get",
			id:   validID,
			setupMock: func(m *MockRepository) {
				m.On("Get", ctx, validID).Return(testTodo, nil)
			},
			wantErr: false,
		},
		{
			name:      "empty id returns error",
			id:        "",
			setupMock: func(m *MockRepository) {},
			wantErr:   true,
			assertErr: func(t *testing.T, err error) {
				require.ErrorIs(t, err, ErrInvalidInput)
				require.Contains(t, err.Error(), "id is required")
			},
		},
		{
			name:      "invalid uuid format returns error",
			id:        "invalid-uuid",
			setupMock: func(m *MockRepository) {},
			wantErr:   true,
			assertErr: func(t *testing.T, err error) {
				require.ErrorIs(t, err, ErrInvalidInput)
				require.Contains(t, err.Error(), "invalid id format")
			},
		},
		{
			name: "not found",
			id:   validID,
			setupMock: func(m *MockRepository) {
				m.On("Get", ctx, validID).Return(nil, ErrNotFound)
			},
			wantErr: true,
			assertErr: func(t *testing.T, err error) {
				require.ErrorIs(t, err, ErrNotFound)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := newTestService(t)
			tt.setupMock(mockRepo)

			todo, err := service.GetTodo(ctx, tt.id)

			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, todo)
				if tt.assertErr != nil {
					tt.assertErr(t, err)
				}
			} else {
				require.NoError(t, err)
				require.NotNil(t, todo)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_UpdateTodo(t *testing.T) {
	ctx := context.Background()
	validID := validUUID()

	newTitle := "New Title"
	newDesc := "New Description"
	emptyTitle := ""

	tests := []struct {
		name       string
		id         string
		update     UpdateTodo
		setupMock  func(*MockRepository)
		wantErr    bool
		assertErr  func(*testing.T, error)
		assertTodo func(*testing.T, *Todo)
	}{
		{
			name: "successful update",
			id:   validID,
			update: UpdateTodo{
				Title:       &newTitle,
				Description: &newDesc,
				Labels:      []string{"new"},
			},
			setupMock: func(m *MockRepository) {
				existingTodo := newValidTodo(t)
				m.On("Get", ctx, validID).Return(existingTodo, nil)
				m.On("Update", ctx, mock.AnythingOfType("*todo.Todo")).Return(nil)
			},
			wantErr: false,
			assertTodo: func(t *testing.T, todo *Todo) {
				require.Equal(t, newTitle, todo.Title)
				require.Equal(t, newDesc, todo.Description)
				require.Equal(t, []string{"new"}, todo.Labels)
			},
		},
		{
			name:      "empty id returns error",
			id:        "",
			update:    UpdateTodo{},
			setupMock: func(m *MockRepository) {},
			wantErr:   true,
			assertErr: func(t *testing.T, err error) {
				require.ErrorIs(t, err, ErrInvalidInput)
			},
		},
		{
			name:      "invalid uuid format returns error",
			id:        "invalid-uuid",
			update:    UpdateTodo{},
			setupMock: func(m *MockRepository) {},
			wantErr:   true,
			assertErr: func(t *testing.T, err error) {
				require.ErrorIs(t, err, ErrInvalidInput)
			},
		},
		{
			name:   "todo not found",
			id:     validID,
			update: UpdateTodo{},
			setupMock: func(m *MockRepository) {
				m.On("Get", ctx, validID).Return(nil, ErrNotFound)
			},
			wantErr: true,
			assertErr: func(t *testing.T, err error) {
				require.ErrorIs(t, err, ErrNotFound)
			},
		},
		{
			name: "invalid update data",
			id:   validID,
			update: UpdateTodo{
				Title: &emptyTitle,
			},
			setupMock: func(m *MockRepository) {
				existingTodo := newValidTodo(t)
				m.On("Get", ctx, validID).Return(existingTodo, nil)
			},
			wantErr: true,
			assertErr: func(t *testing.T, err error) {
				require.ErrorIs(t, err, ErrInvalidInput)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := newTestService(t)
			tt.setupMock(mockRepo)

			todo, err := service.UpdateTodo(ctx, tt.id, tt.update)

			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, todo)
				if tt.assertErr != nil {
					tt.assertErr(t, err)
				}
			} else {
				require.NoError(t, err)
				require.NotNil(t, todo)
				if tt.assertTodo != nil {
					tt.assertTodo(t, todo)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_ChangeStatus(t *testing.T) {
	ctx := context.Background()
	validID := validUUID()

	tests := []struct {
		name          string
		id            string
		newStatus     Status
		initialStatus Status
		setupMock     func(*MockRepository, *Todo)
		wantErr       bool
		assertErr     func(*testing.T, error)
	}{
		{
			name:          "successful status change",
			id:            validID,
			newStatus:     StatusInProgress,
			initialStatus: StatusPending,
			setupMock: func(m *MockRepository, todo *Todo) {
				m.On("Get", ctx, validID).Return(todo, nil)
				m.On("Update", ctx, mock.AnythingOfType("*todo.Todo")).Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "empty id returns error",
			id:        "",
			newStatus: StatusCompleted,
			setupMock: func(m *MockRepository, todo *Todo) {},
			wantErr:   true,
			assertErr: func(t *testing.T, err error) {
				require.ErrorIs(t, err, ErrInvalidInput)
			},
		},
		{
			name:      "invalid uuid format returns error",
			id:        "invalid-uuid",
			newStatus: StatusCompleted,
			setupMock: func(m *MockRepository, todo *Todo) {},
			wantErr:   true,
			assertErr: func(t *testing.T, err error) {
				require.ErrorIs(t, err, ErrInvalidInput)
			},
		},
		{
			name:      "todo not found",
			id:        validID,
			newStatus: StatusCompleted,
			setupMock: func(m *MockRepository, todo *Todo) {
				m.On("Get", ctx, validID).Return(nil, ErrNotFound)
			},
			wantErr: true,
			assertErr: func(t *testing.T, err error) {
				require.ErrorIs(t, err, ErrNotFound)
			},
		},
		{
			name:          "invalid status transition",
			id:            validID,
			newStatus:     StatusBlocked,
			initialStatus: StatusCompleted,
			setupMock: func(m *MockRepository, todo *Todo) {
				m.On("Get", ctx, validID).Return(todo, nil)
			},
			wantErr: true,
			assertErr: func(t *testing.T, err error) {
				require.ErrorIs(t, err, ErrInvalidStatus)
			},
		},
		{
			name:          "invalid status value",
			id:            validID,
			newStatus:     Status("invalid"),
			initialStatus: StatusPending,
			setupMock: func(m *MockRepository, todo *Todo) {
				m.On("Get", ctx, validID).Return(todo, nil)
			},
			wantErr: true,
			assertErr: func(t *testing.T, err error) {
				require.ErrorIs(t, err, ErrInvalidStatus)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := newTestService(t)

			var testTodo *Todo
			if tt.initialStatus != "" {
				testTodo = newValidTodo(t)
				testTodo.Status = tt.initialStatus
			}

			tt.setupMock(mockRepo, testTodo)

			todo, err := service.ChangeStatus(ctx, tt.id, tt.newStatus)

			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, todo)
				if tt.assertErr != nil {
					tt.assertErr(t, err)
				}
			} else {
				require.NoError(t, err)
				require.NotNil(t, todo)
				require.Equal(t, tt.newStatus, todo.Status)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_DeleteTodo(t *testing.T) {
	ctx := context.Background()
	validID := validUUID()

	tests := []struct {
		name      string
		id        string
		setupMock func(*MockRepository)
		wantErr   bool
		assertErr func(*testing.T, error)
	}{
		{
			name: "successful delete",
			id:   validID,
			setupMock: func(m *MockRepository) {
				m.On("Delete", ctx, validID).Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "empty id returns error",
			id:        "",
			setupMock: func(m *MockRepository) {},
			wantErr:   true,
			assertErr: func(t *testing.T, err error) {
				require.ErrorIs(t, err, ErrInvalidInput)
			},
		},
		{
			name:      "invalid uuid format returns error",
			id:        "invalid-uuid",
			setupMock: func(m *MockRepository) {},
			wantErr:   true,
			assertErr: func(t *testing.T, err error) {
				require.ErrorIs(t, err, ErrInvalidInput)
			},
		},
		{
			name: "not found",
			id:   validID,
			setupMock: func(m *MockRepository) {
				m.On("Delete", ctx, validID).Return(ErrNotFound)
			},
			wantErr: true,
			assertErr: func(t *testing.T, err error) {
				require.ErrorIs(t, err, ErrNotFound)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := newTestService(t)
			tt.setupMock(mockRepo)

			err := service.DeleteTodo(ctx, tt.id)

			if tt.wantErr {
				require.Error(t, err)
				if tt.assertErr != nil {
					tt.assertErr(t, err)
				}
			} else {
				require.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_ListTodos(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name         string
		filter       ListFilter
		setupMock    func(*MockRepository)
		wantErr      bool
		assertErr    func(*testing.T, error)
		assertTodos  func(*testing.T, []*Todo)
		assertFilter func(*testing.T, ListFilter) // To verify filter modifications
	}{
		{
			name:   "successful list with defaults",
			filter: ListFilter{},
			setupMock: func(m *MockRepository) {
				expectedTodos := []*Todo{
					{Title: "Todo 1", Status: StatusPending},
					{Title: "Todo 2", Status: StatusCompleted},
				}
				m.On("List", ctx, mock.MatchedBy(func(filter ListFilter) bool {
					return filter.Limit == 50 &&
						filter.SortBy == SortFieldCreateTime &&
						filter.SortOrder == SortOrderDesc
				})).Return(expectedTodos, nil)
			},
			wantErr: false,
			assertTodos: func(t *testing.T, todos []*Todo) {
				require.Len(t, todos, 2)
				require.Equal(t, "Todo 1", todos[0].Title)
			},
		},
		{
			name:   "limit capped at 1000",
			filter: ListFilter{Limit: 5000},
			setupMock: func(m *MockRepository) {
				m.On("List", ctx, mock.MatchedBy(func(filter ListFilter) bool {
					return filter.Limit == 1000
				})).Return([]*Todo{}, nil)
			},
			wantErr: false,
		},
		{
			name: "custom filter applied",
			filter: ListFilter{
				Status:    StatusPending,
				Limit:     100,
				SortBy:    SortFieldTitle,
				SortOrder: SortOrderAsc,
			},
			setupMock: func(m *MockRepository) {
				m.On("List", ctx, mock.MatchedBy(func(filter ListFilter) bool {
					return filter.Status == StatusPending &&
						filter.Limit == 100 &&
						filter.SortBy == SortFieldTitle &&
						filter.SortOrder == SortOrderAsc
				})).Return([]*Todo{}, nil)
			},
			wantErr: false,
		},
		{
			name:      "invalid filter returns error",
			filter:    ListFilter{Status: Status("invalid")},
			setupMock: func(m *MockRepository) {},
			wantErr:   true,
		},
		{
			name:   "repository error",
			filter: ListFilter{},
			setupMock: func(m *MockRepository) {
				m.On("List", ctx, mock.Anything).Return(nil, errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := newTestService(t)
			tt.setupMock(mockRepo)

			todos, err := service.ListTodos(ctx, tt.filter)

			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, todos)
				if tt.assertErr != nil {
					tt.assertErr(t, err)
				}
			} else {
				require.NoError(t, err)
				require.NotNil(t, todos)
				if tt.assertTodos != nil {
					tt.assertTodos(t, todos)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_CountTodos(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		filter      ListFilter
		setupMock   func(*MockRepository)
		wantErr     bool
		assertErr   func(*testing.T, error)
		assertCount func(*testing.T, int)
	}{
		{
			name:   "successful count",
			filter: ListFilter{},
			setupMock: func(m *MockRepository) {
				m.On("Count", ctx, mock.Anything).Return(42, nil)
			},
			wantErr: false,
			assertCount: func(t *testing.T, count int) {
				require.Equal(t, 42, count)
			},
		},
		{
			name:   "count with filter",
			filter: ListFilter{Status: StatusCompleted},
			setupMock: func(m *MockRepository) {
				m.On("Count", ctx, mock.MatchedBy(func(filter ListFilter) bool {
					return filter.Status == StatusCompleted
				})).Return(10, nil)
			},
			wantErr: false,
			assertCount: func(t *testing.T, count int) {
				require.Equal(t, 10, count)
			},
		},
		{
			name:      "invalid filter returns error",
			filter:    ListFilter{Status: Status("invalid")},
			setupMock: func(m *MockRepository) {},
			wantErr:   true,
		},
		{
			name:   "repository error",
			filter: ListFilter{},
			setupMock: func(m *MockRepository) {
				m.On("Count", ctx, mock.Anything).Return(0, errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := newTestService(t)
			tt.setupMock(mockRepo)

			count, err := service.CountTodos(ctx, tt.filter)

			if tt.wantErr {
				require.Error(t, err)
				require.Equal(t, 0, count)
				if tt.assertErr != nil {
					tt.assertErr(t, err)
				}
			} else {
				require.NoError(t, err)
				if tt.assertCount != nil {
					tt.assertCount(t, count)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
