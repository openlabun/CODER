import { User } from '../entities/user.entity';

export interface IUsersRepo {
    save(user: User): Promise<void>;
    findByUsername(username: string): Promise<User | null>;
    findById(id: string): Promise<User | null>;
}
