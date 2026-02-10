import { randomUUID } from 'crypto';

export class Course {
    private constructor(
        public readonly id: string,
        public name: string,
        public code: string,
        public period: string,
        public groupNumber: number,
        public enrollmentCode: string | null,
        public readonly professorId: string,
        public readonly createdAt: Date,
        public updatedAt: Date,
    ) { }

    static create(params: {
        name: string;
        code: string;
        period: string;
        groupNumber: number;
        professorId: string;
        enrollmentCode?: string;
    }) {
        if (!params.name || params.name.trim().length < 3) {
            throw new Error('Course name must be at least 3 characters');
        }
        if (!params.code || params.code.trim().length < 2) {
            throw new Error('Course code is required');
        }
        if (!params.period || params.period.trim().length < 4) {
            throw new Error('Period is required (e.g., 2025-1)');
        }
        if (params.groupNumber < 1) {
            throw new Error('Group number must be at least 1');
        }

        const now = new Date();
        return new Course(
            randomUUID(),
            params.name.trim(),
            params.code.trim(),
            params.period.trim(),
            params.groupNumber,
            params.enrollmentCode || null,
            params.professorId,
            now,
            now,
        );
    }

    static fromPersistence(row: {
        id: string;
        name: string;
        code: string;
        period: string;
        group_number: number;
        enrollment_code?: string | null;
        professor_id: string;
        created_at: Date | string;
        updated_at: Date | string;
    }) {
        return new Course(
            row.id,
            row.name,
            row.code,
            row.period,
            row.group_number,
            row.enrollment_code || null,
            row.professor_id,
            new Date(row.created_at),
            new Date(row.updated_at),
        );
    }

    updateInfo(name: string, code: string, period: string, groupNumber: number) {
        if (name && name.trim().length >= 3) {
            this.name = name.trim();
        }
        if (code && code.trim().length >= 2) {
            this.code = code.trim();
        }
        if (period && period.trim().length >= 4) {
            this.period = period.trim();
        }
        if (groupNumber >= 1) {
            this.groupNumber = groupNumber;
        }
        this.updatedAt = new Date();
    }
}
