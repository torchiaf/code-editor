Todo

- write commands to get code-server password for each users
  #users: echo $(kubectl get secret code-editor-users --namespace code-editor -o jsonpath="{.data.users}" | base64 --decode)
  #password: echo $(kubectl get secret code-editor-config --namespace code-editor -o jsonpath="{.data.code-editor-local-0000_PASSWORD}" | base64 --decode)

- activate/deactivate option tu create/destroy endpoint 
- add dynamic k8s client
- use rule template to create rule
- complete the external login check; should works fine when external API token expires, url is incorrect or API response is not 200
- external user ids should be generate with same algorithm of local users (helm chart)
- error handling
- check log errors on destroyRule() call
- make "code-editor" prefix parametric
- clean helm charts
- add swagger docs
- investigate kubevirt
- investigate, check heartbit file: https://coder.com/docs/code-server/latest/FAQ#where-is-vs-code-configuration-stored
- investigate: https://coder.com/docs/code-server/latest/FAQ#where-is-vs-code-configuration-stored
- explore projetcs section: https://github.com/coder/awesome-code-server
- explore Faq section https://coder.com/docs/code-server/latest/FAQ#where-is-vs-code-configuration-stored
- explore vscode settings. Settings are configurable by setting the file: home/coder/.local/share/code-server/Machine/settings.json
- implement a k8s controller to create a pool of code-server pods to be assigned dynamically to the users, with dynamic authentication 
- implement a k8s controller to keep healthy code-server pods
- refactoring docs
  readme example: https://github.com/andreabenini/podmaster/tree/main/forklift/
