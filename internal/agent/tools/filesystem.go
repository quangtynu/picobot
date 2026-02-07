package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FilesystemTool provides read/write/list operations within the filesystem.
// Args:
// - action: "read" | "write" | "list"
// - path: path to file or directory (relative to workspace)
// - content: for write
type FilesystemTool struct {
	workspaceDir string
}

func NewFilesystemTool(workspaceDir string) *FilesystemTool {
	return &FilesystemTool{workspaceDir: workspaceDir}
}

func (t *FilesystemTool) Name() string        { return "filesystem" }
func (t *FilesystemTool) Description() string { return "Read, write, and list files in the workspace" }

func (t *FilesystemTool) Parameters() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"action": map[string]interface{}{
				"type":        "string",
				"description": "The filesystem operation to perform",
				"enum":        []string{"read", "write", "list"},
			},
			"path": map[string]interface{}{
				"type":        "string",
				"description": "The file or directory path (relative to workspace)",
			},
			"content": map[string]interface{}{
				"type":        "string",
				"description": "Content to write (required when action is 'write')",
			},
		},
		"required": []string{"action", "path"},
	}
}

// resolvePath resolves a path relative to the workspace directory and ensures
// it stays within the workspace boundary.
func (t *FilesystemTool) resolvePath(pathStr string) (string, error) {
	if pathStr == "" {
		return t.workspaceDir, nil
	}
	// Join relative paths with workspace dir; absolute paths are used as-is
	var resolved string
	if filepath.IsAbs(pathStr) {
		resolved = filepath.Clean(pathStr)
	} else {
		resolved = filepath.Join(t.workspaceDir, pathStr)
	}
	// Security: ensure the resolved path is within the workspace
	absResolved, err := filepath.Abs(resolved)
	if err != nil {
		return "", err
	}
	absBase, err := filepath.Abs(t.workspaceDir)
	if err != nil {
		return "", err
	}
	if !strings.HasPrefix(absResolved, absBase) {
		return "", fmt.Errorf("filesystem: path outside workspace not allowed")
	}
	return absResolved, nil
}

func (t *FilesystemTool) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
	actionRaw, ok := args["action"]
	if !ok {
		return "", fmt.Errorf("filesystem: 'action' is required")
	}
	action, ok := actionRaw.(string)
	if !ok {
		return "", fmt.Errorf("filesystem: 'action' must be a string")
	}
	pathRaw, _ := args["path"]
	pathStr := ""
	if pathRaw != nil {
		switch v := pathRaw.(type) {
		case string:
			pathStr = v
		default:
			return "", fmt.Errorf("filesystem: 'path' must be a string")
		}
	}
	resolved, err := t.resolvePath(pathStr)
	if err != nil {
		return "", err
	}
	switch action {
	case "read":
		b, err := os.ReadFile(resolved)
		if err != nil {
			return "", err
		}
		return string(b), nil
	case "write":
		contentRaw, _ := args["content"]
		content := ""
		switch v := contentRaw.(type) {
		case string:
			content = v
		default:
			return "", fmt.Errorf("filesystem: 'content' must be a string")
		}
		// Create parent directories if needed
		dir := filepath.Dir(resolved)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return "", err
		}
		if err := os.WriteFile(resolved, []byte(content), 0o644); err != nil {
			return "", err
		}
		return "written", nil
	case "list":
		entries, err := os.ReadDir(resolved)
		if err != nil {
			return "", err
		}
		out := ""
		for _, e := range entries {
			name := e.Name()
			if e.IsDir() {
				name += "/"
			}
			out += name + "\n"
		}
		return out, nil
	default:
		return "", fmt.Errorf("filesystem: unknown action %s", action)
	}
}
