package models

import (
	"database/sql"
)

type SystemSettings struct {
	WildcardDomain *string `json:"wildcardDomain"`
	MistAppName    string  `json:"mistAppName"`
}

func GetSystemSettings() (*SystemSettings, error) {
	var settings SystemSettings

	var wildcardDomain sql.NullString
	err := db.QueryRow(`SELECT value FROM system_settings WHERE key = ?`, "wildcard_domain").Scan(&wildcardDomain)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if wildcardDomain.Valid && wildcardDomain.String != "" {
		settings.WildcardDomain = &wildcardDomain.String
	}

	var mistAppName string
	err = db.QueryRow(`SELECT value FROM system_settings WHERE key = ?`, "mist_app_name").Scan(&mistAppName)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if mistAppName == "" {
		mistAppName = "mist"
	}
	settings.MistAppName = mistAppName

	return &settings, nil
}

func UpdateSystemSettings(wildcardDomain *string, mistAppName string) (*SystemSettings, error) {
	wildcardValue := ""
	if wildcardDomain != nil {
		wildcardValue = *wildcardDomain
	}
	_, err := db.Exec(`
		INSERT INTO system_settings (key, value, updated_at) 
		VALUES (?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(key) DO UPDATE SET value = ?, updated_at = CURRENT_TIMESTAMP
	`, "wildcard_domain", wildcardValue, wildcardValue)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		INSERT INTO system_settings (key, value, updated_at) 
		VALUES (?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(key) DO UPDATE SET value = ?, updated_at = CURRENT_TIMESTAMP
	`, "mist_app_name", mistAppName, mistAppName)
	if err != nil {
		return nil, err
	}

	return GetSystemSettings()
}

func GenerateAutoDomain(projectName, appName string) (string, error) {
	settings, err := GetSystemSettings()
	if err != nil {
		return "", err
	}

	if settings.WildcardDomain == nil || *settings.WildcardDomain == "" {
		return "", nil
	}

	wildcardDomain := *settings.WildcardDomain
	if len(wildcardDomain) > 0 && wildcardDomain[0] == '*' {
		wildcardDomain = wildcardDomain[1:]
	}
	if len(wildcardDomain) > 0 && wildcardDomain[0] == '.' {
		wildcardDomain = wildcardDomain[1:]
	}

	return projectName + "-" + appName + "." + wildcardDomain, nil
}
