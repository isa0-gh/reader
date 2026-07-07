import { useState, FormEvent } from "react";
import { useNavigate } from "react-router-dom";
import { api } from "../api";
import { useAuth } from "../AuthContext";

export default function AccountPage() {
  const { user, logout } = useAuth();
  const nav = useNavigate();
  const [currentPassword, setCurrentPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [error, setError] = useState("");
  const [done, setDone] = useState(false);

  async function submit(e: FormEvent) {
    e.preventDefault();
    setError("");
    if (newPassword !== confirmPassword) {
      setError("New passwords don't match.");
      return;
    }
    try {
      await api.changePassword(currentPassword, newPassword);
      // Changing the password rotates the server-side session id, so the
      // token this tab is holding is now invalid — send the user back to
      // login instead of leaving them in a broken logged-in state.
      setDone(true);
      setTimeout(() => {
        logout();
        nav("/login");
      }, 1500);
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : "Failed to change password.");
    }
  }

  if (!user) return <div className="container muted">Not logged in.</div>;

  return (
    <div className="auth-wrap">
      <h1>Account</h1>
      <p className="muted" style={{ marginBottom: "1.5rem" }}>{user.email}</p>

      {done ? (
        <p className="muted">Password changed. Redirecting to login…</p>
      ) : (
        <form onSubmit={submit}>
          <div className="form-group">
            <label>Current password</label>
            <input type="password" value={currentPassword} onChange={(e) => setCurrentPassword(e.target.value)} required />
          </div>
          <div className="form-group">
            <label>New password</label>
            <input type="password" value={newPassword} onChange={(e) => setNewPassword(e.target.value)} required />
          </div>
          <div className="form-group">
            <label>Confirm new password</label>
            <input type="password" value={confirmPassword} onChange={(e) => setConfirmPassword(e.target.value)} required />
          </div>
          {error && <div className="error">{error}</div>}
          <button className="btn" type="submit">Change Password</button>
        </form>
      )}
    </div>
  );
}
