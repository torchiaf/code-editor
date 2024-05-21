export interface View {
  Id: string;
  UserId?: string;
  Path: string;
  Status: string;
  VScodeSettings: string;
  Session: string;
  RepoType: string;
  Repo: string;
}

export interface Extension {
  id: string;
  name?: string;
  settings: object;
}

export interface ViewCreateGeneral {
  git: {
    name: string;
    email: string;
  };
  extensions: Extension[];
  vscodeSettings?: string;
}

export interface ViewCreateRepo {
  git: {
    type: string,
    org: string,
    repo: string,
    branch: string,
    commit: string
  }
}

export interface ViewCreate {
  general: ViewCreateGeneral,
  repo?: ViewCreateRepo
}

