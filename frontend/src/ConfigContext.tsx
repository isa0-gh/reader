import { createContext, useContext, useEffect, useState, ReactNode } from "react";

interface AppConfig {
  cdn_url: string;
  register_disabled: boolean;
  login_disabled: boolean;
  maintenance: boolean;
}

const defaults: AppConfig = { cdn_url: "", register_disabled: false, login_disabled: false, maintenance: false };
const Ctx = createContext<AppConfig>(defaults);

export function ConfigProvider({ children }: { children: ReactNode }) {
  const [cfg, setCfg] = useState<AppConfig>(defaults);

  useEffect(() => {
    fetch("/api/v1/config")
      .then(r => r.json())
      .then(setCfg)
      .catch(() => {});
  }, []);

  return <Ctx.Provider value={cfg}>{children}</Ctx.Provider>;
}

export const useConfig = () => useContext(Ctx);
