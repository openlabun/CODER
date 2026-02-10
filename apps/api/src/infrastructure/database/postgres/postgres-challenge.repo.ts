import { Pool } from 'pg';
import { IChallengeRepo } from '../../../core/challenges/interfaces/challenge.repo';
import { Challenge } from '../../../core/challenges/entities/challenge.entity';

export class PostgresChallengeRepo implements IChallengeRepo {
  constructor(private readonly pool: Pool) { }

  async save(ch: Challenge): Promise<void> {
    const sql = `
      INSERT INTO public.challenges (id, title, description, status, time_limit, memory_limit, difficulty, tags, input_format, output_format, constraints, created_at, updated_at)
      VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
      ON CONFLICT (id) DO UPDATE SET
        title = EXCLUDED.title,
        description = EXCLUDED.description,
        status = EXCLUDED.status,
        time_limit = EXCLUDED.time_limit,
        memory_limit = EXCLUDED.memory_limit,
        difficulty = EXCLUDED.difficulty,
        tags = EXCLUDED.tags,
        input_format = EXCLUDED.input_format,
        output_format = EXCLUDED.output_format,
        constraints = EXCLUDED.constraints,
        updated_at = NOW()
    `;
    const values = [
      ch.id,
      ch.title,
      ch.description,
      ch.status,
      ch.timeLimit,
      ch.memoryLimit,
      ch.difficulty,
      ch.tags,
      ch.inputFormat,
      ch.outputFormat,
      ch.constraints,
      ch.createdAt,
      ch.updatedAt
    ];
    await this.pool.query(sql, values);
  }

  async findById(id: string): Promise<Challenge | null> {
    const { rows } = await this.pool.query(
      `SELECT id, title, description, status, time_limit, memory_limit, difficulty, tags, input_format, output_format, constraints, created_at, updated_at
       FROM public.challenges WHERE id = $1 LIMIT 1`, [id],
    );
    if (rows.length === 0) return null;
    return Challenge.fromPersistence(rows[0]);
  }

  async list(): Promise<Challenge[]> {
    const { rows } = await this.pool.query(
      `SELECT id, title, description, status, time_limit, memory_limit, difficulty, tags, input_format, output_format, constraints, created_at, updated_at
       FROM public.challenges ORDER BY created_at DESC`,
    );
    return rows.map(Challenge.fromPersistence);
  }
}
