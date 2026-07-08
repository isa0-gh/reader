import { useState, FormEvent } from "react";
import { Link, useNavigate } from "react-router-dom";
import { api } from "../api";
import { useConfig } from "../ConfigContext";

export default function Register() {
  const nav = useNavigate();
  const { register_disabled, maintenance } = useConfig();
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");

  async function submit(e: FormEvent) {
    e.preventDefault();
    setError("");
    try {
      await api.register(email, password, name);
      nav("/login");
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : "Registration failed.");
    }
  }

  const disabled = register_disabled || maintenance;

  return (
    <div className="auth-wrap">
      <h1>Register</h1>
      {disabled ? (
        <p className="error">
          {maintenance ? "The site is under maintenance. Try again later." : "Registration is currently disabled."}
        </p>
      ) : (
        <form onSubmit={submit}>
          <div className="form-group">
            <label>Name</label>
            <input value={name} onChange={(e) => setName(e.target.value)} required />
          </div>
          <div className="form-group">
            <label>Email</label>
            <input type="email" value={email} onChange={(e) => setEmail(e.target.value)} required />
          </div>
          <div className="form-group">
            <label>Password</label>
            <input type="password" value={password} onChange={(e) => setPassword(e.target.value)} required minLength={8} />
          </div>
          {error && <div className="error">{error}</div>}
          <button className="btn" type="submit">Create account</button>
        </form>
      )}
      <p className="auth-switch">Have an account? <Link to="/login">Login</Link></p>
    </div>
  );
}
