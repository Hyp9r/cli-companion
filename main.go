package main

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"
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

// QUERY u update functionu

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

type Field struct {
	Name string
	Type string
}

type Template struct {
	Filename       string
	ModelName      string
	Module         string
	Path           string
	ModelFields    []Field
	DatabaseFields []string
	FieldSize      int
}

func main() {
	// Clear the screen
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	// Setup scanner
	scanner := bufio.NewScanner(os.Stdin)

	// Setup colors
	info := color.New(color.FgHiYellow).Add(color.Underline)
	question := color.New(color.FgCyan).Add(color.Underline)
	answer := color.New(color.FgGreen).Add(color.Underline)

	info.Println("Welcome database generation cli")

	question.Println("Provide github repository or module")
	scanner.Scan()
	module := scanner.Text()

	question.Println("How would you like to name your model? ")
	scanner.Scan()
	modelName := scanner.Text()
	fileName := toLowerCase(modelName)

	// Initialize folder structure
	repositoryInfraFolder := fmt.Sprintf("%s/%s", REPOSITORY_INFRASTRUCTRUE_PATH, fileName)
	domainModelFolder := fmt.Sprintf("%s/%s/model", DOMAIN_PATH, fileName)

	// Create folders
	err := makeDirectories(repositoryInfraFolder)
	if err != nil {
		panic(err)
	}
	err = makeDirectories(domainModelFolder)
	if err != nil {
		panic(err)
	}

	var modelFields []Field
	var databaseFields []string
	fieldSize := 0

	for {
		question.Println("Enter the name of the field(ENTER to finish): ")
		scanner.Scan()
		operation := scanner.Text()
		// Check if user is done
		if operation == INPUT_DONE {
			break
		}
		question.Println("Enter the type of the field: ")
		scanner.Scan()
		fieldType := scanner.Text()
		field := Field{
			Name: operation,
			Type: fieldType,
		}
		modelFields = append(modelFields, field)
		databaseFields = append(databaseFields, toSnakeCase(operation))
		fieldSize++
	}

	// Create a new template model
	templateModel := Template{
		Filename:       fmt.Sprintf("%s.go", fileName),
		ModelName:      modelName,
		Module:         module,
		Path:           fileName,
		DatabaseFields: databaseFields,
		ModelFields:    modelFields,
		FieldSize:      fieldSize - 1,
	}

	// run the generation
	repoFileContents := generate(templateModel, "templates/repository.tmpl")
	modelFileContents := generate(templateModel, "templates/model.tmpl")
	err = saveFile(repoFileContents, fmt.Sprintf("%s/%s", repositoryInfraFolder, "repository.go"))
	if err != nil {
		panic(err)
	}
	err = saveFile(modelFileContents, fmt.Sprintf("%s/%s.go", domainModelFolder, fileName))
	if err != nil {
		panic(err)
	}

	answer.Printf("%s.go model created\n", fileName)
}

func generate(model Template, pathToTemplate string) string {
	var buf bytes.Buffer
	tmpl, err := template.ParseFiles(pathToTemplate)
	if err != nil {
		panic(err)
	}

	// Execute the template with the person data
	err = tmpl.Execute(&buf, model)
	if err != nil {
		panic(err)
	}

	return buf.String()
}

func saveFile(content string, filePath string) error {
	// Write the content to the file
	err := ioutil.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil
}

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func toLowerCase(str string) string {
	firstChar := rune(str[0])
	if unicode.IsUpper(firstChar) {
		firstChar = unicode.ToLower(firstChar)
	}
	return string(firstChar) + str[1:]
}

func makeDirectories(path string) error {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return err
	}
	return nil
}
