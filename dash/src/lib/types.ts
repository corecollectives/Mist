

export type User = {
  id: string | number;
  username: string;
  email: string;
  role: "owner" | "admin" | "user";
  isAdmin: boolean;
  createdAt?: string;
  updatedAt?: string;
}
