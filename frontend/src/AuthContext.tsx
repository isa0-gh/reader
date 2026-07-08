import { createContext, useContext, useState, ReactNode } from "react";
import { S3Object } from "./api";

interface AuthUser {
  id: number;
  name: string;
  email: string;
  role: string;
  bio?: string;
  avatar?: S3Object | null;
}
interface AuthCtx {
  user: AuthUser | null;
  token: string | null;
  login(token: string, user: AuthUser): void;
  logout(): void;
  updateUser(patch: Partial<AuthUser>): void;
}

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

  // Patches the cached user (e.g. after editing name/bio/avatar on the
  // account page) without a re-login round trip.
  function updateUser(patch: Partial<AuthUser>) {
    setUser((prev) => {
      if (!prev) return prev;
      const next = { ...prev, ...patch };
      localStorage.setItem("user", JSON.stringify(next));
      return next;
    });
  }

  return <Ctx.Provider value={{ user, token, login, logout, updateUser }}>{children}</Ctx.Provider>;
}

export const useAuth = () => useContext(Ctx);
