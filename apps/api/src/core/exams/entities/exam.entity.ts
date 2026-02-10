import { randomUUID } from 'crypto';

export class Exam {
    private constructor(
        public readonly id: string,
        public title: string,
        public description: string,
        public readonly courseId: string,
        public startTime: Date,
        public endTime: Date,
        public durationMinutes: number,
        public readonly createdAt: Date,
        public updatedAt: Date,
    ) { }

    static create(props: {
        title: string;
        description: string;
        courseId: string;
        startTime: Date;
        endTime: Date;
        durationMinutes: number;
    }) {
        if (props.startTime >= props.endTime) {
            throw new Error('Start time must be before end time');
        }
        if (props.durationMinutes <= 0) {
            throw new Error('Duration must be positive');
        }

        const now = new Date();
        return new Exam(
            randomUUID(),
            props.title,
            props.description,
            props.courseId,
            props.startTime,
            props.endTime,
            props.durationMinutes,
            now,
            now,
        );
    }

    isActive(now: Date = new Date()): boolean {
        return now >= this.startTime && now <= this.endTime;
    }

    static fromPersistence(row: any) {
        return new Exam(
            row.id,
            row.title,
            row.description,
            row.course_id,
            new Date(row.start_time),
            new Date(row.end_time),
            row.duration_minutes,
            new Date(row.created_at),
            new Date(row.updated_at),
        );
    }
}
