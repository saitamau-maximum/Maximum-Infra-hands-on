import { createContext, useEffect, useState, ReactNode } from "react";
import { Logout } from "../api/logout";

type User = {
  id: string;
  name: string;
};

type AuthContextType = {
  user: User | null;
  loading: boolean;
  refetch: () => void;
  logout: () => void; 
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

  const logout = () => {
    // クッキーを削除するために、期限切れのクッキーをセット
    document.cookie = "token=; Max-Age=0; path=/; SameSite=None; Secure"; 
    Logout();
    setUser(null); // ユーザー情報のリセット
  };

  useEffect(() => {
    fetchUser();
  }, []);

  return (
    <AuthContext.Provider value={{ user, loading, refetch: fetchUser, logout }}>
      {children}
    </AuthContext.Provider>
  );
};
