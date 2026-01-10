package docker

import (
	"fmt"

	"github.com/corecollectives/mist/models"
)

func GetDeploymentConfigForApp(app *models.App) (int, []string, map[string]string, error) {
	port := 3000
	if app.Port != nil {
		port = int(*app.Port)
	}

	domains, err := models.GetDomainsByAppID(app.ID)
	if err != nil && err.Error() != "sql: no rows in result set" {
		return 0, nil, nil, fmt.Errorf("get domains failed: %w", err)
	}

	var domainStrings []string
	for _, d := range domains {
		domainStrings = append(domainStrings, d.Domain)
	}

	envs, err := models.GetEnvVariablesByAppID(app.ID)
	if err != nil && err.Error() != "sql: no rows in result set" {
		return 0, nil, nil, fmt.Errorf("get env variables failed: %w", err)
	}

	envMap := make(map[string]string)
	for _, env := range envs {
		envMap[env.Key] = env.Value
	}

	return port, domainStrings, envMap, nil
}
