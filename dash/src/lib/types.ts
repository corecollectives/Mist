

export type User = {
  id: string;
  username: string;
  email: string;
  role: "owner" | "admin" | "user";
  isAdmin: boolean;
}
