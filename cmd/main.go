package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/Hyp9r/cli-companion/gen"
	"github.com/Hyp9r/cli-companion/prompt"
	"github.com/fatih/color"
)

const (
	NAMESPACE = "namespace"
	OPERATION = "operation"
)

type CMDArguments struct {
	Operation string
	Namespace string
}

func main() {
	// Clear the screen
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	// Load CMD arguments
	args := loadArguments()
	err := validateArguments(args)
	if err != nil {
		panic(err)
	}

	// Find the path of the executable so we can determine template locations
	executablePath, err := findPathOfExecutable()
	if err != nil {
		panic(err)
	}

	// String kung-fu to slash out the last part of the path
	lastSlashIndex := strings.LastIndex(executablePath, "/")

	// Setup paths of templates
	templatesPath := executablePath[:lastSlashIndex]
	repositoryTemplatePath := fmt.Sprintf("%s/templates/%s", templatesPath, "repository.tmpl")
	modelTemplatePath := fmt.Sprintf("%s/templates/%s", templatesPath, "model.tmpl")

	// Setup scanner
	scanner := bufio.NewScanner(os.Stdin)

	// Create prompter and start it
	prompter := prompt.NewPrompter(scanner, prompt.PrompterType{Operation: args.Operation, Namespace: args.Namespace})
	promptResult, err := prompter.Start()
	if err != nil {
		panic(err)
	}

	// Create generator and generate code
	generator := gen.NewGenerator(repositoryTemplatePath, modelTemplatePath)
	code, err := generator.Generate(promptResult)
	if err != nil {
		panic(err)
	}
	// Prepare folders
	err = generator.PrepareDirectories(promptResult.ModelResult.RepositoryFolder)
	if err != nil {
		panic(err)
	}
	err = generator.PrepareDirectories(promptResult.ModelResult.ModelFolder)
	if err != nil {
		panic(err)
	}

	// run the generation
	err = saveFile(code.RepositoryCode, fmt.Sprintf("%s/%s", promptResult.ModelResult.RepositoryFolder, "repository.go"))
	if err != nil {
		panic(err)
	}
	err = saveFile(code.ModelCode, fmt.Sprintf("%s/%s.go", promptResult.ModelResult.ModelFolder, "model.go"))
	if err != nil {
		panic(err)
	}

	answer := color.New(color.FgGreen).Add(color.Underline)
	answer.Printf("%s.go and repository.go created\n", promptResult.ModelResult.Filename)
}

func loadArguments() CMDArguments {
	operation := flag.String(OPERATION, "list", "Operation which you want to do, list to see possibilities")
	namespace := flag.String(NAMESPACE, "", "Namespace for which you want generation(model, controller)")
	flag.Parse()
	if operation != nil && len(*operation) == 0 {
		panic(fmt.Errorf("no operation provided or provided more then one"))
	}
	if namespace != nil && len(*namespace) == 0 {
		panic(fmt.Errorf("no namespace provided or provided more then one"))
	}
	return CMDArguments{
		Operation: *operation,
		Namespace: *namespace,
	}
}

func findPathOfExecutable() (string, error) {
	executablePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	return executablePath, nil
}

func validateArguments(args CMDArguments) error {
	var ok bool
	switch args.Operation {
	case "list":
	case "make":
		ok = true
		break
	default:
		ok = false
	}
	if !ok {
		return fmt.Errorf("invalid operation")
	}
	switch args.Namespace {
	case "model":
	case "controller":
		ok = true
		break
	default:
		ok = false
	}
	if !ok {
		return fmt.Errorf("invalid namespace")
	}
	return nil
}

func saveFile(content string, filePath string) error {
	// Write the content to the file
	err := ioutil.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil
}
