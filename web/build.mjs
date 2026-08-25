import { mkdir, readFile, writeFile } from "node:fs/promises";

const source = await readFile(new URL("./index.html", import.meta.url), "utf8");
await mkdir(new URL("./dist/", import.meta.url), { recursive: true });
const output = source.replace("__BUILD_MARKER__", "genealogy-story-organizer-web");
await writeFile(new URL("./dist/index.html", import.meta.url), output);
process.stdout.write("built web/dist/index.html\n");
