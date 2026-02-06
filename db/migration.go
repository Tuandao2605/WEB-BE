package db

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jmoiron/sqlx"
)

func RunMigrations(db *sqlx.DB, migrationsDir string) error {
	absDir, err := filepath.Abs(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	fmt.Printf("Running migrations from: %s\n", absDir)

	var files []string

	err = filepath.WalkDir(absDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// 🔥 CHỈ chạy file .up.sql
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".up.sql") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to read migrations dir: %w", err)
	}

	sort.Strings(files)

	for _, file := range files {
		fmt.Printf("Running migration: %s\n", filepath.Base(file))

		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", file, err)
		}

		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("migration failed (%s): %w", file, err)
		}
	}

	return nil
}
