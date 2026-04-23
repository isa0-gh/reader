import { Routes, Route } from "react-router-dom";
import { AuthProvider } from "./AuthContext";
import Nav from "./Nav";
import Home from "./pages/Home";
import SeriesPage from "./pages/SeriesPage";
import ReaderPage from "./pages/ReaderPage";
import Login from "./pages/Login";
import Register from "./pages/Register";

export default function App() {
  return (
    <AuthProvider>
      <Nav />
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/series/:id" element={<SeriesPage />} />
        <Route path="/chapters/:id" element={<ReaderPage />} />
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />
      </Routes>
    </AuthProvider>
  );
}
