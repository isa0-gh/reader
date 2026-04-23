import { jsx as _jsx, jsxs as _jsxs } from "react/jsx-runtime";
import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { api } from "../api";
import { useAuth } from "../AuthContext";
export default function Login() {
    const { login } = useAuth();
    const nav = useNavigate();
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [error, setError] = useState("");
    async function submit(e) {
        e.preventDefault();
        setError("");
        try {
            const { token, user } = await api.login(email, password);
            login(token, user);
            nav("/");
        }
        catch {
            setError("Invalid email or password.");
        }
    }
    return (_jsxs("div", { className: "auth-wrap", children: [_jsx("h1", { children: "Login" }), _jsxs("form", { onSubmit: submit, children: [_jsxs("div", { className: "form-group", children: [_jsx("label", { children: "Email" }), _jsx("input", { type: "email", value: email, onChange: (e) => setEmail(e.target.value), required: true })] }), _jsxs("div", { className: "form-group", children: [_jsx("label", { children: "Password" }), _jsx("input", { type: "password", value: password, onChange: (e) => setPassword(e.target.value), required: true })] }), error && _jsx("div", { className: "error", children: error }), _jsx("button", { className: "btn", type: "submit", children: "Login" })] }), _jsxs("p", { className: "auth-switch", children: ["No account? ", _jsx(Link, { to: "/register", children: "Register" })] })] }));
}
