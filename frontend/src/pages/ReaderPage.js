import { jsx as _jsx, jsxs as _jsxs, Fragment as _Fragment } from "react/jsx-runtime";
/// <reference types="vite/client" />
import { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { api } from "../api";
const CDN = import.meta.env.VITE_CDN_URL ?? "";
export default function ReaderPage() {
    const { id } = useParams();
    const nav = useNavigate();
    const [chapter, setChapter] = useState(null);
    const [error, setError] = useState("");
    useEffect(() => {
        api.getChapter(Number(id))
            .then(setChapter)
            .catch(() => setError("Chapter not found."));
    }, [id]);
    if (error)
        return _jsx("div", { className: "container muted", children: error });
    if (!chapter)
        return _jsx("div", { className: "container muted", children: "Loading\u2026" });
    const pages = [...(chapter.pages ?? [])].sort((a, b) => a.page_number - b.page_number);
    return (_jsxs(_Fragment, { children: [_jsxs("div", { className: "reader-nav", children: [_jsx("button", { onClick: () => nav(`/series/${chapter.series_id}`), children: "\u2190 Back" }), _jsxs("span", { children: ["Ch. ", chapter.number, chapter.title ? ` — ${chapter.title}` : ""] }), _jsxs("div", { style: { display: "flex", gap: "0.5rem" }, children: [_jsx("button", { onClick: () => nav(`/chapters/${Number(id) - 1}`), children: "Prev" }), _jsx("button", { onClick: () => nav(`/chapters/${Number(id) + 1}`), children: "Next" })] })] }), _jsx("div", { className: "reader-pages", children: pages.map((p) => (_jsx("img", { src: `${CDN}/${p.key}`, alt: `Page ${p.page_number}`, loading: "lazy" }, p.id))) })] }));
}
