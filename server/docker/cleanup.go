package docker

import (
	"context"
	"fmt"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func CleanupOldImages(appID int64, keepCount int) error {
	if keepCount < 1 {
		keepCount = 5
	}

	imagePattern := fmt.Sprintf("mist-app-%d-", appID)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	listCmd := exec.CommandContext(ctx, "docker", "images",
		"--filter", fmt.Sprintf("reference=%s*", imagePattern),
		"--format", "{{.Repository}}:{{.Tag}} {{.CreatedAt}}",
	)

	output, err := listCmd.Output()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("listing images timed out")
		}
		return fmt.Errorf("failed to list images: %w", err)
	}

	if len(output) == 0 {
		return nil
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) <= keepCount {
		return nil
	}

	type imageInfo struct {
		name      string
		timestamp string
	}

	var images []imageInfo
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, " ", 2)
		if len(parts) == 2 {
			images = append(images, imageInfo{
				name:      parts[0],
				timestamp: parts[1],
			})
		}
	}

	sort.Slice(images, func(i, j int) bool {
		return images[i].timestamp > images[j].timestamp
	})

	if len(images) > keepCount {
		imagesToRemove := images[keepCount:]

		for _, img := range imagesToRemove {
			rmiCtx, rmiCancel := context.WithTimeout(context.Background(), 1*time.Minute)
			rmiCmd := exec.CommandContext(rmiCtx, "docker", "rmi", "-f", img.name)
			if err := rmiCmd.Run(); err != nil {
				fmt.Printf("Warning: Failed to remove image %s: %v\n", img.name, err)
			}
			rmiCancel()
		}
	}

	return nil
}

func CleanupDanglingImages() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	pruneCmd := exec.CommandContext(ctx, "docker", "image", "prune", "-f")
	if err := pruneCmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("pruning images timed out")
		}
		return fmt.Errorf("failed to prune dangling images: %w", err)
	}

	return nil
}

func CleanupStoppedContainers() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	pruneCmd := exec.CommandContext(ctx, "docker", "container", "prune", "-f")
	if err := pruneCmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("pruning containers timed out")
		}
		return fmt.Errorf("failed to prune stopped containers: %w", err)
	}

	return nil
}

func SystemPrune() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	pruneCmd := exec.CommandContext(ctx, "docker", "system", "prune", "-f")
	output, err := pruneCmd.CombinedOutput()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("system prune timed out")
		}
		return string(output), fmt.Errorf("failed to run system prune: %w", err)
	}

	return string(output), nil
}

func SystemPruneAll() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	pruneCmd := exec.CommandContext(ctx, "docker", "system", "prune", "-a", "-f")
	output, err := pruneCmd.CombinedOutput()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("aggressive system prune timed out")
		}
		return string(output), fmt.Errorf("failed to run aggressive system prune: %w", err)
	}

	return string(output), nil
}
