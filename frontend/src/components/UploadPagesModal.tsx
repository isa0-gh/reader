import { useState, ChangeEvent } from "react";
import { PiArrowUp, PiArrowDown, PiX } from "react-icons/pi";
import Modal from "../Modal";
import { api, uploadFile } from "../api";

interface FileEntry { file: File; pageNumber: number; preview: string; }

export default function UploadPagesModal({ chapterId, onClose, onDone }: { chapterId: number; onClose(): void; onDone(): void }) {
  const [entries, setEntries] = useState<FileEntry[]>([]);
  const [progress, setProgress] = useState("");
  const [error, setError] = useState("");

  function onFilesChange(e: ChangeEvent<HTMLInputElement>) {
    const files = Array.from(e.target.files ?? []).sort((a, b) => a.name.localeCompare(b.name));
    setEntries(files.map((file, i) => ({
      file,
      pageNumber: i + 1,
      preview: URL.createObjectURL(file),
    })));
  }

  function moveUp(i: number) {
    if (i === 0) return;
    setEntries(prev => {
      const next = [...prev];
      [next[i - 1], next[i]] = [next[i], next[i - 1]];
      return next.map((e, idx) => ({ ...e, pageNumber: idx + 1 }));
    });
  }

  function moveDown(i: number) {
    setEntries(prev => {
      if (i === prev.length - 1) return prev;
      const next = [...prev];
      [next[i], next[i + 1]] = [next[i + 1], next[i]];
      return next.map((e, idx) => ({ ...e, pageNumber: idx + 1 }));
    });
  }

  function setPage(i: number, val: string) {
    setEntries(prev => prev.map((e, idx) => idx === i ? { ...e, pageNumber: parseInt(val) || idx + 1 } : e));
  }

  function remove(i: number) {
    setEntries(prev => prev.filter((_, idx) => idx !== i).map((e, idx) => ({ ...e, pageNumber: idx + 1 })));
  }

  async function submit() {
    if (!entries.length) return;
    setError("");
    try {
      const pages: { key: string; bucket: string; page_number: number }[] = [];
      for (let i = 0; i < entries.length; i++) {
        setProgress(`Uploading ${i + 1} / ${entries.length}…`);
        const { key, bucket } = await uploadFile(entries[i].file, `chapters/${chapterId}`);
        pages.push({ key, bucket, page_number: entries[i].pageNumber });
      }
      setProgress("Saving…");
      await api.uploadChapterPages(chapterId, pages);
      onDone();
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : "Upload failed");
      setProgress("");
    }
  }

  return (
    <Modal title="Upload Pages" onClose={onClose}>
      <div className="modal-form">
        <div className="form-group">
          <label>Select images</label>
          <input type="file" accept="image/*" multiple onChange={onFilesChange} />
        </div>

        {entries.length > 0 && (
          <div style={{ display: "flex", flexDirection: "column", gap: "0.5rem", marginBottom: "1rem", maxHeight: "340px", overflowY: "auto" }}>
            {entries.map((e, i) => (
              <div key={e.preview} className="page-row">
                <img src={e.preview} className="page-row-thumb" />
                <span className="page-row-name">{e.file.name}</span>
                <input
                  type="number" min={1} value={e.pageNumber}
                  onChange={ev => setPage(i, ev.target.value)}
                  className="page-row-num"
                />
                <button onClick={() => moveUp(i)} disabled={i === 0} className="page-row-btn" aria-label="Move up"><PiArrowUp /></button>
                <button onClick={() => moveDown(i)} disabled={i === entries.length - 1} className="page-row-btn" aria-label="Move down"><PiArrowDown /></button>
                <button onClick={() => remove(i)} className="page-row-btn page-row-btn--danger" aria-label="Remove"><PiX /></button>
              </div>
            ))}
          </div>
        )}

        {progress && <div className="muted" style={{ marginBottom: "0.5rem" }}>{progress}</div>}
        {error && <div className="error">{error}</div>}
        <button className="btn" onClick={submit} disabled={!entries.length || !!progress}>
          {progress || `Upload ${entries.length > 0 ? `${entries.length} page${entries.length > 1 ? "s" : ""}` : ""}`}
        </button>
      </div>
    </Modal>
  );
}
