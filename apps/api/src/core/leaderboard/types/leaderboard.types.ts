export interface LeaderboardEntry {
    rank: number;
    userId: string;
    username: string;
    score: number;
    timeMs?: number;
    submittedAt?: Date;
}

export interface CourseLeaderboardEntry {
    rank: number;
    userId: string;
    username: string;
    totalScore: number;
    challengesSolved: number;
    totalTimeMs: number;
}
