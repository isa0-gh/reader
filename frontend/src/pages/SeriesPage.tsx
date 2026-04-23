import { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { api, Series } from "../api";

export default function SeriesPage() {
  const { id } = useParams<{ id: string }>();
  const nav = useNavigate();
  const [series, setSeries] = useState<Series | null>(null);
  const [error, setError] = useState("");

  useEffect(() => {
    api.getSeries(Number(id))
      .then(setSeries)
      .catch(() => setError("Series not found."));
  }, [id]);

  if (error) return <div className="container muted">{error}</div>;
  if (!series) return <div className="container muted">Loading…</div>;

  const chapters = [...(series.chapters ?? [])].sort((a, b) => b.number - a.number);

  return (
    <div className="container">
      <div className="series-detail">
        <img src={series.cover_image || ""} alt={series.title} />
        <div className="series-meta">
          <h1>{series.title}</h1>
          <div className="meta-row">
            {series.author && <span>Author: {series.author}</span>}
            {series.artist && series.artist !== series.author && <span> · Art: {series.artist}</span>}
            <span> · {series.status}</span>
          </div>
          {series.description && <p>{series.description}</p>}
        </div>
      </div>

      <div className="chapter-list">
        <h2>Chapters ({chapters.length})</h2>
        {chapters.map((ch) => (
          <div key={ch.id} className="chapter-item" onClick={() => nav(`/chapters/${ch.id}`)}>
            <span className="ch-num">Ch. {ch.number}{ch.title ? ` — ${ch.title}` : ""}</span>
          </div>
        ))}
      </div>
    </div>
  );
}
