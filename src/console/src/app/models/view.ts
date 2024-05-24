export interface View {
  Id: string;
  Name: string;
  UserId?: string;
  Path: string;
  Query: string;
  Status: string;
  VScodeSettings: string;
  GitAuth: boolean;
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
  name: string;
  git: {
    name: string;
    email: string;
  };
  extensions: Extension[];
  vscodeSettings?: string;
  sshKey: string;
}

export interface ViewCreateRepo {
  git: {
    type: string | null,
    org: string | null,
    repo: string | null,
    branch: string | null,
    commit: string | null
  }
}

export interface ViewCreate {
  general: ViewCreateGeneral,
  repo: ViewCreateRepo
}

