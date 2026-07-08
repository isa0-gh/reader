/// <reference types="vite/client" />

import { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { api, Chapter, Comment, Series } from "../api";
import { useAuth } from "../AuthContext";
import { useConfig } from "../ConfigContext";
import { useToast } from "../ToastContext";
import { useConfirm } from "../ConfirmContext";
import CommentItem from "../components/CommentItem";

const CAN_CREATE = ["uploader", "moderator", "admin"];
const CAN_MODERATE = ["moderator", "admin"];

export default function ReaderPage() {
  const { id: seriesId, chapterId } = useParams<{ id: string; chapterId: string }>();
  const nav = useNavigate();
  const { user } = useAuth();
  const { cdn_url } = useConfig();
  const toast = useToast();
  const confirm = useConfirm();
  const [chapter, setChapter] = useState<Chapter | null>(null);
  const [series, setSeries] = useState<Series | null>(null);
  const [error, setError] = useState("");
  const [comments, setComments] = useState<Comment[]>([]);
  const [commentBody, setCommentBody] = useState("");
  const [commentError, setCommentError] = useState("");

  useEffect(() => {
    api.getChapter(Number(chapterId)).then(setChapter).catch(() => setError("Chapter not found."));
  }, [chapterId]);

  useEffect(() => {
    api.getSeries(Number(seriesId)).then(setSeries).catch(() => {});
  }, [seriesId]);

  useEffect(() => {
    if (!chapterId) return;
    api.listComments(Number(chapterId)).then((res) => setComments(res.items)).catch(() => {});
  }, [chapterId]);

  async function handleAddComment(e: React.FormEvent) {
    e.preventDefault();
    if (!commentBody.trim()) return;
    try {
      const c = await api.createComment(Number(chapterId), commentBody);
      setComments((cs) => [...cs, c]);
      setCommentBody("");
      setCommentError("");
    } catch (e: any) { setCommentError(e.message); }
  }

  async function handleDeleteComment(id: number) {
    const ok = await confirm("Delete this comment?", { danger: true });
    if (!ok) return;
    try {
      await api.deleteComment(id);
      setComments((cs) => cs.filter((c) => c.id !== id));
      toast.show("Comment deleted.", "success");
    } catch (e: any) { toast.show(e.message, "error"); }
  }

  async function handleSuspend(userId: number, duration: string) {
    try {
      await api.suspendComments(userId, duration);
      toast.show(duration && duration !== "none" ? "Commenting suspended." : "Suspension cleared.", "success");
    } catch (e: any) { toast.show(e.message, "error"); }
  }

  if (error) return <div className="container muted">{error}</div>;
  if (!chapter) return <div className="container muted">Loading…</div>;

  const pages = [...(chapter.pages ?? [])].sort((a, b) => a.page_number - b.page_number);
  const canEdit = user && CAN_CREATE.includes(user.role);

  const siblings = [...(series?.chapters ?? [])].sort((a, b) => a.number - b.number);
  const idx = siblings.findIndex(c => c.id === chapter.id);
  const prevId = idx > 0 ? siblings[idx - 1].id : null;
  const nextId = idx !== -1 && idx < siblings.length - 1 ? siblings[idx + 1].id : null;

  return (
    <>
      <div className="reader-nav">
        <button onClick={() => nav(`/series/${seriesId}`)}>← Back</button>
        <span>Ch. {chapter.number}{chapter.title ? ` — ${chapter.title}` : ""}</span>
        <div style={{ display: "flex", gap: "0.5rem" }}>
          {canEdit && <button onClick={() => nav(`/series/${seriesId}/${chapterId}/edit`)}>Edit</button>}
          <button disabled={!prevId} onClick={() => prevId && nav(`/series/${seriesId}/${prevId}`)}>Prev</button>
          <button disabled={!nextId} onClick={() => nextId && nav(`/series/${seriesId}/${nextId}`)}>Next</button>
        </div>
      </div>

      <div className="reader-pages">
        {pages.map((p) => (
          <img key={p.id} src={`${cdn_url}/${p.key}`} alt={`Page ${p.page_number}`} loading="lazy" />
        ))}
      </div>

      <div className="container comment-section" style={{ maxWidth: "40rem" }}>
        <h2 className="comment-section-header">Comments ({comments.length})</h2>

        <div className="comment-list">
          {comments.map((c) => (
            <CommentItem
              key={c.id}
              comment={c}
              canDelete={!!user && (c.user_id === user.id || CAN_MODERATE.includes(user.role))}
              canModerate={!!user && CAN_MODERATE.includes(user.role) && c.user_id !== user.id}
              onDelete={handleDeleteComment}
              onSuspend={handleSuspend}
            />
          ))}
        </div>

        {user ? (
          <form className="comment-form form-group" onSubmit={handleAddComment}>
            <textarea
              value={commentBody}
              onChange={(e) => setCommentBody(e.target.value)}
              placeholder="Add a comment…"
              rows={3}
            />
            {commentError && <p className="error">{commentError}</p>}
            <button className="btn" type="submit">Post comment</button>
          </form>
        ) : (
          <p className="comment-login-prompt">Log in to comment.</p>
        )}
      </div>
    </>
  );
}
