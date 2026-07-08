import { useEffect, useState } from "react";
import { PiTrash, PiCaretLeft, PiCaretRight } from "react-icons/pi";
import { api, User } from "../api";
import { useAuth } from "../AuthContext";
import { useToast } from "../ToastContext";
import { useConfirm } from "../ConfirmContext";

const PAGE_SIZE = 50;
const ROLES = ["reader", "uploader", "moderator", "admin"];

export default function AdminUsersPage() {
  const { user: me } = useAuth();
  const toast = useToast();
  const confirm = useConfirm();
  const [users, setUsers] = useState<User[]>([]);
  const [cursor, setCursor] = useState<number | undefined>();
  const [history, setHistory] = useState<number[]>([]);
  const [error, setError] = useState("");

  async function load(after?: number) {
    try {
      setUsers(await api.listUsers({ limit: PAGE_SIZE, after }));
    } catch (e: any) { setError(e.message); }
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
    const p = h.pop();
    setHistory(h);
    setCursor(p);
    load(p === 0 ? undefined : p);
  }

  async function handleRoleChange(id: number, role: string) {
    try {
      await api.updateUserRole(id, role);
      setUsers((u) => u.map((x) => x.id === id ? { ...x, role } : x));
      toast.show("Role updated.", "success");
    } catch (e: any) { toast.show(e.message, "error"); }
  }

  async function handleDelete(id: number) {
    const ok = await confirm("Delete this user? This can't be undone.", { danger: true });
    if (!ok) return;
    try {
      await api.deleteUser(id);
      setUsers((u) => u.filter((x) => x.id !== id));
      toast.show("User deleted.", "success");
    } catch (e: any) { toast.show(e.message, "error"); }
  }

  if (error) return <p className="error">{error}</p>;

  return (
    <>
      <h2>Users</h2>
      <div className="table-scroll">
        <table className="admin-table">
          <thead>
            <tr>{["ID", "Name", "Email", "Role", "Created", ""].map((h, i) => <th key={i}>{h}</th>)}</tr>
          </thead>
          <tbody>
            {users.map((u) => (
              <tr key={u.id}>
                <td>{u.id}</td>
                <td>{u.name}</td>
                <td>{u.email}</td>
                <td>
                  <select className="role-select" value={u.role} disabled={u.id === me?.id}
                    onChange={(e) => handleRoleChange(u.id, e.target.value)}>
                    {ROLES.map((r) => <option key={r} value={r}>{r}</option>)}
                  </select>
                </td>
                <td>{new Date(u.created_at).toLocaleDateString()}</td>
                <td>
                  <button className="btn-outline" style={{ color: "var(--danger)" }}
                    disabled={u.id === me?.id} onClick={() => handleDelete(u.id)}>
                    <PiTrash /> Delete
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      <div className="admin-pagination">
        <button className="btn-outline" onClick={prev} disabled={!history.length}><PiCaretLeft /> Prev</button>
        <button className="btn-outline" onClick={next} disabled={users.length < PAGE_SIZE}>Next <PiCaretRight /></button>
      </div>
    </>
  );
}
