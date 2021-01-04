package repository

import (
	"context"

	"github.com/dannypaul/go-skeleton/internal/primitive"
)

type Copier interface {
	Copy(destination interface{}) error
}

type ListCopier interface {
	CopyAll(ctx context.Context, destination interface{}) error
}

type Creator interface {
	Create(ctx context.Context, model interface{}) (Copier, error)
}

type Setter interface {
	Set(ctx context.Context, filters []Filter, Key string, Value interface{}) error
	SetById(ctx context.Context, id primitive.Id, Key string, Value interface{}) error
	SetAll(ctx context.Context, filters []Filter, setters []KeyValue) error
	SetAllById(ctx context.Context, id primitive.Id, setters []KeyValue) error
	UnSet(ctx context.Context, filters []Filter, Key string) error
}

type Patcher interface {
	Patch(ctx context.Context, id primitive.Id, patches []Patch) error
}

type Counter interface {
	Count(ctx context.Context, filters []Filter) (int64, error)
}

type Incrementer interface {
	IncrementById(ctx context.Context, id primitive.Id, Key string, incrementBy int) error
}

type Deleter interface {
	Delete(ctx context.Context, id primitive.Id) (int64, error)
}

type Adder interface {
	Add(ctx context.Context, filters []Filter, key string) (float64, error)
}

type Finder interface {
	FindById(ctx context.Context, id primitive.Id) (Copier, error)
	FindSingle(ctx context.Context, filters []Filter) (Copier, error)
}

type Replacer interface {
	Replace(ctx context.Context, id primitive.Id, Value interface{}) (Copier, error)
}

type Patch struct {
	Action string
	Key    string
	Value  interface{}
}

type KeyValue struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

type Filter struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

type Page struct {
	Skip  int64
	Limit int64
	Total int64
}
