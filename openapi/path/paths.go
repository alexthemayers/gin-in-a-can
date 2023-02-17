package path

import (
	"github.com/sasswart/gin-in-a-can/openapi/errors"
	"github.com/sasswart/gin-in-a-can/openapi/operation"
	"github.com/sasswart/gin-in-a-can/openapi/parameter"
	"github.com/sasswart/gin-in-a-can/tree"
	"path/filepath"
)

var _ tree.NodeTraverser = &Item{}

// PathItem is a programmatic representation of the Path Item object defined here: https://swagger.io/specification/#path-item-object
type Item struct {
	tree.Node
	Ref         string `yaml:"$ref"` // must be defined in the format of a PathItem object
	Summary     string
	Description string
	Get         *operation.Operation
	Post        *operation.Operation
	Patch       *operation.Operation
	Delete      *operation.Operation
	Parameters  []parameter.Parameter
}

func (p *Item) GetChildren() map[string]tree.NodeTraverser {
	return p.Operations()
}

func (p *Item) GetRef() string {
	return p.Ref
}

func (p *Item) GetPath() string {
	name := p.GetName()
	return name
}

func (p *Item) GetBasePath() string {
	if p.GetParent() == nil {
		panic("PathItem should never be at the root of the tree")
	}
	// TODO: Deal with absolute paths for both of these parameters
	// For now both of these params are assumed relative
	basePath := filepath.Join(p.GetParent().GetBasePath(), filepath.Dir(p.Ref))
	return basePath
}

func (p *Item) Operations() map[string]tree.NodeTraverser {
	operations := map[string]tree.NodeTraverser{}
	if p.Get != nil {
		operations["get"] = p.Get
	}
	if p.Post != nil {
		operations["post"] = p.Post
	}
	if p.Patch != nil {
		operations["patch"] = p.Patch
	}
	if p.Delete != nil {
		operations["delete"] = p.Delete
	}
	return operations
}

func (p *Item) SetChild(i string, child tree.NodeTraverser) {
	if op, ok := child.(*operation.Operation); ok {
		p.Operations()[i] = op
		return

	}
	panic("(p *PathItem) setChild(): " + errors.ErrCastFail)
}