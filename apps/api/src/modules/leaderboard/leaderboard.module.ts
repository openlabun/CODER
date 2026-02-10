import { Module } from '@nestjs/common';
import { LeaderboardController } from './leaderboard.controller';
import { GetChallengeLeaderboardUseCase } from '../../core/leaderboard/use-cases/get-challenge-leaderboard.use-case';
import { GetCourseLeaderboardUseCase } from '../../core/leaderboard/use-cases/get-course-leaderboard.use-case';
import { DatabaseModule } from '../../infrastructure/database/database.module';

@Module({
    imports: [DatabaseModule],
    controllers: [LeaderboardController],
    providers: [
        {
            provide: GetChallengeLeaderboardUseCase,
            useFactory: (repo) => new GetChallengeLeaderboardUseCase(repo),
            inject: ['SubmissionRepo'],
        },
        {
            provide: GetCourseLeaderboardUseCase,
            useFactory: (repo) => new GetCourseLeaderboardUseCase(repo),
            inject: ['SubmissionRepo'],
        },
    ],
})
export class LeaderboardModule { }
