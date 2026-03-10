import { Module } from '@nestjs/common';
import { HealthController } from './modules/health/health.controller';
import { DbController } from './modules/database/db.controller';
import { CacheController } from './modules/cache/cache.controller';
import { REDIS_CLIENT, createRedisClient } from './infrastructure/cache/redis.provider';
import { AuthModule } from './modules/auth/auth.module';
import { ChallengesModule } from './modules/challenges/challenges.module';
import { DatabaseModule } from './infrastructure/database/database.module';
import { SubmissionsModule } from './modules/submissions/submissions.module';
import { CacheModule } from './infrastructure/cache/cache.module';
import { CoursesModule } from './modules/courses/courses.module';
import { LeaderboardModule } from './modules/leaderboard/leaderboard.module';
import { MetricsModule } from './modules/metrics/metrics.module';
import { ExamsModule } from './modules/exams/exams.module';
import { AIModule } from './modules/ai/ai.module';



@Module({
  imports: [
    DatabaseModule,
    CacheModule,
    AuthModule,
    ChallengesModule,
    SubmissionsModule,
    CoursesModule,
    LeaderboardModule,
    MetricsModule,
    ExamsModule,
    AIModule,
  ],
  controllers: [HealthController, DbController, CacheController],
  providers: [
    { provide: REDIS_CLIENT, useFactory: () => createRedisClient() },
  ],
})
export class AppModule { }
