package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

func ExecCommand(ctx context.Context, command string) error {
	cmd := exec.CommandContext(ctx, "/usr/bin/env", "bash", "-c", command)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func ExecComposeReset(ctx context.Context) error {
	return ExecCommand(ctx, "docker compose down -v && docker compose up -d")
}

func ExecGitlabRunnerRegister(ctx context.Context, gitlabURL, token string) error {
	return ExecCommand(ctx, fmt.Sprintf(`docker compose exec gitlab-runner \
  gitlab-runner register --non-interactive \
  --url "%s" \
  --token "%s" \
  --name GitClassrooms \
  --executor docker \
  --docker-image "ubuntu:24"
 `, gitlabURL, token))
}
