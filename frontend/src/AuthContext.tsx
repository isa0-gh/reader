import { createContext, useContext, useState, ReactNode } from "react";

interface AuthUser { id: number; name: string; role: string; }
interface AuthCtx { user: AuthUser | null; token: string | null; login(token: string, user: AuthUser): void; logout(): void; }

const Ctx = createContext<AuthCtx>(null!);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [token, setToken] = useState(() => localStorage.getItem("token"));
  const [user, setUser] = useState<AuthUser | null>(() => {
    const u = localStorage.getItem("user");
    return u ? JSON.parse(u) : null;
  });

  function login(t: string, u: AuthUser) {
    localStorage.setItem("token", t);
    localStorage.setItem("user", JSON.stringify(u));
    setToken(t); setUser(u);
  }

  function logout() {
    localStorage.removeItem("token");
    localStorage.removeItem("user");
    setToken(null); setUser(null);
  }

  return <Ctx.Provider value={{ user, token, login, logout }}>{children}</Ctx.Provider>;
}

export const useAuth = () => useContext(Ctx);
