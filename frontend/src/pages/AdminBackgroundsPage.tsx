import { useEffect, useState, ChangeEvent } from "react";
import { PiTrash } from "react-icons/pi";
import { api, uploadFile, Background } from "../api";
import { useConfig } from "../ConfigContext";
import { useToast } from "../ToastContext";
import { useConfirm } from "../ConfirmContext";

export default function AdminBackgroundsPage() {
  const { cdn_url } = useConfig();
  const toast = useToast();
  const confirm = useConfirm();
  const [backgrounds, setBackgrounds] = useState<Background[]>([]);
  const [uploading, setUploading] = useState(false);

  function load() {
    api.listBackgrounds().then(setBackgrounds).catch(() => {});
  }

  useEffect(load, []);

  async function onFilesChange(e: ChangeEvent<HTMLInputElement>) {
    const files = Array.from(e.target.files ?? []);
    if (!files.length) return;
    e.target.value = "";
    setUploading(true);
    try {
      for (const file of files) {
        const { key, bucket } = await uploadFile(file, "backgrounds");
        await api.addBackground(key, bucket);
      }
      load();
      toast.show(`Added ${files.length} background${files.length > 1 ? "s" : ""}.`, "success");
    } catch (err: unknown) {
      toast.show(err instanceof Error ? err.message : "Failed to upload background.", "error");
    } finally {
      setUploading(false);
    }
  }

  async function handleDelete(bg: Background) {
    const ok = await confirm("Remove this background? It'll stop showing on the site immediately.", { danger: true });
    if (!ok) return;
    try {
      await api.deleteBackground(bg.id);
      setBackgrounds((bs) => bs.filter((b) => b.id !== bg.id));
      toast.show("Background removed.", "success");
    } catch (err: unknown) {
      toast.show(err instanceof Error ? err.message : "Failed to remove background.", "error");
    }
  }

  return (
    <>
      <h2>Backgrounds</h2>
      <p className="admin-subtext">
        Wallpaper images shown behind the whole site, cycling every few seconds when more than one is set.
        With none set, the site falls back to a plain gradient.
      </p>

      <div className="form-group" style={{ maxWidth: 320 }}>
        <label>Add background(s)</label>
        <input type="file" accept="image/*" multiple onChange={onFilesChange} disabled={uploading} />
      </div>

      {uploading && <p className="muted">Uploading…</p>}

      {backgrounds.length === 0 ? (
        <p className="muted">No backgrounds set — using the default gradient.</p>
      ) : (
        <div className="background-grid">
          {backgrounds.map((bg) => (
            <div key={bg.id} className="background-tile">
              <img src={`${cdn_url}/${bg.image.key}`} alt="" />
              <button className="btn-outline" style={{ color: "var(--danger)" }} onClick={() => handleDelete(bg)}>
                <PiTrash /> Remove
              </button>
            </div>
          ))}
        </div>
      )}
    </>
  );
}
