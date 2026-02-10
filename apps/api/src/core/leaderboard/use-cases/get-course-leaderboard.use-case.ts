import { ISubmissionRepo } from '../../Submission/interfaces/submission.repo';
import { CourseLeaderboardEntry } from '../types/leaderboard.types';

export class GetCourseLeaderboardUseCase {
    constructor(private submissionRepo: ISubmissionRepo) { }

    async execute(courseId: string): Promise<CourseLeaderboardEntry[]> {
        const results = await (this.submissionRepo as any).getBestByCourse(courseId);

        return results.map((r: any, index: number) => ({
            rank: index + 1,
            userId: r.user_id,
            username: r.username,
            totalScore: parseInt(r.total_score) || 0,
            challengesSolved: parseInt(r.challenges_solved) || 0,
            totalTimeMs: parseInt(r.total_time_ms) || 0,
        }));
    }
}
