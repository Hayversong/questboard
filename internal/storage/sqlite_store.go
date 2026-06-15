package storage

import (
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Hayversong/questboard/internal/model"
	_ "modernc.org/sqlite"
)

const (
	EnvSQLiteFile     = "QUESTBOARD_DB_FILE"
	defaultSQLiteFile = "questboard.db"
)

type ProjectRepositorySQLite struct {
	db *sql.DB
}

type CardRepositorySQLite struct {
	db *sql.DB
}

type ActivityRepositorySQLite struct {
	db *sql.DB
}

type projectRecord struct {
	ID          string
	Name        string
	Description string
	CreatedAt   string
	UpdatedAt   string
}

type cardRecord struct {
	ID          string
	ProjectID   string
	Title       string
	Description string
	Status      string
	Rarity      string
	Deadline    string
	Order       int
	CreatedAt   string
	UpdatedAt   string
}

type activityRecord struct {
	ID          string
	ProjectID   string
	CardID      string
	Description string
	CreatedAt   string
}

func SQLiteFilePath() string {
	if path := os.Getenv(EnvSQLiteFile); path != "" {
		return path
	}

	dir := os.Getenv(EnvDataDir)
	if dir == "" {
		dir = defaultDataDir
	}

	return filepath.Join(dir, defaultSQLiteFile)
}

func OpenSQLite() (*sql.DB, error) {
	path := SQLiteFilePath()

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		db.Close()
		return nil, err
	}

	if err := MigrateSQLite(db); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func MigrateSQLite(db *sql.DB) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS projects (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS cards (
			id TEXT PRIMARY KEY,
			project_id TEXT NOT NULL,
			title TEXT NOT NULL,
			description TEXT NOT NULL DEFAULT '',
			status TEXT NOT NULL,
			rarity TEXT NOT NULL DEFAULT 'common',
			deadline TEXT NOT NULL DEFAULT '',
			order_index INTEGER NOT NULL DEFAULT 0,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL,
			FOREIGN KEY(project_id) REFERENCES projects(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS activities (
			id TEXT PRIMARY KEY,
			project_id TEXT NOT NULL,
			card_id TEXT,
			description TEXT NOT NULL,
			created_at TEXT NOT NULL,
			FOREIGN KEY(project_id) REFERENCES projects(id) ON DELETE CASCADE,
			FOREIGN KEY(card_id) REFERENCES cards(id) ON DELETE SET NULL
		)`,
		`CREATE INDEX IF NOT EXISTS idx_cards_project_id ON cards(project_id)`,
		`CREATE INDEX IF NOT EXISTS idx_activities_project_id ON activities(project_id)`,
		`CREATE INDEX IF NOT EXISTS idx_activities_card_id ON activities(card_id)`,
	}

	for _, statement := range statements {
		if _, err := db.Exec(statement); err != nil {
			return err
		}
	}

	return nil
}

func LoadProjectsSQLite() ([]model.Project, error) {
	db, err := OpenSQLite()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	projectRepository := NewProjectRepositorySQLite(db)
	cardRepository := NewCardRepositorySQLite(db)
	activityRepository := NewActivityRepositorySQLite(db)

	records, err := projectRepository.List()
	if err != nil {
		return nil, err
	}

	projects := make([]model.Project, 0, len(records))

	for _, record := range records {
		cards, err := cardRepository.ListByProject(record.ID)
		if err != nil {
			return nil, err
		}

		activities, err := activityRepository.ListByProject(record.ID)
		if err != nil {
			return nil, err
		}

		project := model.Project{
			ID:         record.ID,
			Name:       record.Name,
			Cards:      make([]model.Card, 0, len(cards)),
			Activities: make([]model.Activity, 0, len(activities)),
		}

		for _, card := range cards {
			project.Cards = append(project.Cards, model.Card{
				ID:          card.ID,
				Title:       card.Title,
				Description: card.Description,
				Status:      card.Status,
				Rarity:      card.Rarity,
				Deadline:    card.Deadline,
				Order:       card.Order,
			})
		}

		for _, activity := range activities {
			project.Activities = append(project.Activities, model.Activity{
				Message: activity.Description,
				Time:    activity.CreatedAt,
			})
		}

		projects = append(projects, project)
	}

	return projects, nil
}

func SaveProjectsSQLite(projects []model.Project) error {
	db, err := OpenSQLite()
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if err := replaceSnapshot(tx, projects); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func NewProjectRepositorySQLite(db *sql.DB) ProjectRepositorySQLite {
	return ProjectRepositorySQLite{db: db}
}

func NewCardRepositorySQLite(db *sql.DB) CardRepositorySQLite {
	return CardRepositorySQLite{db: db}
}

func NewActivityRepositorySQLite(db *sql.DB) ActivityRepositorySQLite {
	return ActivityRepositorySQLite{db: db}
}

func (r ProjectRepositorySQLite) Create(project projectRecord) error {
	now := nowSQLite()
	if project.CreatedAt == "" {
		project.CreatedAt = now
	}
	if project.UpdatedAt == "" {
		project.UpdatedAt = now
	}

	_, err := r.db.Exec(
		`INSERT INTO projects (id, name, description, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)`,
		project.ID,
		project.Name,
		project.Description,
		project.CreatedAt,
		project.UpdatedAt,
	)

	return err
}

func (r ProjectRepositorySQLite) Update(project projectRecord) error {
	result, err := r.db.Exec(
		`UPDATE projects
		SET name = ?, description = ?, updated_at = ?
		WHERE id = ?`,
		project.Name,
		project.Description,
		nowSQLite(),
		project.ID,
	)
	if err != nil {
		return err
	}

	return requireRowsAffected(result)
}

func (r ProjectRepositorySQLite) Delete(id string) error {
	result, err := r.db.Exec(`DELETE FROM projects WHERE id = ?`, id)
	if err != nil {
		return err
	}

	return requireRowsAffected(result)
}

func (r ProjectRepositorySQLite) FindByID(id string) (projectRecord, error) {
	var project projectRecord

	err := r.db.QueryRow(
		`SELECT id, name, description, created_at, updated_at
		FROM projects
		WHERE id = ?`,
		id,
	).Scan(
		&project.ID,
		&project.Name,
		&project.Description,
		&project.CreatedAt,
		&project.UpdatedAt,
	)
	if err != nil {
		return projectRecord{}, err
	}

	return project, nil
}

func (r ProjectRepositorySQLite) List() ([]projectRecord, error) {
	rows, err := r.db.Query(
		`SELECT id, name, description, created_at, updated_at
		FROM projects
		ORDER BY created_at, id`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []projectRecord

	for rows.Next() {
		var project projectRecord

		if err := rows.Scan(
			&project.ID,
			&project.Name,
			&project.Description,
			&project.CreatedAt,
			&project.UpdatedAt,
		); err != nil {
			return nil, err
		}

		projects = append(projects, project)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}

func (r CardRepositorySQLite) Create(card cardRecord) error {
	now := nowSQLite()
	if card.CreatedAt == "" {
		card.CreatedAt = now
	}
	if card.UpdatedAt == "" {
		card.UpdatedAt = now
	}

	_, err := r.db.Exec(
		`INSERT INTO cards (
			id, project_id, title, description, status, rarity,
			deadline, order_index, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		card.ID,
		card.ProjectID,
		card.Title,
		card.Description,
		card.Status,
		card.Rarity,
		card.Deadline,
		card.Order,
		card.CreatedAt,
		card.UpdatedAt,
	)

	return err
}

func (r CardRepositorySQLite) Update(card cardRecord) error {
	result, err := r.db.Exec(
		`UPDATE cards
		SET title = ?, description = ?, status = ?, rarity = ?,
			deadline = ?, order_index = ?, updated_at = ?
		WHERE id = ?`,
		card.Title,
		card.Description,
		card.Status,
		card.Rarity,
		card.Deadline,
		card.Order,
		nowSQLite(),
		card.ID,
	)
	if err != nil {
		return err
	}

	return requireRowsAffected(result)
}

func (r CardRepositorySQLite) Delete(id string) error {
	result, err := r.db.Exec(`DELETE FROM cards WHERE id = ?`, id)
	if err != nil {
		return err
	}

	return requireRowsAffected(result)
}

func (r CardRepositorySQLite) FindByID(id string) (cardRecord, error) {
	var card cardRecord

	err := r.db.QueryRow(
		`SELECT
			id, project_id, title, description, status, rarity,
			deadline, order_index, created_at, updated_at
		FROM cards
		WHERE id = ?`,
		id,
	).Scan(
		&card.ID,
		&card.ProjectID,
		&card.Title,
		&card.Description,
		&card.Status,
		&card.Rarity,
		&card.Deadline,
		&card.Order,
		&card.CreatedAt,
		&card.UpdatedAt,
	)
	if err != nil {
		return cardRecord{}, err
	}

	return card, nil
}

func (r CardRepositorySQLite) ListByProject(projectID string) ([]cardRecord, error) {
	rows, err := r.db.Query(
		`SELECT
			id, project_id, title, description, status, rarity,
			deadline, order_index, created_at, updated_at
		FROM cards
		WHERE project_id = ?
		ORDER BY status, order_index, created_at, id`,
		projectID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []cardRecord

	for rows.Next() {
		var card cardRecord

		if err := rows.Scan(
			&card.ID,
			&card.ProjectID,
			&card.Title,
			&card.Description,
			&card.Status,
			&card.Rarity,
			&card.Deadline,
			&card.Order,
			&card.CreatedAt,
			&card.UpdatedAt,
		); err != nil {
			return nil, err
		}

		cards = append(cards, card)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return cards, nil
}

func (r ActivityRepositorySQLite) Create(activity activityRecord) error {
	if activity.ID == "" {
		activity.ID = strconv.FormatInt(time.Now().UnixNano(), 10)
	}
	if activity.CreatedAt == "" {
		activity.CreatedAt = nowSQLite()
	}

	cardID := nullableString(activity.CardID)

	_, err := r.db.Exec(
		`INSERT INTO activities (id, project_id, card_id, description, created_at)
		VALUES (?, ?, ?, ?, ?)`,
		activity.ID,
		activity.ProjectID,
		cardID,
		activity.Description,
		activity.CreatedAt,
	)

	return err
}

func (r ActivityRepositorySQLite) Update(activity activityRecord) error {
	cardID := nullableString(activity.CardID)

	result, err := r.db.Exec(
		`UPDATE activities
		SET project_id = ?, card_id = ?, description = ?
		WHERE id = ?`,
		activity.ProjectID,
		cardID,
		activity.Description,
		activity.ID,
	)
	if err != nil {
		return err
	}

	return requireRowsAffected(result)
}

func (r ActivityRepositorySQLite) Delete(id string) error {
	result, err := r.db.Exec(`DELETE FROM activities WHERE id = ?`, id)
	if err != nil {
		return err
	}

	return requireRowsAffected(result)
}

func (r ActivityRepositorySQLite) FindByID(id string) (activityRecord, error) {
	var activity activityRecord
	var cardID sql.NullString

	err := r.db.QueryRow(
		`SELECT id, project_id, card_id, description, created_at
		FROM activities
		WHERE id = ?`,
		id,
	).Scan(
		&activity.ID,
		&activity.ProjectID,
		&cardID,
		&activity.Description,
		&activity.CreatedAt,
	)
	if err != nil {
		return activityRecord{}, err
	}

	activity.CardID = cardID.String

	return activity, nil
}

func (r ActivityRepositorySQLite) ListByProject(projectID string) ([]activityRecord, error) {
	rows, err := r.db.Query(
		`SELECT id, project_id, card_id, description, created_at
		FROM activities
		WHERE project_id = ?
		ORDER BY created_at DESC, id DESC`,
		projectID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []activityRecord

	for rows.Next() {
		var activity activityRecord
		var cardID sql.NullString

		if err := rows.Scan(
			&activity.ID,
			&activity.ProjectID,
			&cardID,
			&activity.Description,
			&activity.CreatedAt,
		); err != nil {
			return nil, err
		}

		activity.CardID = cardID.String
		activities = append(activities, activity)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return activities, nil
}

func replaceSnapshot(tx *sql.Tx, projects []model.Project) error {
	now := nowSQLite()

	if _, err := tx.Exec(`DELETE FROM activities`); err != nil {
		return err
	}
	if _, err := tx.Exec(`DELETE FROM cards`); err != nil {
		return err
	}
	if _, err := tx.Exec(`DELETE FROM projects`); err != nil {
		return err
	}

	for _, project := range projects {
		if _, err := tx.Exec(
			`INSERT INTO projects (id, name, description, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?)`,
			project.ID,
			project.Name,
			"",
			now,
			now,
		); err != nil {
			return err
		}

		for _, card := range project.Cards {
			if _, err := tx.Exec(
				`INSERT INTO cards (
					id, project_id, title, description, status, rarity,
					deadline, order_index, created_at, updated_at
				) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
				card.ID,
				project.ID,
				card.Title,
				card.Description,
				card.Status,
				card.Rarity,
				card.Deadline,
				card.Order,
				now,
				now,
			); err != nil {
				return err
			}
		}

		for index, activity := range project.Activities {
			createdAt := activity.Time
			if createdAt == "" {
				createdAt = now
			}

			if _, err := tx.Exec(
				`INSERT INTO activities (id, project_id, card_id, description, created_at)
				VALUES (?, ?, ?, ?, ?)`,
				activityID(project.ID, index, activity),
				project.ID,
				nil,
				activity.Message,
				createdAt,
			); err != nil {
				return err
			}
		}
	}

	return nil
}

func activityID(projectID string, index int, activity model.Activity) string {
	return projectID + "-activity-" + strconv.Itoa(index)
}

func nowSQLite() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func nullableString(value string) sql.NullString {
	return sql.NullString{
		String: value,
		Valid:  value != "",
	}
}

func requireRowsAffected(result sql.Result) error {
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func IsNotFound(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
