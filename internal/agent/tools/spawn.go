package tools

import (
	"context"
	"fmt"
)

// SpawnTool creates a background subagent; stubbed for v0.
// Args: {"agent": "name", "task": "..."}

type SpawnTool struct{}

func NewSpawnTool() *SpawnTool { return &SpawnTool{} }

func (t *SpawnTool) Name() string        { return "spawn" }
func (t *SpawnTool) Description() string { return "Spawn a background subagent (stub)" }

func (t *SpawnTool) Parameters() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"agent": map[string]interface{}{
				"type":        "string",
				"description": "The name of the agent to spawn",
			},
			"task": map[string]interface{}{
				"type":        "string",
				"description": "The task description for the spawned agent",
			},
		},
		"required": []string{},
	}
}

func (t *SpawnTool) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
	agentName, _ := args["agent"].(string)
	task, _ := args["task"].(string)
	if agentName == "" && task == "" {
		return "", fmt.Errorf("spawn: 'agent' or 'task' required")
	}
	// For v0 we simply return an acknowledgement
	return fmt.Sprintf("spawned: agent=%s task=%s", agentName, task), nil
}
