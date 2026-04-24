/// <reference types="vite/client" />

import { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { api, Chapter, Series } from "../api";
import { useAuth } from "../AuthContext";
import { useConfig } from "../ConfigContext";

const CAN_CREATE = ["uploader", "moderator", "admin"];

export default function ReaderPage() {
  const { id: seriesId, chapterId } = useParams<{ id: string; chapterId: string }>();
  const nav = useNavigate();
  const { user } = useAuth();
  const { cdn_url } = useConfig();
  const [chapter, setChapter] = useState<Chapter | null>(null);
  const [series, setSeries] = useState<Series | null>(null);
  const [error, setError] = useState("");

  useEffect(() => {
    api.getChapter(Number(chapterId)).then(setChapter).catch(() => setError("Chapter not found."));
  }, [chapterId]);

  useEffect(() => {
    api.getSeries(Number(seriesId)).then(setSeries).catch(() => {});
  }, [seriesId]);

  if (error) return <div className="container muted">{error}</div>;
  if (!chapter) return <div className="container muted">Loading…</div>;

  const pages = [...(chapter.pages ?? [])].sort((a, b) => a.page_number - b.page_number);
  const canEdit = user && CAN_CREATE.includes(user.role);

  const siblings = [...(series?.chapters ?? [])].sort((a, b) => a.number - b.number);
  const idx = siblings.findIndex(c => c.id === chapter.id);
  const prevId = idx > 0 ? siblings[idx - 1].id : null;
  const nextId = idx !== -1 && idx < siblings.length - 1 ? siblings[idx + 1].id : null;

  return (
    <>
      <div className="reader-nav">
        <button onClick={() => nav(`/series/${seriesId}`)}>← Back</button>
        <span>Ch. {chapter.number}{chapter.title ? ` — ${chapter.title}` : ""}</span>
        <div style={{ display: "flex", gap: "0.5rem" }}>
          {canEdit && <button onClick={() => nav(`/series/${seriesId}/${chapterId}/edit`)}>Edit</button>}
          <button disabled={!prevId} onClick={() => prevId && nav(`/series/${seriesId}/${prevId}`)}>Prev</button>
          <button disabled={!nextId} onClick={() => nextId && nav(`/series/${seriesId}/${nextId}`)}>Next</button>
        </div>
      </div>

      <div className="reader-pages">
        {pages.map((p) => (
          <img key={p.id} src={`${cdn_url}/${p.key}`} alt={`Page ${p.page_number}`} loading="lazy" />
        ))}
      </div>
    </>
  );
}
