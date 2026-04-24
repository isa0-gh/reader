import { useEffect, useState } from "react";
import { api, S3Object } from "../api";

export default function AdminS3CleanPage() {
  const [objects, setObjects] = useState<S3Object[]>([]);
  const [result, setResult] = useState<{ deleted: string[]; failed: string[] } | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  async function load() {
    try {
      setObjects(await api.listOrphanedObjects());
    } catch (e: any) {
      setError(e.message);
    }
  }

  useEffect(() => { load(); }, []);

  async function handlePurge() {
    if (!confirm(`Delete ${objects.length} orphaned object(s) from S3?`)) return;
    setLoading(true);
    try {
      const res = await api.purgeOrphanedObjects();
      setResult(res);
      setObjects([]);
    } catch (e: any) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  }

  if (error) return <p className="error">{error}</p>;

  return (
    <main style={{ padding: "1rem" }}>
      <h2>S3 Cleanup</h2>
      <p>{objects.length} orphaned object(s) found (pages whose chapter has been deleted).</p>

      {objects.length > 0 && (
        <>
          <table style={{ width: "100%", borderCollapse: "collapse", marginBottom: "1rem" }}>
            <thead>
              <tr>
                {["ID", "Key", "Bucket", "Chapter ID", "Page #"].map((h) => (
                  <th key={h} style={{ textAlign: "left", padding: "0.5rem", borderBottom: "1px solid #333" }}>{h}</th>
                ))}
              </tr>
            </thead>
            <tbody>
              {objects.map((o) => (
                <tr key={o.id}>
                  <td style={{ padding: "0.5rem" }}>{o.id}</td>
                  <td style={{ padding: "0.5rem", fontFamily: "monospace", fontSize: "0.85em" }}>{o.key}</td>
                  <td style={{ padding: "0.5rem" }}>{o.bucket}</td>
                  <td style={{ padding: "0.5rem" }}>{o.chapter_id}</td>
                  <td style={{ padding: "0.5rem" }}>{o.page_number}</td>
                </tr>
              ))}
            </tbody>
          </table>
          <button onClick={handlePurge} disabled={loading} style={{ color: "red" }}>
            {loading ? "Purging…" : `Purge ${objects.length} object(s)`}
          </button>
        </>
      )}

      {result && (
        <div style={{ marginTop: "1rem" }}>
          <p>✓ Deleted: {result.deleted.length}</p>
          {result.failed.length > 0 && <p style={{ color: "red" }}>✗ Failed: {result.failed.join(", ")}</p>}
        </div>
      )}
    </main>
  );
}
