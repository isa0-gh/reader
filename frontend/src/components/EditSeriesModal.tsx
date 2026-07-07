import { useState, FormEvent } from "react";
import Modal from "../Modal";
import { api, uploadFile, Series } from "../api";

export default function EditSeriesModal({ series, onClose, onSave }: { series: Series; onClose(): void; onSave(s: Series): void }) {
  const [title, setTitle] = useState(series.title);
  const [slug, setSlug] = useState(series.slug);
  const [description, setDescription] = useState(series.description);
  const [author, setAuthor] = useState(series.author);
  const [artist, setArtist] = useState(series.artist);
  const [status, setStatus] = useState(series.status);
  const [cover, setCover] = useState<File | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  async function submit(e: FormEvent) {
    e.preventDefault();
    setLoading(true);
    setError("");
    try {
      let cover_key = "";
      let cover_bucket = "";
      if (cover) {
        const res = await uploadFile(cover, "covers");
        cover_key = res.key;
        cover_bucket = res.bucket;
      }
      const updated = await api.updateSeries(series.id, {
        title,
        slug,
        description,
        author,
        artist,
        status,
        cover_key,
        cover_bucket,
      });
      onSave(updated);
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : "Failed to update series");
    } finally {
      setLoading(false);
    }
  }

  return (
    <Modal title="Edit Series" onClose={onClose}>
      <form onSubmit={submit} className="modal-form">
        <div className="form-group">
          <label>Title</label>
          <input value={title} onChange={(e) => setTitle(e.target.value)} required />
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
          <label>Replace Cover Image</label>
          <input type="file" accept="image/*" onChange={(e) => setCover(e.target.files?.[0] ?? null)} />
        </div>
        {error && <div className="error">{error}</div>}
        <button className="btn" type="submit" disabled={loading}>{loading ? "Saving…" : "Save Changes"}</button>
      </form>
    </Modal>
  );
}
