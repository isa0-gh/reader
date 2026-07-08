import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { PiPlus, PiMagnifyingGlass } from "react-icons/pi";
import { useAuth } from "../AuthContext";
import { useConfig } from "../ConfigContext";
import { useFavorites } from "../FavoritesContext";
import CreateSeriesModal from "../components/CreateSeriesModal";
import { api, Series } from "../api";

const CAN_CREATE = ["uploader", "moderator", "admin"];
const PAGE_SIZE = 24;

export default function Home() {
  const nav = useNavigate();
  const { user } = useAuth();
  const { cdn_url } = useConfig();
  const { favorites } = useFavorites();
  const [tab, setTab] = useState<"browse" | "favorites">("browse");
  const [seriesList, setSeriesList] = useState<Series[]>([]);
  const [total, setTotal] = useState(0);
  const [query, setQuery] = useState("");
  const [loading, setLoading] = useState(true);
  const [loadingMore, setLoadingMore] = useState(false);
  const [showCreate, setShowCreate] = useState(false);
  const [visibleFavorites, setVisibleFavorites] = useState(PAGE_SIZE);

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
        <p className="page-title">{tab === "browse" ? "Browse" : "My Favorites"}</p>
        {tab === "browse" && canCreate && <button className="btn-outline" onClick={() => setShowCreate(true)}><PiPlus /> New Series</button>}
      </div>

      {user && (
        <div className="home-tabs">
          <button className={"home-tab" + (tab === "browse" ? " home-tab--active" : "")} onClick={() => setTab("browse")}>Browse</button>
          <button className={"home-tab" + (tab === "favorites" ? " home-tab--active" : "")} onClick={() => setTab("favorites")}>
            My Favorites{favorites.length > 0 ? ` (${favorites.length})` : ""}
          </button>
        </div>
      )}

      {tab === "browse" ? (
        <>
          <div className="search-input-wrap">
            <PiMagnifyingGlass className="search-input-icon" aria-hidden="true" />
            <input
              className="search-input"
              placeholder="Search by title…"
              value={query}
              onChange={(e) => setQuery(e.target.value)}
            />
          </div>

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
                    {s.cover_image
                      ? <img src={`${cdn_url}/${s.cover_image.key}`} alt={s.title} />
                      : <div className="cover-placeholder" aria-hidden="true" />}
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
        </>
      ) : (
        favorites.length === 0 ? (
          <p className="muted">No favorites yet. Heart a series to track it here.</p>
        ) : (
          <>
            <div className="series-grid">
              {favorites.slice(0, visibleFavorites).map((f) => {
                const s = f.series;
                if (!s) return null;
                // Base on last_read_chapter (the preloaded object), not the raw
                // _id: if that chapter's since been deleted, the id is stale but
                // the object comes back null, matching the "Not started" label.
                const target = f.last_read_chapter
                  ? `/series/${s.id}/${f.last_read_chapter.id}`
                  : `/series/${s.id}`;
                return (
                  <div key={f.id} className="series-card" role="button" tabIndex={0}
                    onClick={() => nav(target)}
                    onKeyDown={(e) => e.key === "Enter" && nav(target)}>
                    {s.cover_image
                      ? <img src={`${cdn_url}/${s.cover_image.key}`} alt={s.title} />
                      : <div className="cover-placeholder" aria-hidden="true" />}
                    <div className="card-info">
                      <div className="card-title">{s.title}</div>
                      <div className="card-status">
                        {f.last_read_chapter ? `Continue: Ch. ${f.last_read_chapter.number}` : "Not started"}
                      </div>
                    </div>
                  </div>
                );
              })}
            </div>
            {favorites.length > visibleFavorites && (
              <button className="btn-outline" style={{ marginTop: "1.5rem" }}
                onClick={() => setVisibleFavorites((n) => n + PAGE_SIZE)}>
                Show more ({favorites.length - visibleFavorites} remaining)
              </button>
            )}
          </>
        )
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
