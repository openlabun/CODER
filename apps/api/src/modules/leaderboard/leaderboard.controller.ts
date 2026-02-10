import { Controller, Get, Param } from '@nestjs/common';
import { ApiTags, ApiOperation, ApiResponse, ApiParam } from '@nestjs/swagger';
import { GetChallengeLeaderboardUseCase } from '../../core/leaderboard/use-cases/get-challenge-leaderboard.use-case';
import { GetCourseLeaderboardUseCase } from '../../core/leaderboard/use-cases/get-course-leaderboard.use-case';

@ApiTags('leaderboard')
@Controller('leaderboard')
export class LeaderboardController {
    constructor(
        private getChallengeLeaderboard: GetChallengeLeaderboardUseCase,
        private getCourseLeaderboard: GetCourseLeaderboardUseCase,
    ) { }

    @Get('challenge/:id')
    @ApiOperation({ summary: 'Get leaderboard for a specific challenge' })
    @ApiParam({ name: 'id', description: 'Challenge ID' })
    @ApiResponse({ status: 200, description: 'Returns challenge leaderboard with rankings' })
    async getChallenge(@Param('id') id: string) {
        const entries = await this.getChallengeLeaderboard.execute(id);
        return {
            challengeId: id,
            entries,
        };
    }

    @Get('course/:id')
    @ApiOperation({ summary: 'Get leaderboard for a specific course' })
    @ApiParam({ name: 'id', description: 'Course ID' })
    @ApiResponse({ status: 200, description: 'Returns course leaderboard with total scores' })
    async getCourse(@Param('id') id: string) {
        const entries = await this.getCourseLeaderboard.execute(id);
        return {
            courseId: id,
            entries,
        };
    }
}

