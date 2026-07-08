import { useEffect, useState } from "react";
import { PiTrash, PiCheckCircle, PiXCircle } from "react-icons/pi";
import { api, S3Object } from "../api";
import { useConfirm } from "../ConfirmContext";

export default function AdminS3CleanPage() {
  const confirm = useConfirm();
  const [objects, setObjects] = useState<S3Object[]>([]);
  const [result, setResult] = useState<{ deleted: string[]; failed: string[] } | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  async function load() {
    try { setObjects((await api.listOrphanedObjects()) ?? []); }
    catch (e: any) { setError(e.message); }
  }

  useEffect(() => { load(); }, []);

  async function handlePurge() {
    const ok = await confirm(`Delete ${objects.length} orphaned object(s) from S3? This can't be undone.`, { danger: true, confirmLabel: "Purge" });
    if (!ok) return;
    setLoading(true);
    try {
      const res = await api.purgeOrphanedObjects();
      setResult(res);
      setObjects([]);
    } catch (e: any) { setError(e.message); }
    finally { setLoading(false); }
  }

  if (error) return <p className="error">{error}</p>;

  return (
    <>
      <h2>S3 Cleanup</h2>
      <p className="admin-subtext">{objects.length} orphaned object(s) — pages whose chapter or series has been deleted.</p>

      {objects.length > 0 && (
        <>
          <div className="table-scroll">
            <table className="admin-table">
              <thead>
                <tr>{["ID", "Key", "Bucket", "Chapter ID", "Page #"].map((h) => <th key={h}>{h}</th>)}</tr>
              </thead>
              <tbody>
                {objects.map((o) => (
                  <tr key={o.id}>
                    <td>{o.id}</td>
                    <td className="mono">{o.key}</td>
                    <td>{o.bucket}</td>
                    <td>{o.chapter_id}</td>
                    <td>{o.page_number}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          <button className="btn-outline" onClick={handlePurge} disabled={loading}
            style={{ color: "var(--danger)", borderColor: "var(--danger)", width: "auto" }}>
            <PiTrash /> {loading ? "Purging…" : `Purge ${objects.length} object(s)`}
          </button>
        </>
      )}

      {result && (
        <div className="admin-result">
          <p style={{ display: "flex", alignItems: "center", gap: "0.4rem" }}>
            <PiCheckCircle style={{ color: "#2f9e44" }} /> Deleted from S3 and DB: {result.deleted?.length ?? 0}
          </p>
          {result.failed?.length > 0 && (
            <div className="fail-list">
              <p style={{ display: "flex", alignItems: "center", gap: "0.4rem" }}>
                <PiXCircle /> Failed ({result.failed.length}) — DB rows kept. Check server logs.
              </p>
              <ul>{result.failed.map((k) => <li key={k}>{k}</li>)}</ul>
            </div>
          )}
        </div>
      )}
    </>
  );
}
