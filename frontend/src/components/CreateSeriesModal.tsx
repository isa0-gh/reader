import { useState, FormEvent } from "react";
import Modal from "../Modal";
import { api, uploadFile, Series } from "../api";

export default function CreateSeriesModal({ onClose, onCreate }: { onClose(): void; onCreate(s: Series): void }) {
  const [title, setTitle] = useState("");
  const [slug, setSlug] = useState("");
  const [description, setDescription] = useState("");
  const [author, setAuthor] = useState("");
  const [artist, setArtist] = useState("");
  const [status, setStatus] = useState("ongoing");
  const [cover, setCover] = useState<File | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  async function submit(e: FormEvent) {
    e.preventDefault();
    setLoading(true);
    setError("");
    try {
      let cover_image = "";
      if (cover) {
        const { public_url } = await uploadFile(cover, "covers");
        cover_image = public_url;
      }
      const series = await api.createSeries({ title, slug, description, author, artist, status, cover_image });
      onCreate(series);
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : "Failed to create series");
    } finally {
      setLoading(false);
    }
  }

  return (
    <Modal title="New Series" onClose={onClose}>
      <form onSubmit={submit} className="modal-form">
        <div className="form-group">
          <label>Title</label>
          <input value={title} onChange={(e) => { setTitle(e.target.value); setSlug(e.target.value.toLowerCase().replace(/\s+/g, "-").replace(/[^a-z0-9-]/g, "")); }} required />
        </div>
        <div className="form-group">
          <label>Slug</label>
          <input value={slug} onChange={(e) => setSlug(e.target.value)} required />
        </div>
        <div className="form-group">
          <label>Description</label>
          <textarea value={description} onChange={(e) => setDescription(e.target.value)} rows={3} />
        </div>
        <div className="form-row">
          <div className="form-group">
            <label>Author</label>
            <input value={author} onChange={(e) => setAuthor(e.target.value)} />
          </div>
          <div className="form-group">
            <label>Artist</label>
            <input value={artist} onChange={(e) => setArtist(e.target.value)} />
          </div>
        </div>
        <div className="form-group">
          <label>Status</label>
          <select value={status} onChange={(e) => setStatus(e.target.value)}>
            <option value="ongoing">Ongoing</option>
            <option value="completed">Completed</option>
            <option value="hiatus">Hiatus</option>
          </select>
        </div>
        <div className="form-group">
          <label>Cover Image</label>
          <input type="file" accept="image/*" onChange={(e) => setCover(e.target.files?.[0] ?? null)} />
        </div>
        {error && <div className="error">{error}</div>}
        <button className="btn" type="submit" disabled={loading}>{loading ? "Creating…" : "Create Series"}</button>
      </form>
    </Modal>
  );
}
