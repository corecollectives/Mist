

export type User = {
  id: string | number;
  username: string;
  email: string;
  role: "owner" | "admin" | "user";
  isAdmin?: boolean;
  createdAt?: string;
  updatedAt?: string;
}

export type Project = {
  id: number;
  name: string;
  description: string;
  tags?: string[];
  ownerId: string | number;
  owner?: User
  projectMembers: User[];
  createdAt?: string;
  updatedAt?: string

}
