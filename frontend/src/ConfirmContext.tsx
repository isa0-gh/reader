import { createContext, useCallback, useContext, useRef, useState, ReactNode } from "react";
import Modal from "./Modal";

interface ConfirmOptions { title?: string; danger?: boolean; confirmLabel?: string; }
interface ConfirmCtx { confirm(message: string, options?: ConfirmOptions): Promise<boolean>; }

const Ctx = createContext<ConfirmCtx>(null!);

interface PendingConfirm extends ConfirmOptions { message: string; }

export function ConfirmProvider({ children }: { children: ReactNode }) {
  const [pending, setPending] = useState<PendingConfirm | null>(null);
  const resolver = useRef<((v: boolean) => void) | undefined>(undefined);

  const confirm = useCallback((message: string, options?: ConfirmOptions) => {
    setPending({ message, ...options });
    return new Promise<boolean>((resolve) => { resolver.current = resolve; });
  }, []);

  function settle(value: boolean) {
    resolver.current?.(value);
    setPending(null);
  }

  return (
    <Ctx.Provider value={{ confirm }}>
      {children}
      {pending && (
        <Modal title={pending.title ?? "Are you sure?"} onClose={() => settle(false)}>
          <p className="confirm-message">{pending.message}</p>
          <div className="confirm-actions">
            <button className="btn-outline" onClick={() => settle(false)}>Cancel</button>
            <button
              className={pending.danger ? "btn btn--danger" : "btn"}
              style={{ width: "auto" }}
              onClick={() => settle(true)}
              autoFocus
            >
              {pending.confirmLabel ?? (pending.danger ? "Delete" : "Confirm")}
            </button>
          </div>
        </Modal>
      )}
    </Ctx.Provider>
  );
}

export const useConfirm = () => useContext(Ctx).confirm;
