package render_test

import (
	"github.com/sasswart/gin-in-a-can/config"
	"github.com/sasswart/gin-in-a-can/openapi"
	"github.com/sasswart/gin-in-a-can/openapi/operation"
	"github.com/sasswart/gin-in-a-can/openapi/path"
	"github.com/sasswart/gin-in-a-can/openapi/request"
	"github.com/sasswart/gin-in-a-can/render"
	golang "github.com/sasswart/gin-in-a-can/render/go"
	"github.com/sasswart/gin-in-a-can/tree"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"os"
	"strings"
	"testing"
	"text/template"
)

var (
	md       = tree.Metadata{"package": "testPackage", "some": "metadata"}
	toRender = buildTestSpec()
)

func buildTestSpec() *openapi.OpenAPI {

	root := openapi.OpenAPI{
		Node: tree.Node{Name: "openapi"},
	}
	root.SetMetadata(md)
	p := path.Item{
		Node: tree.Node{Name: "pathitem"},
	}
	root.SetChild("/endpoint", &p)
	p.SetParent(&root)

	get := operation.Operation{
		Node: tree.Node{Name: "operation"},
	}
	p.SetChild("get", &get)
	get.SetParent(&p)

	post := operation.Operation{
		Node: tree.Node{Name: "pathitem2"},
	}
	p.SetChild("post", &post)
	post.SetParent(&p)

	requestBody := request.Body{
		Node: tree.Node{Name: "requestbody"},
	}
	post.SetChild("Body", &requestBody)
	requestBody.SetParent(&post)

	requestBody2 := request.Body{
		Node: tree.Node{Name: "requestbody2"},
	}
	get.SetChild("Body", &requestBody2)
	requestBody2.SetParent(&get)

	return &root
}

func Test_Render_Render(t *testing.T) {
	tempFolder, err := os.MkdirTemp(os.TempDir(), "CanTestArtifacts")
	defer os.RemoveAll(tempFolder)

	cfg := newTestConfig()
	err = cfg.Load()
	if err != nil {
		t.Errorf(err.Error())
	}
	cfg.OutputPath = tempFolder
	e := render.Engine{}.With(&golang.Renderer{Base: &render.Base{}}, cfg)
	r := *e.GetRenderer()
	r.SetTemplateFuncMap(&template.FuncMap{
		"ToUpper": strings.ToUpper,
		"ToTitle": func(s string) string {
			caser := cases.Title(language.English)
			return caser.String(s)
		},
		"SanitiseName": r.SanitiseName,
		"SanitiseType": r.SanitiseType,
	})
	_, err = tree.Traverse(toRender, e.BuildRenderNode())
	if err != nil {
		t.Errorf(err.Error())
	}
}
func newTestConfig() config.Data {
	config.ConfigFilePath = "../config/config_test.yml"
	config.Debug = true
	return config.Data{
		Generator: config.Generator{},
		Template: config.Template{
			Name: "go-gin",
		},
		TemplatesDir: "../templates",
		OpenAPIFile:  "../openapi/test/fixtures/validation_no_refs.yaml",
		OutputPath:   ".",
	}
}
