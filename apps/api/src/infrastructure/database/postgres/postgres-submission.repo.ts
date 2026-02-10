import { Pool } from "pg";
import { Submission } from "../../../core/Submission/entities/submission.entity";
import { ISubmissionRepo } from "../../../core/Submission/interfaces/submission.repo";

export class PostgresSubmissionRepo implements ISubmissionRepo {
    constructor(private readonly pool: Pool) { }

    async save(sub: Submission): Promise<void> {
        const sql = `
      INSERT INTO public.submissions (id, challenge_id, user_id, code, language, status, score, time_ms_total, created_at, updated_at)
      VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
      ON CONFLICT (id) DO UPDATE SET
        challenge_id = EXCLUDED.challenge_id,
        user_id      = EXCLUDED.user_id,
        code         = EXCLUDED.code,
        language     = EXCLUDED.language,
        status       = EXCLUDED.status,
        score        = EXCLUDED.score,
        time_ms_total = EXCLUDED.time_ms_total,
        updated_at   = NOW()
    `;
        const vals = [
            sub.id,
            sub.challengeId,
            sub.userId,
            sub.code,
            sub.language,
            sub.status,
            sub.score,
            sub.timeMsTotal,
            sub.createdAt,
            sub.updatedAt,
        ];
        await this.pool.query(sql, vals);
    }

    async findById(id: string): Promise<Submission | null> {
        const { rows } = await this.pool.query(
            `SELECT id, challenge_id, user_id, code, language, status, score, time_ms_total, created_at, updated_at
         FROM public.submissions
        WHERE id = $1
        LIMIT 1`,
            [id],
        );
        if (rows.length === 0) return null;
        return Submission.fromPersistence(rows[0]);
    }

    async list(params: {
        challengeId?: string;
        userId?: string;
        status?: string;
        limit?: number;
        offset?: number;
    }): Promise<Submission[]> {
        const conditions: string[] = [];
        const values: any[] = [];
        let paramIndex = 1;

        if (params.challengeId) {
            conditions.push(`challenge_id = $${paramIndex++}`);
            values.push(params.challengeId);
        }
        if (params.userId) {
            conditions.push(`user_id = $${paramIndex++}`);
            values.push(params.userId);
        }
        if (params.status) {
            conditions.push(`status = $${paramIndex++}`);
            values.push(params.status);
        }

        const whereClause = conditions.length > 0 ? `WHERE ${conditions.join(' AND ')}` : '';
        const limit = params.limit ?? 100;
        const offset = params.offset ?? 0;

        const sql = `
            SELECT id, challenge_id, user_id, code, language, status, score, time_ms_total, created_at, updated_at
            FROM public.submissions
            ${whereClause}
            ORDER BY created_at DESC
            LIMIT $${paramIndex++} OFFSET $${paramIndex++}
        `;
        values.push(limit, offset);

        const { rows } = await this.pool.query(sql, values);
        return rows.map((r) => Submission.fromPersistence(r));
    }

    async getBestByChallenge(challengeId: string): Promise<any[]> {
        const query = `
            SELECT DISTINCT ON (s.user_id)
                s.user_id,
                u.username,
                s.score,
                s.time_ms_total,
                s.created_at
            FROM submissions s
            INNER JOIN users u ON s.user_id = u.id
            WHERE s.challenge_id = $1 AND s.status = 'accepted'
            ORDER BY s.user_id, s.score DESC, s.time_ms_total ASC, s.created_at ASC
        `;
        const result = await this.pool.query(query, [challengeId]);
        return result.rows;
    }

    async getBestByCourse(courseId: string): Promise<any[]> {
        const query = `
            SELECT 
                s.user_id,
                u.username,
                SUM(best.score) as total_score,
                COUNT(DISTINCT best.challenge_id) as challenges_solved,
                SUM(best.time_ms_total) as total_time_ms
            FROM (
                SELECT DISTINCT ON (s2.user_id, s2.challenge_id)
                    s2.user_id,
                    s2.challenge_id,
                    s2.score,
                    s2.time_ms_total
                FROM submissions s2
                INNER JOIN course_challenges cc ON s2.challenge_id = cc.challenge_id
                WHERE cc.course_id = $1 AND s2.status = 'accepted'
                ORDER BY s2.user_id, s2.challenge_id, s2.score DESC, s2.time_ms_total ASC
            ) best
            INNER JOIN submissions s ON best.user_id = s.user_id
            INNER JOIN users u ON s.user_id = u.id
            INNER JOIN course_students cs ON s.user_id = cs.student_id AND cs.course_id = $1
            GROUP BY s.user_id, u.username
            ORDER BY total_score DESC, challenges_solved DESC, total_time_ms ASC
        `;
        const result = await this.pool.query(query, [courseId]);
        return result.rows;
    }

    async count(params: {
        challengeId?: string;
        userId?: string;
        status?: string;
    }): Promise<number> {
        const conditions: string[] = [];
        const values: any[] = [];
        let paramIndex = 1;

        if (params.challengeId) {
            conditions.push(`challenge_id = $${paramIndex++}`);
            values.push(params.challengeId);
        }
        if (params.userId) {
            conditions.push(`user_id = $${paramIndex++}`);
            values.push(params.userId);
        }
        if (params.status) {
            conditions.push(`status = $${paramIndex++}`);
            values.push(params.status);
        }

        const whereClause = conditions.length > 0 ? `WHERE ${conditions.join(' AND ')}` : '';
        const sql = `SELECT COUNT(*) FROM submissions ${whereClause}`;

        const result = await this.pool.query(sql, values);
        return parseInt(result.rows[0].count);
    }
}
