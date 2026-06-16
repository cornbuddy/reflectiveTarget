const development = process.env.NODE_ENV == "development"
const sourcemap = development
    ? "inline"
    : "none";
const minify = !development;

const result = await Bun.build({
    minify,
    sourcemap,
    entrypoints: ["./js/index.js"],
    outdir: "./dist",
});

console.log(result);
