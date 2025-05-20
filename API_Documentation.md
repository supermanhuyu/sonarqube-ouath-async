# API Documentation

This project provides APIs for synchronizing GitLab user/group permissions to SonarQube, as well as querying related data. All APIs are prefixed with `/api/v1`.

---

## 1. Synchronization APIs

### 1.1 Trigger Full Async Sync
- **Method:** POST
- **Endpoint:** `/api/v1/async`
- **Auth:** Required (SonarQube Token in header)
- **Request Body:** None
- **Response:**
  - `200 OK`
  - Example:
    ```json
    { "data": "ok" }
    ```

### 1.2 Sync All Projects
- **Method:** POST
- **Endpoint:** `/api/v1/sync/all`
- **Auth:** Required
- **Request Body:** None
- **Response:**
  - `200 OK`
  - Example:
    ```json
    { "success": true, "synced": 123 }
    ```

### 1.3 Sync Single Project
- **Method:** POST
- **Endpoint:** `/api/v1/sync/project/:projectId`
- **Auth:** Required
- **Path Parameter:**
  - `projectId` (string, required): GitLab project ID
- **Request Body:** None
- **Response:**
  - `200 OK`
  - Example:
    ```json
    { "success": true, "projectId": "123" }
    ```

### 1.4 Sync User Permissions
- **Method:** POST
- **Endpoint:** `/api/v1/sync/user/:username`
- **Auth:** Required
- **Path Parameter:**
  - `username` (string, required): GitLab username
- **Request Body:** None
- **Response:**
  - `200 OK`
  - Example:
    ```json
    { "success": true, "username": "johndoe" }
    ```

---

## 2. SonarQube Data Query APIs

### 2.1 List All SonarQube Projects
- **Method:** GET
- **Endpoint:** `/api/v1/sonar/projects/list`
- **Auth:** Required
- **Query Parameters:** None
- **Response:**
  - `200 OK`
  - Example:
    ```json
    [
      { "kee": "project-key", "name": "Project Name", ... },
      ...
    ]
    ```

### 2.2 List All SonarQube Users
- **Method:** GET
- **Endpoint:** `/api/v1/sonar/users/list`
- **Auth:** Required
- **Query Parameters:** None
- **Response:**
  - `200 OK`
  - Example:
    ```json
    [
      { "login": "user1", "name": "User One", ... },
      ...
    ]
    ```

---

## 3. GitLab Data Query APIs

### 3.1 List All GitLab Projects
- **Method:** GET
- **Endpoint:** `/api/v1/gitlab/projects/list`
- **Auth:** Required
- **Query Parameters:** None
- **Response:**
  - `200 OK`
  - Example:
    ```json
    [
      { "id": 123, "name": "Project Name", ... },
      ...
    ]
    ```

### 3.2 List Members of a GitLab Project
- **Method:** GET
- **Endpoint:** `/api/v1/gitlab/users/:projectId`
- **Auth:** Required
- **Path Parameter:**
  - `projectId` (string, required): GitLab project ID
- **Response:**
  - `200 OK`
  - Example:
    ```json
    [
      { "id": 1, "username": "johndoe", ... },
      ...
    ]
    ```

---

## Auth
All endpoints require authentication via a SonarQube API token, passed in the `Authorization` header.

## Error Response Example
```json
{
  "err": "error message"
}
```

---
For more details, see `README.md` or source code comments.
