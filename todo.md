Todo

- move initContainers to separate template files
- activate/deactivate option tu create/destroy endpoint 
- add dynamic k8s client
- use rule template to create rule
- complete the external login check; should works fine when external API token expires, url is incorrect or API response is not 200
- external user ids should be generate with same algorithm of local users (helm chart)
- resolve TODOs in the code
- errors handling
- check log errors on destroyRule() call
- make "code-editor" prefix parametric
- vs-code-settings ~~and extensions~~ must be installed via API
- Write license
- Add a saved postman endpoints for dev
- Clean helm charts
- Add swagger
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
