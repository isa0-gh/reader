import { NavLink, Outlet } from "react-router-dom";
import "../admin.css";

const TABS = [
  { to: "/admin/users", label: "Users" },
  { to: "/admin/s3", label: "S3 Cleanup" },
];

export default function AdminLayout() {
  return (
    <div className="admin-shell">
      <aside className="admin-sidebar">
        <div className="admin-sidebar-label">Admin</div>
        <div className="admin-sidebar-nav">
          {TABS.map((t) => (
            <NavLink
              key={t.to}
              to={t.to}
              className={({ isActive }) => "admin-sidebar-link" + (isActive ? " admin-sidebar-link--active" : "")}
            >
              {t.label}
            </NavLink>
          ))}
        </div>
      </aside>
      <div className="admin-content">
        <Outlet />
      </div>
    </div>
  );
}
