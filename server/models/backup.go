package models

import (
	"time"

	"github.com/corecollectives/mist/utils"
)

type BackupType string
type BackupStatus string
type StorageType string

const (
	BackupTypeManual        BackupType = "manual"
	BackupTypeScheduled     BackupType = "scheduled"
	BackupTypePreDeployment BackupType = "pre_deployment"
	BackupTypeAutomatic     BackupType = "automatic"

	BackupStatusPending    BackupStatus = "pending"
	BackupStatusInProgress BackupStatus = "in_progress"
	BackupStatusCompleted  BackupStatus = "completed"
	BackupStatusFailed     BackupStatus = "failed"
	BackupStatusDeleted    BackupStatus = "deleted"

	StorageTypeLocal StorageType = "local"
	StorageTypeS3    StorageType = "s3"
	StorageTypeGCS   StorageType = "gcs"
	StorageTypeAzure StorageType = "azure"
	StorageTypeFTP   StorageType = "ftp"
)

type Backup struct {
	ID                int64        `db:"id" json:"id"`
	AppID             int64        `db:"app_id" json:"appId"`
	BackupType        BackupType   `db:"backup_type" json:"backupType"`
	BackupName        string       `db:"backup_name" json:"backupName"`
	FilePath          string       `db:"file_path" json:"filePath"`
	FileSize          *int64       `db:"file_size" json:"fileSize,omitempty"`
	CompressionType   string       `db:"compression_type" json:"compressionType"`
	DatabaseType      *string      `db:"database_type" json:"databaseType,omitempty"`
	DatabaseVersion   *string      `db:"database_version" json:"databaseVersion,omitempty"`
	StorageType       StorageType  `db:"storage_type" json:"storageType"`
	StoragePath       *string      `db:"storage_path" json:"storagePath,omitempty"`
	Status            BackupStatus `db:"status" json:"status"`
	Progress          int          `db:"progress" json:"progress"`
	ErrorMessage      *string      `db:"error_message" json:"errorMessage,omitempty"`
	Checksum          *string      `db:"checksum" json:"checksum,omitempty"`
	ChecksumAlgorithm string       `db:"checksum_algorithm" json:"checksumAlgorithm"`
	IsVerified        bool         `db:"is_verified" json:"isVerified"`
	VerifiedAt        *time.Time   `db:"verified_at" json:"verifiedAt,omitempty"`
	CanRestore        bool         `db:"can_restore" json:"canRestore"`
	LastRestoreAt     *time.Time   `db:"last_restore_at" json:"lastRestoreAt,omitempty"`
	RestoreCount      int          `db:"restore_count" json:"restoreCount"`
	RetentionDays     *int         `db:"retention_days" json:"retentionDays,omitempty"`
	AutoDeleteAt      *time.Time   `db:"auto_delete_at" json:"autoDeleteAt,omitempty"`
	CreatedBy         *int64       `db:"created_by" json:"createdBy,omitempty"`
	CreatedAt         time.Time    `db:"created_at" json:"createdAt"`
	CompletedAt       *time.Time   `db:"completed_at" json:"completedAt,omitempty"`
	Duration          *int         `db:"duration" json:"duration,omitempty"`
	Notes             *string      `db:"notes" json:"notes,omitempty"`
}

func (b *Backup) ToJson() map[string]interface{} {
	return map[string]interface{}{
		"id":                b.ID,
		"appId":             b.AppID,
		"backupType":        b.BackupType,
		"backupName":        b.BackupName,
		"filePath":          b.FilePath,
		"fileSize":          b.FileSize,
		"compressionType":   b.CompressionType,
		"databaseType":      b.DatabaseType,
		"databaseVersion":   b.DatabaseVersion,
		"storageType":       b.StorageType,
		"storagePath":       b.StoragePath,
		"status":            b.Status,
		"progress":          b.Progress,
		"errorMessage":      b.ErrorMessage,
		"checksum":          b.Checksum,
		"checksumAlgorithm": b.ChecksumAlgorithm,
		"isVerified":        b.IsVerified,
		"verifiedAt":        b.VerifiedAt,
		"canRestore":        b.CanRestore,
		"lastRestoreAt":     b.LastRestoreAt,
		"restoreCount":      b.RestoreCount,
		"retentionDays":     b.RetentionDays,
		"autoDeleteAt":      b.AutoDeleteAt,
		"createdBy":         b.CreatedBy,
		"createdAt":         b.CreatedAt,
		"completedAt":       b.CompletedAt,
		"duration":          b.Duration,
		"notes":             b.Notes,
	}
}

func (b *Backup) InsertInDB() error {
	id := utils.GenerateRandomId()
	b.ID = id
	query := `
	INSERT INTO backups (
		id, app_id, backup_type, backup_name, file_path,
		file_size, compression_type, database_type, database_version,
		storage_type, storage_path, status, checksum_algorithm,
		retention_days, auto_delete_at, created_by, notes
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	RETURNING created_at
	`
	err := db.QueryRow(query, b.ID, b.AppID, b.BackupType, b.BackupName, b.FilePath,
		b.FileSize, b.CompressionType, b.DatabaseType, b.DatabaseVersion,
		b.StorageType, b.StoragePath, b.Status, b.ChecksumAlgorithm,
		b.RetentionDays, b.AutoDeleteAt, b.CreatedBy, b.Notes).Scan(&b.CreatedAt)
	return err
}

func GetBackupsByAppID(appID int64) ([]Backup, error) {
	var backups []Backup
	query := `
	SELECT id, app_id, backup_type, backup_name, file_path, file_size,
	       compression_type, database_type, database_version,
	       storage_type, storage_path, status, progress, error_message,
	       checksum, checksum_algorithm, is_verified, verified_at,
	       can_restore, last_restore_at, restore_count,
	       retention_days, auto_delete_at, created_by, created_at,
	       completed_at, duration, notes
	FROM backups
	WHERE app_id = ? AND status != 'deleted'
	ORDER BY created_at DESC
	`
	rows, err := db.Query(query, appID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var backup Backup
		err := rows.Scan(
			&backup.ID, &backup.AppID, &backup.BackupType, &backup.BackupName,
			&backup.FilePath, &backup.FileSize, &backup.CompressionType,
			&backup.DatabaseType, &backup.DatabaseVersion, &backup.StorageType,
			&backup.StoragePath, &backup.Status, &backup.Progress, &backup.ErrorMessage,
			&backup.Checksum, &backup.ChecksumAlgorithm, &backup.IsVerified, &backup.VerifiedAt,
			&backup.CanRestore, &backup.LastRestoreAt, &backup.RestoreCount,
			&backup.RetentionDays, &backup.AutoDeleteAt, &backup.CreatedBy, &backup.CreatedAt,
			&backup.CompletedAt, &backup.Duration, &backup.Notes,
		)
		if err != nil {
			return nil, err
		}
		backups = append(backups, backup)
	}

	return backups, rows.Err()
}

func GetBackupByID(backupID int64) (*Backup, error) {
	var backup Backup
	query := `
	SELECT id, app_id, backup_type, backup_name, file_path, file_size,
	       compression_type, database_type, database_version,
	       storage_type, storage_path, status, progress, error_message,
	       checksum, checksum_algorithm, is_verified, verified_at,
	       can_restore, last_restore_at, restore_count,
	       retention_days, auto_delete_at, created_by, created_at,
	       completed_at, duration, notes
	FROM backups
	WHERE id = ?
	`
	err := db.QueryRow(query, backupID).Scan(
		&backup.ID, &backup.AppID, &backup.BackupType, &backup.BackupName,
		&backup.FilePath, &backup.FileSize, &backup.CompressionType,
		&backup.DatabaseType, &backup.DatabaseVersion, &backup.StorageType,
		&backup.StoragePath, &backup.Status, &backup.Progress, &backup.ErrorMessage,
		&backup.Checksum, &backup.ChecksumAlgorithm, &backup.IsVerified, &backup.VerifiedAt,
		&backup.CanRestore, &backup.LastRestoreAt, &backup.RestoreCount,
		&backup.RetentionDays, &backup.AutoDeleteAt, &backup.CreatedBy, &backup.CreatedAt,
		&backup.CompletedAt, &backup.Duration, &backup.Notes,
	)
	if err != nil {
		return nil, err
	}
	return &backup, nil
}

func (b *Backup) UpdateStatus(status BackupStatus, errorMsg *string) error {
	query := `
	UPDATE backups
	SET status = ?, error_message = ?, completed_at = CURRENT_TIMESTAMP
	WHERE id = ?
	`
	_, err := db.Exec(query, status, errorMsg, b.ID)
	return err
}

func (b *Backup) UpdateProgress(progress int) error {
	query := `UPDATE backups SET progress = ? WHERE id = ?`
	_, err := db.Exec(query, progress, b.ID)
	return err
}

func (b *Backup) MarkAsRestored() error {
	query := `
	UPDATE backups
	SET last_restore_at = CURRENT_TIMESTAMP,
	    restore_count = restore_count + 1
	WHERE id = ?
	`
	_, err := db.Exec(query, b.ID)
	return err
}

func DeleteExpiredBackups() error {
	query := `
	UPDATE backups
	SET status = 'deleted'
	WHERE auto_delete_at IS NOT NULL AND auto_delete_at < CURRENT_TIMESTAMP AND status != 'deleted'
	`
	_, err := db.Exec(query)
	return err
}
