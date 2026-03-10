// apps/api/src/db.controller.ts
import { Controller, Get, Inject } from '@nestjs/common';
import { ApiTags, ApiOperation, ApiResponse } from '@nestjs/swagger';
import { Pool } from 'pg';
import { PG_POOL } from '../../infrastructure/database/postgres.provider';

@ApiTags('health')
@Controller('db')
export class DbController {
  constructor(@Inject(PG_POOL) private readonly pool: Pool) {}

  @Get('health')
  @ApiOperation({ summary: 'PostgreSQL database health check' })
  @ApiResponse({ status: 200, description: 'Returns database connection status and latency' })
  async health() {
    const start = Date.now();
    const result = await this.pool.query('SELECT 1 as ok;');
    const ms = Date.now() - start;
    return {
      ok: result.rows[0]?.ok === 1,
      durationMs: ms,
    };
  }
}
