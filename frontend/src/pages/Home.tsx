import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../AuthContext";
import { useConfig } from "../ConfigContext";
import CreateSeriesModal from "../components/CreateSeriesModal";
import { api, Series } from "../api";

const CAN_CREATE = ["uploader", "moderator", "admin"];
const PAGE_SIZE = 24;

export default function Home() {
  const nav = useNavigate();
  const { user } = useAuth();
  const { cdn_url } = useConfig();
  const [seriesList, setSeriesList] = useState<Series[]>([]);
  const [total, setTotal] = useState(0);
  const [query, setQuery] = useState("");
  const [loading, setLoading] = useState(true);
  const [loadingMore, setLoadingMore] = useState(false);
  const [showCreate, setShowCreate] = useState(false);

  function load(q: string, offset: number, append: boolean) {
    (append ? setLoadingMore : setLoading)(true);
    api.listSeries({ limit: PAGE_SIZE, offset, q: q || undefined })
      .then((res) => {
        setSeriesList((prev) => (append ? [...prev, ...res.items] : res.items));
        setTotal(res.total);
      })
      .finally(() => (append ? setLoadingMore : setLoading)(false));
  }

  useEffect(() => {
    const handle = setTimeout(() => load(query, 0, false), 300);
    return () => clearTimeout(handle);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [query]);

  const canCreate = user && CAN_CREATE.includes(user.role);
  const canLoadMore = seriesList.length < total;

  return (
    <div className="container">
      <div className="page-header">
        <p className="page-title">Browse</p>
        {canCreate && <button className="btn-outline" onClick={() => setShowCreate(true)}>+ New Series</button>}
      </div>

      <input
        className="search-input"
        placeholder="Search by title…"
        value={query}
        onChange={(e) => setQuery(e.target.value)}
        style={{ marginBottom: "1rem", width: "100%", maxWidth: 320, padding: "0.5rem 0.75rem", border: "1px solid var(--border)", borderRadius: 4 }}
      />

      {loading ? (
        <p className="muted">Loading…</p>
      ) : seriesList.length === 0 ? (
        <p className="muted">{query ? "No series match your search." : "No series yet."}</p>
      ) : (
        <>
          <div className="series-grid">
            {seriesList.map((s) => (
            <div key={s.id} className="series-card" role="button" tabIndex={0}
              onClick={() => nav(`/series/${s.id}`)}
              onKeyDown={(e) => e.key === "Enter" && nav(`/series/${s.id}`)}>
                <img src={s.cover_image ? `${cdn_url}/${s.cover_image.key}` : ""} alt={s.title} />
                <div className="card-info">
                  <div className="card-title">{s.title}</div>
                  <div className="card-status">{s.status}</div>
                </div>
              </div>
            ))}
          </div>

          {canLoadMore && (
            <button
              className="btn-outline"
              style={{ marginTop: "1.5rem" }}
              onClick={() => load(query, seriesList.length, true)}
              disabled={loadingMore}
            >
              {loadingMore ? "Loading…" : `Load more (${total - seriesList.length} remaining)`}
            </button>
          )}
        </>
      )}

      {showCreate && (
        <CreateSeriesModal
          onClose={() => setShowCreate(false)}
          onCreate={(s) => {
            setSeriesList((prev) => [s, ...prev]);
            setTotal((t) => t + 1);
            setShowCreate(false);
            nav(`/series/${s.id}`);
          }}
        />
      )}
    </div>
  );
}
