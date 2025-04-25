// src/shared/auth/AuthProvider.tsx
import { createContext, useEffect, useState, ReactNode } from "react";

type User = {
  id: string;
  name: string;
};

type AuthContextType = {
  user: User | null;
  loading: boolean;
  refetch: () => void;
};

export const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  const fetchUser = async () => {
    try {
      setLoading(true);
      const res = await fetch("http://localhost:8080/api/user/me", {// TODO: 環境変数化
        method: "GET",
        credentials: "include", // ← Cookieを送信
      });
      if (res.ok) {
        const data = await res.json();
        setUser(data);
      } else {
        setUser(null);
      }
    } catch {
      setUser(null);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchUser();
  }, []);

  return (
    <AuthContext.Provider value={{ user, loading, refetch: fetchUser }}>
      {children}
    </AuthContext.Provider>
  );
};
