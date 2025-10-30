import React, { createContext, useContext, useEffect, useState } from "react";
import type { User } from "../lib/types";
import { toast } from "sonner";


type AuthContextType = {
  setupRequired: boolean | null;
  setSetupRequired: React.Dispatch<React.SetStateAction<boolean | null>>;
  user: User | null;
  setUser: React.Dispatch<React.SetStateAction<User | null>>;
  logout: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | null>(null);


export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [setupRequired, setSetupRequired] = React.useState<boolean | null>(null);
  const [user, setUser] = useState<User | null>(null)


  const logout = async () => {
    try {
      const response = await fetch("/api/auth/logout", {
        method: "POST",
        credentials: "include"
      });
      const data = await response.json();
      if (data.success) {
        setUser(null);
        toast.success("Logged out successfully");
      } else {
        console.error("Logout failed:", data.error);
      }
    }
    catch (error) {
      console.error("Error during logout:", error);
    }
  }

  useEffect(() => {
    const fetchAuth = async () => {
      try {
        const response = await fetch("/api/auth/me", {
          method: "GET",
          credentials: "include"
        });
        const data = await response.json();
        if (data.success) {
          if (data.data.setupRequired === true) {
            setSetupRequired(true);
          }
          else {
            setSetupRequired(false);
            if (data.data.user) {
              setUser({ ...data.data.user, isAdmin: data.data.user.role === "owner" || data.data.user.role === "admin" });
            }
          }
        }
      }
      catch (error) {
        console.error("Error fetching auth data:", error);
      }
    }
    fetchAuth();
  }, []);

  return (
    <AuthContext.Provider value={{ setupRequired, setSetupRequired, user, setUser, logout }}>
      {children}
    </AuthContext.Provider>
  );

}

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === null) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}
