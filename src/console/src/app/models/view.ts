export interface View {
  Id: string;
  UserId?: string;
  Path: string;
  Status: string;
  VScodeSettings: string;
}

export interface Extension {
  id: string;
  name?: string;
  settings: object;
}

export interface ViewCreate {
  git: {
    name: string;
    email: string;
  };
  extensions: Extension[];
  vscodeSettings?: string;
}

