package docker

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/docker/docker/pkg/archive"
	"github.com/gorilla/websocket"

	// "github.com/moby/go-archive"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/network"

	// "github.com/moby/moby/api/types/strslice"
	"github.com/moby/moby/client"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func DeployHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer ws.Close()
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("could not initialise docker client"+err.Error()))
		return
	}

	defer cli.Close()

	repoDir, _ := filepath.Abs("../../test/")
	fmt.Println("repoDir: ", repoDir)

	tar, err := archive.TarWithOptions(repoDir, &archive.TarOptions{})
	if err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("could not create tar"+err.Error()))
		return
	}

	imageTag := "testimage1"

	buildResponse, err := cli.ImageBuild(ctx, tar, client.ImageBuildOptions{
		Tags:       []string{imageTag},
		Dockerfile: "Dockerfile",
		Remove:     true,
	})
	if err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("could not build image"+err.Error()))
	}
	defer buildResponse.Body.Close()

	scanner := bufio.NewScanner(buildResponse.Body)
	for scanner.Scan() {
		ws.WriteMessage(websocket.TextMessage, []byte(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("errorreading build response: "+err.Error()))
		return
	}

	contName := "myapp_container"
	timeout := 10
	_ = cli.ContainerStop(ctx, contName, client.ContainerStopOptions{
		Timeout: &timeout,
	})
	_ = cli.ContainerRemove(ctx, contName, client.ContainerRemoveOptions{
		Force: true,
	})

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageTag,
		Tty:   false,
	}, &container.HostConfig{
		PortBindings: map[network.Port][]network.PortBinding{},
	}, nil, nil, contName)

	if err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("error container create failed: "+err.Error()))
		return
	}
	if err := cli.ContainerStart(ctx, resp.ID, client.ContainerStartOptions{}); err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("error container start failed: "+err.Error()))
		return
	}
	ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("container %s started\n", resp.ID)))

	logOpts := client.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Timestamps: true,
	}

	logReader, err := cli.ContainerLogs(ctx, resp.ID, logOpts)
	if err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("error container logs failed: "+err.Error()))
		return
	}
	defer logReader.Close()
	scanner = bufio.NewScanner(logReader)
	for scanner.Scan() {
		ws.WriteMessage(websocket.TextMessage, []byte(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("error reading logs: "+err.Error()))
		return
	}

	ws.WriteMessage(websocket.TextMessage, []byte("deployment doneeee"))

}
