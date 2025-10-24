import React, { createContext, useContext, useEffect, useState } from "react";
import type { User } from "../lib/types";


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
      } else {
        console.error("Logout failed:", data.error);
      }
    }
    catch (error) {
      console.error("Error during logout:", error);
    }
  }

  useEffect(() => {
    const checkSetupRequired = async () => {
      try {
        const response = await fetch("/api/auth/check-setup-status");
        const data = await response.json();
        setSetupRequired(data.setupRequired);
      } catch (error) {
        console.error("Error checking setup requirement:", error);
      }
    };
    checkSetupRequired();

    const fetchUser = async () => {
      try {
        const response = await fetch("/api/auth/me", {
          credentials: "include"
        });
        const data = await response.json();
        if (data.success) {
          setUser({ ...data.data, isAdmin: data.data.role === "owner" || data.data.role === "admin" });
        } else {
          setUser(null);
        }
      } catch (error) {
        console.error("Error fetching user:", error);
        setUser(null);
      }
    }
    fetchUser();
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
