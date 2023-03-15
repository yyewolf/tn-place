import { resolve } from "path";
/** @type {import('vite').UserConfig} */
module.exports = {
  root: "src",
  build: {
    outDir: "../dist",
    emptyOutDir: true,
    rollupOptions: {
      input: {
        main: resolve(__dirname, "src/index.html"),
        place: resolve(__dirname, "src/place.html"),
      },
    },
  },
};
