// apps/api/src/redis.provider.ts
import Redis from 'ioredis';

export const REDIS_CLIENT = 'REDIS_CLIENT';

export function createRedisClient() {
  const host = process.env.REDIS_HOST || '127.0.0.1';
  const port = Number(process.env.REDIS_PORT || 6379);

  const client = new Redis({
    host,
    port,
    lazyConnect: false,      // conecta de una vez
    maxRetriesPerRequest: 2, // no se quede colgado
  });

  client.on('connect', () => console.log('[redis] conectado'));
  client.on('error', (err) => console.error('[redis] error:', err.message));

  return client;
}
