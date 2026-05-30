import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../AuthContext";
import { useConfig } from "../ConfigContext";
import CreateSeriesModal from "../components/CreateSeriesModal";
import { api, Series } from "../api";

const CAN_CREATE = ["uploader", "moderator", "admin"];

export default function Home() {
  const nav = useNavigate();
  const { user } = useAuth();
  const { cdn_url } = useConfig();
  const [seriesList, setSeriesList] = useState<Series[]>([]);
  const [loading, setLoading] = useState(true);
  const [showCreate, setShowCreate] = useState(false);

  useEffect(() => {
    api.listSeries().then(setSeriesList).finally(() => setLoading(false));
  }, []);

  const canCreate = user && CAN_CREATE.includes(user.role);

  return (
    <div className="container">
      <div className="page-header">
        <p className="page-title">Browse</p>
        {canCreate && <button className="btn-outline" onClick={() => setShowCreate(true)}>+ New Series</button>}
      </div>

      {loading ? (
        <p className="muted">Loading…</p>
      ) : seriesList.length === 0 ? (
        <p className="muted">No series yet.</p>
      ) : (
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
      )}

      {showCreate && (
        <CreateSeriesModal
          onClose={() => setShowCreate(false)}
          onCreate={(s) => { setSeriesList((prev) => [s, ...prev]); setShowCreate(false); nav(`/series/${s.id}`); }}
        />
      )}
    </div>
  );
}
