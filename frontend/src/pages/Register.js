import { jsx as _jsx, jsxs as _jsxs } from "react/jsx-runtime";
import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { api } from "../api";
export default function Register() {
    const nav = useNavigate();
    const [name, setName] = useState("");
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [error, setError] = useState("");
    async function submit(e) {
        e.preventDefault();
        setError("");
        try {
            await api.register(email, password, name);
            nav("/login");
        }
        catch {
            setError("Registration failed. Email may already be in use.");
        }
    }
    return (_jsxs("div", { className: "auth-wrap", children: [_jsx("h1", { children: "Register" }), _jsxs("form", { onSubmit: submit, children: [_jsxs("div", { className: "form-group", children: [_jsx("label", { children: "Name" }), _jsx("input", { value: name, onChange: (e) => setName(e.target.value), required: true })] }), _jsxs("div", { className: "form-group", children: [_jsx("label", { children: "Email" }), _jsx("input", { type: "email", value: email, onChange: (e) => setEmail(e.target.value), required: true })] }), _jsxs("div", { className: "form-group", children: [_jsx("label", { children: "Password" }), _jsx("input", { type: "password", value: password, onChange: (e) => setPassword(e.target.value), required: true, minLength: 8 })] }), error && _jsx("div", { className: "error", children: error }), _jsx("button", { className: "btn", type: "submit", children: "Create account" })] }), _jsxs("p", { className: "auth-switch", children: ["Have an account? ", _jsx(Link, { to: "/login", children: "Login" })] })] }));
}
