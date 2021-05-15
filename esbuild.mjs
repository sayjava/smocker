import esbuild from "esbuild";
import sassPlugin from "esbuild-plugin-sass";
import htmlPlugin from "@chialab/esbuild-plugin-html";

esbuild
  .build({
    entryPoints: ["client/index.html"],
    outdir: "build",
    bundle: true,
    minify: true,
    sourcemap: true,
    loader: {
      ".png": "file",
    },
    plugins: [htmlPlugin(), sassPlugin()],
  })
  .catch((e) => console.error(e.message));
