import { Link, useNavigate } from "react-router-dom";
import { useAuth } from "./AuthContext";

export default function Nav() {
  const { user, logout } = useAuth();
  const nav = useNavigate();

  function handleLogout() {
    logout();
    nav("/");
  }

  return (
    <nav>
      <Link to="/" className="logo">Reader</Link>
      <div className="nav-links">
        {user ? (
          <>
            <span>{user.name}</span>
            <a href="#" onClick={handleLogout}>Logout</a>
          </>
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
