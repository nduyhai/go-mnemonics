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
	FullName    string   `json:"full_name"`
	HTMLURL     string   `json:"html_url"`
	Description string   `json:"description"`
	Language    string   `json:"language"`
	Topics      []string `json:"topics"`
}

type SimpleRepo struct {
	Name        string   `json:"name"`
	URL         string   `json:"url"`
	Description string   `json:"description"`
	Language    string   `json:"language"`
	Topics      []string `json:"topics"`
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
			// Default to "Unknown" if the language is not specified
			language := r.Language
			if language == "" {
				language = "Unknown"
			}

			result = append(result, SimpleRepo{
				Name:        r.FullName,
				URL:         r.HTMLURL,
				Description: r.Description,
				Language:    language,
				Topics:      r.Topics,
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

// formatTopics formats repository topics as clickable links (limited to 3)
func formatTopics(topics []string) string {
	var topicLinks []string
	// Limit to maximum 3 topics
	topicsToShow := topics
	if len(topicsToShow) > 3 {
		topicsToShow = topicsToShow[:3]
	}
	for _, topic := range topicsToShow {
		topicLinks = append(topicLinks, fmt.Sprintf("[%s](https://github.com/topics/%s)", topic, topic))
	}
	return strings.Join(topicLinks, ", ")
}

// generateMarkdownContent creates markdown content for a list of repositories
func generateMarkdownContent(repos []SimpleRepo, title string, username string, timestamp string, includeLanguage bool) string {
	var sb strings.Builder

	// Add header
	sb.WriteString(fmt.Sprintf("# 🌟 %s Repositories Starred by [@%s](https://github.com/%s)\n\n", title, username, username))
	sb.WriteString(fmt.Sprintf("Auto-generated on %s\n\n", timestamp))

	// Add a table header based on whether a language column is included
	if includeLanguage {
		sb.WriteString("| Name | Description | Language | Topics |\n|------|-------------|----------|-------|\n")
	} else {
		sb.WriteString("| Name | Description | Topics |\n|------|-------------|-------|\n")
	}

	// Add each repository to the markdown
	for _, r := range repos {
		desc := strings.ReplaceAll(r.Description, "\n", " ")
		topicsFormatted := formatTopics(r.Topics)

		if includeLanguage {
			sb.WriteString(fmt.Sprintf("| [%s](%s) | %s | %s | %s |\n", r.Name, r.URL, desc, r.Language, topicsFormatted))
		} else {
			sb.WriteString(fmt.Sprintf("| [%s](%s) | %s | %s |\n", r.Name, r.URL, desc, topicsFormatted))
		}
	}

	return sb.String()
}

func writeMarkdown(repos []SimpleRepo, username string, defaultLangs map[string]bool) error {
	// Organize repositories by language
	reposByLanguage := organizeReposByLanguage(repos)

	// Create the base stars directory
	starsDir := "stars"
	if err := ensureDirectoryExists(starsDir); err != nil {
		return err
	}

	// Generate timestamp for all files
	timestamp := time.Now().Format(time.RFC3339)

	// Create a map to hold repositories for the "unknown" category
	unknownRepos := make([]SimpleRepo, 0)

	// Create a Markdown file for each default language
	for language, langRepos := range reposByLanguage {
		lowerLang := strings.ToLower(language)

		// If not a default language, add to unknown
		if !defaultLangs[lowerLang] && lowerLang != "unknown" {
			unknownRepos = append(unknownRepos, langRepos...)
			continue
		}

		// Create a language directory
		langDir := filepath.Join(starsDir, lowerLang)
		if err := ensureDirectoryExists(langDir); err != nil {
			return err
		}

		// Generate markdown content and write to a file
		mdContent := generateMarkdownContent(langRepos, language, username, timestamp, false)
		mdFilePath := filepath.Join(langDir, "starred.md")
		if err := os.WriteFile(mdFilePath, []byte(mdContent), 0644); err != nil {
			return err
		}
	}

	// Add existing "unknown" repositories to our collected non-default language repos
	if unknownLangRepos, exists := reposByLanguage["Unknown"]; exists {
		unknownRepos = append(unknownRepos, unknownLangRepos...)
	}

	// Sort unknown repos by name
	sort.Slice(unknownRepos, func(i, j int) bool {
		return unknownRepos[i].Name < unknownRepos[j].Name
	})

	// If we have any unknown repos, create the unknown directory and markdown file
	if len(unknownRepos) > 0 {
		unknownDir := filepath.Join(starsDir, "unknown")
		if err := ensureDirectoryExists(unknownDir); err != nil {
			return err
		}

		// Generate markdown content and write to a file
		mdContent := generateMarkdownContent(unknownRepos, "Unknown", username, timestamp, true)
		mdFilePath := filepath.Join(unknownDir, "starred.md")
		if err := os.WriteFile(mdFilePath, []byte(mdContent), 0644); err != nil {
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

	// Get default languages from environment variable
	defaultLangsStr := os.Getenv("DEFAULT_LANGUAGES")
	if defaultLangsStr == "" {
		defaultLangsStr = "go,rust,java" // Default if isn't specified
	}

	// Parse default languages
	defaultLangs := make(map[string]bool)
	for _, lang := range strings.Split(defaultLangsStr, ",") {
		defaultLangs[strings.ToLower(strings.TrimSpace(lang))] = true
	}

	repos, err := fetchStarred(username, token)
	if err != nil {
		panic(err)
	}

	if err := writeJSON(repos); err != nil {
		panic(err)
	}

	if err := writeMarkdown(repos, username, defaultLangs); err != nil {
		panic(err)
	}
}
