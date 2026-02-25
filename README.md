# toomani

The amount of repositories we have to manage is too damn high!
Pulling all repositories from a specific group in your workspace can be tedious.  
[Mani](https://manicli.com/) is an excellent tool for pulling and managing multiple repositories.
However, setting up the initial configuration file for Mani can be time-consuming.

This is where **toomani** helps.
It generates an initial Mani configuration file, or optionally a shell script, for a specific group or organization on **GitLab** or **GitHub**.

## Demo

Run the following command:

```shell
export GITLAB_TOKEN="<your GitLab personal access token>"

go run github.com/alex0ptr/toomani@latest gitlab \
  --group gitlab-org/architecture \
  > mani.yml
```

This will generate a Mani configuration file for the `gitlab-org/architecture` group:

```yaml
projects:
    design-doc:
        url: git@gitlab.com:gitlab-org/architecture/autoflow/design-doc.git
        path: autoflow/design-doc
    design-docs:
        url: git@gitlab.com:gitlab-org/architecture/design-docs.git
        path: design-docs
    staff-plus-engineer-role-questions:
        url: git@gitlab.com:gitlab-org/architecture/staff-plus-engineer-role-questions.git
        path: staff-plus-engineer-role-questions
tasks:
    fetch-all:
        cmd: git fetch --all
        description: Fetch all remotes.
    pre-commit:
        cmd: pre-commit install --allow-missing-config
        description: Install pre-commit hooks.
    pull-main:
        cmd: git switch main && git pull
        description: Switch to main and pull.
```

You can now run:

```shell
mani sync
mani run --all "<task>"
```

## Advanced Usage

Use `go run github.com/alex0ptr/toomani@latest gitlab -h` or `go run github.com/alex0ptr/toomani@latest github -h` to discover all options.

### Self-Hosted GitLab

Use the `--host` flag or `GITLAB_HOST` env var to point to your self-hosted instance:

```shell
export GITLAB_HOST="gitlab.mycompany.com"
export GITLAB_TOKEN="<your token>"

go run github.com/alex0ptr/toomani@latest gitlab \
  --group my-team \
  > mani.yml
```

### GitHub

toomani also supports generating a Mani configuration file from a GitHub organization or user account:

```shell
export GITHUB_TOKEN="<your GitHub personal access token>"

go run github.com/alex0ptr/toomani@latest github \
  --owner my-org \
  > mani.yml
```

### Exclude Sub-groups or Prefixes

```shell
# GitLab
go run github.com/alex0ptr/toomani@latest gitlab \
  --group my-team \
  --exclude-prefix my-team/playground,my-team/archive \
  > mani.yml

# GitHub
go run github.com/alex0ptr/toomani@latest github \
  --owner my-org \
  --exclude-prefix my-org/playground,my-org/archive \
  > mani.yml
```

### Include Only Repositories with a Specific Prefix

```shell
# GitLab
go run github.com/alex0ptr/toomani@latest gitlab \
  --group my-team \
  --match-prefix my-team/java-,my-team/important- \
  > mani.yml

# GitHub
go run github.com/alex0ptr/toomani@latest github \
  --owner my-org \
  --match-prefix my-org/java-,my-org/important- \
  > mani.yml
```