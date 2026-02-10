import { createPgPool } from "../infrastructure/database/postgres.provider";
import { createRedisClient } from "../infrastructure/cache/redis.provider";
import { PostgresSubmissionRepo } from "../infrastructure/database/postgres/postgres-submission.repo";
import { PostgresTestCaseRepo } from "../infrastructure/database/postgres/postgres-test-case.repo";
import { Submission } from "../core/Submission/entities/submission.entity";
import { spawnSync } from "child_process";
import fs from "fs";
import path from "path";
import os from "os";
import { createLogger } from "../infrastructure/logging/logger";
import { metricsCollector } from "../infrastructure/metrics/metrics";

const logger = createLogger("SubmissionWorker");

async function runRunnerForAllTests(
    subId: string,
    challengeId: string,
    testCaseRepo: PostgresTestCaseRepo,
    pool: any,
) {
    // Get test cases from database
    const testCases = await testCaseRepo.findByChallengeId(challengeId);
    console.log(
        `[worker] Found ${testCases.length} test cases for challenge ${challengeId}`,
    );
    if (testCases.length === 0) {
        return {
            status: "error",
            message: "No test cases found for this challenge",
        };
    }

    const hostSubmissionsDir = process.env.HOST_SUBMISSIONS_DIR || "";
    console.log(`[worker] HOST_SUBMISSIONS_DIR: ${hostSubmissionsDir}`);

    // Determine paths for test cases
    // We must use a directory that is shared between the worker container and the host
    // because we are using the host's docker socket.
    // We'll use 'temp_tests' inside the project root.

    // Container path (where worker writes files)
    // Assumes worker is running in /usr/src/app/src/workers or similar, and project root is /usr/src/app
    const projectRoot = path.resolve(__dirname, "../..");
    const tempDir = path.join(projectRoot, "temp_tests", `sub-${subId}`);
    console.log(`[worker] Container tempDir: ${tempDir}`);

    // Host path (what we pass to docker -v)
    let hostTempDir = tempDir;
    if (hostSubmissionsDir) {
        // HOST_SUBMISSIONS_DIR points to .../src/core/Submission
        // We want .../temp_tests
        // We are in a Linux container, so we use path.resolve.
        // We assume HOST_SUBMISSIONS_DIR is an absolute path on the host.
        // If it's passed correctly from docker-compose, we can just resolve relative to it.
        // However, we can't use path.resolve(hostSubmissionsDir, ...) if we are inside the container
        // and hostSubmissionsDir is a path on the host system (which might be different).
        // But here, we just want to construct the string for the volume mount.

        // If HOST_SUBMISSIONS_DIR is /path/to/apps/api/src/core/Submission
        // We want /path/to/temp_tests/sub-ID

        // Go up 3 levels: .../src/core/Submission -> .../src/core -> .../src -> .../apps/api -> NO wait
        // The previous code did: ..\..\..\temp_tests
        // .../src/core/Submission (1) -> .../src/core (2) -> .../src (3) -> .../apps/api (4) ??
        // Let's look at the structure:
        // apps/api/src/core/Submission
        // apps/api/temp_tests (This seems to be where it wants to go, based on projectRoot)

        // projectRoot in worker is .../apps/api (derived from __dirname/../..)
        // tempDir is projectRoot/temp_tests/sub-ID

        // So we need to go up from Submission to apps/api.
        // src/core/Submission -> src/core -> src -> apps/api. That is 3 '..' if we are in Submission.
        // Let's just use string manipulation to be safe since we are manipulating a path that might not exist inside the container.

        const submissionDirParent = path.dirname(
            path.dirname(path.dirname(hostSubmissionsDir)),
        );
        // If hostSubmissionsDir is /.../src/core/Submission
        // dirname -> /.../src/core
        // dirname -> /.../src
        // dirname -> /.../apps/api (This is where temp_tests should be?)

        // Wait, the previous code was: path.win32.resolve(hostSubmissionsDir, "..\\..\\..\\temp_tests", `sub-${subId}`);
        // That means it went up 3 levels.

        hostTempDir = path.join(
            submissionDirParent,
            "temp_tests",
            `sub-${subId}`,
        );
    }
    console.log(`[worker] Host tempDir: ${hostTempDir}`);

    if (!fs.existsSync(tempDir)) {
        fs.mkdirSync(tempDir, { recursive: true });
    }

    // Write test cases to temporary files
    testCases.forEach((tc, index) => {
        const inputFile = path.join(tempDir, `input${index + 1}.in`);
        const outputFile = path.join(tempDir, `output${index + 1}.out`);
        fs.writeFileSync(inputFile, tc.input);
        fs.writeFileSync(outputFile, tc.expectedOutput);
        console.log(`[worker] Wrote test case ${index + 1} to ${inputFile}`);
    });

    const submissionBase = path.resolve(__dirname, "../core/Submission", subId);

    let hostSubmissionBase = submissionBase;
    if (hostSubmissionsDir) {
        hostSubmissionBase = path.join(hostSubmissionsDir, subId);
    }

    // Read meta.json
    if (!fs.existsSync(submissionBase)) {
        cleanupTempDir(tempDir);
        return null;
    }
    const metaPath = path.join(submissionBase, "meta.json");
    if (!fs.existsSync(metaPath)) {
        cleanupTempDir(tempDir);
        return null;
    }
    const meta = JSON.parse(fs.readFileSync(metaPath, "utf8"));
    const lang = meta.language;
    const codeRel = meta.codeFile;

    // Determine image by language
    let image = "";
    if (lang === "python") image = "juez_runner_python:local";
    else if (lang === "node" || lang === "javascript")
        image = "juez_runner_node:local";
    else if (lang === "cpp" || lang === "c++") image = "juez_runner_cpp:local";
    else if (lang === "java") image = "juez_runner_java:local";
    else {
        cleanupTempDir(tempDir);
        return { status: "error", message: `unsupported language ${lang}` };
    }

    // Get time limit from challenge
    const challengeResult = await pool.query(
        "SELECT * FROM challenges WHERE id = $1",
        [challengeId],
    );
    const timeLimit = challengeResult.rows[0]?.time_limit || 1500;

    // Build payload
    const payload = {
        source_file: "/code/" + path.basename(codeRel),
        time_limit_ms: timeLimit,
    };

    // Mount code and test files
    const mounts: string[] = [];
    if (hostSubmissionsDir) {
        let hostCode = path.join(hostSubmissionBase, "code");

        console.log(`[worker] Mounting host code path: ${hostCode}`);
        mounts.push("-v", `${hostCode}:/code:ro`);
    } else {
        const mountCode = path.join(submissionBase, "code");
        mounts.push("-v", `${mountCode}:/code:ro`);
    }

    // Mount temporary test directory
    mounts.push("-v", `${hostTempDir}:/tests:ro`);

    const args = [
        "run",
        "-i", // <--- ADDED THIS
        "--rm",
        "--network",
        "none",
        "--memory",
        "512m",
        "--cpus",
        "0.5",
        ...mounts,
        image,
    ];

    console.log("[worker] running docker", args.join(" "));
    const proc = spawnSync("docker", args, {
        input: JSON.stringify(payload),
        encoding: "utf8",
        maxBuffer: 50 * 1024 * 1024,
    });

    // Cleanup temporary directory
    // cleanupTempDir(tempDir); // 👈 DISABLED FOR DEBUGGING

    if (proc.error) {
        console.error("[worker] docker run error", proc.error);
        return { status: "error", message: String(proc.error) };
    }

    // Log raw output for debugging
    console.log("[worker] runner stdout:", proc.stdout);
    console.log("[worker] runner stderr:", proc.stderr);

    if (proc.status !== 0) {
        console.error("[worker] runner exit code", proc.status, proc.stderr);
        try {
            return JSON.parse(proc.stdout);
        } catch (e) {
            return { status: "error", stderr: proc.stderr || proc.stdout };
        }
    }
    try {
        const result = JSON.parse(proc.stdout);
        // Calculate score based on passed test cases
        if (result.cases && Array.isArray(result.cases)) {
            const totalPoints = testCases.reduce(
                (sum, tc) => sum + tc.points,
                0,
            );
            const earnedPoints = result.cases.reduce(
                (sum: number, c: any, idx: number) => {
                    return (
                        sum + (c.status === "OK" ? testCases[idx].points : 0)
                    );
                },
                0,
            );
            result.score =
                totalPoints > 0
                    ? Math.round((earnedPoints / totalPoints) * 100)
                    : 0;
        }
        return result;
    } catch (e) {
        console.error("[worker] invalid runner output", proc.stdout);
        return {
            status: "error",
            message: "invalid runner output",
            raw: proc.stdout,
        };
    }
}

function cleanupTempDir(dir: string) {
    try {
        if (fs.existsSync(dir)) {
            fs.rmSync(dir, { recursive: true, force: true });
        }
    } catch (err) {
        console.error("[worker] cleanup error:", err);
    }
}

async function main() {
    const pool = createPgPool();
    const redis = createRedisClient();
    const submissionRepo = new PostgresSubmissionRepo(pool);
    const testCaseRepo = new PostgresTestCaseRepo(pool);

    logger.info("Worker started, waiting for jobs");

    const shutdown = async () => {
        logger.info("Worker shutting down");
        await redis.quit();
        await pool.end();
        process.exit(0);
    };
    process.on("SIGINT", shutdown);
    process.on("SIGTERM", shutdown);

    while (true) {
        try {
            const res = await redis.brpop("queue:submissions", 0);
            if (!res) continue;
            const subId = res[1];

            logger.info("Picked submission job", { submissionId: subId });
            metricsCollector.incrementSubmissionsTotal();

            const current = await submissionRepo.findById(subId);
            if (!current) {
                logger.warn("Submission not found", { submissionId: subId });
                continue;
            }

            current.start();
            await submissionRepo.save(current);

            const startTime = Date.now();

            // Real execution using runner images with DB test cases
            const runnerResult = await runRunnerForAllTests(
                subId,
                current.challengeId,
                testCaseRepo,
                pool,
            );

            const duration = Date.now() - startTime;

            if (!runnerResult) {
                current.fail();
                await submissionRepo.save(current);
                metricsCollector.incrementSubmissionsFailed();
                logger.warn("No runner result", {
                    submissionId: subId,
                    duration,
                });
                continue;
            }

            // Update submission based on result
            if (runnerResult.status === "ACCEPTED") {
                current.accept();
                metricsCollector.incrementSubmissionsAccepted();
            } else if (
                runnerResult.status === "PARTIAL" ||
                runnerResult.status === "WRONG_ANSWER"
            ) {
                current.reject();
                metricsCollector.incrementSubmissionsRejected();
            } else if (runnerResult.status === "TIME_LIMIT_EXCEEDED") {
                current.fail();
                metricsCollector.incrementSubmissionsFailed();
            } else if (runnerResult.status === "ERROR") {
                current.fail();
                metricsCollector.incrementSubmissionsFailed();
            } else {
                current.fail();
                metricsCollector.incrementSubmissionsFailed();
            }

            // Update score and time
            if (runnerResult.score !== undefined) {
                current.updateScore(
                    runnerResult.score,
                    runnerResult.timeMsTotal || 0,
                );
                metricsCollector.recordExecutionTime(
                    runnerResult.timeMsTotal || 0,
                );
            }

            await submissionRepo.save(current);

            logger.info("Submission processed", {
                submissionId: subId,
                challengeId: current.challengeId,
                userId: current.userId,
                status: current.status,
                score: runnerResult.score || 0,
                duration,
            });
        } catch (err: any) {
            logger.error("Worker loop error", err, { error: err.message });
            await new Promise((r) => setTimeout(r, 500));
        }
    }
}

main().catch((err) => {
    logger.error("Worker fatal error", err);
    process.exit(1);
});
