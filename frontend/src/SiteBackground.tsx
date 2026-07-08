import { useEffect, useState } from "react";
import { PiImage, PiImageBroken } from "react-icons/pi";
import { api, Background } from "./api";
import { useConfig } from "./ConfigContext";

const ROTATE_MS = 12000;
const STORAGE_KEY = "wallpaper_enabled";

// Mounted once at the app root, behind everything (fixed, pointer-events:none).
// Renders nothing — falling back to index.css's plain gradient mesh — until
// an admin has uploaded at least one wallpaper via the Backgrounds admin panel.
// Opt-in, not opt-out: defaults off until a visitor explicitly turns it on
// (plain localStorage preference, no account needed).
export default function SiteBackground() {
  const { cdn_url } = useConfig();
  const [backgrounds, setBackgrounds] = useState<Background[]>([]);
  const [active, setActive] = useState(0);
  const [enabled, setEnabled] = useState(() => localStorage.getItem(STORAGE_KEY) === "1");

  useEffect(() => {
    api.listBackgrounds().then(setBackgrounds).catch(() => {});
  }, []);

  // Muted/secondary text tokens are tuned for a flat --bg and read as
  // washed-out over a photo; boost their contrast site-wide (CSS custom
  // properties cascade) only while a wallpaper is actually showing.
  useEffect(() => {
    const showing = backgrounds.length > 0 && enabled;
    document.body.classList.toggle("has-wallpaper", showing);
    return () => document.body.classList.remove("has-wallpaper");
  }, [backgrounds.length, enabled]);

  useEffect(() => {
    if (backgrounds.length < 2 || !enabled) return;
    if (window.matchMedia("(prefers-reduced-motion: reduce)").matches) return;
    const id = setInterval(() => setActive((i) => (i + 1) % backgrounds.length), ROTATE_MS);
    return () => clearInterval(id);
  }, [backgrounds.length, enabled]);

  function toggle() {
    setEnabled((e) => {
      const next = !e;
      localStorage.setItem(STORAGE_KEY, next ? "1" : "0");
      return next;
    });
  }

  if (backgrounds.length === 0) return null;

  return (
    <>
      {enabled && (
        <div className="site-background" aria-hidden="true">
          {backgrounds.map((bg, i) => (
            <div
              key={bg.id}
              className={"site-background-layer" + (i === active ? " site-background-layer--active" : "")}
              style={{ backgroundImage: `url(${cdn_url}/${bg.image.key})` }}
            />
          ))}
          <div className="site-background-scrim" />
        </div>
      )}
      <button
        className="wallpaper-toggle"
        onClick={toggle}
        title={enabled ? "Turn background image off" : "Turn background image on"}
      >
        {enabled ? <PiImage /> : <PiImageBroken />}
        {enabled ? "Wallpaper on" : "Wallpaper off"}
      </button>
    </>
  );
}
