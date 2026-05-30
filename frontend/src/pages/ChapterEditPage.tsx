import { useEffect, useState, ChangeEvent } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { api, Chapter, Page, uploadFile } from "../api";
import { useConfig } from "../ConfigContext";

interface FileEntry { file: File; pageNumber: number; preview: string; }

export default function ChapterEditPage() {
  const { id: seriesId, chapterId } = useParams<{ id: string; chapterId: string }>();
  const nav = useNavigate();
  const { cdn_url } = useConfig();
  const [chapter, setChapter] = useState<Chapter | null>(null);
  const [newFiles, setNewFiles] = useState<FileEntry[]>([]);
  const [progress, setProgress] = useState("");
  const [error, setError] = useState("");

  function load() {
    api.getChapter(Number(chapterId)).then(setChapter).catch(() => {});
  }

  useEffect(load, [chapterId]);

  function onFilesChange(e: ChangeEvent<HTMLInputElement>) {
    const files = Array.from(e.target.files ?? []).sort((a, b) => a.name.localeCompare(b.name));
    const maxPage = Math.max(0, ...(chapter?.pages ?? []).map(p => p.page_number));
    setNewFiles(prev => [
      ...prev,
      ...files.map((file, i) => ({
        file,
        pageNumber: maxPage + prev.length + i + 1,
        preview: URL.createObjectURL(file),
      }))
    ]);
  }

  function moveUp(i: number, isNew: boolean) {
    if (isNew) {
      if (i === 0) return;
      setNewFiles(prev => {
        const next = [...prev];
        [next[i - 1], next[i]] = [next[i], next[i - 1]];
        return next;
      });
    } else {
      if (!chapter || i === 0) return;
      setChapter(prev => {
        if (!prev) return prev;
        const pages = [...prev.pages];
        [pages[i - 1], pages[i]] = [pages[i], pages[i - 1]];
        return { ...prev, pages };
      });
    }
  }

  function moveDown(i: number, isNew: boolean) {
    if (isNew) {
      setNewFiles(prev => {
        if (i === prev.length - 1) return prev;
        const next = [...prev];
        [next[i], next[i + 1]] = [next[i + 1], next[i]];
        return next;
      });
    } else {
      setChapter(prev => {
        if (!prev || i === prev.pages.length - 1) return prev;
        const pages = [...prev.pages];
        [pages[i], pages[i + 1]] = [pages[i + 1], pages[i]];
        return { ...prev, pages };
      });
    }
  }

  function setPageNum(i: number, val: string, isNew: boolean) {
    const num = parseInt(val) || 1;
    if (isNew) {
      setNewFiles(prev => prev.map((e, idx) => idx === i ? { ...e, pageNumber: num } : e));
    } else {
      setChapter(prev => {
        if (!prev) return prev;
        return { ...prev, pages: prev.pages.map((p, idx) => idx === i ? { ...p, page_number: num } : p) };
      });
    }
  }

  function removeNew(i: number) {
    setNewFiles(prev => prev.filter((_, idx) => idx !== i));
  }

  async function deleteExisting(page: Page) {
    if (!confirm(`Delete page ${page.page_number}?`)) return;
    try {
      await api.deleteChapterPage(Number(chapterId), page.id);
      setChapter(prev => prev ? { ...prev, pages: prev.pages.filter(p => p.id !== page.id) } : prev);
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : "Delete failed");
    }
  }

  async function save() {
    if (!newFiles.length) return;
    setError("");
    try {
      const pages: { key: string; bucket: string; page_number: number }[] = [];
      for (let i = 0; i < newFiles.length; i++) {
        setProgress(`Uploading ${i + 1} / ${newFiles.length}…`);
        const { key, bucket } = await uploadFile(newFiles[i].file, `chapters/${chapterId}`);
        pages.push({ key, bucket, page_number: newFiles[i].pageNumber });
      }
      setProgress("Saving…");
      await api.uploadChapterPages(Number(chapterId), pages);
      setNewFiles([]);
      setProgress("");
      load();
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : "Upload failed");
      setProgress("");
    }
  }

  if (!chapter) return <div className="container muted">Loading…</div>;

  const pages = [...(chapter.pages ?? [])].sort((a, b) => a.page_number - b.page_number);

  return (
    <div className="container">
      <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", marginBottom: "1.5rem" }}>
        <h1 style={{ fontSize: "1.3rem", margin: 0 }}>Edit Ch. {chapter.number}</h1>
        <div style={{ display: "flex", gap: "0.5rem" }}>
          <button className="btn-outline" onClick={() => nav(`/series/${seriesId}/${chapterId}`)}>View</button>
          <button className="btn-outline" onClick={() => nav(`/series/${seriesId}`)}>← Back</button>
        </div>
      </div>

      <h2 style={{ fontSize: "1rem", marginBottom: "0.75rem" }}>Existing Pages ({pages.length})</h2>
      <div style={{ display: "flex", flexDirection: "column", gap: "0.5rem", marginBottom: "2rem" }}>
        {pages.length === 0
          ? <div className="muted">No pages yet.</div>
          : pages.map((p, i) => (
          <div key={p.id} style={rowStyle}>
            <img src={`${cdn_url}/${p.key}`} style={thumbStyle} />
            <span style={{ flex: 1, fontSize: "0.85rem", overflow: "hidden", textOverflow: "ellipsis" }}>{p.key.split("/").pop()}</span>
            <input type="number" min={1} value={p.page_number} onChange={e => setPageNum(i, e.target.value, false)} style={inputStyle} />
            <button onClick={() => moveUp(i, false)} disabled={i === 0} style={btnStyle}>↑</button>
            <button onClick={() => moveDown(i, false)} disabled={i === pages.length - 1} style={btnStyle}>↓</button>
            <button onClick={() => deleteExisting(p)} style={{ ...btnStyle, color: "#c00" }}>✕</button>
          </div>
        ))}
      </div>

      <h2 style={{ fontSize: "1rem", marginBottom: "0.75rem" }}>Add New Pages</h2>
      <input type="file" accept="image/*" multiple onChange={onFilesChange} style={{ marginBottom: "1rem" }} />
      {newFiles.length > 0 && (
        <div style={{ display: "flex", flexDirection: "column", gap: "0.5rem", marginBottom: "1rem" }}>
          {newFiles.map((e, i) => (
            <div key={e.preview} style={rowStyle}>
              <img src={e.preview} style={thumbStyle} />
              <span style={{ flex: 1, fontSize: "0.85rem", overflow: "hidden", textOverflow: "ellipsis" }}>{e.file.name}</span>
              <input type="number" min={1} value={e.pageNumber} onChange={ev => setPageNum(i, ev.target.value, true)} style={inputStyle} />
              <button onClick={() => moveUp(i, true)} disabled={i === 0} style={btnStyle}>↑</button>
              <button onClick={() => moveDown(i, true)} disabled={i === newFiles.length - 1} style={btnStyle}>↓</button>
              <button onClick={() => removeNew(i)} style={{ ...btnStyle, color: "#c00" }}>✕</button>
            </div>
          ))}
        </div>
      )}

      {progress && <div className="muted" style={{ marginBottom: "0.5rem" }}>{progress}</div>}
      {error && <div className="error" style={{ marginBottom: "0.5rem" }}>{error}</div>}
      <button className="btn" onClick={save} disabled={!newFiles.length || !!progress}>
        {progress || `Upload ${newFiles.length} new page${newFiles.length !== 1 ? "s" : ""}`}
      </button>
    </div>
  );
}

const rowStyle: React.CSSProperties = {
  display: "flex", alignItems: "center", gap: "0.5rem",
  borderBottom: "1px solid var(--border)", paddingBottom: "0.5rem",
};

const thumbStyle: React.CSSProperties = {
  width: 40, height: 54, objectFit: "cover", borderRadius: 3, flexShrink: 0,
};

const inputStyle: React.CSSProperties = {
  width: 50, padding: "0.25rem 0.4rem", border: "1px solid var(--border)",
  borderRadius: 3, fontSize: "0.85rem",
};

const btnStyle: React.CSSProperties = {
  background: "none", border: "1px solid var(--border)", borderRadius: 3,
  cursor: "pointer", padding: "0.2rem 0.4rem", fontSize: "0.85rem",
};
