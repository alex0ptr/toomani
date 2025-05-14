# toomani

The amount of repositories we have to manage is too damn high!
Pulling all repositories from a specific group in your workspace can be tedious.  
[Mani](https://manicli.com/) is an excellent tool for pulling and managing multiple repositories.
However, setting up the initial configuration file for Mani can be time-consuming.

This is where **toomani** helps.  
It generates an initial Mani configuration file, or optionally a shell script, for a specific group on your GitLab server.

### Usage

Run the following command:

```shell
export GITLAB_HOST="<your self-hosted GitLab instance, if applicable>"
export GITLAB_TOKEN="<your GitLab personal access token>"

go run github.com/alex0ptr/toomani@latest gitlab --group gitlab-org/architecture -o mani > mani.yml
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

Note that you also have the option to filter which repositories you want to include in the mani configuration file.
Use `go run github.com/alex0ptr/toomani@latest gitlab -h` to discover options.