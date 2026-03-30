package junie

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/gentleman-programming/gentle-ai/internal/model"
	"github.com/gentleman-programming/gentle-ai/internal/system"
)

func TestDetect(t *testing.T) {
	tests := []struct {
		name            string
		stat            statResult
		wantInstalled   bool
		wantConfigPath  string
		wantConfigFound bool
		wantErr         bool
	}{
		{
			name:            "config directory found",
			stat:            statResult{isDir: true},
			wantInstalled:   true,
			wantConfigPath:  filepath.Join("/tmp/home", ".junie"),
			wantConfigFound: true,
		},
		{
			name:            "config missing",
			stat:            statResult{err: os.ErrNotExist},
			wantInstalled:   false,
			wantConfigPath:  filepath.Join("/tmp/home", ".junie"),
			wantConfigFound: false,
		},
		{
			name:    "stat error bubbles up",
			stat:    statResult{err: errors.New("permission denied")},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Adapter{
				statPath: func(string) statResult {
					return tt.stat
				},
			}

			installed, _, configPath, configFound, err := a.Detect(context.Background(), "/tmp/home")
			if (err != nil) != tt.wantErr {
				t.Fatalf("Detect() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			if installed != tt.wantInstalled {
				t.Fatalf("Detect() installed = %v, want %v", installed, tt.wantInstalled)
			}

			if configPath != tt.wantConfigPath {
				t.Fatalf("Detect() configPath = %q, want %q", configPath, tt.wantConfigPath)
			}

			if configFound != tt.wantConfigFound {
				t.Fatalf("Detect() configFound = %v, want %v", configFound, tt.wantConfigFound)
			}
		})
	}
}

func TestConfigPathsCrossPlatform(t *testing.T) {
	a := NewAdapter()
	home := "/tmp/home"

	if got := a.GlobalConfigDir(home); got != filepath.Join(home, ".junie") {
		t.Fatalf("GlobalConfigDir() = %q, want %q", got, filepath.Join(home, ".junie"))
	}

	if got := a.SkillsDir(home); got != filepath.Join(home, ".junie", "skills") {
		t.Fatalf("SkillsDir() = %q, want %q", got, filepath.Join(home, ".junie", "skills"))
	}

	if got := a.MCPConfigPath(home, "ctx7"); got != filepath.Join(home, ".junie", "mcp", "mcp.json") {
		t.Fatalf("MCPConfigPath() = %q, want %q", got, filepath.Join(home, ".junie", "mcp", "mcp.json"))
	}

	if got := a.SystemPromptFile(home); got != filepath.Join(home, ".junie", "AGENTS.md") {
		t.Fatalf("SystemPromptFile() = %q, want %q", got, filepath.Join(home, ".junie", "AGENTS.md"))
	}

	if got := a.SubAgentsDir(home); got != filepath.Join(home, ".junie", "agents") {
		t.Fatalf("SubAgentsDir() = %q, want %q", got, filepath.Join(home, ".junie", "agents"))
	}
}

func TestStrategies(t *testing.T) {
	a := NewAdapter()

	if got := a.SystemPromptStrategy(); got != model.StrategyFileReplace {
		t.Fatalf("SystemPromptStrategy() = %v, want %v", got, model.StrategyFileReplace)
	}

	if got := a.MCPStrategy(); got != model.StrategyMCPConfigFile {
		t.Fatalf("MCPStrategy() = %v, want %v", got, model.StrategyMCPConfigFile)
	}
}

func TestCapabilities(t *testing.T) {
	a := NewAdapter()

	if !a.SupportsSkills() {
		t.Fatal("Junie should support skills")
	}

	if !a.SupportsSystemPrompt() {
		t.Fatal("Junie should support system prompt")
	}

	if !a.SupportsMCP() {
		t.Fatal("Junie should support MCP")
	}

	if !a.SupportsSubAgents() {
		t.Fatal("Junie should support sub-agents")
	}

	if a.SupportsOutputStyles() {
		t.Fatal("Junie should not support output styles")
	}

	if a.SupportsSlashCommands() {
		t.Fatal("Junie should not support slash commands")
	}
}

func TestAutoInstall(t *testing.T) {
	a := NewAdapter()

	if !a.SupportsAutoInstall() {
		t.Fatal("Junie should support auto-install (CLI available)")
	}

	// Unix install
	cmds, err := a.InstallCommand(system.PlatformProfile{OS: "darwin"})
	if err != nil {
		t.Fatalf("InstallCommand(darwin) error: %v", err)
	}
	if len(cmds) != 1 {
		t.Fatalf("InstallCommand(darwin) = %d commands, want 1", len(cmds))
	}

	// Windows install
	cmds, err = a.InstallCommand(system.PlatformProfile{OS: "windows"})
	if err != nil {
		t.Fatalf("InstallCommand(windows) error: %v", err)
	}
	if len(cmds) != 1 {
		t.Fatalf("InstallCommand(windows) = %d commands, want 1", len(cmds))
	}
}

func TestIdentity(t *testing.T) {
	a := NewAdapter()

	if a.Agent() != model.AgentJunie {
		t.Fatalf("Agent() = %q, want %q", a.Agent(), model.AgentJunie)
	}

	if a.Tier() != model.TierFull {
		t.Fatalf("Tier() = %q, want %q", a.Tier(), model.TierFull)
	}
}

func TestEmbeddedSubAgentsDir(t *testing.T) {
	a := NewAdapter()

	if got := a.EmbeddedSubAgentsDir(); got != "junie/agents" {
		t.Fatalf("EmbeddedSubAgentsDir() = %q, want %q", got, "junie/agents")
	}
}
