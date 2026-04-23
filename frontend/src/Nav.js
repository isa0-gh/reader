import { jsx as _jsx, Fragment as _Fragment, jsxs as _jsxs } from "react/jsx-runtime";
import { Link, useNavigate } from "react-router-dom";
import { useAuth } from "./AuthContext";
export default function Nav() {
    const { user, logout } = useAuth();
    const nav = useNavigate();
    function handleLogout() {
        logout();
        nav("/");
    }
    return (_jsxs("nav", { children: [_jsx(Link, { to: "/", className: "logo", children: "Reader" }), _jsx("div", { className: "nav-links", children: user ? (_jsxs(_Fragment, { children: [_jsx("span", { children: user.name }), _jsx("a", { href: "#", onClick: handleLogout, children: "Logout" })] })) : (_jsxs(_Fragment, { children: [_jsx(Link, { to: "/login", children: "Login" }), _jsx(Link, { to: "/register", children: "Register" })] })) })] }));
}
