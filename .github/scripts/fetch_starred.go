package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Repo struct {
	FullName    string `json:"full_name"`
	HTMLURL     string `json:"html_url"`
	Description string `json:"description"`
	Language    string `json:"language"`
}

type SimpleRepo struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	Description string `json:"description"`
	Language    string `json:"language"`
	Category    string `json:"category"`
}

func fetchStarred(username, token string) ([]SimpleRepo, error) {
	client := &http.Client{}
	var result []SimpleRepo
	page := 1
	perPage := 100

	for {
		url := fmt.Sprintf("https://api.github.com/users/%s/starred?per_page=%d&page=%d", username, perPage, page)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "token "+token)
		req.Header.Set("Accept", "application/vnd.github.v3+json")

		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(resp.Body)

		body, _ := io.ReadAll(resp.Body)
		var repos []Repo
		if err := json.Unmarshal(body, &repos); err != nil {
			return nil, err
		}
		if len(repos) == 0 {
			break
		}

		for _, r := range repos {
			// Default to "Unknown" if language is not specified
			language := r.Language
			if language == "" {
				language = "Unknown"
			}

			result = append(result, SimpleRepo{
				Name:        r.FullName,
				URL:         r.HTMLURL,
				Description: r.Description,
				Language:    language,
				Category:    "Uncategorized",
			})
		}
		page++
	}

	return result, nil
}

// organizeReposByLanguage groups repositories by language and sorts them alphabetically by name
func organizeReposByLanguage(repos []SimpleRepo) map[string][]SimpleRepo {
	// Create a map to hold repositories grouped by language
	reposByLanguage := make(map[string][]SimpleRepo)

	// Group repositories by language
	for _, repo := range repos {
		reposByLanguage[repo.Language] = append(reposByLanguage[repo.Language], repo)
	}

	// Sort repositories within each language alphabetically by name
	for language, langRepos := range reposByLanguage {
		sort.Slice(langRepos, func(i, j int) bool {
			return langRepos[i].Name < langRepos[j].Name
		})
		reposByLanguage[language] = langRepos
	}

	return reposByLanguage
}

// ensureDirectoryExists creates a directory if it doesn't exist
func ensureDirectoryExists(path string) error {
	return os.MkdirAll(path, 0755)
}

func writeJSON(repos []SimpleRepo) error {
	// Write the complete JSON file with all repositories
	file, err := os.Create("starred.json")
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(repos)
}

func writeMarkdown(repos []SimpleRepo, username string) error {
	// Organize repositories by language
	reposByLanguage := organizeReposByLanguage(repos)

	// Create the base stars directory
	starsDir := "stars"
	if err := ensureDirectoryExists(starsDir); err != nil {
		return err
	}

	// Generate timestamp for all files
	timestamp := time.Now().Format(time.RFC3339)

	// Create a Markdown file for each language
	for language, langRepos := range reposByLanguage {
		// Create a language directory
		langDir := filepath.Join(starsDir, strings.ToLower(language))
		if err := ensureDirectoryExists(langDir); err != nil {
			return err
		}

		// Build markdown content
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("# 🌟 %s Repositories Starred by @%s\n\n", language, username))
		sb.WriteString(fmt.Sprintf("Auto-generated on %s\n\n", timestamp))
		sb.WriteString("| Name | Description |\n|------|-------------|\n")

		// Add each repository to the markdown
		for _, r := range langRepos {
			desc := strings.ReplaceAll(r.Description, "\n", " ")
			sb.WriteString(fmt.Sprintf("| [%s](%s) | %s |\n", r.Name, r.URL, desc))
		}

		// Write the Markdown file
		mdFilePath := filepath.Join(langDir, "starred.md")
		if err := os.WriteFile(mdFilePath, []byte(sb.String()), 0644); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	username := os.Getenv("GITHUB_USERNAME")
	token := os.Getenv("GITHUB_TOKEN")

	if username == "" || token == "" {
		fmt.Println("Missing GITHUB_USERNAME or GITHUB_TOKEN env vars")
		os.Exit(1)
	}

	repos, err := fetchStarred(username, token)
	if err != nil {
		panic(err)
	}

	if err := writeJSON(repos); err != nil {
		panic(err)
	}

	if err := writeMarkdown(repos, username); err != nil {
		panic(err)
	}
}
