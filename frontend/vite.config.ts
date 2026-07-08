import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      // Plain string target doesn't set X-Forwarded-For, and the backend's
      // ClientIPFromXFFTrustedProxies(1) (main.go) needs exactly one hop's
      // worth of XFF to resolve a client IP for rate limiting — without it,
      // /auth/login and /auth/register 500 with "client ip not resolved" in
      // dev. xfwd adds X-Forwarded-For/Host/Port, matching the one nginx hop
      // assumed in production.
      "/api": {
        target: "http://localhost:8080",
        xfwd: true,
      },
    },
  },
});
