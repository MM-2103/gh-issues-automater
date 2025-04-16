package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
)

type Config struct {
	Repo     string              `json:"repo"`
	Keywords map[string]CrudOpts `json:"keywords"`
}

type CrudOpts struct {
	Create bool `json:"create"`
	Read   bool `json:"read"`
	Update bool `json:"update"`
	Delete bool `json:"delete"`
}

type OpenRouterModel struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type OpenRouterResponse struct {
	Data []OpenRouterModel `json:"data"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: issue-creator <config-file.json> or --list-models openrouter/")
		os.Exit(1)
	}

	// Check for command line flags
	if os.Args[1] == "--list-models" && len(os.Args) > 2 && os.Args[2] == "openrouter/" {
		listOpenRouterModels()
		return
	}

	configFile := os.Args[1]
	config, err := loadConfig(configFile)
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	for keyword, crudOpts := range config.Keywords {
		if hasCRUD(crudOpts) {
			createCRUDIssues(config.Repo, keyword, crudOpts)
		} else {
			createBaseIssue(config.Repo, keyword)
		}
	}
}

func loadConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(file, &config); err != nil {
		return nil, err
	}

	if config.Repo == "" {
		return nil, fmt.Errorf("repository not specified in config")
	}

	return &config, nil
}

func hasCRUD(opts CrudOpts) bool {
	return opts.Create || opts.Read || opts.Update || opts.Delete
}

func createCRUDIssues(repo, keyword string, opts CrudOpts) {
	crudOperations := map[string]bool{
		"aanmaken":    opts.Create,
		"lezen":       opts.Read,
		"updaten":     opts.Update,
		"verwijderen": opts.Delete,
	}

	for operation, enabled := range crudOperations {
		if enabled {
			title := fmt.Sprintf("Als dev-team wil ik %s kunnen %s", keyword, operation)
			createIssue(repo, title)
		}
	}
}

func createBaseIssue(repo, keyword string) {
	title := fmt.Sprintf("Als dev-team wil ik %s hebben", keyword)
	createIssue(repo, title)
}

func createIssue(repo, title string) {
	body := ".github/ISSUE_TEMPLATE/issue-template.md"

	cmd := exec.Command("gh", "issue", "create",
		"--repo", repo,
		"--title", title,
		"-F", body,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error creating issue '%s': %v\nOutput: %s\n", title, err, string(output))
		return
	}

	fmt.Printf("Successfully created issue: %s\n", title)
}

func listOpenRouterModels() {
	url := "https://openrouter.ai/api/v1/models"
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}
	
	// Add API key if available
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey != "" {
		req.Header.Add("Authorization", "Bearer "+apiKey)
	}
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error fetching models: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: received status code %d\n", resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		fmt.Println(string(body))
		return
	}
	
	var openRouterResp OpenRouterResponse
	if err := json.NewDecoder(resp.Body).Decode(&openRouterResp); err != nil {
		fmt.Printf("Error decoding response: %v\n", err)
		return
	}
	
	fmt.Println("Available OpenRouter Models:")
	fmt.Println("===========================")
	for _, model := range openRouterResp.Data {
		fmt.Printf("ID: %s\nName: %s\nDescription: %s\n\n", model.ID, model.Name, model.Description)
	}
}
