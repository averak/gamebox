package main

import (
	"github.com/averak/gamebox/cmd/protoc-gen-gamebox-server/handler"
	"github.com/averak/gamebox/protobuf/custom_option"
	"google.golang.org/protobuf/compiler/protogen"
)

func main() {
	generator := handler.NewGenerator[*custom_option.MethodOption](handler.Config{
		MethodOptExt:      custom_option.E_MethodOption,
		MethodOptIdent:    protogen.GoImportPath("github.com/averak/gamebox/app/infrastructure/connect/aop").Ident("MethodOption"),
		MethodOptExtIdent: protogen.GoImportPath("github.com/averak/gamebox/protobuf/custom_option").Ident("E_MethodOption"),
		MethodErrDefIdent: protogen.GoImportPath("github.com/averak/gamebox/app/infrastructure/connect/aop").Ident("MethodErrDefinition"),
		ProxyIdent:        protogen.GoImportPath("github.com/averak/gamebox/app/infrastructure/connect/aop").Ident("Proxy"),
	})
	protogen.Options{}.Run(func(plugin *protogen.Plugin) error {
		for _, file := range plugin.Files {
			if file.Desc.Package() != "api" && file.Desc.Package() != "api.debug" {
				continue
			}
			err := generator.Generate(plugin, file)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
