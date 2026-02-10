import { Module, Global } from '@nestjs/common';
import { Pool } from 'pg';
import { PG_POOL, createPgPool } from './postgres.provider';
import { PostgresChallengeRepo } from './postgres/postgres-challenge.repo';
import { PostgresSubmissionRepo } from './postgres/postgres-submission.repo';
import { PostgresUsersRepo } from './postgres/postgres-users.repo';
import { PostgresTestCaseRepo } from './postgres/postgres-test-case.repo';
import { PostgresCourseRepo } from './postgres/postgres-course.repo';

@Global()
@Module({
  providers: [
    {
      provide: PG_POOL,
      useFactory: () => createPgPool(),
    },
    {
      provide: 'ChallengeRepo',
      useFactory: (pool: Pool) => new PostgresChallengeRepo(pool),
      inject: [PG_POOL],
    },
    {
      provide: 'SubmissionRepo',
      useFactory: (pool: Pool) => new PostgresSubmissionRepo(pool),
      inject: [PG_POOL],
    },
    {
      provide: 'UsersRepo',
      useFactory: (pool: Pool) => new PostgresUsersRepo(pool),
      inject: [PG_POOL],
    },
    {
      provide: 'TestCaseRepo',
      useFactory: (pool: Pool) => new PostgresTestCaseRepo(pool),
      inject: [PG_POOL],
    },
    {
      provide: 'CourseRepo',
      useFactory: (pool: Pool) => new PostgresCourseRepo(pool),
      inject: [PG_POOL],
    },
  ],
  exports: [PG_POOL, 'ChallengeRepo', 'SubmissionRepo', 'UsersRepo', 'TestCaseRepo', 'CourseRepo'],
})
export class DatabaseModule { }


