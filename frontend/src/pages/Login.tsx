import { useState, FormEvent } from "react";
import { Link, useNavigate } from "react-router-dom";
import { api } from "../api";
import { useAuth } from "../AuthContext";
import { useConfig } from "../ConfigContext";

export default function Login() {
  const { login } = useAuth();
  const { login_disabled, maintenance } = useConfig();
  const nav = useNavigate();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");

  async function submit(e: FormEvent) {
    e.preventDefault();
    setError("");
    try {
      const { token, user } = await api.login(email, password);
      login(token, user);
      nav("/");
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : "Login failed.");
    }
  }

  const disabled = login_disabled || maintenance;

  return (
    <div className="auth-wrap">
      <h1>Login</h1>
      {disabled ? (
        <p className="error">
          {maintenance ? "The site is under maintenance. Try again later." : "Login is currently disabled."}
        </p>
      ) : (
        <form onSubmit={submit}>
          <div className="form-group">
            <label>Email</label>
            <input type="email" value={email} onChange={(e) => setEmail(e.target.value)} required />
          </div>
          <div className="form-group">
            <label>Password</label>
            <input type="password" value={password} onChange={(e) => setPassword(e.target.value)} required />
          </div>
          {error && <div className="error">{error}</div>}
          <button className="btn" type="submit">Login</button>
        </form>
      )}
      <p className="auth-switch">No account? <Link to="/register">Register</Link></p>
    </div>
  );
}
