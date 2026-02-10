import { Pool } from 'pg';
import { User } from '../../../modules/auth/entities/user.entity';
import { IUsersRepo } from '../../../modules/auth/interfaces/users.repo';

export class PostgresUsersRepo implements IUsersRepo {
    constructor(private readonly pool: Pool) { }

    async save(user: User): Promise<void> {
        const sql = `
      INSERT INTO public.users (id, username, password, role, created_at, updated_at)
      VALUES ($1, $2, $3, $4, $5, $6)
      ON CONFLICT (id) DO UPDATE SET
        username   = EXCLUDED.username,
        password   = EXCLUDED.password,
        role       = EXCLUDED.role,
        updated_at = NOW()
    `;
        const vals = [
            user.id,
            user.username,
            user.passwordHash,
            user.role,
            user.createdAt,
            user.updatedAt,
        ];
        await this.pool.query(sql, vals);
    }

    async findByUsername(username: string): Promise<User | null> {
        const { rows } = await this.pool.query(
            `SELECT id, username, password, role, created_at, updated_at
       FROM public.users
       WHERE username = $1
       LIMIT 1`,
            [username],
        );
        if (!rows.length) return null;
        return User.fromPersistence(rows[0]);
    }

    async findById(id: string): Promise<User | null> {
        const { rows } = await this.pool.query(
            `SELECT id, username, password, role, created_at, updated_at
       FROM public.users
       WHERE id = $1
       LIMIT 1`,
            [id],
        );
        if (!rows.length) return null;
        return User.fromPersistence(rows[0]);
    }
}
