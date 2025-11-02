package dockerdeploy

import (
	"fmt"
)

func (d *Deployer) DeployerMain(Id int64) (string, error) {

	// var id int64
	// err := d.DB.QueryRow("SELECT id FROM deployments WHERE id = ?", Id).Scan(&id)
	// if err == sql.ErrNoRows {
	// 	fmt.Println("No deployment found with id:", Id)
	// } else if err != nil {
	// 	fmt.Println("Query error:", err)
	// } else {
	// 	fmt.Println("Deployment found with id:", id)
	// }
	// fmt.Println("Rows: ", rows)
	fmt.Println("Deploying deployment", Id)
	dep, err := d.loadDeployment(Id)
	// if err != nil {
	// 	http.Error(w, "Deployment not found", 404)
	// 	return
	// }

	if err != nil {
		fmt.Printf("Error loading deployment: %v\n", err)

		return "", err
	}

	fmt.Printf("Loaded deployment: %+v\n", dep)
	appContextPath := "../../test"
	imageTag := dep.CommitHash
	containerName := fmt.Sprintf("app-%d", Id)

	go func() {
		err := d.DeployApp(dep, appContextPath, imageTag, containerName)
		if err != nil {
			// http.Error(w, err.Error(), 500)
			fmt.Printf("Error deploying app: %v\n", err)
			dep.Status = "failed"
			d.UpdateDeployment(dep)
		}
	}()

	return "Deployment started", nil
}
