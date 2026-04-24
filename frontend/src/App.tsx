import { Routes, Route } from "react-router-dom";
import { AuthProvider } from "./AuthContext";
import { ConfigProvider } from "./ConfigContext";
import Nav from "./Nav";
import Home from "./pages/Home";
import SeriesPage from "./pages/SeriesPage";
import ReaderPage from "./pages/ReaderPage";
import ChapterEditPage from "./pages/ChapterEditPage";
import Login from "./pages/Login";
import Register from "./pages/Register";

import AdminUsersPage from "./pages/AdminUsersPage";
import AdminS3CleanPage from "./pages/AdminS3CleanPage";

export default function App() {
  return (
    <ConfigProvider>
      <AuthProvider>
        <Nav />
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/series/:id" element={<SeriesPage />} />
          <Route path="/series/:id/:chapterId" element={<ReaderPage />} />
          <Route path="/series/:id/:chapterId/edit" element={<ChapterEditPage />} />
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
          <Route path="/admin/users" element={<AdminUsersPage />} />
          <Route path="/admin/s3" element={<AdminS3CleanPage />} />
        </Routes>
      </AuthProvider>
    </ConfigProvider>
  );
}
