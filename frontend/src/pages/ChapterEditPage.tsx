import { useEffect, useState, ChangeEvent } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { PiArrowLeft, PiArrowUp, PiArrowDown, PiX } from "react-icons/pi";
import { api, Chapter, Page, uploadFile } from "../api";
import { useConfig } from "../ConfigContext";
import { useConfirm } from "../ConfirmContext";

interface FileEntry { file: File; pageNumber: number; preview: string; }

export default function ChapterEditPage() {
  const { id: seriesId, chapterId } = useParams<{ id: string; chapterId: string }>();
  const nav = useNavigate();
  const { cdn_url } = useConfig();
  const confirm = useConfirm();
  const [chapter, setChapter] = useState<Chapter | null>(null);
  const [initialPages, setInitialPages] = useState<Page[]>([]);
  const [newFiles, setNewFiles] = useState<FileEntry[]>([]);
  const [progress, setProgress] = useState("");
  const [error, setError] = useState("");
  const [chNumber, setChNumber] = useState("");
  const [chTitle, setChTitle] = useState("");
  const [savingInfo, setSavingInfo] = useState(false);
  const [savingOrder, setSavingOrder] = useState(false);

  function load() {
    api.getChapter(Number(chapterId)).then(c => {
      setChapter(c);
      setInitialPages(c.pages ?? []);
      setChNumber(String(c.number));
      setChTitle(c.title ?? "");
    }).catch(() => {});
  }

  useEffect(load, [chapterId]);

  async function saveInfo() {
    setError("");
    setSavingInfo(true);
    try {
      const updated = await api.updateChapter(Number(chapterId), { number: parseFloat(chNumber) || 0, title: chTitle });
      setChapter(prev => prev ? { ...prev, number: updated.number, title: updated.title } : prev);
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : "Failed to update chapter");
    } finally {
      setSavingInfo(false);
    }
  }

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
      return;
    }
    swapExisting(i, i - 1);
  }

  function moveDown(i: number, isNew: boolean) {
    if (isNew) {
      setNewFiles(prev => {
        if (i === prev.length - 1) return prev;
        const next = [...prev];
        [next[i], next[i + 1]] = [next[i + 1], next[i]];
        return next;
      });
      return;
    }
    swapExisting(i, i + 1);
  }

  // Existing pages are ordered by page_number, not array position, so
  // swapping must operate on the sorted view and write page_number back.
  function swapExisting(i: number, j: number) {
    setChapter(prev => {
      if (!prev) return prev;
      const sorted = [...prev.pages].sort((a, b) => a.page_number - b.page_number);
      if (i < 0 || j < 0 || i >= sorted.length || j >= sorted.length) return prev;
      const a = sorted[i], b = sorted[j];
      const swappedNum = a.page_number;
      return {
        ...prev,
        pages: prev.pages.map(p => {
          if (p.id === a.id) return { ...p, page_number: b.page_number };
          if (p.id === b.id) return { ...p, page_number: swappedNum };
          return p;
        }),
      };
    });
  }

  function setPageNum(i: number, val: string, isNew: boolean) {
    const num = parseInt(val) || 1;
    if (isNew) {
      setNewFiles(prev => prev.map((e, idx) => idx === i ? { ...e, pageNumber: num } : e));
      return;
    }
    setChapter(prev => {
      if (!prev) return prev;
      const sorted = [...prev.pages].sort((a, b) => a.page_number - b.page_number);
      const target = sorted[i];
      if (!target) return prev;
      return { ...prev, pages: prev.pages.map(p => p.id === target.id ? { ...p, page_number: num } : p) };
    });
  }

  function removeNew(i: number) {
    setNewFiles(prev => prev.filter((_, idx) => idx !== i));
  }

  async function deleteExisting(page: Page) {
    const ok = await confirm(`Delete page ${page.page_number}?`, { danger: true });
    if (!ok) return;
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

  async function saveOrder() {
    if (!chapter) return;
    setError("");
    setSavingOrder(true);
    try {
      await api.reorderChapterPages(Number(chapterId), chapter.pages.map(p => ({ id: p.id, page_number: p.page_number })));
      load();
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : "Failed to save page order");
    } finally {
      setSavingOrder(false);
    }
  }

  if (!chapter) return <div className="container muted">Loading…</div>;

  const pages = [...(chapter.pages ?? [])].sort((a, b) => a.page_number - b.page_number);
  const pagesDirty = JSON.stringify(pages.map(p => [p.id, p.page_number]))
    !== JSON.stringify([...initialPages].sort((a, b) => a.page_number - b.page_number).map(p => [p.id, p.page_number]));

  return (
    <div className="container">
      <div className="page-header">
        <h1 className="page-title">Edit Ch. {chapter.number}</h1>
        <div style={{ display: "flex", gap: "0.5rem" }}>
          <button className="btn-outline" onClick={() => nav(`/series/${seriesId}/${chapterId}`)}>View</button>
          <button className="btn-outline" onClick={() => nav(`/series/${seriesId}`)}><PiArrowLeft /> Back</button>
        </div>
      </div>

      <div className="form-row" style={{ alignItems: "flex-end", marginBottom: "2rem" }}>
        <div className="form-group" style={{ flex: "0 0 90px" }}>
          <label>Number</label>
          <input type="number" step="any" value={chNumber} onChange={e => setChNumber(e.target.value)} />
        </div>
        <div className="form-group" style={{ flex: 1 }}>
          <label>Title</label>
          <input value={chTitle} onChange={e => setChTitle(e.target.value)} />
        </div>
        <button className="btn-outline" onClick={saveInfo} disabled={savingInfo} style={{ marginBottom: "1rem" }}>
          {savingInfo ? "Saving…" : "Save Info"}
        </button>
      </div>

      <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", marginBottom: "0.75rem" }}>
        <h2 style={{ fontSize: "1rem", margin: 0 }}>Existing Pages ({pages.length})</h2>
        {pagesDirty && <button className="btn-outline" onClick={saveOrder} disabled={savingOrder}>{savingOrder ? "Saving…" : "Save Order"}</button>}
      </div>
      <div style={{ display: "flex", flexDirection: "column", gap: "0.5rem", marginBottom: "2rem" }}>
        {pages.length === 0
          ? <div className="muted">No pages yet.</div>
          : pages.map((p, i) => (
          <div key={p.id} className="page-row">
            <img src={`${cdn_url}/${p.key}`} className="page-row-thumb" />
            <span className="page-row-name">{p.key.split("/").pop()}</span>
            <input type="number" min={1} value={p.page_number} onChange={e => setPageNum(i, e.target.value, false)} className="page-row-num" />
            <button onClick={() => moveUp(i, false)} disabled={i === 0} className="page-row-btn" aria-label="Move up"><PiArrowUp /></button>
            <button onClick={() => moveDown(i, false)} disabled={i === pages.length - 1} className="page-row-btn" aria-label="Move down"><PiArrowDown /></button>
            <button onClick={() => deleteExisting(p)} className="page-row-btn page-row-btn--danger" aria-label="Delete"><PiX /></button>
          </div>
        ))}
      </div>

      <h2 style={{ fontSize: "1rem", marginBottom: "0.75rem" }}>Add New Pages</h2>
      <input type="file" accept="image/*" multiple onChange={onFilesChange} style={{ marginBottom: "1rem" }} />
      {newFiles.length > 0 && (
        <div style={{ display: "flex", flexDirection: "column", gap: "0.5rem", marginBottom: "1rem" }}>
          {newFiles.map((e, i) => (
            <div key={e.preview} className="page-row">
              <img src={e.preview} className="page-row-thumb" />
              <span className="page-row-name">{e.file.name}</span>
              <input type="number" min={1} value={e.pageNumber} onChange={ev => setPageNum(i, ev.target.value, true)} className="page-row-num" />
              <button onClick={() => moveUp(i, true)} disabled={i === 0} className="page-row-btn" aria-label="Move up"><PiArrowUp /></button>
              <button onClick={() => moveDown(i, true)} disabled={i === newFiles.length - 1} className="page-row-btn" aria-label="Move down"><PiArrowDown /></button>
              <button onClick={() => removeNew(i)} className="page-row-btn page-row-btn--danger" aria-label="Remove"><PiX /></button>
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
