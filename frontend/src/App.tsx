import { Routes, Route, Navigate } from "react-router-dom";
import { AuthProvider } from "./AuthContext";
import { ConfigProvider } from "./ConfigContext";
import { ToastProvider } from "./ToastContext";
import { ConfirmProvider } from "./ConfirmContext";
import { FavoritesProvider } from "./FavoritesContext";
import SiteBackground from "./SiteBackground";
import Nav from "./Nav";
import Home from "./pages/Home";
import SeriesPage from "./pages/SeriesPage";
import ReaderPage from "./pages/ReaderPage";
import ChapterEditPage from "./pages/ChapterEditPage";
import Login from "./pages/Login";
import Register from "./pages/Register";
import AccountPage from "./pages/AccountPage";

import AdminLayout from "./pages/AdminLayout";
import AdminUsersPage from "./pages/AdminUsersPage";
import AdminS3CleanPage from "./pages/AdminS3CleanPage";
import AdminBackgroundsPage from "./pages/AdminBackgroundsPage";

export default function App() {
  return (
    <ConfigProvider>
      <SiteBackground />
      <AuthProvider>
        <ToastProvider>
          <ConfirmProvider>
            <FavoritesProvider>
              <Nav />
              <Routes>
                <Route path="/" element={<Home />} />
                <Route path="/series/:id" element={<SeriesPage />} />
                <Route path="/series/:id/:chapterId" element={<ReaderPage />} />
                <Route path="/series/:id/:chapterId/edit" element={<ChapterEditPage />} />
                <Route path="/login" element={<Login />} />
                <Route path="/register" element={<Register />} />
                <Route path="/account" element={<AccountPage />} />
                <Route path="/admin" element={<AdminLayout />}>
                  <Route index element={<Navigate to="users" replace />} />
                  <Route path="users" element={<AdminUsersPage />} />
                  <Route path="backgrounds" element={<AdminBackgroundsPage />} />
                  <Route path="s3" element={<AdminS3CleanPage />} />
                </Route>
              </Routes>
            </FavoritesProvider>
          </ConfirmProvider>
        </ToastProvider>
      </AuthProvider>
    </ConfigProvider>
  );
}
