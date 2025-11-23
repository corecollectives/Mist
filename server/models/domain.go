package models

import (
	"time"

	"github.com/corecollectives/mist/utils"
)

type Domain struct {
	ID        int64     `json:"id"`
	AppID     int64     `json:"appId"`
	Domain    string    `json:"domain"`
	SslStatus string    `json:"sslStatus"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func CreateDomain(appID int64, domain string) (*Domain, error) {
	id := utils.GenerateRandomId()
	query := `
		INSERT INTO domains (id, app_id, domain_name, ssl_status, created_at, updated_at)
		VALUES (?, ?, ?, 'pending', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, app_id, domain_name, ssl_status, created_at, updated_at
	`
	var d Domain
	err := db.QueryRow(query, id, appID, domain).Scan(
		&d.ID, &d.AppID, &d.Domain, &d.SslStatus, &d.CreatedAt, &d.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func GetDomainsByAppID(appID int64) ([]Domain, error) {
	query := `
		SELECT id, app_id, domain_name, ssl_status, created_at, updated_at
		FROM domains
		WHERE app_id = ?
		ORDER BY created_at ASC
	`
	rows, err := db.Query(query, appID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var domains []Domain
	for rows.Next() {
		var d Domain
		err := rows.Scan(&d.ID, &d.AppID, &d.Domain, &d.SslStatus, &d.CreatedAt, &d.UpdatedAt)
		if err != nil {
			return nil, err
		}
		domains = append(domains, d)
	}
	return domains, nil
}

func GetPrimaryDomainByAppID(appID int64) (*Domain, error) {
	query := `
		SELECT id, app_id, domain_name, ssl_status, created_at, updated_at
		FROM domains
		WHERE app_id = ?
		ORDER BY created_at ASC
		LIMIT 1
	`
	var d Domain
	err := db.QueryRow(query, appID).Scan(&d.ID, &d.AppID, &d.Domain, &d.SslStatus, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func UpdateDomain(id int64, domain string) error {
	query := `
		UPDATE domains
		SET domain_name = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	_, err := db.Exec(query, domain, id)
	return err
}

func DeleteDomain(id int64) error {
	query := `DELETE FROM domains WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}

func GetDomainByID(id int64) (*Domain, error) {
	query := `
		SELECT id, app_id, domain_name, ssl_status, created_at, updated_at
		FROM domains
		WHERE id = ?
	`
	var d Domain
	err := db.QueryRow(query, id).Scan(&d.ID, &d.AppID, &d.Domain, &d.SslStatus, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &d, nil
}
