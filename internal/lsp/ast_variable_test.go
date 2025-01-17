package lsp

import (
	"context"
	"testing"

	"github.com/mrjosh/helm-ls/internal/tree-sitter/gotemplate"
	sitter "github.com/smacker/go-tree-sitter"
)

func TestGetVariableDefinitionDirectDecleration(t *testing.T) {

	var template = `
{{ $variable := "text" }}
{{ $variable }}
	`

	node, err := sitter.ParseCtx(context.Background(), []byte(template), gotemplate.GetLanguage())

	if err != nil {
		t.Errorf("Parsing did not work")
	}

	usageNode := node.NamedChild(3)
	definitionNode := GetVariableDefinitionOfNode(usageNode, template)

	if definitionNode == nil {
		t.Errorf("Could not get definitionNode")
	} else if definitionNode.Content([]byte(template)) != "$variable := \"text\"" {
		t.Errorf("Definition did not match but was %s", definitionNode.Content([]byte(template)))
	}
}

func TestGetVariableDefinitionOtherDecleration(t *testing.T) {
	var template = `
{{ $variable := "text" }}
{{ $someOther := "text" }}
{{ $variable }}
	`

	node, err := sitter.ParseCtx(context.Background(), []byte(template), gotemplate.GetLanguage())

	if err != nil {
		t.Errorf("Parsing did not work")
	}

	usageNode := node.NamedChild(5)
	definitionNode := GetVariableDefinitionOfNode(usageNode, template)

	if definitionNode == nil {
		t.Errorf("Could not get definitionNode")
	} else if definitionNode.Content([]byte(template)) != "$variable := \"text\"" {
		t.Errorf("Definition did not match but was %s", definitionNode.Content([]byte(template)))
	}

}

func TestGetVariableDefinitionRange(t *testing.T) {
	var template = `
{{ range $index, $element := pipeline }}{{ $index }}{{ $element }}{{ end }}
	`

	node, err := sitter.ParseCtx(context.Background(), []byte(template), gotemplate.GetLanguage())

	if err != nil {
		t.Errorf("Parsing did not work")
	}

	elementUsageNode := node.NamedChild(1).NamedChild(2)
	logger.Println(elementUsageNode.Content([]byte(template))) // $element
	definitionNode := GetVariableDefinitionOfNode(elementUsageNode, template)

	if definitionNode == nil {
		t.Errorf("Could not get definitionNode")
	} else if definitionNode.Content([]byte(template)) != "$element" {
		t.Errorf("Definition did not match but was %s", definitionNode.Content([]byte(template)))
	}

	indexUsageNode := node.NamedChild(1).NamedChild(1)
	logger.Println(elementUsageNode.Content([]byte(template))) // $index
	definitionNode = GetVariableDefinitionOfNode(indexUsageNode, template)

	if definitionNode == nil {
		t.Errorf("Could not get definitionNode")
	} else if definitionNode.Content([]byte(template)) != "$index" {
		t.Errorf("Definition did not match but was %s", definitionNode.Content([]byte(template)))
	}
}
