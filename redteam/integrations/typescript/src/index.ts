import { spawnSync } from "node:child_process";
import { readFileSync } from "node:fs";
import { resolve } from "node:path";

export type RedteamFormat = "json" | "html";

export interface RunRedteamOptions {
  /** Path to redteam.yaml */
  configPath: string;
  /** Output report path */
  outputPath: string;
  format?: RedteamFormat;
  /** Use CLI mock invoker (no HTTP to target) */
  mock?: boolean;
  /** Exit non-zero when any finding is present */
  failOnFindings?: boolean;
  /** Extra YAML attack directories (each passed as -attacks-dir) */
  attacksDirs?: string[];
  /** Skip embedded library (only -attacks-dir packs) */
  noBuiltin?: boolean;
  /** Override path to redteam binary */
  binaryPath?: string;
  cwd?: string;
}

function resolveBinary(explicit?: string): string {
  if (explicit) return resolve(explicit);
  const env = process.env.REDTEAM_BIN;
  if (env) return resolve(env);
  return "redteam";
}

/**
 * Runs the redteam CLI synchronously. Requires the binary on PATH or REDTEAM_BIN.
 */
export function runRedteamSync(opts: RunRedteamOptions): {
  status: number | null;
  stderr: string;
  stdout: string;
} {
  const bin = resolveBinary(opts.binaryPath);
  const args = [
    "run",
    "--config",
    opts.configPath,
    "--output",
    opts.outputPath,
    "--format",
    opts.format ?? "json",
  ];
  if (opts.mock) args.push("--mock");
  if (opts.failOnFindings) args.push("--fail-on-findings");
  if (opts.noBuiltin) args.push("--no-builtin");
  for (const d of opts.attacksDirs ?? []) {
    args.push("--attacks-dir", d);
  }

  const r = spawnSync(bin, args, {
    encoding: "utf8",
    cwd: opts.cwd,
    env: process.env,
  });
  return {
    status: r.status,
    stderr: r.stderr ?? "",
    stdout: r.stdout ?? "",
  };
}

/** Parse JSON report after a successful run. */
export function readReportJson(reportPath: string): unknown {
  return JSON.parse(readFileSync(reportPath, "utf8"));
}
