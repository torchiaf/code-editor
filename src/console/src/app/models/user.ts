export type Role = 'Admin' | 'User';

export interface UserDetails extends User {
  role: Role;
}

export interface User {
  id: string;
  name?: string;
  username?: string;
}
