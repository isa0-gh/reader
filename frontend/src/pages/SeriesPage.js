import { jsx as _jsx, jsxs as _jsxs } from "react/jsx-runtime";
import { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { api } from "../api";
export default function SeriesPage() {
    const { id } = useParams();
    const nav = useNavigate();
    const [series, setSeries] = useState(null);
    const [error, setError] = useState("");
    useEffect(() => {
        api.getSeries(Number(id))
            .then(setSeries)
            .catch(() => setError("Series not found."));
    }, [id]);
    if (error)
        return _jsx("div", { className: "container muted", children: error });
    if (!series)
        return _jsx("div", { className: "container muted", children: "Loading\u2026" });
    const chapters = [...(series.chapters ?? [])].sort((a, b) => b.number - a.number);
    return (_jsxs("div", { className: "container", children: [_jsxs("div", { className: "series-detail", children: [_jsx("img", { src: series.cover_image || "", alt: series.title }), _jsxs("div", { className: "series-meta", children: [_jsx("h1", { children: series.title }), _jsxs("div", { className: "meta-row", children: [series.author && _jsxs("span", { children: ["Author: ", series.author] }), series.artist && series.artist !== series.author && _jsxs("span", { children: [" \u00B7 Art: ", series.artist] }), _jsxs("span", { children: [" \u00B7 ", series.status] })] }), series.description && _jsx("p", { children: series.description })] })] }), _jsxs("div", { className: "chapter-list", children: [_jsxs("h2", { children: ["Chapters (", chapters.length, ")"] }), chapters.map((ch) => (_jsx("div", { className: "chapter-item", onClick: () => nav(`/chapters/${ch.id}`), children: _jsxs("span", { className: "ch-num", children: ["Ch. ", ch.number, ch.title ? ` — ${ch.title}` : ""] }) }, ch.id)))] })] }));
}
