import { Module } from '@nestjs/common';
import { AuthModule } from '../auth/auth.module';              
import { SubmissionsController } from './submissions.controller';
import { SubmissionsService } from './submissions.service';

import { CreateSubmissionUseCase } from '../../core/Submission/use-cases/create-submission.use-case';
import { GetSubmissionUseCase } from '../../core/Submission/use-cases/get-submission.use-case';

import { PG_POOL } from '../../infrastructure/database/postgres.provider';
import { REDIS_CLIENT } from '../../infrastructure/cache/redis.provider';

import { PostgresSubmissionRepo } from '../../infrastructure/database/postgres/postgres-submission.repo';
import { RedisSubmissionQueue, SUBMISSION_QUEUE } from '../../infrastructure/queue/redis-submission.queue';
import { PostgresChallengeRepo } from '../../infrastructure/database/postgres/postgres-challenge.repo';

@Module({
  imports: [AuthModule], 

  controllers: [SubmissionsController],
  providers: [
    SubmissionsService,

    // Repo Submissions (PG)
    { provide: 'SubmissionRepo', useFactory: (pool) => new PostgresSubmissionRepo(pool), inject: [PG_POOL] },

    // Repo Challenges (PG) solo lectura para validar challengeId
    { provide: 'ChallengeRepoForSub', useFactory: (pool) => new PostgresChallengeRepo(pool), inject: [PG_POOL] },

    // Cola Redis
    { provide: SUBMISSION_QUEUE, useFactory: (redis) => new RedisSubmissionQueue(redis), inject: [REDIS_CLIENT] },

  {
    provide: 'ISubmissionRepo',           // alias para inyectar en el service
    useExisting: 'SubmissionRepo',
    },
    // Use cases
    {
      provide: CreateSubmissionUseCase,
      useFactory: (subRepo, queue, challengeRepo) => new CreateSubmissionUseCase(subRepo, queue, challengeRepo),
      inject: ['SubmissionRepo', SUBMISSION_QUEUE, 'ChallengeRepoForSub'],
    },
    { provide: GetSubmissionUseCase, useFactory: (repo) => new GetSubmissionUseCase(repo), inject: ['SubmissionRepo'] },
  ],
})
export class SubmissionsModule {}
