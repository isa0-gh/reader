import { useState, useEffect, FormEvent, ChangeEvent } from "react";
import { useNavigate } from "react-router-dom";
import { PiTrash, PiXCircle } from "react-icons/pi";
import { api, uploadFile, Comment } from "../api";
import { useAuth } from "../AuthContext";
import { useConfig } from "../ConfigContext";
import { useToast } from "../ToastContext";
import { useConfirm } from "../ConfirmContext";
import { useFavorites } from "../FavoritesContext";

const COMMENTS_PAGE_SIZE = 20;
const FAVORITES_PAGE_SIZE = 12;

export default function AccountPage() {
  const { user, logout, updateUser } = useAuth();
  const { cdn_url } = useConfig();
  const toast = useToast();
  const confirm = useConfirm();
  const { favorites, toggleFavorite } = useFavorites();
  const nav = useNavigate();

  const [name, setName] = useState(user?.name ?? "");
  const [bio, setBio] = useState(user?.bio ?? "");
  const [avatarFile, setAvatarFile] = useState<File | null>(null);
  const [avatarPreview, setAvatarPreview] = useState<string | null>(null);
  const [savingProfile, setSavingProfile] = useState(false);

  const [comments, setComments] = useState<Comment[]>([]);
  const [commentsTotal, setCommentsTotal] = useState(0);
  const [loadingMoreComments, setLoadingMoreComments] = useState(false);
  const [visibleFavorites, setVisibleFavorites] = useState(FAVORITES_PAGE_SIZE);

  const [newEmail, setNewEmail] = useState(user?.email ?? "");
  const [emailPassword, setEmailPassword] = useState("");
  const [emailError, setEmailError] = useState("");
  const [emailDone, setEmailDone] = useState(false);

  const [currentPassword, setCurrentPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [passwordError, setPasswordError] = useState("");
  const [passwordDone, setPasswordDone] = useState(false);

  function loadComments(offset: number, append: boolean) {
    if (append) setLoadingMoreComments(true);
    api.listMyComments({ limit: COMMENTS_PAGE_SIZE, offset }).then((res) => {
      setComments((prev) => (append ? [...prev, ...res.items] : res.items));
      setCommentsTotal(res.total);
    }).catch(() => {}).finally(() => { if (append) setLoadingMoreComments(false); });
  }

  useEffect(() => { loadComments(0, false); }, []);

  if (!user) return <div className="container muted">Not logged in.</div>;

  function onAvatarChange(e: ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0];
    if (!file) return;
    setAvatarFile(file);
    setAvatarPreview(URL.createObjectURL(file));
  }

  async function saveProfile(e: FormEvent) {
    e.preventDefault();
    setSavingProfile(true);
    try {
      if (avatarFile) {
        const uploaded = await uploadFile(avatarFile, "avatars");
        const withAvatar = await api.setAvatar(uploaded.key, uploaded.bucket);
        updateUser({ avatar: withAvatar.avatar });
        setAvatarFile(null);
      }
      const updated = await api.updateProfile(name, bio);
      updateUser({ name: updated.name, bio: updated.bio });
      toast.show("Profile updated.", "success");
    } catch (err: unknown) {
      toast.show(err instanceof Error ? err.message : "Failed to update profile.", "error");
    } finally {
      setSavingProfile(false);
    }
  }

  async function handleRemovePhoto() {
    if (avatarFile) {
      setAvatarFile(null);
      setAvatarPreview(null);
      return;
    }
    try {
      const updated = await api.clearAvatar();
      updateUser({ avatar: updated.avatar });
      toast.show("Profile picture removed.", "success");
    } catch (err: unknown) {
      toast.show(err instanceof Error ? err.message : "Failed to remove photo.", "error");
    }
  }

  async function handleRemoveFavorite(seriesId: number) {
    try {
      await toggleFavorite(seriesId);
    } catch (err: unknown) {
      toast.show(err instanceof Error ? err.message : "Failed to remove favorite.", "error");
    }
  }

  async function handleDeleteComment(id: number) {
    const ok = await confirm("Delete this comment?", { danger: true });
    if (!ok) return;
    try {
      await api.deleteComment(id);
      setComments((cs) => cs.filter((c) => c.id !== id));
      setCommentsTotal((t) => t - 1);
      toast.show("Comment deleted.", "success");
    } catch (err: unknown) {
      toast.show(err instanceof Error ? err.message : "Failed to delete comment.", "error");
    }
  }

  async function submitEmail(e: FormEvent) {
    e.preventDefault();
    setEmailError("");
    if (newEmail === user!.email) {
      setEmailError("That's already your email.");
      return;
    }
    try {
      await api.changeEmail(emailPassword, newEmail);
      // Changing the email rotates the server-side session id, same as
      // password change — send the user back to login instead of leaving
      // them in a broken logged-in state.
      setEmailDone(true);
      setTimeout(() => {
        logout();
        nav("/login");
      }, 1500);
    } catch (err: unknown) {
      setEmailError(err instanceof Error ? err.message : "Failed to change email.");
    }
  }

  async function submitPassword(e: FormEvent) {
    e.preventDefault();
    setPasswordError("");
    if (newPassword !== confirmPassword) {
      setPasswordError("New passwords don't match.");
      return;
    }
    try {
      await api.changePassword(currentPassword, newPassword);
      // Changing the password rotates the server-side session id, so the
      // token this tab is holding is now invalid — send the user back to
      // login instead of leaving them in a broken logged-in state.
      setPasswordDone(true);
      setTimeout(() => {
        logout();
        nav("/login");
      }, 1500);
    } catch (err: unknown) {
      setPasswordError(err instanceof Error ? err.message : "Failed to change password.");
    }
  }

  const avatarSrc = avatarPreview ?? (user.avatar ? `${cdn_url}/${user.avatar.key}` : null);

  return (
    <div className="container account-page">
      <div className="page-header">
        <p className="page-title">Account</p>
      </div>

      <section className="account-section">
        <h2>Profile</h2>
        <form onSubmit={saveProfile}>
          <div className="account-avatar-row">
            <label className="avatar-upload" title="Change photo">
              {avatarSrc
                ? <img src={avatarSrc} alt={name} className="avatar-lg" />
                : <div className="avatar-lg avatar-lg--placeholder">{name.charAt(0).toUpperCase()}</div>}
              <span className="avatar-upload-hint">Change photo</span>
              <input type="file" accept="image/*" onChange={onAvatarChange} hidden />
            </label>
            <div style={{ flex: 1 }}>
              <div className="form-group">
                <label>Name</label>
                <input value={name} onChange={(e) => setName(e.target.value)} required />
              </div>
              <p className="muted">{user.email}</p>
              {avatarSrc && (
                <button type="button" className="btn-outline" style={{ marginTop: "0.5rem" }} onClick={handleRemovePhoto}>
                  <PiXCircle /> Remove photo
                </button>
              )}
            </div>
          </div>
          <div className="form-group">
            <label>Bio</label>
            <textarea value={bio} onChange={(e) => setBio(e.target.value)} rows={3} placeholder="Tell people a bit about yourself…" />
          </div>
          <button className="btn" style={{ width: "auto" }} type="submit" disabled={savingProfile}>
            {savingProfile ? "Saving…" : "Save Profile"}
          </button>
        </form>
      </section>

      <section className="account-section">
        <h2>My Favorites {favorites.length > 0 ? `(${favorites.length})` : ""}</h2>
        {favorites.length === 0 ? (
          <p className="muted">No favorites yet. Heart a series to track it here.</p>
        ) : (
          <>
            <div className="favorite-list">
              {favorites.slice(0, visibleFavorites).map((f) => {
                const s = f.series;
                if (!s) return null;
                // Base on last_read_chapter (the preloaded object), not the raw
                // _id: if that chapter's since been deleted, the id is stale but
                // the object comes back null, matching the "Not started" label.
                const target = f.last_read_chapter ? `/series/${s.id}/${f.last_read_chapter.id}` : `/series/${s.id}`;
                return (
                  <div key={f.id} className="favorite-row">
                    {s.cover_image
                      ? <img src={`${cdn_url}/${s.cover_image.key}`} alt={s.title} className="favorite-row-thumb" />
                      : <div className="favorite-row-thumb cover-placeholder" aria-hidden="true" />}
                    <div className="favorite-row-info" role="button" tabIndex={0}
                      onClick={() => nav(target)} onKeyDown={(e) => e.key === "Enter" && nav(target)}>
                      <div className="comment-author">{s.title}</div>
                      <div className="muted">{f.last_read_chapter ? `Continue: Ch. ${f.last_read_chapter.number}` : "Not started"}</div>
                    </div>
                    <button className="btn-outline" onClick={() => handleRemoveFavorite(s.id)}><PiXCircle /> Remove</button>
                  </div>
                );
              })}
            </div>
            {favorites.length > visibleFavorites && (
              <button className="btn-outline" style={{ marginTop: "0.75rem" }}
                onClick={() => setVisibleFavorites((n) => n + FAVORITES_PAGE_SIZE)}>
                Show more ({favorites.length - visibleFavorites} remaining)
              </button>
            )}
          </>
        )}
      </section>

      <section className="account-section">
        <h2>My Comments {commentsTotal > 0 ? `(${commentsTotal})` : ""}</h2>
        {comments.length === 0 ? (
          <p className="muted">You haven't posted any comments yet.</p>
        ) : (
          <>
            <div className="comment-list">
              {comments.map((c) => (
                <div key={c.id} className="comment">
                  <div className="comment-body">
                    <div className="comment-header-row">
                      {c.chapter
                        ? <button className="dropdown-item" style={{ width: "auto", padding: 0 }}
                            onClick={() => nav(`/series/${c.chapter!.series_id}/${c.chapter!.id}`)}>
                            Ch. {c.chapter.number}{c.chapter.title ? ` — ${c.chapter.title}` : ""}
                          </button>
                        : <span className="comment-author">Deleted chapter</span>}
                      <span className="comment-time">{new Date(c.created_at).toLocaleString()}</span>
                    </div>
                    <p className="comment-text">{c.body}</p>
                  </div>
                  <button className="btn-outline" style={{ color: "var(--danger)" }} onClick={() => handleDeleteComment(c.id)}>
                    <PiTrash /> Delete
                  </button>
                </div>
              ))}
            </div>
            {comments.length < commentsTotal && (
              <button className="btn-outline" style={{ marginTop: "0.75rem" }}
                onClick={() => loadComments(comments.length, true)} disabled={loadingMoreComments}>
                {loadingMoreComments ? "Loading…" : `Load more (${commentsTotal - comments.length} remaining)`}
              </button>
            )}
          </>
        )}
      </section>

      <section className="account-section">
        <h2>Email</h2>
        {emailDone ? (
          <p className="muted">Email changed. Redirecting to login…</p>
        ) : (
          <form onSubmit={submitEmail}>
            <div className="form-group">
              <label>New email</label>
              <input type="email" value={newEmail} onChange={(e) => setNewEmail(e.target.value)} required />
            </div>
            <div className="form-group">
              <label>Current password</label>
              <input type="password" value={emailPassword} onChange={(e) => setEmailPassword(e.target.value)} required />
            </div>
            {emailError && <div className="error">{emailError}</div>}
            <button className="btn" style={{ width: "auto" }} type="submit">Change Email</button>
          </form>
        )}
      </section>

      <section className="account-section">
        <h2>Password</h2>
        {passwordDone ? (
          <p className="muted">Password changed. Redirecting to login…</p>
        ) : (
          <form onSubmit={submitPassword}>
            <div className="form-group">
              <label>Current password</label>
              <input type="password" value={currentPassword} onChange={(e) => setCurrentPassword(e.target.value)} required />
            </div>
            <div className="form-group">
              <label>New password</label>
              <input type="password" value={newPassword} onChange={(e) => setNewPassword(e.target.value)} required />
            </div>
            <div className="form-group">
              <label>Confirm new password</label>
              <input type="password" value={confirmPassword} onChange={(e) => setConfirmPassword(e.target.value)} required />
            </div>
            {passwordError && <div className="error">{passwordError}</div>}
            <button className="btn" style={{ width: "auto" }} type="submit">Change Password</button>
          </form>
        )}
      </section>
    </div>
  );
}
