import { Global, Module } from '@nestjs/common';
import { REDIS_CLIENT, createRedisClient } from './redis.provider';

@Global() // 👈 esto lo hace visible en todos los módulos
@Module({
  providers: [
    { provide: REDIS_CLIENT, useFactory: () => createRedisClient() },
  ],
  exports: [REDIS_CLIENT],
})
export class CacheModule {}
