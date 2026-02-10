import { Module } from '@nestjs/common';
import { ChallengesController } from './challenges.controller';
import { TestCasesController } from './test-cases.controller';
import { ChallengesService } from './challenges.service';
import { CreateChallengeUseCase } from '../../core/challenges/use-cases/create-challenge.use-case';
import { ListChallengesUseCase } from '../../core/challenges/use-cases/list-challenges.use-case';
import { GetChallengeUseCase } from '../../core/challenges/use-cases/get-challenge.use-case';
import { PublishChallengeUseCase } from '../../core/challenges/use-cases/publish-challenge.use-case';
import { ArchiveChallengeUseCase } from '../../core/challenges/use-cases/archive-challenge.use-case';
import { CreateTestCaseUseCase } from '../../core/challenges/use-cases/create-test-case.use-case';
import { ListTestCasesUseCase } from '../../core/challenges/use-cases/list-test-cases.use-case';
import { DeleteTestCaseUseCase } from '../../core/challenges/use-cases/delete-test-case.use-case';
import { PG_POOL } from '../../infrastructure/database/postgres.provider';
import { PostgresChallengeRepo } from '../../infrastructure/database/postgres/postgres-challenge.repo';
import { DatabaseModule } from '../../infrastructure/database/database.module';

@Module({
  imports: [DatabaseModule],
  controllers: [ChallengesController, TestCasesController],
  providers: [
    ChallengesService,
    {
      provide: 'ChallengeRepo',
      useFactory: (pool) => new PostgresChallengeRepo(pool),
      inject: [PG_POOL],
    },
    {
      provide: CreateChallengeUseCase,
      useFactory: (repo) => new CreateChallengeUseCase(repo),
      inject: ['ChallengeRepo'],
    },
    {
      provide: ListChallengesUseCase,
      useFactory: (repo) => new ListChallengesUseCase(repo),
      inject: ['ChallengeRepo'],
    },
    {
      provide: GetChallengeUseCase,
      useFactory: (repo) => new GetChallengeUseCase(repo),
      inject: ['ChallengeRepo'],
    },
    {
      provide: PublishChallengeUseCase,
      useFactory: (repo) => new PublishChallengeUseCase(repo),
      inject: ['ChallengeRepo'],
    },
    {
      provide: ArchiveChallengeUseCase,
      useFactory: (repo) => new ArchiveChallengeUseCase(repo),
      inject: ['ChallengeRepo'],
    },
    {
      provide: CreateTestCaseUseCase,
      useFactory: (repo) => new CreateTestCaseUseCase(repo),
      inject: ['TestCaseRepo'],
    },
    {
      provide: ListTestCasesUseCase,
      useFactory: (repo) => new ListTestCasesUseCase(repo),
      inject: ['TestCaseRepo'],
    },
    {
      provide: DeleteTestCaseUseCase,
      useFactory: (repo) => new DeleteTestCaseUseCase(repo),
      inject: ['TestCaseRepo'],
    },
  ],
})
export class ChallengesModule { }

