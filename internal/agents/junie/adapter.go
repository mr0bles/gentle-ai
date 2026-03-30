package junie

import (
	"context"
	"os"
	"path/filepath"

	"github.com/gentleman-programming/gentle-ai/internal/model"
	"github.com/gentleman-programming/gentle-ai/internal/system"
)

type statResult struct {
	isDir bool
	err   error
}

type Adapter struct {
	statPath func(string) statResult
}

func NewAdapter() *Adapter {
	return &Adapter{
		statPath: defaultStat,
	}
}

// --- Identity ---

func (a *Adapter) Agent() model.AgentID {
	return model.AgentJunie
}

func (a *Adapter) Tier() model.SupportTier {
	return model.TierFull
}

// --- Detection ---

func (a *Adapter) Detect(_ context.Context, homeDir string) (bool, string, string, bool, error) {
	configPath := filepath.Join(homeDir, ".junie")

	stat := a.statPath(configPath)
	if stat.err != nil {
		if os.IsNotExist(stat.err) {
			return false, "", configPath, false, nil
		}
		return false, "", "", false, stat.err
	}

	// Junie can be used as IDE plugin or CLI — if config dir exists, it's installed.
	return stat.isDir, "", configPath, stat.isDir, nil
}

// --- Installation ---

func (a *Adapter) SupportsAutoInstall() bool {
	return true
}

func (a *Adapter) InstallCommand(profile system.PlatformProfile) ([][]string, error) {
	if profile.OS == "windows" {
		return [][]string{
			{"powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", "iex (irm 'https://junie.jetbrains.com/install.ps1')"},
		}, nil
	}
	return [][]string{
		{"bash", "-c", "curl -fsSL https://junie.jetbrains.com/install.sh | bash"},
	}, nil
}

// --- Config paths ---

func (a *Adapter) GlobalConfigDir(homeDir string) string {
	return filepath.Join(homeDir, ".junie")
}

func (a *Adapter) SystemPromptDir(homeDir string) string {
	return filepath.Join(homeDir, ".junie")
}

func (a *Adapter) SystemPromptFile(homeDir string) string {
	return filepath.Join(homeDir, ".junie", "AGENTS.md")
}

func (a *Adapter) SkillsDir(homeDir string) string {
	return filepath.Join(homeDir, ".junie", "skills")
}

func (a *Adapter) SettingsPath(homeDir string) string {
	return filepath.Join(homeDir, ".junie", "mcp", "mcp.json")
}

// --- Config strategies ---

func (a *Adapter) SystemPromptStrategy() model.SystemPromptStrategy {
	return model.StrategyFileReplace
}

func (a *Adapter) MCPStrategy() model.MCPStrategy {
	return model.StrategyMCPConfigFile
}

// --- MCP ---

func (a *Adapter) MCPConfigPath(homeDir string, _ string) string {
	return filepath.Join(homeDir, ".junie", "mcp", "mcp.json")
}

// --- Optional capabilities ---

func (a *Adapter) SupportsOutputStyles() bool {
	return false
}

func (a *Adapter) OutputStyleDir(_ string) string {
	return ""
}

func (a *Adapter) SupportsSlashCommands() bool {
	return false
}

func (a *Adapter) CommandsDir(_ string) string {
	return ""
}

func (a *Adapter) SupportsSkills() bool {
	return true
}

func (a *Adapter) SupportsSystemPrompt() bool {
	return true
}

func (a *Adapter) SupportsMCP() bool {
	return true
}

// --- Sub-agent support (Junie custom subagents in ~/.junie/agents/) ---

func (a *Adapter) SupportsSubAgents() bool {
	return true
}

func (a *Adapter) SubAgentsDir(homeDir string) string {
	return filepath.Join(homeDir, ".junie", "agents")
}

func (a *Adapter) EmbeddedSubAgentsDir() string {
	return "junie/agents"
}

func defaultStat(path string) statResult {
	info, err := os.Stat(path)
	if err != nil {
		return statResult{err: err}
	}

	return statResult{isDir: info.IsDir()}
}
