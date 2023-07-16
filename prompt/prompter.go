package prompt

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/fatih/color"
)

const (
	INPUT_DONE                     = ""
	REPOSITORY_INFRASTRUCTRUE_PATH = "./infrastructure"
	DOMAIN_PATH                    = "./domain"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

type Prompter struct {
	scanner      *bufio.Scanner
	prompterType PrompterType
}

func NewPrompter(scanner *bufio.Scanner, prompterType PrompterType) *Prompter {
	return &Prompter{
		scanner:      scanner,
		prompterType: prompterType,
	}
}

func (p *Prompter) Start() (PromptResult, error) {
	// Setup colors
	info := color.New(color.FgHiYellow).Add(color.Underline)
	question := color.New(color.FgCyan).Add(color.Underline)

	info.Println("Welcome database generation cli")

	question.Println("Provide github repository or module")
	p.scanner.Scan()
	module := p.scanner.Text()

	switch p.prompterType.Operation {
	case "make":
		return p.ResolveNamespace(p.prompterType.Namespace, module)
	case "list":
		return PromptResult{}, nil
	default:
		return PromptResult{}, fmt.Errorf("unrecognized prompt operation")
	}
}

func (p *Prompter) ResolveNamespace(namespace string, module string) (PromptResult, error) {
	switch namespace {
	case "model":
		templateModel, err := p.ModelPrompt(module)
		if err != nil {
			return PromptResult{}, nil
		}
		return templateModel, nil
	default:
		return PromptResult{}, fmt.Errorf("unrecognized namespace")
	}
}

func (p *Prompter) ModelPrompt(module string) (PromptResult, error) {
	// Setup colors
	question := color.New(color.FgCyan).Add(color.Underline)
	question.Println("How would you like to name your model? ")
	p.scanner.Scan()
	modelName := p.scanner.Text()
	fileName := toLowerCase(modelName)

	// Initialize folder structure
	repositoryInfraFolder := fmt.Sprintf("%s/%s", REPOSITORY_INFRASTRUCTRUE_PATH, fileName)
	domainModelFolder := fmt.Sprintf("%s/%s/model", DOMAIN_PATH, fileName)

	var modelFields []Field
	var databaseFields []string
	fieldSize := 0

	for {
		question.Println("Enter the name of the field(ENTER to finish): ")
		p.scanner.Scan()
		operation := p.scanner.Text()
		// Check if user is done
		if operation == INPUT_DONE {
			break
		}
		question.Println("Enter the type of the field: ")
		p.scanner.Scan()
		fieldType := p.scanner.Text()
		valid := allowedType(fieldType)
		if !valid {
			return PromptResult{}, fmt.Errorf("invalid field type")
		}
		field := Field{
			Name: operation,
			Type: fieldType,
		}
		modelFields = append(modelFields, field)
		databaseFields = append(databaseFields, toSnakeCase(operation))
		fieldSize++
	}

	// Create a new template model
	templateModel := PromptResult{
		ModelResult: &ModelResult{
			Filename:         fmt.Sprintf("%s.go", fileName),
			ModelName:        modelName,
			Module:           module,
			Path:             fileName,
			DatabaseFields:   databaseFields,
			ModelFields:      modelFields,
			FieldSize:        fieldSize - 1,
			RepositoryFolder: repositoryInfraFolder,
			ModelFolder:      domainModelFolder,
		},
		ControllerResult: nil,
	}
	return templateModel, nil
}

func (p *Prompter) ListPrompt() {
	fmt.Println("These operations are available")
}

func toLowerCase(str string) string {
	firstChar := rune(str[0])
	if unicode.IsUpper(firstChar) {
		firstChar = unicode.ToLower(firstChar)
	}
	return string(firstChar) + str[1:]
}

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
