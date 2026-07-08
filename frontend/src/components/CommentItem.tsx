import { useEffect, useRef, useState } from "react";
import { PiDotsThree } from "react-icons/pi";
import { Comment } from "../api";
import Modal from "../Modal";
import { useConfig } from "../ConfigContext";

const CUSTOM_UNITS: [string, number][] = [
  ["Hours", 1],
  ["Days", 24],
  ["Weeks", 24 * 7],
];

export default function CommentItem({
  comment, canDelete, canModerate, onDelete, onSuspend,
}: {
  comment: Comment;
  canDelete: boolean;
  canModerate: boolean;
  onDelete(id: number): void;
  onSuspend(userId: number, duration: string): void;
}) {
  const { cdn_url } = useConfig();
  const [open, setOpen] = useState(false);
  const [customOpen, setCustomOpen] = useState(false);
  const [amount, setAmount] = useState(3);
  const [unitHours, setUnitHours] = useState(24);
  const ref = useRef<HTMLDivElement>(null);
  const name = comment.user?.name ?? `User #${comment.user_id}`;
  const avatar = comment.user?.avatar;
  const suspendedUntil = comment.user?.comment_suspended_until;
  const isSuspended = !!suspendedUntil && new Date(suspendedUntil) > new Date();

  useEffect(() => {
    function handler(e: MouseEvent) {
      if (ref.current && !ref.current.contains(e.target as Node)) setOpen(false);
    }
    document.addEventListener("mousedown", handler);
    return () => document.removeEventListener("mousedown", handler);
  }, []);

  function act(fn: () => void) {
    fn();
    setOpen(false);
  }

  function submitCustom(e: React.FormEvent) {
    e.preventDefault();
    if (amount > 0) onSuspend(comment.user_id, `${Math.round(amount * unitHours)}h`);
    setCustomOpen(false);
  }

  const hasMenu = canDelete || canModerate;

  return (
    <div className="comment">
      {avatar
        ? <img src={`${cdn_url}/${avatar.key}`} alt="" className="comment-avatar" />
        : <div className="comment-avatar comment-avatar--placeholder" aria-hidden="true">{name.charAt(0).toUpperCase()}</div>}
      <div className="comment-body">
        <div className="comment-header-row">
          <span className="comment-author">{name}</span>
          <span className="comment-time">{new Date(comment.created_at).toLocaleString()}</span>
        </div>
        <p className="comment-text">{comment.body}</p>
      </div>

      {hasMenu && (
        <div className="comment-menu" ref={ref}>
          <button
            className="comment-menu-trigger"
            onClick={() => setOpen((o) => !o)}
            aria-haspopup="true"
            aria-expanded={open}
            aria-label="Comment actions"
          >
            <PiDotsThree />
          </button>
          {open && (
            <div className="account-dropdown" role="menu">
              {canDelete && (
                <button
                  className="dropdown-item dropdown-item--danger"
                  onClick={() => act(() => onDelete(comment.id))}
                  role="menuitem"
                >
                  Delete comment
                </button>
              )}
              {canDelete && canModerate && <div className="dropdown-divider" />}
              {canModerate && (
                <>
                  <div className="dropdown-label">Commenting</div>
                  {isSuspended && (
                    <div className="dropdown-status">
                      Suspended until {new Date(suspendedUntil!).toLocaleString()}
                    </div>
                  )}
                  <button
                    className="dropdown-item"
                    onClick={() => act(() => setCustomOpen(true))}
                    role="menuitem"
                  >
                    Suspend…
                  </button>
                  {isSuspended && (
                    <button
                      className="dropdown-item"
                      onClick={() => act(() => onSuspend(comment.user_id, "none"))}
                      role="menuitem"
                    >
                      Remove suspension
                    </button>
                  )}
                </>
              )}
            </div>
          )}
        </div>
      )}

      {customOpen && (
        <Modal title={`Suspend ${name}`} onClose={() => setCustomOpen(false)}>
          <form onSubmit={submitCustom} className="modal-form">
            <div className="form-row">
              <div className="form-group">
                <label>Amount</label>
                <input
                  type="number"
                  min={1}
                  value={amount}
                  onChange={(e) => setAmount(parseInt(e.target.value) || 1)}
                  autoFocus
                />
              </div>
              <div className="form-group">
                <label>Unit</label>
                <select value={unitHours} onChange={(e) => setUnitHours(Number(e.target.value))}>
                  {CUSTOM_UNITS.map(([label, hours]) => (
                    <option key={label} value={hours}>{label}</option>
                  ))}
                </select>
              </div>
            </div>
            <button className="btn" type="submit">Suspend commenting</button>
          </form>
        </Modal>
      )}
    </div>
  );
}
