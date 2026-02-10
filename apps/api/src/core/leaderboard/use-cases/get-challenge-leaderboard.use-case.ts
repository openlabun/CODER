import { ISubmissionRepo } from '../../Submission/interfaces/submission.repo';
import { LeaderboardEntry } from '../types/leaderboard.types';

export class GetChallengeLeaderboardUseCase {
    constructor(private submissionRepo: ISubmissionRepo) { }

    async execute(challengeId: string): Promise<LeaderboardEntry[]> {
        const results = await (this.submissionRepo as any).getBestByChallenge(challengeId);

        return results
            .sort((a: any, b: any) => {
                if (b.score !== a.score) return b.score - a.score;
                if (a.time_ms_total !== b.time_ms_total) return a.time_ms_total - b.time_ms_total;
                return new Date(a.created_at).getTime() - new Date(b.created_at).getTime();
            })
            .map((r: any, index: number) => ({
                rank: index + 1,
                userId: r.user_id,
                username: r.username,
                score: r.score,
                timeMs: r.time_ms_total,
                submittedAt: new Date(r.created_at),
            }));
    }
}
