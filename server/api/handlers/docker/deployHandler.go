package docker

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
)

func (d *Deployer) DeployHandler(w http.ResponseWriter, r *http.Request) {
	depIDStr := r.URL.Query().Get("id")
	depID, err := strconv.ParseInt(depIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid deployment ID", 400)
		return
	}
	var id int64
	err = d.DB.QueryRow("SELECT id FROM deployments WHERE id = ?", depID).Scan(&id)
	if err == sql.ErrNoRows {
		fmt.Println("No deployment found with id:", depID)
	} else if err != nil {
		fmt.Println("Query error:", err)
	} else {
		fmt.Println("Deployment found with id:", id)
	}
	// fmt.Println("Rows: ", rows)
	fmt.Println("Deploying deployment", depID)
	dep, err := d.loadDeployment(depID)
	// if err != nil {
	// 	http.Error(w, "Deployment not found", 404)
	// 	return
	// }

	if err != nil {
		fmt.Printf("Error loading deployment: %v\n", err)
		http.Error(w, "Deployment not found", 404)
		return
	}

	fmt.Printf("Loaded deployment: %+v\n", dep)
	appContextPath := "../../test"
	imageTag := dep.CommitHash
	containerName := fmt.Sprintf("app-%d", depID)

	go func() {
		err := d.DeployApp(dep, appContextPath, imageTag, containerName)
		if err != nil {
			// http.Error(w, err.Error(), 500)
			fmt.Printf("Error deploying app: %v\n", err)
			dep.Status = "failed"
			d.UpdateDeployment(dep)
		}
	}()

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Deployment started"))
}
