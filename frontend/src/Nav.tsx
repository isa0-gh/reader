import { useState, useRef, useEffect } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useAuth } from "./AuthContext";
import { useConfig } from "./ConfigContext";

export default function Nav() {
  const { user, logout } = useAuth();
  const { cdn_url } = useConfig();
  const nav = useNavigate();
  const [open, setOpen] = useState(false);
  const ref = useRef<HTMLDivElement>(null);

  useEffect(() => {
    function handler(e: MouseEvent) {
      if (ref.current && !ref.current.contains(e.target as Node)) setOpen(false);
    }
    document.addEventListener("mousedown", handler);
    return () => document.removeEventListener("mousedown", handler);
  }, []);

  function handleLogout() {
    logout();
    setOpen(false);
    nav("/");
  }

  return (
    <nav>
      <Link to="/" className="logo">Reader</Link>
      <div className="nav-links">
        {user ? (
          <div className="account-menu" ref={ref}>
            <button className="account-trigger" onClick={() => setOpen((o) => !o)} aria-haspopup="true" aria-expanded={open}>
              {user.avatar
                ? <img src={`${cdn_url}/${user.avatar.key}`} alt="" className="avatar-xs" />
                : <span className="avatar-xs avatar-xs--placeholder" aria-hidden="true">{user.name.charAt(0).toUpperCase()}</span>}
              <span>{user.name}</span>
              <svg width="12" height="12" viewBox="0 0 12 12" fill="none" aria-hidden="true" style={{ opacity: 0.5, transition: "transform 0.15s", transform: open ? "rotate(180deg)" : "none" }}>
                <path d="M2 4l4 4 4-4" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"/>
              </svg>
            </button>
            {open && (
              <div className="account-dropdown" role="menu">
                <div className="dropdown-label">{user.email ?? user.name}</div>
                <Link to="/account" className="dropdown-item" onClick={() => setOpen(false)} role="menuitem">Account</Link>
                {user.role === "admin" && <>
                  <Link to="/admin" className="dropdown-item" onClick={() => setOpen(false)} role="menuitem">Admin dashboard</Link>
                  <div className="dropdown-divider" />
                </>}
                <button className="dropdown-item dropdown-item--danger" onClick={handleLogout} role="menuitem">Logout</button>
              </div>
            )}
          </div>
        ) : (
          <>
            <Link to="/login">Login</Link>
            <Link to="/register">Register</Link>
          </>
        )}
      </div>
    </nav>
  );
}
