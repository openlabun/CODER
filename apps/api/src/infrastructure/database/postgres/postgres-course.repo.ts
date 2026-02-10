import { Pool } from 'pg';
import { Course } from '../../../core/courses/entities/course.entity';
import { ICourseRepo } from '../../../core/courses/interfaces/course.repo';

export class PostgresCourseRepo implements ICourseRepo {
    constructor(private pool: Pool) { }

    async save(course: Course): Promise<void> {
        const query = `
      INSERT INTO courses (id, name, code, period, group_number, enrollment_code, professor_id, created_at, updated_at)
      VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
      ON CONFLICT (id) DO UPDATE SET
        name = EXCLUDED.name,
        code = EXCLUDED.code,
        period = EXCLUDED.period,
        group_number = EXCLUDED.group_number,
        enrollment_code = EXCLUDED.enrollment_code,
        updated_at = NOW()
    `;
        await this.pool.query(query, [
            course.id,
            course.name,
            course.code,
            course.period,
            course.groupNumber,
            course.enrollmentCode,
            course.professorId,
            course.createdAt,
            course.updatedAt,
        ]);
    }

    async update(course: Course): Promise<void> {
        return this.save(course);
    }

    async findById(id: string): Promise<Course | null> {
        const query = 'SELECT * FROM courses WHERE id = $1';
        const result = await this.pool.query(query, [id]);
        if (result.rows.length === 0) return null;
        return Course.fromPersistence(result.rows[0]);
    }

    async findAll(): Promise<Course[]> {
        const query = 'SELECT * FROM courses ORDER BY created_at DESC';
        const result = await this.pool.query(query);
        return result.rows.map((row) => Course.fromPersistence(row));
    }

    async list(): Promise<Course[]> {
        const query = 'SELECT * FROM courses ORDER BY created_at DESC';
        const result = await this.pool.query(query);
        return result.rows.map((row) => Course.fromPersistence(row));
    }

    async findByProfessor(professorId: string): Promise<Course[]> {
        const query = 'SELECT * FROM courses WHERE professor_id = $1 ORDER BY created_at DESC';
        const result = await this.pool.query(query, [professorId]);
        return result.rows.map((row) => Course.fromPersistence(row));
    }

    async findByStudent(studentId: string): Promise<Course[]> {
        const query = `
      SELECT c.* FROM courses c
      INNER JOIN course_students cs ON c.id = cs.course_id
      WHERE cs.student_id = $1
      ORDER BY c.created_at DESC
    `;
        const result = await this.pool.query(query, [studentId]);
        return result.rows.map((row) => Course.fromPersistence(row));
    }

    async addStudent(courseId: string, studentId: string): Promise<void> {
        const query = `
      INSERT INTO course_students (course_id, student_id)
      VALUES ($1, $2)
      ON CONFLICT (course_id, student_id) DO NOTHING
    `;
        await this.pool.query(query, [courseId, studentId]);
    }

    async removeStudent(courseId: string, studentId: string): Promise<void> {
        const query = 'DELETE FROM course_students WHERE course_id = $1 AND student_id = $2';
        await this.pool.query(query, [courseId, studentId]);
    }

    async getStudents(courseId: string): Promise<string[]> {
        const query = 'SELECT student_id FROM course_students WHERE course_id = $1';
        const result = await this.pool.query(query, [courseId]);
        return result.rows.map((row) => row.student_id);
    }

    async addChallenge(courseId: string, challengeId: string): Promise<void> {
        const query = `
      INSERT INTO course_challenges (course_id, challenge_id)
      VALUES ($1, $2)
      ON CONFLICT (course_id, challenge_id) DO NOTHING
    `;
        await this.pool.query(query, [courseId, challengeId]);
    }

    async removeChallenge(courseId: string, challengeId: string): Promise<void> {
        const query = 'DELETE FROM course_challenges WHERE course_id = $1 AND challenge_id = $2';
        await this.pool.query(query, [courseId, challengeId]);
    }

    async getChallenges(courseId: string): Promise<string[]> {
        const query = 'SELECT challenge_id FROM course_challenges WHERE course_id = $1';
        const result = await this.pool.query(query, [courseId]);
        return result.rows.map((row) => row.challenge_id);
    }
}
