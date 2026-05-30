import { useState, FormEvent } from "react";
import Modal from "../Modal";
import { api, Chapter } from "../api";

export default function CreateChapterModal({ seriesId, onClose, onCreate }: { seriesId: number; onClose(): void; onCreate(c: Chapter): void }) {
  const [number, setNumber] = useState("");
  const [title, setTitle] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  async function submit(e: FormEvent) {
    e.preventDefault();
    setLoading(true);
    setError("");
    try {
      const chapter = await api.createChapter({ series_id: seriesId, number: parseFloat(number), title });
      onCreate(chapter);
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : "Failed to create chapter");
    } finally {
      setLoading(false);
    }
  }

  return (
    <Modal title="New Chapter" onClose={onClose}>
      <form onSubmit={submit} className="modal-form">
        <div className="form-group">
          <label>Chapter Number</label>
          <input type="number" step="0.1" value={number} onChange={(e) => setNumber(e.target.value)} required />
        </div>
        <div className="form-group">
          <label>Title (optional)</label>
          <input value={title} onChange={(e) => setTitle(e.target.value)} />
        </div>
        {error && <div className="error">{error}</div>}
        <button className="btn" type="submit" disabled={loading}>{loading ? "Creating…" : "Create Chapter"}</button>
      </form>
    </Modal>
  );
}
