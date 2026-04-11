package postgres

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type searchDocumentQueryer interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

func ComposeProjectEmbeddingText(ctx context.Context, db searchDocumentQueryer, projectID uuid.UUID) (string, error) {
	var contentText string
	err := db.QueryRow(ctx, `SELECT COALESCE(compose_project_embedding_text($1), '')`, projectID).Scan(&contentText)
	if err != nil {
		return "", fmt.Errorf("compose project embedding text: %w", err)
	}

	return contentText, nil
}

func RefreshProjectSearchDocument(ctx context.Context, db searchDocumentQueryer, projectID uuid.UUID) (string, bool, error) {
	contentText, err := ComposeProjectEmbeddingText(ctx, db, projectID)
	if err != nil {
		return "", false, err
	}

	newHash := fmt.Sprintf("%x", sha256.Sum256([]byte(contentText)))

	var storedHash *string
	err = db.QueryRow(ctx, `SELECT search_content_hash FROM project_search_documents WHERE project_id = $1`, projectID).Scan(&storedHash)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return "", false, fmt.Errorf("load stored search content hash: %w", err)
	}

	if storedHash != nil && *storedHash == newHash {
		return contentText, false, nil
	}

	_, err = db.Exec(ctx, `
		INSERT INTO project_search_documents (
			project_id,
			search_document,
			search_trgm,
			search_content_hash,
			search_composed_at
		) VALUES (
			$1,
			compose_project_search_doc($1),
			compose_project_search_trgm($1),
			$2,
			NOW()
		)
		ON CONFLICT (project_id) DO UPDATE
		SET search_document = compose_project_search_doc($1),
		    search_trgm = compose_project_search_trgm($1),
		    search_content_hash = $2,
		    search_composed_at = NOW()`,
		projectID,
		newHash,
	)
	if err != nil {
		return "", false, fmt.Errorf("upsert project search document: %w", err)
	}

	return contentText, true, nil
}

func UpdateProjectSearchEmbedding(ctx context.Context, db searchDocumentQueryer, projectID uuid.UUID, embeddingVec []float32) error {
	if len(embeddingVec) == 0 {
		return fmt.Errorf("empty embedding vector")
	}

	vecStrs := make([]string, len(embeddingVec))
	for i, v := range embeddingVec {
		vecStrs[i] = fmt.Sprintf("%f", v)
	}

	_, err := db.Exec(
		ctx,
		"UPDATE project_search_documents SET search_embedding = $1::vector WHERE project_id = $2",
		"["+strings.Join(vecStrs, ",")+"]",
		projectID,
	)
	if err != nil {
		return fmt.Errorf("update search embedding: %w", err)
	}

	return nil
}
