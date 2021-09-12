package main

import (
	"fmt"
	"github.com/yandex-cloud/cq-provider-yandex/gen/util/modelfromproto"
	"path/filepath"

	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
	"github.com/yandex-cloud/cq-provider-yandex/gen/util"
)

func generate(resource string, pathToProto string) {
	out := filepath.Join(util.ResourcesDir, "resourcemanager_"+strcase.ToSnake(inflection.Plural(resource))+".go")

	b := modelfromproto.TableBuilder{
		Service: "ResourceManager",
	}

	err := b.WithMessageFromProto(resource, pathToProto, "cloudapi", "cloudapi/third_party/googleapis")
	if err != nil {
		return
	}

	tableModel, err := b.Build()
	if err != nil {
		return
	}

	tableModel.Multiplex = fmt.Sprintf("client.MultiplexBy(client.%s)", inflection.Plural(resource))

	util.SilentExecute(util.TemplatesDir{
		MainFile: "resource_manager.go.tmpl",
		Path:     "templates",
	}, map[string]interface{}{
		"resource": resource,
		"table":    tableModel,
	}, out)
}

func main() {
	generate("Cloud", "yandex/cloud/resourcemanager/v1/cloud.proto")
	generate("Folder", "yandex/cloud/resourcemanager/v1/folder.proto")
}
