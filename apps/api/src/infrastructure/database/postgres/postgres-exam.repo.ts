import { Injectable, Inject } from '@nestjs/common';
import { Pool } from 'pg';
import { IExamRepo } from '../../../core/exams/interfaces/exam.repo';
import { Exam } from '../../../core/exams/entities/exam.entity';

@Injectable()
export class PostgresExamRepo implements IExamRepo {
    constructor(@Inject('PG_POOL') private readonly pool: Pool) { }

    async save(exam: Exam): Promise<void> {
        const query = `
      INSERT INTO exams (id, title, description, course_id, start_time, end_time, duration_minutes, created_at, updated_at)
      VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
      ON CONFLICT (id) DO UPDATE SET
        title = EXCLUDED.title,
        description = EXCLUDED.description,
        course_id = EXCLUDED.course_id,
        start_time = EXCLUDED.start_time,
        end_time = EXCLUDED.end_time,
        duration_minutes = EXCLUDED.duration_minutes,
        updated_at = EXCLUDED.updated_at
    `;
        await this.pool.query(query, [
            exam.id,
            exam.title,
            exam.description,
            exam.courseId,
            exam.startTime,
            exam.endTime,
            exam.durationMinutes,
            exam.createdAt,
            exam.updatedAt,
        ]);
    }

    async findById(id: string): Promise<Exam | null> {
        const res = await this.pool.query('SELECT * FROM exams WHERE id = $1', [id]);
        if (res.rows.length === 0) return null;
        return Exam.fromPersistence(res.rows[0]);
    }

    async findByCourseId(courseId: string): Promise<Exam[]> {
        const res = await this.pool.query('SELECT * FROM exams WHERE course_id = $1 ORDER BY created_at DESC', [courseId]);
        return res.rows.map((row) => Exam.fromPersistence(row));
    }

    async addChallengeToExam(examId: string, challengeId: string, points: number, order: number): Promise<void> {
        await this.pool.query(
            'INSERT INTO exam_challenges (exam_id, challenge_id, points, "order") VALUES ($1, $2, $3, $4)',
            [examId, challengeId, points, order],
        );
    }

    async getExamChallenges(examId: string): Promise<any[]> {
        const query = `
      SELECT c.*, ec.points, ec."order"
      FROM challenges c
      JOIN exam_challenges ec ON c.id = ec.challenge_id
      WHERE ec.exam_id = $1
      ORDER BY ec."order" ASC
    `;
        const res = await this.pool.query(query, [examId]);
        return res.rows.map(row => ({
            id: row.id,
            title: row.title,
            difficulty: row.difficulty,
            points: row.points,
            order: row.order
        }));
    }
}
