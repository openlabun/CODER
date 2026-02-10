import { randomUUID } from 'crypto';

export type UserRole = 'student' | 'professor';

export class User {
    constructor(
        public readonly id: string,
        public readonly username: string,
        public readonly passwordHash: string,
        public readonly role: UserRole,
        public readonly createdAt: Date,
        public updatedAt: Date,
    ) { }

    static create(props: { username: string; passwordHash: string; role: UserRole }) {
        const now = new Date();
        return new User(
            randomUUID(),
            props.username,
            props.passwordHash,
            props.role,
            now,
            now,
        );
    }

    static fromPersistence(row: any) {
        return new User(
            row.id,
            row.username,
            row.password,
            row.role,
            new Date(row.created_at),
            new Date(row.updated_at),
        );
    }
}
