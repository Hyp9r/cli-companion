package gen

import (
	"bytes"
	"html/template"
	"os"

	"github.com/Hyp9r/cli-companion/prompt"
)

type Generator struct {
	repositoryTemplatePath string
	modelTemplatePath      string
}

func NewGenerator(repPath, modelPath string) *Generator {
	return &Generator{
		repositoryTemplatePath: repPath,
		modelTemplatePath:      modelPath,
	}
}

func (g *Generator) Generate(model prompt.PromptResult) (Code, error) {
	var repBuf bytes.Buffer
	var modelBuf bytes.Buffer
	// Parse repository template
	tmpl, err := template.ParseFiles(g.repositoryTemplatePath)
	if err != nil {
		return Code{}, err
	}

	// Parse model template
	mTmpl, err := template.ParseFiles(g.modelTemplatePath)
	if err != nil {
		return Code{}, err
	}

	// Execute the template with the repository data
	err = tmpl.Execute(&repBuf, model.ModelResult)
	if err != nil {
		return Code{}, err
	}

	// Execute the template with the repository data
	err = mTmpl.Execute(&modelBuf, model.ModelResult)
	if err != nil {
		return Code{}, err
	}

	return Code{
		RepositoryCode: repBuf.String(),
		ModelCode:      modelBuf.String(),
	}, nil
}

func (g *Generator) PrepareDirectories(path string) error {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return err
	}
	return nil
}
