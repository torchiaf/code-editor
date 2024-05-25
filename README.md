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
  POST https://localhost/code-editor/api/v1/login
  ```
- Create `code-server` instance:
  ```
  POST https://localhost/code-editor/api/v1/views

  response:
  {
    "viewId": "0001",
    "path": path of the code-server instance UI,
    "query-param": path to the cloned repo provisioned in the code-server instance,
    "code-server-session": code-server instance token to use in the browser
  }
  ```

- Go to https://localhost/code-editor/$path/?folder=$query-param and enjoy the power of VSCode in your web browser!

### Console

  The web-gui console is accessible here: https://localhost/code-editor/console .
  - Admin users can assign new Code Editor instances to the standard users.
  - Standard Users can use it to access to their code instances.

  Admin console

  ![image](https://github.com/torchiaf/code-editor/assets/26394656/daeeeca8-269d-439a-8549-863943329ed7)


  Create Page

  ![image](https://github.com/torchiaf/code-editor/assets/26394656/327b275c-5954-4d7c-89ac-12c0cd99bc86)


  Users console

  ![image](https://github.com/torchiaf/code-editor/assets/26394656/dcbcdd78-f83c-4cfd-88b7-b289f636c98a)


  VSCode instance
  
  ![image](https://github.com/torchiaf/code-editor/assets/26394656/4b36a843-2253-4af3-8a66-69783277a3a3)





