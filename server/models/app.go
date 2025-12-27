package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/corecollectives/mist/utils"
)

type DeploymentStrategy string
type AppStatus string
type AppType string
type RestartPolicy string

const (
	DeploymentAuto   DeploymentStrategy = "auto"
	DeploymentManual DeploymentStrategy = "manual"

	StatusStopped   AppStatus = "stopped"
	StatusRunning   AppStatus = "running"
	StatusError     AppStatus = "error"
	StatusBuilding  AppStatus = "building"
	StatusDeploying AppStatus = "deploying"

	AppTypeWeb      AppType = "web"
	AppTypeService  AppType = "service"
	AppTypeDatabase AppType = "database"

	RestartPolicyNo            RestartPolicy = "no"
	RestartPolicyAlways        RestartPolicy = "always"
	RestartPolicyOnFailure     RestartPolicy = "on-failure"
	RestartPolicyUnlessStopped RestartPolicy = "unless-stopped"
)

type App struct {
	ID          int64   `db:"id" json:"id"`
	ProjectID   int64   `db:"project_id" json:"project_id"`
	CreatedBy   int64   `db:"created_by" json:"created_by"`
	Name        string  `db:"name" json:"name"`
	Description *string `db:"description" json:"description,omitempty"`

	AppType      AppType `db:"app_type" json:"app_type"`
	TemplateName *string `db:"template_name" json:"template_name,omitempty"`

	GitProviderID *int64  `db:"git_provider_id" json:"git_provider_id,omitempty"`
	GitRepository *string `db:"git_repository" json:"git_repository,omitempty"`
	GitBranch     string  `db:"git_branch" json:"git_branch,omitempty"`
	GitCloneURL   *string `db:"git_clone_url" json:"git_clone_url,omitempty"`

	DeploymentStrategy DeploymentStrategy `db:"deployment_strategy" json:"deployment_strategy"`
	Port               *int64             `db:"port" json:"port,omitempty"`
	RootDirectory      string             `db:"root_directory" json:"root_directory,omitempty"`
	BuildCommand       *string            `db:"build_command" json:"build_command,omitempty"`
	StartCommand       *string            `db:"start_command" json:"start_command,omitempty"`
	DockerfilePath     *string            `db:"dockerfile_path" json:"dockerfile_path,omitempty"`

	CPULimit      *float64      `db:"cpu_limit" json:"cpu_limit,omitempty"`
	MemoryLimit   *int          `db:"memory_limit" json:"memory_limit,omitempty"`
	RestartPolicy RestartPolicy `db:"restart_policy" json:"restart_policy"`

	HealthcheckPath     *string   `db:"healthcheck_path" json:"healthcheck_path,omitempty"`
	HealthcheckInterval int       `db:"healthcheck_interval" json:"healthcheck_interval"`
	HealthcheckTimeout  int       `db:"healthcheck_timeout" json:"healthcheck_timeout"`
	HealthcheckRetries  int       `db:"healthcheck_retries" json:"healthcheck_retries"`
	Status              AppStatus `db:"status" json:"status"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

func (a *App) ToJson() map[string]interface{} {
	return map[string]interface{}{
		"id":                  a.ID,
		"projectId":           a.ProjectID,
		"createdBy":           a.CreatedBy,
		"name":                a.Name,
		"description":         a.Description,
		"appType":             a.AppType,
		"templateName":        a.TemplateName,
		"gitProviderId":       a.GitProviderID,
		"gitRepository":       a.GitRepository,
		"gitBranch":           a.GitBranch,
		"gitCloneUrl":         a.GitCloneURL,
		"deploymentStrategy":  a.DeploymentStrategy,
		"port":                a.Port,
		"rootDirectory":       a.RootDirectory,
		"buildCommand":        a.BuildCommand,
		"startCommand":        a.StartCommand,
		"dockerfilePath":      a.DockerfilePath,
		"cpuLimit":            a.CPULimit,
		"memoryLimit":         a.MemoryLimit,
		"restartPolicy":       a.RestartPolicy,
		"healthcheckPath":     a.HealthcheckPath,
		"healthcheckInterval": a.HealthcheckInterval,
		"healthcheckTimeout":  a.HealthcheckTimeout,
		"healthcheckRetries":  a.HealthcheckRetries,
		"status":              a.Status,
		"createdAt":           a.CreatedAt,
		"updatedAt":           a.UpdatedAt,
	}
}

func (a *App) InsertInDB() error {
	id := utils.GenerateRandomId()
	a.ID = id

	if a.AppType == "" {
		a.AppType = AppTypeWeb
	}
	if a.RestartPolicy == "" {
		a.RestartPolicy = RestartPolicyUnlessStopped
	}
	if a.DeploymentStrategy == "" {
		a.DeploymentStrategy = DeploymentAuto
	}
	if a.Status == "" {
		a.Status = StatusStopped
	}

	query := `
	INSERT INTO apps (
		id, name, description, project_id, created_by, app_type, template_name,
		port, deployment_strategy, restart_policy, git_provider_id, git_repository,
		git_branch, git_clone_url, root_directory, build_command, start_command,
		dockerfile_path, cpu_limit, memory_limit, healthcheck_path,
		healthcheck_interval, healthcheck_timeout, healthcheck_retries, status
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	RETURNING 
		created_at, updated_at
	`
	err := db.QueryRow(query,
		a.ID, a.Name, a.Description, a.ProjectID, a.CreatedBy, a.AppType, a.TemplateName,
		a.Port, a.DeploymentStrategy, a.RestartPolicy, a.GitProviderID, a.GitRepository,
		a.GitBranch, a.GitCloneURL, a.RootDirectory, a.BuildCommand, a.StartCommand,
		a.DockerfilePath, a.CPULimit, a.MemoryLimit, a.HealthcheckPath,
		a.HealthcheckInterval, a.HealthcheckTimeout, a.HealthcheckRetries, a.Status,
	).Scan(&a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func GetApplicationByProjectID(projectId int64) ([]App, error) {
	var apps []App
	query := `
	SELECT id, project_id, created_by, name, description, app_type, template_name,
	       git_provider_id, git_repository, git_branch, git_clone_url, 
	       deployment_strategy, port, root_directory, build_command, start_command,
	       dockerfile_path, cpu_limit, memory_limit, restart_policy,
	       healthcheck_path, healthcheck_interval, healthcheck_timeout, healthcheck_retries,
	       status, created_at, updated_at
	FROM apps
	WHERE project_id = ?
	`
	rows, err := db.Query(query, projectId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var app App
		err := rows.Scan(
			&app.ID, &app.ProjectID, &app.CreatedBy, &app.Name, &app.Description,
			&app.AppType, &app.TemplateName, &app.GitProviderID, &app.GitRepository,
			&app.GitBranch, &app.GitCloneURL, &app.DeploymentStrategy, &app.Port,
			&app.RootDirectory, &app.BuildCommand, &app.StartCommand, &app.DockerfilePath,
			&app.CPULimit, &app.MemoryLimit, &app.RestartPolicy,
			&app.HealthcheckPath, &app.HealthcheckInterval, &app.HealthcheckTimeout,
			&app.HealthcheckRetries, &app.Status, &app.CreatedAt, &app.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		apps = append(apps, app)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return apps, nil
}

func GetApplicationByID(appId int64) (*App, error) {
	var app App
	query := `
	SELECT id, project_id, created_by, name, description, app_type, template_name,
	       git_provider_id, git_repository, git_branch, git_clone_url,
	       deployment_strategy, port, root_directory, build_command, start_command,
	       dockerfile_path, cpu_limit, memory_limit, restart_policy,
	       healthcheck_path, healthcheck_interval, healthcheck_timeout, healthcheck_retries,
	       status, created_at, updated_at
	FROM apps
	WHERE id = ?
	`
	err := db.QueryRow(query, appId).Scan(
		&app.ID, &app.ProjectID, &app.CreatedBy, &app.Name, &app.Description,
		&app.AppType, &app.TemplateName, &app.GitProviderID, &app.GitRepository,
		&app.GitBranch, &app.GitCloneURL, &app.DeploymentStrategy, &app.Port,
		&app.RootDirectory, &app.BuildCommand, &app.StartCommand, &app.DockerfilePath,
		&app.CPULimit, &app.MemoryLimit, &app.RestartPolicy,
		&app.HealthcheckPath, &app.HealthcheckInterval, &app.HealthcheckTimeout,
		&app.HealthcheckRetries, &app.Status, &app.CreatedAt, &app.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func (app *App) UpdateApplication() error {
	query := `
	UPDATE apps
	SET 
		name = ?,
		description = ?,
		app_type = ?,
		template_name = ?,
		git_provider_id = ?,
		git_repository = ?,
		git_branch = ?,
		git_clone_url = ?,
		deployment_strategy = ?,
		port = ?,
		root_directory = ?,
		build_command = ?,
		start_command = ?,
		dockerfile_path = ?,
		cpu_limit = ?,
		memory_limit = ?,
		restart_policy = ?,
		healthcheck_path = ?,
		healthcheck_interval = ?,
		healthcheck_timeout = ?,
		healthcheck_retries = ?,
		status = ?,
		updated_at = CURRENT_TIMESTAMP
	WHERE id = ?
	`
	_, err := db.Exec(query,
		app.Name, app.Description, app.AppType, app.TemplateName,
		app.GitProviderID, app.GitRepository, app.GitBranch, app.GitCloneURL,
		app.DeploymentStrategy, app.Port, app.RootDirectory,
		app.BuildCommand, app.StartCommand, app.DockerfilePath,
		app.CPULimit, app.MemoryLimit, app.RestartPolicy,
		app.HealthcheckPath, app.HealthcheckInterval, app.HealthcheckTimeout,
		app.HealthcheckRetries, app.Status, app.ID,
	)
	return err
}

func IsUserApplicationOwner(userId int64, appId int64) (bool, error) {
	var createdBy int64
	err := db.QueryRow(`
		SELECT created_by FROM apps WHERE id = ?
	`, appId).Scan(&createdBy)
	if err != nil {
		return false, err
	}
	return createdBy == userId, nil
}

func FindApplicationIDByGitRepoAndBranch(gitRepo string, gitBranch string) (int64, error) {
	var appId int64
	err := db.QueryRow(`
		SELECT id FROM apps WHERE git_repository = ? AND git_branch = ?
	`, gitRepo, gitBranch).Scan(&appId)
	if err != nil {
		return 0, err
	}
	return appId, nil
}

func GetUserIDByAppID(appID int64) (*int64, error) {
	query := `
		SELECT created_by FROM apps WHERE id = ?
	`
	var userID int64
	err := db.QueryRow(query, appID).Scan(&userID)
	if err != nil {
		return nil, err
	}
	return &userID, nil
}

func GetAppIDByDeploymentID(depID int64) (int64, error) {
	query := `
		SELECT app_id FROM deployments WHERE id = ?
	`
	var appID int64
	err := db.QueryRow(query, depID).Scan(&appID)
	if err != nil {
		return 0, err
	}
	return appID, nil
}

func GetAppRepoInfo(appId int64) (string, string, int64, string, error) {
	var repo sql.NullString
	var branch sql.NullString
	var name string
	var projectId int64

	err := db.QueryRow(`
		SELECT git_repository, git_branch, project_id, name
		FROM apps WHERE id = ?
	`, appId).Scan(&repo, &branch, &projectId, &name)
	if err != nil {
		return "", "", 0, "", err
	}

	repoStr := ""
	if repo.Valid {
		repoStr = repo.String
	}

	branchStr := ""
	if branch.Valid {
		branchStr = branch.String
	}

	return repoStr, branchStr, projectId, name, nil
}

func GetAppRepoAndBranch(appID int64) (string, string, error) {
	var repoName sql.NullString
	var branch string
	err := db.QueryRow(`SELECT git_repository, COALESCE(git_branch, 'main') FROM apps WHERE id = ?`, appID).
		Scan(&repoName, &branch)
	if err != nil {
		return "", "", err
	}
	if !repoName.Valid || repoName.String == "" {
		return "", "", fmt.Errorf("app has no git repository configured")
	}
	return repoName.String, branch, nil
}

func DeleteApplication(appID int64) error {
	query := `DELETE FROM apps WHERE id = ?`
	_, err := db.Exec(query, appID)
	return err
}

// if git_clone_url is set, it uses that
// if not, it falls back to git_repository
// assuming github, becuase in v1.0.0, we only had github as a git provider and we were only storing git_repository and appending the https://github.com/ manually, bad design decision but whatever
func GetAppCloneURL(appID int64, userID int64) (string, string, bool, error) {
	gitProviderID, gitRepository, _, gitCloneURL, _, _, err := GetAppGitInfo(appID)
	if err != nil {
		return "", "", false, fmt.Errorf("failed to get app git info: %w", err)
	}

	// case 1 - git_clone_url is set, use it
	if gitCloneURL != nil && *gitCloneURL != "" {
		var accessToken string
		if gitProviderID != nil {
			token, _, needsRefresh, err := GetGitProviderAccessToken(*gitProviderID)
			if err != nil {
				// try to get GitHub access token as fallback
				_, err := GetInstallationIDByUserID(userID)
				if err == nil {
					installation, err := GetInstallationByUserID(int(userID))
					if err == nil {
						accessToken = installation.AccessToken
					}
				}
			} else {
				// Check if token needs refresh
				if needsRefresh {
					// Try to refresh the token
					newToken, err := RefreshGitProviderToken(*gitProviderID)
					if err == nil {
						token = newToken
					}
					// If refresh fails, continue with expired token - deployment will fail but at least we tried
				}
				accessToken = token
			}
		}
		return *gitCloneURL, accessToken, false, nil
	}

	// Case 2 - git_clone_url missing but git_repository present (legacy github apps)
	if gitRepository != nil && *gitRepository != "" {
		// assume GitHub and construct the clone URL
		cloneURL := fmt.Sprintf("https://github.com/%s.git", *gitRepository)

		// ary to get access token from GitHub installation
		var accessToken string
		_, err := GetInstallationIDByUserID(userID)
		if err == nil {
			installation, err := GetInstallationByUserID(int(userID))
			if err == nil {
				accessToken = installation.AccessToken
			}
		}

		// mark that we should migrate this app's data
		return cloneURL, accessToken, true, nil
	}

	return "", "", false, fmt.Errorf("app has no git repository or clone URL configured")
}
