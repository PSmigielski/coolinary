// Bun's built-in Redis client. redis-cache may evict; redis-state never does.
import { RedisClient } from "bun";
import { config } from "../config";

export const redisCache = new RedisClient(config.redisCacheUrl, { enableOfflineQueue: false });
export const redisState = new RedisClient(config.redisStateUrl, { enableOfflineQueue: false });

export async function ping(client: RedisClient): Promise<string> {
  try {
    await client.send("PING", []);
    return "ok";
  } catch (err) {
    return `down: ${(err as Error).message}`;
  }
}
