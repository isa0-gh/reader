import { useEffect, useRef, useState } from "react";
import { Comment } from "../api";

const SUSPEND_PRESETS: [string, string][] = [
  ["1h", "Suspend 1 hour"],
  ["1d", "Suspend 1 day"],
  ["1y", "Suspend 1 year"],
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
  const [open, setOpen] = useState(false);
  const ref = useRef<HTMLDivElement>(null);
  const name = comment.user?.name ?? `User #${comment.user_id}`;

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

  const hasMenu = canDelete || canModerate;

  return (
    <div className="comment">
      <div className="comment-avatar" aria-hidden="true">{name.charAt(0).toUpperCase()}</div>
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
            ⋯
          </button>
          {open && (
            <div className="account-dropdown" role="menu">
              {canDelete && (
                <button
                  className="dropdown-item dropdown-item--danger"
                  onClick={() => act(() => onDelete(comment.id))}
                  role="menuitem"
                >
                  Delete
                </button>
              )}
              {canDelete && canModerate && <div className="dropdown-divider" />}
              {canModerate && (
                <>
                  <div className="dropdown-label">Suspend commenting</div>
                  {SUSPEND_PRESETS.map(([value, label]) => (
                    <button
                      key={value}
                      className="dropdown-item"
                      onClick={() => act(() => onSuspend(comment.user_id, value))}
                      role="menuitem"
                    >
                      {label}
                    </button>
                  ))}
                  <button
                    className="dropdown-item"
                    onClick={() => act(() => {
                      const d = prompt("Custom duration (Go duration syntax, e.g. 72h30m):");
                      if (d) onSuspend(comment.user_id, d);
                    })}
                    role="menuitem"
                  >
                    Custom duration…
                  </button>
                  <button
                    className="dropdown-item"
                    onClick={() => act(() => onSuspend(comment.user_id, "none"))}
                    role="menuitem"
                  >
                    Remove suspension
                  </button>
                </>
              )}
            </div>
          )}
        </div>
      )}
    </div>
  );
}
