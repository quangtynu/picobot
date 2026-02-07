package tools

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestSkillManager_CreateSkill(t *testing.T) {
	tmpDir := t.TempDir()
	mgr := NewSkillManager(tmpDir)

	err := mgr.CreateSkill("test-skill", "Test description", "# Test\n\nTest content")
	if err != nil {
		t.Fatalf("CreateSkill failed: %v", err)
	}

	// Verify file was created
	skillPath := filepath.Join(tmpDir, "skills", "test-skill", "SKILL.md")
	content, err := os.ReadFile(skillPath)
	if err != nil {
		t.Fatalf("Failed to read created skill: %v", err)
	}

	strContent := string(content)
	if !containsString(strContent, "name: test-skill") {
		t.Error("Skill content missing name in frontmatter")
	}
	if !containsString(strContent, "Test content") {
		t.Error("Skill content missing body content")
	}
}

func TestSkillManager_ListSkills(t *testing.T) {
	tmpDir := t.TempDir()
	mgr := NewSkillManager(tmpDir)

	// Create test skills
	mgr.CreateSkill("skill1", "Description 1", "Content 1")
	mgr.CreateSkill("skill2", "Description 2", "Content 2")

	skills, err := mgr.ListSkills()
	if err != nil {
		t.Fatalf("ListSkills failed: %v", err)
	}

	if len(skills) != 2 {
		t.Errorf("Expected 2 skills, got %d", len(skills))
	}
}

func TestSkillManager_GetSkill(t *testing.T) {
	tmpDir := t.TempDir()
	mgr := NewSkillManager(tmpDir)

	mgr.CreateSkill("test-skill", "Test description", "# Test\n\nTest content")

	content, err := mgr.GetSkill("test-skill")
	if err != nil {
		t.Fatalf("GetSkill failed: %v", err)
	}

	if !containsString(content, "Test content") {
		t.Error("Skill content not returned correctly")
	}
}

func TestSkillManager_DeleteSkill(t *testing.T) {
	tmpDir := t.TempDir()
	mgr := NewSkillManager(tmpDir)

	mgr.CreateSkill("test-skill", "Test description", "Content")

	err := mgr.DeleteSkill("test-skill")
	if err != nil {
		t.Fatalf("DeleteSkill failed: %v", err)
	}

	// Verify skill is gone
	skillPath := filepath.Join(tmpDir, "skills", "test-skill")
	if _, err := os.Stat(skillPath); !os.IsNotExist(err) {
		t.Error("Skill directory still exists after deletion")
	}
}

func TestCreateSkillTool_Execute(t *testing.T) {
	tmpDir := t.TempDir()
	mgr := NewSkillManager(tmpDir)
	tool := NewCreateSkillTool(mgr)

	args := map[string]interface{}{
		"name":        "test-skill",
		"description": "Test description",
		"content":     "# Test Content",
	}

	result, err := tool.Execute(context.Background(), args)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if !containsString(result, "created successfully") {
		t.Errorf("Unexpected result: %s", result)
	}

	// Verify skill was created
	skillPath := filepath.Join(tmpDir, "skills", "test-skill", "SKILL.md")
	if _, err := os.Stat(skillPath); os.IsNotExist(err) {
		t.Error("Skill file was not created")
	}
}

func TestListSkillsTool_Execute(t *testing.T) {
	tmpDir := t.TempDir()
	mgr := NewSkillManager(tmpDir)
	mgr.CreateSkill("skill1", "Description 1", "Content 1")
	mgr.CreateSkill("skill2", "Description 2", "Content 2")

	tool := NewListSkillsTool(mgr)
	result, err := tool.Execute(context.Background(), map[string]interface{}{})
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if !containsString(result, "skill1") || !containsString(result, "skill2") {
		t.Errorf("Expected skills in result, got: %s", result)
	}
}

func TestReadSkillTool_Execute(t *testing.T) {
	tmpDir := t.TempDir()
	mgr := NewSkillManager(tmpDir)
	mgr.CreateSkill("test-skill", "Test description", "# Test Content")

	tool := NewReadSkillTool(mgr)
	args := map[string]interface{}{"name": "test-skill"}
	result, err := tool.Execute(context.Background(), args)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if !containsString(result, "Test Content") {
		t.Errorf("Expected skill content in result, got: %s", result)
	}
}

func TestDeleteSkillTool_Execute(t *testing.T) {
	tmpDir := t.TempDir()
	mgr := NewSkillManager(tmpDir)
	mgr.CreateSkill("test-skill", "Test description", "Content")

	tool := NewDeleteSkillTool(mgr)
	args := map[string]interface{}{"name": "test-skill"}
	result, err := tool.Execute(context.Background(), args)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if !containsString(result, "deleted successfully") {
		t.Errorf("Unexpected result: %s", result)
	}

	// Verify skill is gone
	skillPath := filepath.Join(tmpDir, "skills", "test-skill")
	if _, err := os.Stat(skillPath); !os.IsNotExist(err) {
		t.Error("Skill directory still exists after deletion")
	}
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			containsAt(s, substr))))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
