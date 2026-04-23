/// <reference types="vite/client" />

import { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { api, Chapter } from "../api";

const CDN = import.meta.env.VITE_CDN_URL ?? "";

export default function ReaderPage() {
  const { id } = useParams<{ id: string }>();
  const nav = useNavigate();
  const [chapter, setChapter] = useState<Chapter | null>(null);
  const [error, setError] = useState("");

  useEffect(() => {
    api.getChapter(Number(id))
      .then(setChapter)
      .catch(() => setError("Chapter not found."));
  }, [id]);

  if (error) return <div className="container muted">{error}</div>;
  if (!chapter) return <div className="container muted">Loading…</div>;

  const pages = [...(chapter.pages ?? [])].sort((a, b) => a.page_number - b.page_number);

  return (
    <>
      <div className="reader-nav">
        <button onClick={() => nav(`/series/${chapter.series_id}`)}>← Back</button>
        <span>Ch. {chapter.number}{chapter.title ? ` — ${chapter.title}` : ""}</span>
        <div style={{ display: "flex", gap: "0.5rem" }}>
          <button onClick={() => nav(`/chapters/${Number(id) - 1}`)}>Prev</button>
          <button onClick={() => nav(`/chapters/${Number(id) + 1}`)}>Next</button>
        </div>
      </div>

      <div className="reader-pages">
        {pages.map((p) => (
          <img key={p.id} src={`${CDN}/${p.key}`} alt={`Page ${p.page_number}`} loading="lazy" />
        ))}
      </div>
    </>
  );
}
