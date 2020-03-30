// Code generated by sysl DO NOT EDIT.
package simple

import (
	"context"

	"github.com/anz-bank/sysl/demo/examples/Code-Generation/gen/jsonplaceholder"
)

// DefaultSimpleImpl  ...
type DefaultSimpleImpl struct {
}

// NewDefaultSimpleImpl for simple
func NewDefaultSimpleImpl() *DefaultSimpleImpl {
	return &DefaultSimpleImpl{}
}

// Get Client
type GetClient struct {
}

// GetFoobarList Client
type GetFoobarListClient struct {
	GetTodos func(ctx context.Context, req *jsonplaceholder.GetTodosRequest) (*jsonplaceholder.TodosResponse, error)
}

// ServiceInterface for simple
type ServiceInterface struct {
	Get           func(ctx context.Context, req *GetRequest, client GetClient) (*Welcome, error)
	GetFoobarList func(ctx context.Context, req *GetFoobarListRequest, client GetFoobarListClient) (*jsonplaceholder.TodosResponse, error)
}