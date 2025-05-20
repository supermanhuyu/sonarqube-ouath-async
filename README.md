# Sync GitLab user permissions to SonarQube

> **Note:** Most of this project was created with the help of VSCode Copilot.

This service synchronizes GitLab user and group permissions to SonarQube, and provides HTTP APIs for manual or automated sync and data queries.

## Build

```shell
go build -o sonarqube-ouath-async main.go
```

## Run

```shell
./sonarqube-ouath-async -gitlabAddr https://gitlab.example.com -gitlabToken "your_gitlab_token" -sonarAddr http://sonar.example.com:9000 -sonarToken "your_sonarqube_token"
```

## API Endpoints

All APIs are prefixed with `/api/v1`.

### Sync APIs

- **POST /api/v1/async**

  Trigger a full async sync from GitLab to SonarQube (runs in background).

  ```shell
  curl -X POST http://localhost:8080/api/v1/async
  ```

- **POST /api/v1/sync/all**

  Sync all GitLab project members and permissions to SonarQube.

  ```shell
  curl -X POST http://localhost:8080/api/v1/sync/all
  ```

- **POST /api/v1/sync/project/:projectId**

  Sync a single GitLab project to SonarQube.

  ```shell
  curl -X POST http://localhost:8080/api/v1/sync/project/123
  ```

- **POST /api/v1/sync/user/:username**

  Sync all permissions for a specific GitLab user to SonarQube.

  ```shell
  curl -X POST http://localhost:8080/api/v1/sync/user/johndoe
  ```

### SonarQube Data Query APIs

- **GET /api/v1/sonar/projects/list**

  Get all SonarQube projects (from the database).

  ```shell
  curl http://localhost:8080/api/v1/sonar/projects/list
  ```

- **GET /api/v1/sonar/users/list**

  Get all SonarQube users.

  ```shell
  curl http://localhost:8080/api/v1/sonar/users/list
  ```

### GitLab Data Query APIs

- **GET /api/v1/gitlab/projects/list**

  Get all GitLab projects.

  ```shell
  curl http://localhost:8080/api/v1/gitlab/projects/list
  ```

- **GET /api/v1/gitlab/users/:projectId**

  Get all members of a specific GitLab project.

  ```shell
  curl http://localhost:8080/api/v1/gitlab/users/123
  ```

## Configuration

You can pass configuration via command-line flags:

- `-gitlabAddr` : GitLab server address (e.g., https://gitlab.example.com)
- `-gitlabToken` : GitLab personal access token
- `-sonarAddr` : SonarQube server address (e.g., http://sonar.example.com:9000)
- `-sonarToken` : SonarQube API token

## Logging

Logs are printed to stdout and include info and error messages for all sync and API operations.

## License

MIT
