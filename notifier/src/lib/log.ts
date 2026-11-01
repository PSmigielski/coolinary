// JSON lines on stdout (docs/service-conventions.md).
import { config } from "../config";

const levels = { debug: 10, info: 20, warn: 30, error: 40 } as const;
type Level = keyof typeof levels;
const min = levels[(config.logLevel as Level) in levels ? (config.logLevel as Level) : "info"];

function write(level: Level, msg: string, fields: Record<string, unknown> = {}) {
  if (levels[level] < min) return;
  const line = JSON.stringify({ time: new Date().toISOString(), level, service: config.serviceName, msg, ...fields });
  (level === "error" ? process.stderr : process.stdout).write(line + "\n");
}

export const log = {
  debug: (msg: string, f?: Record<string, unknown>) => write("debug", msg, f),
  info: (msg: string, f?: Record<string, unknown>) => write("info", msg, f),
  warn: (msg: string, f?: Record<string, unknown>) => write("warn", msg, f),
  error: (msg: string, f?: Record<string, unknown>) => write("error", msg, f),
};
