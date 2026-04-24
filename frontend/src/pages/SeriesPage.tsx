import { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { api, Series, Chapter } from "../api";
import { useAuth } from "../AuthContext";
import CreateChapterModal from "../components/CreateChapterModal";

const CAN_CREATE = ["uploader", "moderator", "admin"];
const CAN_DELETE = ["moderator", "admin"];

export default function SeriesPage() {
  const { id } = useParams<{ id: string }>();
  const nav = useNavigate();
  const { user } = useAuth();
  const [series, setSeries] = useState<Series | null>(null);
  const [error, setError] = useState("");
  const [showCreate, setShowCreate] = useState(false);

  useEffect(() => {
    api.getSeries(Number(id)).then(setSeries).catch(() => setError("Series not found."));
  }, [id]);

  if (error) return <div className="container muted">{error}</div>;
  if (!series) return <div className="container muted">Loading…</div>;

  const chapters = [...(series.chapters ?? [])].sort((a, b) => b.number - a.number);
  const canCreate = user && CAN_CREATE.includes(user.role);
  const canDelete = user && CAN_DELETE.includes(user.role);

  async function handleDeleteSeries() {
    if (!confirm("Delete this series?")) return;
    await api.deleteSeries(series!.id);
    nav("/");
  }

  async function handleDeleteChapter(e: React.MouseEvent, chId: number) {
    e.stopPropagation();
    if (!confirm("Delete this chapter?")) return;
    await api.deleteChapter(chId);
    setSeries((prev) => prev ? { ...prev, chapters: prev.chapters.filter((c) => c.id !== chId) } : prev);
  }

  function onChapterCreated(c: Chapter) {
    setSeries((prev) => prev ? { ...prev, chapters: [c, ...(prev.chapters ?? [])] } : prev);
    setShowCreate(false);
    nav(`/series/${series!.id}/${c.id}`);
  }

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
          {canDelete && <button className="btn-outline" style={{ color: "red" }} onClick={handleDeleteSeries}>Delete Series</button>}
        </div>
      </div>

      <div className="chapter-list">
        <div className="chapter-list-header">
          <h2>Chapters ({chapters.length})</h2>
          {canCreate && <button className="btn-outline" onClick={() => setShowCreate(true)}>+ New Chapter</button>}
        </div>
        {chapters.map((ch) => (
          <div key={ch.id} className="chapter-item" onClick={() => nav(`/series/${series.id}/${ch.id}`)}>
            <span className="ch-num">Ch. {ch.number}{ch.title ? ` — ${ch.title}` : ""}</span>
            {canDelete && <button className="btn-outline" style={{ color: "red" }} onClick={(e) => handleDeleteChapter(e, ch.id)}>Delete</button>}
          </div>
        ))}
      </div>

      {showCreate && (
        <CreateChapterModal
          seriesId={series.id}
          onClose={() => setShowCreate(false)}
          onCreate={onChapterCreated}
        />
      )}
    </div>
  );
}
