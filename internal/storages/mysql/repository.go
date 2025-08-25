package mysql

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	models "go-platform/internal/models/dogs"
	"go-platform/pkg/db/mysql"
)

type MySQLRepository struct {
	mysql *mysql.MySQLClient
}

func NewMySQLRepository(mysql *mysql.MySQLClient) *MySQLRepository {
	return &MySQLRepository{mysql: mysql}
}

// InsertDog inserts a dog into MySQL
func (r *MySQLRepository) InsertDog(ctx context.Context, dog *models.Dog) (string, error) {
	query := `
		INSERT INTO dogs (breed, image_url, created_at)
		VALUES (?, ?, ?)`

	dog.CreatedAt = time.Now()

	result, err := r.mysql.DB().ExecContext(ctx, query, dog.Breed, dog.ImageURL, dog.CreatedAt)
	if err != nil {
		slog.Error("Failed to insert dog into MySQL", "error", err)
		return "", fmt.Errorf("failed to insert dog into MySQL: %w", err)
	}

	// Get the inserted ID
	lastID, err := result.LastInsertId()
	if err != nil {
		slog.Error("Failed to get last insert ID from MySQL", "error", err)
		return "", fmt.Errorf("failed to get last insert ID from MySQL: %w", err)
	}

	slog.Info("Dog inserted into MySQL", "id", lastID)
	return fmt.Sprintf("%d", lastID), nil
}
