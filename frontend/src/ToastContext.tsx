import { createContext, useCallback, useContext, useRef, useState, ReactNode } from "react";

type ToastKind = "success" | "error" | "info";
interface ToastItem { id: number; kind: ToastKind; message: string; }

interface ToastCtx { show(message: string, kind?: ToastKind): void; }

const Ctx = createContext<ToastCtx>(null!);

export function ToastProvider({ children }: { children: ReactNode }) {
  const [toasts, setToasts] = useState<ToastItem[]>([]);
  const nextId = useRef(0);

  const dismiss = useCallback((id: number) => {
    setToasts((t) => t.filter((x) => x.id !== id));
  }, []);

  const show = useCallback((message: string, kind: ToastKind = "info") => {
    const id = nextId.current++;
    setToasts((t) => [...t, { id, kind, message }]);
    setTimeout(() => dismiss(id), 4500);
  }, [dismiss]);

  return (
    <Ctx.Provider value={{ show }}>
      {children}
      <div className="toast-viewport" role="region" aria-label="Notifications">
        {toasts.map((t) => (
          <div key={t.id} className={`toast toast--${t.kind}`} role="status" onClick={() => dismiss(t.id)}>
            <span className="toast-dot" aria-hidden="true" />
            <span className="toast-message">{t.message}</span>
          </div>
        ))}
      </div>
    </Ctx.Provider>
  );
}

export const useToast = () => useContext(Ctx);
