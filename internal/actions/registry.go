package actions

import "context"

import "github.com/ecuware/pmgwire/internal/pmg"

type Data map[string]interface{}

type Params map[string]interface{}

type Action interface {
	Name() string
	Execute(ctx context.Context, client *pmg.Client, input Data, params Params, filters map[string]string) (Data, error)
}

var registry = make(map[string]Action)

func Register(a Action) {
	registry[a.Name()] = a
}

func Get(name string) (Action, bool) {
	a, ok := registry[name]
	return a, ok
}

func All() map[string]Action {
	return registry
}