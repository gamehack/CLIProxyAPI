package pluginhost

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestAdaptQuotaStorageJSON_NativeCopilot(t *testing.T) {
	raw := []byte(`{"type":"copilot","account_id":"123","username":"test-user","access_token":"copilot-inference-token","refresh_token":"github-oauth-token","base_url":"https://api.individual.githubcopilot.com"}`)

	adapted := adaptQuotaStorageJSON("copilot", raw)
	var got map[string]any
	if err := json.Unmarshal(adapted, &got); err != nil {
		t.Fatalf("adapted JSON is invalid: %v", err)
	}

	if got["github_access_token"] != "github-oauth-token" {
		t.Fatalf("github_access_token = %v, want GitHub OAuth token from refresh_token", got["github_access_token"])
	}
	if got["github_base_url"] != nativeCopilotGitHubBaseURL {
		t.Fatalf("github_base_url = %v, want %q", got["github_base_url"], nativeCopilotGitHubBaseURL)
	}
	if got["access_token"] != "copilot-inference-token" {
		t.Fatalf("Copilot inference token was changed: %v", got["access_token"])
	}
	if got["refresh_token"] != "github-oauth-token" {
		t.Fatalf("native refresh_token was changed: %v", got["refresh_token"])
	}
	if got["base_url"] != "https://api.individual.githubcopilot.com" {
		t.Fatalf("Copilot base_url was changed: %v", got["base_url"])
	}
	if bytes.Equal(adapted, raw) {
		t.Fatal("expected native Copilot quota storage to be adapted")
	}
}

func TestAdaptQuotaStorageJSON_PreservesExistingPluginCredential(t *testing.T) {
	raw := []byte(`{"type":"copilot","github_access_token":"plugin-github-token","github_base_url":"https://github.com","refresh_token":"other-token"}`)

	got := adaptQuotaStorageJSON("copilot", raw)
	if !bytes.Equal(got, raw) {
		t.Fatalf("existing plugin-owned Copilot credential should remain unchanged: %s", got)
	}
}

func TestAdaptQuotaStorageJSON_DoesNotUseInferenceTokenAsGitHubToken(t *testing.T) {
	raw := []byte(`{"type":"copilot","access_token":"copilot-inference-token","base_url":"https://api.individual.githubcopilot.com"}`)

	got := adaptQuotaStorageJSON("copilot", raw)
	if !bytes.Equal(got, raw) {
		t.Fatalf("credential without a GitHub OAuth token should remain unchanged: %s", got)
	}
}

func TestAdaptQuotaStorageJSON_OnlyAppliesToCopilot(t *testing.T) {
	raw := []byte(`{"type":"codex","access_token":"access","refresh_token":"refresh"}`)

	got := adaptQuotaStorageJSON("codex", raw)
	if !bytes.Equal(got, raw) {
		t.Fatalf("non-Copilot credential should remain unchanged: %s", got)
	}
}
