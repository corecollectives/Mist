import React, { createContext, useContext, useEffect, useState } from "react";
import { User } from "../lib/types";



type AuthContextType = {
  setupRequired: boolean | null;
  setSetupRequired: React.Dispatch<React.SetStateAction<boolean | null>>;
}

const AuthContext = createContext<AuthContextType | null>(null);


export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [setupRequired, setSetupRequired] = React.useState<boolean | null>(null);
  const [user, SetUser] = useState<User | null>(null)

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
  }, []);

  return (
    <AuthContext.Provider value={{ setupRequired, setSetupRequired }}>
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
