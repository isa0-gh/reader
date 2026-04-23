// Placeholder home page — shows a grid of series.
// Replace MOCK_SERIES with a real API call once a /series list endpoint exists.
import { useNavigate } from "react-router-dom";

const MOCK_SERIES = [
  { id: 1, title: "One Piece", status: "ongoing", cover_image: "" },
  { id: 2, title: "Berserk", status: "hiatus", cover_image: "" },
  { id: 3, title: "Vagabond", status: "hiatus", cover_image: "" },
];

export default function Home() {
  const nav = useNavigate();
  return (
    <div className="container">
      <p className="page-title">Browse</p>
      <div className="series-grid">
        {MOCK_SERIES.map((s) => (
          <div key={s.id} className="series-card" onClick={() => nav(`/series/${s.id}`)}>
            <img src={s.cover_image || ""} alt={s.title} />
            <div className="card-info">
              <div className="card-title">{s.title}</div>
              <div className="card-status">{s.status}</div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
