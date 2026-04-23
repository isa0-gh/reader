import { jsx as _jsx } from "react/jsx-runtime";
import { createContext, useContext, useState } from "react";
const Ctx = createContext(null);
export function AuthProvider({ children }) {
    const [token, setToken] = useState(() => localStorage.getItem("token"));
    const [user, setUser] = useState(() => {
        const u = localStorage.getItem("user");
        return u ? JSON.parse(u) : null;
    });
    function login(t, u) {
        localStorage.setItem("token", t);
        localStorage.setItem("user", JSON.stringify(u));
        setToken(t);
        setUser(u);
    }
    function logout() {
        localStorage.removeItem("token");
        localStorage.removeItem("user");
        setToken(null);
        setUser(null);
    }
    return _jsx(Ctx.Provider, { value: { user, token, login, logout }, children: children });
}
export const useAuth = () => useContext(Ctx);
