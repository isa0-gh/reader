import { jsx as _jsx, jsxs as _jsxs } from "react/jsx-runtime";
import { Routes, Route } from "react-router-dom";
import { AuthProvider } from "./AuthContext";
import Nav from "./Nav";
import Home from "./pages/Home";
import SeriesPage from "./pages/SeriesPage";
import ReaderPage from "./pages/ReaderPage";
import Login from "./pages/Login";
import Register from "./pages/Register";
export default function App() {
    return (_jsxs(AuthProvider, { children: [_jsx(Nav, {}), _jsxs(Routes, { children: [_jsx(Route, { path: "/", element: _jsx(Home, {}) }), _jsx(Route, { path: "/series/:id", element: _jsx(SeriesPage, {}) }), _jsx(Route, { path: "/chapters/:id", element: _jsx(ReaderPage, {}) }), _jsx(Route, { path: "/login", element: _jsx(Login, {}) }), _jsx(Route, { path: "/register", element: _jsx(Register, {}) })] })] }));
}
