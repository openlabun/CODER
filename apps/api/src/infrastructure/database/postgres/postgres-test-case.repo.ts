import { Pool } from 'pg';
import { TestCase } from '../../../core/challenges/entities/test-case.entity';
import { ITestCaseRepo } from '../../../core/challenges/interfaces/test-case.repo';

export class PostgresTestCaseRepo implements ITestCaseRepo {
    constructor(private pool: Pool) { }

    async save(testCase: TestCase): Promise<void> {
        const query = `
      INSERT INTO test_cases (id, challenge_id, name, input, expected_output, is_sample, points, created_at)
      VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
      ON CONFLICT (id) DO UPDATE SET
        name = EXCLUDED.name,
        input = EXCLUDED.input,
        expected_output = EXCLUDED.expected_output,
        is_sample = EXCLUDED.is_sample,
        points = EXCLUDED.points
    `;
        await this.pool.query(query, [
            testCase.id,
            testCase.challengeId,
            testCase.name,
            testCase.input,
            testCase.expectedOutput,
            testCase.isSample,
            testCase.points,
            testCase.createdAt,
        ]);
    }

    async findByChallengeId(challengeId: string): Promise<TestCase[]> {
        const query = `
      SELECT * FROM test_cases 
      WHERE challenge_id = $1 
      ORDER BY is_sample DESC, name ASC
    `;
        const result = await this.pool.query(query, [challengeId]);
        return result.rows.map((row) => TestCase.fromPersistence(row));
    }

    async findSamplesByChallengeId(challengeId: string): Promise<TestCase[]> {
        const query = `
      SELECT * FROM test_cases 
      WHERE challenge_id = $1 AND is_sample = true
      ORDER BY name ASC
    `;
        const result = await this.pool.query(query, [challengeId]);
        return result.rows.map((row) => TestCase.fromPersistence(row));
    }

    async findById(id: string): Promise<TestCase | null> {
        const query = 'SELECT * FROM test_cases WHERE id = $1';
        const result = await this.pool.query(query, [id]);
        if (result.rows.length === 0) return null;
        return TestCase.fromPersistence(result.rows[0]);
    }

    async deleteById(id: string): Promise<void> {
        const query = 'DELETE FROM test_cases WHERE id = $1';
        await this.pool.query(query, [id]);
    }
}
