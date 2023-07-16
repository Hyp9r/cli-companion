package prompt

type PromptResult struct {
	ModelResult      *ModelResult
	ControllerResult *ControllerResult
}

type ModelResult struct {
	Filename         string
	ModelName        string
	Module           string
	Path             string
	ModelFields      []Field
	DatabaseFields   []string
	FieldSize        int
	RepositoryFolder string
	ModelFolder      string
}

type ControllerResult struct {
	Filename string
	Module   string
}
