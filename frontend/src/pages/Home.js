import { jsx as _jsx, jsxs as _jsxs } from "react/jsx-runtime";
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
    return (_jsxs("div", { className: "container", children: [_jsx("p", { className: "page-title", children: "Browse" }), _jsx("div", { className: "series-grid", children: MOCK_SERIES.map((s) => (_jsxs("div", { className: "series-card", onClick: () => nav(`/series/${s.id}`), children: [_jsx("img", { src: s.cover_image || "", alt: s.title }), _jsxs("div", { className: "card-info", children: [_jsx("div", { className: "card-title", children: s.title }), _jsx("div", { className: "card-status", children: s.status })] })] }, s.id))) })] }));
}
