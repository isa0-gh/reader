import { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { PiHeart, PiHeartFill, PiPencilSimple, PiTrash, PiPlus } from "react-icons/pi";
import { api, Series, Chapter } from "../api";
import { useAuth } from "../AuthContext";
import { useConfig } from "../ConfigContext";
import { useConfirm } from "../ConfirmContext";
import { useFavorites } from "../FavoritesContext";
import CreateChapterModal from "../components/CreateChapterModal";
import EditSeriesModal from "../components/EditSeriesModal";

const CAN_CREATE = ["uploader", "moderator", "admin"];
const CAN_DELETE = ["moderator", "admin"];

export default function SeriesPage() {
  const { id } = useParams<{ id: string }>();
  const nav = useNavigate();
  const { user } = useAuth();
  const { cdn_url } = useConfig();
  const confirm = useConfirm();
  const { isFavorite, toggleFavorite } = useFavorites();
  const [series, setSeries] = useState<Series | null>(null);
  const [error, setError] = useState("");
  const [showCreate, setShowCreate] = useState(false);
  const [showEdit, setShowEdit] = useState(false);

  useEffect(() => {
    api.getSeries(Number(id)).then(setSeries).catch(() => setError("Series not found."));
  }, [id]);

  if (error) return <div className="container muted">{error}</div>;
  if (!series) return <div className="container muted">Loading…</div>;

  const chapters = [...(series.chapters ?? [])].sort((a, b) => b.number - a.number);
  const canCreate = user && CAN_CREATE.includes(user.role);
  const canDelete = user && CAN_DELETE.includes(user.role);

  async function handleDeleteSeries() {
    const ok = await confirm("Delete this series? All its chapters go with it.", { danger: true });
    if (!ok) return;
    await api.deleteSeries(series!.id);
    nav("/");
  }

  async function handleDeleteChapter(e: React.MouseEvent, chId: number) {
    e.stopPropagation();
    const ok = await confirm("Delete this chapter?", { danger: true });
    if (!ok) return;
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
        {series.cover_image
          ? <img src={`${cdn_url}/${series.cover_image.key}`} alt={series.title} />
          : <div className="cover-placeholder" aria-hidden="true" />}
        <div className="series-meta">
          <div style={{ display: "flex", alignItems: "center", gap: "0.6rem" }}>
            <h1>{series.title}</h1>
            {user && (
              <button
                className={"favorite-toggle" + (isFavorite(series.id) ? " favorite-toggle--active" : "")}
                onClick={() => toggleFavorite(series.id, series)}
                aria-pressed={isFavorite(series.id)}
                aria-label={isFavorite(series.id) ? "Remove from favorites" : "Add to favorites"}
                title={isFavorite(series.id) ? "Remove from favorites" : "Add to favorites"}
              >
                {isFavorite(series.id) ? <PiHeartFill /> : <PiHeart />}
              </button>
            )}
          </div>
          <div className="meta-row">
            {series.author && <span>Author: {series.author}</span>}
            {series.artist && series.artist !== series.author && <span> · Art: {series.artist}</span>}
            <span> · {series.status}</span>
          </div>
          {series.description && <p>{series.description}</p>}
          {canDelete && (
            <div style={{ display: "flex", gap: "0.5rem" }}>
              <button className="btn-outline" onClick={() => setShowEdit(true)}><PiPencilSimple /> Edit Series</button>
              <button className="btn-outline" style={{ color: "var(--danger)" }} onClick={handleDeleteSeries}><PiTrash /> Delete Series</button>
            </div>
          )}
        </div>
      </div>

      <div className="chapter-list">
        <div className="chapter-list-header">
          <h2>Chapters ({chapters.length})</h2>
          {canCreate && <button className="btn-outline" onClick={() => setShowCreate(true)}><PiPlus /> New Chapter</button>}
        </div>
        {chapters.map((ch) => (
          <div key={ch.id} className="chapter-item" role="button" tabIndex={0}
            onClick={() => nav(`/series/${series.id}/${ch.id}`)}
            onKeyDown={(e) => e.key === "Enter" && nav(`/series/${series.id}/${ch.id}`)}>
            <span className="ch-num">Ch. {ch.number}{ch.title ? ` — ${ch.title}` : ""}</span>
            {canDelete && <button className="btn-outline" style={{ color: "var(--danger)" }} onClick={(e) => handleDeleteChapter(e, ch.id)}><PiTrash /> Delete</button>}
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

      {showEdit && (
        <EditSeriesModal
          series={series}
          onClose={() => setShowEdit(false)}
          onSave={(updated) => { setSeries(updated); setShowEdit(false); }}
        />
      )}
    </div>
  );
}
