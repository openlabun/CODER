// apps/api/src/cache.controller.ts
import { Controller, Get, Inject } from '@nestjs/common';
import { ApiTags, ApiOperation, ApiResponse } from '@nestjs/swagger';
import { REDIS_CLIENT } from '../../infrastructure/cache/redis.provider';
import type Redis from 'ioredis';

@ApiTags('health')
@Controller('cache')
export class CacheController {
  constructor(@Inject(REDIS_CLIENT) private readonly redis: Redis) {}

  @Get('health')
  @ApiOperation({ summary: 'Redis cache health check' })
  @ApiResponse({ status: 200, description: 'Returns Redis connection status and latency' })
  async health() {
    const start = Date.now();
    try {
      const pong = await this.redis.ping();
      const ms = Date.now() - start;
      return { ok: pong === 'PONG', durationMs: ms };
    } catch (err: any) {
      const ms = Date.now() - start;
      return { ok: false, durationMs: ms, error: err?.message ?? 'unknown error' };
    }
  }
}
