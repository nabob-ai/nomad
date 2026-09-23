package memory

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewFilePreferenceStorage(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "test-preferences")
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Errorf("Failed to remove temp dir: %v", err)
		}
	}()

	storage, err := NewFilePreferenceStorage(tmpDir)
	if err != nil {
		t.Fatalf("NewFilePreferenceStorage failed: %v", err)
	}

	if storage == nil {
		t.Fatal("storage should not be nil")
	}

	// 验证目录是否创建
	if _, err := os.Stat(tmpDir); os.IsNotExist(err) {
		t.Error("storage directory should be created")
	}
}

func TestFilePreferenceStorage_Save(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "test-preferences-save")
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Errorf("Failed to remove temp dir: %v", err)
		}
	}()

	storage, _ := NewFilePreferenceStorage(tmpDir)

	preferences := []*Preference{
		{
			ID:         "pref-1",
			UserID:     "user-1",
			Category:   CategoryUI,
			Key:        "theme",
			Value:      "dark",
			Strength:   0.8,
			Confidence: 0.9,
		},
	}

	err := storage.Save(context.Background(), "user-1", preferences)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// 验证文件是否创建
	filePath := filepath.Join(tmpDir, "user-1.json")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Error("preference file should be created")
	}
}

func TestFilePreferenceStorage_Load(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "test-preferences-load")
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Errorf("Failed to remove temp dir: %v", err)
		}
	}()

	storage, _ := NewFilePreferenceStorage(tmpDir)

	// 保存偏好
	original := []*Preference{
		{
			ID:         "pref-1",
			UserID:     "user-1",
			Category:   CategoryUI,
			Key:        "theme",
			Value:      "dark",
			Strength:   0.8,
			Confidence: 0.9,
		},
	}
	if err := storage.Save(context.Background(), "user-1", original); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// 加载偏好
	loaded, err := storage.Load(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if len(loaded) != 1 {
		t.Errorf("loaded %d preferences, want 1", len(loaded))
	}

	if loaded[0].ID != "pref-1" {
		t.Errorf("ID = %s, want 'pref-1'", loaded[0].ID)
	}

	if loaded[0].Value != "dark" {
		t.Errorf("Value = %s, want 'dark'", loaded[0].Value)
	}
}

func TestFilePreferenceStorage_Load_NotExist(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "test-preferences-not-exist")
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Errorf("Failed to remove temp dir: %v", err)
		}
	}()

	storage, _ := NewFilePreferenceStorage(tmpDir)

	// 加载不存在的用户
	loaded, err := storage.Load(context.Background(), "non-existent-user")
	if err != nil {
		t.Fatalf("Load should not error for non-existent user: %v", err)
	}

	if len(loaded) != 0 {
		t.Errorf("loaded %d preferences, want 0", len(loaded))
	}
}

func TestFilePreferenceStorage_Delete(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "test-preferences-delete")
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Errorf("Failed to remove temp dir: %v", err)
		}
	}()

	storage, _ := NewFilePreferenceStorage(tmpDir)

	// 保存偏好
	preferences := []*Preference{
		{ID: "pref-1", UserID: "user-1", Category: CategoryUI, Key: "theme", Value: "dark"},
	}
	if err := storage.Save(context.Background(), "user-1", preferences); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// 删除
	err := storage.Delete(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// 验证文件已删除
	filePath := filepath.Join(tmpDir, "user-1.json")
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		t.Error("file should be deleted")
	}
}

func TestFilePreferenceStorage_List(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "test-preferences-list")
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Errorf("Failed to remove temp dir: %v", err)
		}
	}()

	storage, _ := NewFilePreferenceStorage(tmpDir)

	// 保存多个用户的偏好
	users := []string{"user-1", "user-2", "user-3"}
	for _, userID := range users {
		prefs := []*Preference{
			{ID: "pref-1", UserID: userID, Category: CategoryUI, Key: "theme", Value: "dark"},
		}
		if err := storage.Save(context.Background(), userID, prefs); err != nil {
			t.Fatalf("Save failed: %v", err)
		}
	}

	// 列出所有用户
	userIDs, err := storage.List(context.Background())
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(userIDs) != 3 {
		t.Errorf("List returned %d users, want 3", len(userIDs))
	}
}

func TestNewPersistentPreferenceManager(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "test-persistent-manager")
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Errorf("Failed to remove temp dir: %v", err)
		}
	}()

	storage, _ := NewFilePreferenceStorage(tmpDir)
	config := DefaultPreferenceManagerConfig()

	manager := NewPersistentPreferenceManager(config, storage)

	if manager == nil {
		t.Fatal("NewPersistentPreferenceManager returned nil")
	}

	if manager.storage == nil {
		t.Error("storage should not be nil")
	}
}

func TestPersistentPreferenceManager_SaveLoadUser(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "test-persistent-save-load")
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Errorf("Failed to remove temp dir: %v", err)
		}
	}()

	storage, _ := NewFilePreferenceStorage(tmpDir)
	manager := NewPersistentPreferenceManager(DefaultPreferenceManagerConfig(), storage)

	// 添加偏好
	pref := &Preference{
		UserID:     "user-1",
		Category:   CategoryUI,
		Key:        "theme",
		Value:      "dark",
		Strength:   0.8,
		Confidence: 0.9,
	}
	if err := manager.AddPreference(context.Background(), pref); err != nil {
		t.Fatalf("AddPreference failed: %v", err)
	}

	// 保存
	err := manager.SaveUser(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("SaveUser failed: %v", err)
	}

	// 创建新的管理器
	newManager := NewPersistentPreferenceManager(DefaultPreferenceManagerConfig(), storage)

	// 加载
	err = newManager.LoadUser(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("LoadUser failed: %v", err)
	}

	// 验证加载的偏好
	loaded, err := newManager.GetPreference(context.Background(), "user-1", CategoryUI, "theme")
	if err != nil {
		t.Fatalf("GetPreference failed: %v", err)
	}

	if loaded.Value != "dark" {
		t.Errorf("Value = %s, want 'dark'", loaded.Value)
	}
}

func TestPersistentPreferenceManager_DeleteUser(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "test-persistent-delete")
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Errorf("Failed to remove temp dir: %v", err)
		}
	}()

	storage, _ := NewFilePreferenceStorage(tmpDir)
	manager := NewPersistentPreferenceManager(DefaultPreferenceManagerConfig(), storage)

	// 添加偏好
	pref := &Preference{
		UserID:     "user-1",
		Category:   CategoryUI,
		Key:        "theme",
		Value:      "dark",
		Strength:   0.8,
		Confidence: 0.9,
	}
	if err := manager.AddPreference(context.Background(), pref); err != nil {
		t.Fatalf("AddPreference failed: %v", err)
	}
	if err := manager.SaveUser(context.Background(), "user-1"); err != nil {
		t.Fatalf("SaveUser failed: %v", err)
	}

	// 删除用户
	err := manager.DeleteUser(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("DeleteUser failed: %v", err)
	}

	// 验证内存中已删除
	_, err = manager.GetPreference(context.Background(), "user-1", CategoryUI, "theme")
	if err == nil {
		t.Error("preference should be deleted from memory")
	}

	// 验证文件已删除
	filePath := filepath.Join(tmpDir, "user-1.json")
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		t.Error("file should be deleted")
	}
}

func TestPersistentPreferenceManager_SaveAll(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "test-persistent-save-all")
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Errorf("Failed to remove temp dir: %v", err)
		}
	}()

	storage, _ := NewFilePreferenceStorage(tmpDir)
	manager := NewPersistentPreferenceManager(DefaultPreferenceManagerConfig(), storage)

	// 添加多个用户的偏好
	users := []string{"user-1", "user-2", "user-3"}
	for _, userID := range users {
		pref := &Preference{
			UserID:     userID,
			Category:   CategoryUI,
			Key:        "theme",
			Value:      "dark",
			Strength:   0.8,
			Confidence: 0.9,
		}
		if err := manager.AddPreference(context.Background(), pref); err != nil {
			t.Fatalf("AddPreference failed: %v", err)
		}
	}

	// 保存所有
	err := manager.SaveAll(context.Background())
	if err != nil {
		t.Fatalf("SaveAll failed: %v", err)
	}

	// 验证所有文件都创建了
	for _, userID := range users {
		filePath := filepath.Join(tmpDir, userID+".json")
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("file for %s should exist", userID)
		}
	}
}

func TestPersistentPreferenceManager_LoadAll(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "test-persistent-load-all")
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Errorf("Failed to remove temp dir: %v", err)
		}
	}()

	storage, _ := NewFilePreferenceStorage(tmpDir)
	manager := NewPersistentPreferenceManager(DefaultPreferenceManagerConfig(), storage)

	// 添加并保存
	users := []string{"user-1", "user-2"}
	for _, userID := range users {
		pref := &Preference{
			UserID:     userID,
			Category:   CategoryUI,
			Key:        "theme",
			Value:      "dark",
			Strength:   0.8,
			Confidence: 0.9,
		}
		if err := manager.AddPreference(context.Background(), pref); err != nil {
			t.Fatalf("AddPreference failed: %v", err)
		}
	}
	if err := manager.SaveAll(context.Background()); err != nil {
		t.Fatalf("SaveAll failed: %v", err)
	}

	// 创建新的管理器并加载所有
	newManager := NewPersistentPreferenceManager(DefaultPreferenceManagerConfig(), storage)
	err := newManager.LoadAll(context.Background())
	if err != nil {
		t.Fatalf("LoadAll failed: %v", err)
	}

	// 验证所有用户的偏好都加载了
	for _, userID := range users {
		_, err := newManager.GetPreference(context.Background(), userID, CategoryUI, "theme")
		if err != nil {
			t.Errorf("preference for %s should be loaded", userID)
		}
	}
}

func TestPersistentPreferenceManager_AutoSave(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "test-persistent-autosave")

	storage, _ := NewFilePreferenceStorage(tmpDir)
	manager := NewPersistentPreferenceManager(DefaultPreferenceManagerConfig(), storage)

	// 添加偏好
	pref := &Preference{
		UserID:     "user-1",
		Category:   CategoryUI,
		Key:        "theme",
		Value:      "dark",
		Strength:   0.8,
		Confidence: 0.9,
	}
	if err := manager.AddPreference(context.Background(), pref); err != nil {
		t.Fatalf("AddPreference failed: %v", err)
	}

	// 启动自动保存（1 秒间隔）
	ctx, cancel := context.WithCancel(context.Background())

	go manager.AutoSave(ctx, 1)

	// 等待自动保存执行
	time.Sleep(1500 * time.Millisecond)

	// 验证文件是否创建
	filePath := filepath.Join(tmpDir, "user-1.json")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Error("auto-save should have created the file")
	}

	// 先取消 context，等待 AutoSave goroutine 完成最后一次保存
	cancel()
	time.Sleep(100 * time.Millisecond)

	// 然后清理临时目录
	if err := os.RemoveAll(tmpDir); err != nil {
		t.Errorf("Failed to remove temp dir: %v", err)
	}
}
