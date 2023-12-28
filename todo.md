Todo

- errors handling
- add user ID
- make "code-editor" prefix parametric
- Replace username with user ID in routing and object names
- vs-code-settings and extensions must be installed via APis, not hardcoded -> remove them from src/templates
- Write license
- Create an external authentication mode:
  -1 DONE pre-requisite: routes dynamic creation (services, deployments, ingressroutes)
  -2 helm param to switch to external authentication
  -3 helm param to save the url of external login
  -4 implement a new /register-user endpoint to make the login to the external authentication and register the user into code-editor
- Add a saved postman endpoints for dev
- Clean helm charts
- Add swagger
- disable telemetry: --disable-telemetry option
- Investigate kubevirt
- Investigate, check heartbit file: https://coder.com/docs/code-server/latest/FAQ#where-is-vs-code-configuration-stored
- Investigate: https://coder.com/docs/code-server/latest/FAQ#where-is-vs-code-configuration-stored
- Explore projetcs section: https://github.com/coder/awesome-code-server
- Explore Faq section https://coder.com/docs/code-server/latest/FAQ#where-is-vs-code-configuration-stored
- Explore vscode settings. Settings are configurable by setting the file: home/coder/.local/share/code-server/Machine/settings.json
- Implement a k8s controller to create a pool of code-server pods to be assigned dynamically to the users, with dynamic authentication 
- Implement a k8s controller to keep healthy code-server pods
- Refactoring docs
  readme example: https://github.com/andreabenini/podmaster/tree/main/forklift/
