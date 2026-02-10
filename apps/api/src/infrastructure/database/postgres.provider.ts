// apps/api/src/postgres.provider.ts
import { Pool } from 'pg';

export const PG_POOL = 'PG_POOL';

export function createPgPool() {
  const pool = new Pool({
    host: process.env.POSTGRES_HOST || 'localhost',
    port: Number(process.env.POSTGRES_PORT || 5432),
    user: process.env.POSTGRES_USER || 'juez',
    password: process.env.POSTGRES_PASSWORD || 'secret',
    database: process.env.POSTGRES_DB || 'juezdb',
    max: 10, // conexiones en el pool
    idleTimeoutMillis: 10_000,
  });
  return pool;
}
