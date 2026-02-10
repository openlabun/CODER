import fs from "fs";
import { spawnSync } from "child_process";
import { glob } from "glob";

// === Leer configuración desde stdin ===
const raw = fs.readFileSync(0, "utf8");
const payload = JSON.parse(raw || "{}");

const src = payload.source_file || "/code/solution.js";
const timeLimit = payload.time_limit_ms || payload.time_limit || 1500;

// === 1️⃣ Verificar errores de sintaxis (como compilación) ===
const check = spawnSync("node", ["--check", src], { encoding: "utf8" });
if (check.status !== 0) {
    console.log(
        JSON.stringify({
            status: "COMPILATION_ERROR",
            timeMsTotal: 0,
            cases: [],
            stderr: (check.stderr || "").trim(),
        }),
    );
    process.exit(0);
}

// === 2️⃣ Buscar archivos de test ===
const inputs = glob.sync("/tests/*.in").sort();
const cases = [];
let totalTime = 0;
let correct = 0;
let runtimeErrors = 0;
let wrongAnswers = 0;

// === 3️⃣ Ejecutar cada test ===
for (let i = 0; i < inputs.length; i++) {
    const inFile = inputs[i];
    const base = inFile.match(/input(\d+)\.in$/)?.[1] || "";
    const outFile = `/tests/output${base}.out`;

    const input = fs.existsSync(inFile) ? fs.readFileSync(inFile, "utf8") : "";

    const expected = fs.existsSync(outFile)
        ? fs.readFileSync(outFile, "utf8").trim()
        : "";

    const start = Date.now();
    let proc;

    try {
        proc = spawnSync("node", [src], {
            input,
            timeout: timeLimit,
            encoding: "utf8",
            maxBuffer: 10 * 1024 * 1024,
            shell: false,
        });
    } catch (e) {
        // Error lanzado directamente por spawnSync
        cases.push({
            caseId: i + 1,
            status: "RUNTIME_ERROR",
            timeMs: 0,
            stderr: e.message || String(e),
        });
        runtimeErrors++;
        continue;
    }

    const elapsed = Date.now() - start;
    totalTime += elapsed;

    // Timeout (TLE)
    if (proc.error && proc.error.code === "ETIMEDOUT") {
        cases.push({
            caseId: i + 1,
            status: "TLE",
            timeMs: timeLimit,
            stderr: "",
        });
        runtimeErrors++;
        continue;
    }

    // Proceso con error (exit code != 0)
    if (proc.status !== 0) {
        const stderrMsg =
            (proc.stderr || "").trim() || (proc.error?.message ?? "");
        cases.push({
            caseId: i + 1,
            status: "RUNTIME_ERROR",
            timeMs: elapsed,
            stderr: stderrMsg,
        });
        runtimeErrors++;
        continue;
    }

    // Comparar salidas
    const actual = (proc.stdout || "").trim();
    const ok = actual === expected;

    cases.push({
        caseId: i + 1,
        status: ok ? "OK" : "WRONG_ANSWER",
        timeMs: elapsed,
        stderr: ok ? "" : `Expected: "${expected}", Got: "${actual}"`,
    });

    if (ok) correct++;
    else wrongAnswers++;
}

// === 4️⃣ Determinar estado general ===
let status = "ACCEPTED";
if (runtimeErrors > 0) status = "RUNTIME_ERROR";
else if (wrongAnswers > 0) status = "WRONG_ANSWER";

console.log(JSON.stringify({ status, timeMsTotal: totalTime, cases }));
