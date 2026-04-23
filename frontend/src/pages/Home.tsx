import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../AuthContext";
import CreateSeriesModal from "../components/CreateSeriesModal";
import { Series } from "../api";

const CAN_CREATE = ["uploader", "moderator", "admin"];

// Replace MOCK_SERIES with a real API call once a /series list endpoint exists.
const INITIAL: Series[] = [
  { id: 1, title: "One Piece", slug: "one-piece", status: "ongoing", cover_image: "", description: "", author: "", artist: "", chapters: [] },
  { id: 2, title: "Berserk", slug: "berserk", status: "hiatus", cover_image: "", description: "", author: "", artist: "", chapters: [] },
];

export default function Home() {
  const nav = useNavigate();
  const { user } = useAuth();
  const [seriesList, setSeriesList] = useState<Series[]>(INITIAL);
  const [showCreate, setShowCreate] = useState(false);

  const canCreate = user && CAN_CREATE.includes(user.role);

  return (
    <div className="container">
      <div className="page-header">
        <p className="page-title">Browse</p>
        {canCreate && <button className="btn-outline" onClick={() => setShowCreate(true)}>+ New Series</button>}
      </div>

      <div className="series-grid">
        {seriesList.map((s) => (
          <div key={s.id} className="series-card" onClick={() => nav(`/series/${s.id}`)}>
            <img src={s.cover_image || ""} alt={s.title} />
            <div className="card-info">
              <div className="card-title">{s.title}</div>
              <div className="card-status">{s.status}</div>
            </div>
          </div>
        ))}
      </div>

      {showCreate && (
        <CreateSeriesModal
          onClose={() => setShowCreate(false)}
          onCreate={(s) => { setSeriesList((prev) => [s, ...prev]); setShowCreate(false); nav(`/series/${s.id}`); }}
        />
      )}
    </div>
  );
}
