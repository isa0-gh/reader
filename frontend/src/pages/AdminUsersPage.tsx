import { useEffect, useState } from "react";
import { api, User } from "../api";
import { useAuth } from "../AuthContext";

const PAGE_SIZE = 50;
const ROLES = ["reader", "uploader", "moderator", "admin"];

export default function AdminUsersPage() {
  const { user: me } = useAuth();
  const [users, setUsers] = useState<User[]>([]);
  const [cursor, setCursor] = useState<number | undefined>();
  const [history, setHistory] = useState<number[]>([]);
  const [error, setError] = useState("");

  async function load(after?: number) {
    try {
      const data = await api.listUsers({ limit: PAGE_SIZE, after });
      setUsers(data);
    } catch (e: any) {
      setError(e.message);
    }
  }

  useEffect(() => { load(); }, []);

  function next() {
    if (users.length < PAGE_SIZE) return;
    const last = users[users.length - 1].id;
    setHistory((h) => [...h, cursor ?? 0]);
    setCursor(last);
    load(last);
  }

  function prev() {
    const h = [...history];
    const prev = h.pop();
    setHistory(h);
    setCursor(prev);
    load(prev === 0 ? undefined : prev);
  }

  async function handleRoleChange(id: number, role: string) {
    try {
      await api.updateUserRole(id, role);
      setUsers((u) => u.map((x) => x.id === id ? { ...x, role } : x));
    } catch (e: any) {
      alert(e.message);
    }
  }

  async function handleDelete(id: number) {
    if (!confirm("Delete this user?")) return;
    try {
      await api.deleteUser(id);
      setUsers((u) => u.filter((x) => x.id !== id));
    } catch (e: any) {
      alert(e.message);
    }
  }

  if (error) return <p className="error">{error}</p>;

  return (
    <main style={{ padding: "1rem" }}>
      <h2>Users</h2>
      <table style={{ width: "100%", borderCollapse: "collapse" }}>
        <thead>
          <tr>
            {["ID", "Name", "Email", "Role", "Created", ""].map((h, i) => (
              <th key={i} style={{ textAlign: "left", padding: "0.5rem", borderBottom: "1px solid #333" }}>{h}</th>
            ))}
          </tr>
        </thead>
        <tbody>
          {users.map((u) => (
            <tr key={u.id}>
              <td style={{ padding: "0.5rem" }}>{u.id}</td>
              <td style={{ padding: "0.5rem" }}>{u.name}</td>
              <td style={{ padding: "0.5rem" }}>{u.email}</td>
              <td style={{ padding: "0.5rem" }}>
                <select
                  value={u.role}
                  disabled={u.id === me?.id}
                  onChange={(e) => handleRoleChange(u.id, e.target.value)}
                >
                  {ROLES.map((r) => <option key={r} value={r}>{r}</option>)}
                </select>
              </td>
              <td style={{ padding: "0.5rem" }}>{new Date(u.created_at).toLocaleDateString()}</td>
              <td style={{ padding: "0.5rem" }}>
                <button
                  disabled={u.id === me?.id}
                  onClick={() => handleDelete(u.id)}
                  style={{ color: "red" }}
                >
                  Delete
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
      <div style={{ marginTop: "1rem", display: "flex", gap: "0.5rem" }}>
        <button onClick={prev} disabled={!history.length}>← Prev</button>
        <button onClick={next} disabled={users.length < PAGE_SIZE}>Next →</button>
      </div>
    </main>
  );
}
