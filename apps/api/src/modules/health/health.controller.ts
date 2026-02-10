// apps/api/src/health.controller.ts
import { Controller, Get } from '@nestjs/common';

@Controller('health')
export class HealthController {
  @Get()
  health() {
    return { status: 'ok', ts: new Date().toISOString() };
  }
}
