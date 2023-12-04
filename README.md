# code-editor
code-server running on Kubernetes with Rest API authentication and multi-user support.

## Usage

- Install `code-editor`
  
  ```bash
  npm install
  npm run cluster:create
  kubectl apply -f https://raw.githubusercontent.com/traefik/traefik/v2.10/docs/content/reference/dynamic-configuration/kubernetes-crd-definition-v1.yml
  kubectl apply -f https://raw.githubusercontent.com/traefik/traefik/v2.10/docs/content/reference/dynamic-configuration/kubernetes-crd-rbac.yml
  npm run code-editor:install
  ```
- Authenticate a user:
  ```
  POST http://localhost/code-editor/api/v1/auth

  req:
  {
      "user": "user1",
      "password": "password1"
  }
  resp:
  {
    "code-server-session": "some-string",
    "path": "code-editor/some-path-for-user1"
  }
  ```
- Go to http://localhost/code-editor/some-path-for-user1/?folder=/git/code-editor
