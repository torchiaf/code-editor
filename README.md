# code-editor
[code-server](https://github.com/coder/code-server) running on Kubernetes with Rest API authentication and multi-user support.

## Requirements

- Node
- Helm
- k3d
- kubectl
- yq

## Usage

- Install `code-editor`
  
  ```bash
  npm install
  npm run cluster:create
  npm run traefik:install
  npm run code-editor:install
  ```
- Authenticate a user:
  ```
  POST http://localhost/code-editor/api/v1/login

  body:
  {
      "user": "user1",
      "password": "password1"
  }
  resp:
  {
    "token: ...
  }
  ```
- Enable `code-server` instance:
  ```
  POST http://localhost/code-editor/api/v1/enable

  header:
  {
      "token": ...
  }
  resp:
  {
    "status": "enabled",
    "code-server-session": some-token,
    "path": some-string,
  }
  ```
- Set vscode configs:
  ```
  POST http://localhost/code-editor/api/v1/enable

  header:
  {
      "token": ...
  }
  body:
  {
    "git": {
      org,
      repo,
      branch
      ...
    },
  }
  resp:
  {
    "status": "configured",
		"query": some query params to access the git repository
  }
  ```
- Go to http://localhost/code-editor/$path-for-user1/?folder=$query
