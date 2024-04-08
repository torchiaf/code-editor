# code-editor
code-editor is an open source application to deploy [code-server](https://github.com/coder/code-server) on Kubernetes, adding a JWT authentication and multi-user support.

## Requirements

- Node
- Helm
- k3d
- kubectl
- yq

## Quick Start

### Installation
  
    # Provide dev tools
    npm install

    # Create a new cluster using k3d
    npm run cluster:create

    # Install code-editor using Helm charts
    npm run code-editor:install

### Using code-editor

  `code-editor` provides http APIs to manage users and code-server instances. Consult the full API list here [Openapi spec](https://github.com/torchiaf/code-editor/blob/main/docs/openapi.yaml)

#### Workflow
  
- Login:
  ```
  POST http://localhost/code-editor/api/v1/login
  ```
- Enable `code-server` instance:
  ```
  POST http://localhost/code-editor/api/v1/enable

  response:
  {
    "code-server-session": code-server instance token to use in the browser,
    "path": path for the code-server instance,
  }
  ```
- Add `code-server` configs:
  ```
  POST http://localhost/code-editor/api/v1/enable

  response:
  {
    "query-param": some query params to access the git repository
  }
  ```

- Go to https://localhost/code-editor/$path/?folder=$query-param and enjoy the power of VSCode in your web browser!
