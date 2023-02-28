package path_test

import (
	"github.com/sasswart/gin-in-a-can/openapi"
	"github.com/sasswart/gin-in-a-can/openapi/operation"
	"github.com/sasswart/gin-in-a-can/openapi/path"
	"github.com/sasswart/gin-in-a-can/tree"
	"strings"
	"testing"
)

func TestOpenAPI_PathItem_GetAndSetName(t *testing.T) {
	want := "testName"
	i := path.Item{Node: tree.Node{Name: "this should be overwritten by a call to SetName()"}}
	i.SetName(want)
	got := strings.Join(i.GetName(), "")
	if got != want {
		t.Fail()
	}
}

func Test_Path_GetBasePath(t *testing.T) {
	i := path.Item{
		Ref:  "./relative/reference/path.yaml",
		Node: tree.Node{},
	}
	o := openapi.OpenAPI{
		Node: tree.Node{},
	}
	o.SetBasePath("/this/base/path")
	o.SetChild("/testEndpoint", &i)
	i.SetParent(&o)
	want := "/this/base/path/relative/reference"
	got := i.GetBasePath()
	if got != want {
		t.Fail()
	}

}

func TestOpenAPI_PathItem_GetAndSetChildren(t *testing.T) {
	i := path.Item{}
	i.SetChild("post", &operation.Operation{})
	i.SetChild("get", &operation.Operation{})
	i.SetChild("patch", &operation.Operation{})
	i.SetChild("delete", &operation.Operation{})
	for method, op := range i.GetChildren() {
		switch method {
		case "post":
		case "get":
		case "patch":
		case "delete":
			o, ok := op.(*operation.Operation)
			if !ok {
				t.Errorf("PathItem.operations() is not successfully returning *operations")
			}
			t.Logf(o.OperationId)
		}
	}
}

func TestOpenAPI_PathItem_GetRef(t *testing.T) {
	want := "testRef"
	i := path.Item{Ref: want}
	got := i.GetRef()
	if got != want {
		t.Fail()
	}
}
