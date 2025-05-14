package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"

	"os"
	"path/filepath"
	"sort"
	"strings"

	"train-backend/internal/infrastructure/database"
)

func main() {
	// Загружаем .env
	if err := godotenv.Load(); err != nil {
		fmt.Println("⚠️ Не удалось загрузить .env:", err)
	}

	ctx := context.Background()
	dbpool, err := database.NewPool(ctx)
	if err != nil {
		fmt.Printf("❌ не удалось подключиться к БД: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	migrationsDir := "migrations"
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		fmt.Printf("❌ не удалось прочитать директорию миграций: %v\n", err)
		os.Exit(1)
	}

	// сортировка по имени файла
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}

		path := filepath.Join(migrationsDir, file.Name())
		content, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("🔴 Не удалось прочитать %s: %v\n", file.Name(), err)
			continue
		}

		fmt.Printf("🟡 Выполняется миграция: %s\n", file.Name())
		_, err = dbpool.Exec(ctx, string(content))
		if err != nil {
			fmt.Printf("🔴 Ошибка при выполнении %s: %v\n", file.Name(), err)
			continue
		}

		fmt.Printf("✅ Миграция выполнена: %s\n", file.Name())
	}
}
