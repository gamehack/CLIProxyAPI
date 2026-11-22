package pluginhost

import (
	"encoding/json"
	"strings"
)

const nativeCopilotGitHubBaseURL = "https://github.com"

// adaptQuotaStorageJSON maps the native Copilot credential fields to the names
// expected by quota plugins such as cpa-subscription-tools.
//
// Native CLIProxyAPI Copilot auth stores the GitHub OAuth token in refresh_token;
// access_token is a separate short-lived Copilot inference token and must not be
// used for GitHub account or quota APIs. The returned JSON is passed only to the
// quota provider and is never persisted back to the credential file.
func adaptQuotaStorageJSON(provider string, raw []byte) []byte {
	if normalizeProviderID(provider) != "copilot" || len(raw) == 0 {
		return raw
	}

	var storage map[string]json.RawMessage
	if err := json.Unmarshal(raw, &storage); err != nil || storage == nil {
		return raw
	}

	var authType string
	if err := json.Unmarshal(storage["type"], &authType); err != nil ||
		!strings.EqualFold(strings.TrimSpace(authType), "copilot") {
		return raw
	}

	var githubAccessToken string
	if value, ok := storage["github_access_token"]; ok {
		_ = json.Unmarshal(value, &githubAccessToken)
	}
	if strings.TrimSpace(githubAccessToken) != "" {
		return raw
	}

	var githubOAuthToken string
	if value, ok := storage["refresh_token"]; !ok ||
		json.Unmarshal(value, &githubOAuthToken) != nil ||
		strings.TrimSpace(githubOAuthToken) == "" {
		return raw
	}

	storage["github_access_token"], _ = json.Marshal(strings.TrimSpace(githubOAuthToken))

	var githubBaseURL string
	if value, ok := storage["github_base_url"]; ok {
		_ = json.Unmarshal(value, &githubBaseURL)
	}
	if strings.TrimSpace(githubBaseURL) == "" {
		storage["github_base_url"], _ = json.Marshal(nativeCopilotGitHubBaseURL)
	}

	adapted, err := json.Marshal(storage)
	if err != nil {
		return raw
	}
	return adapted
}
