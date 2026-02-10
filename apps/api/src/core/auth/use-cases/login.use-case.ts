import { UserEntity } from '../entities/user.entity';

export class LoginUseCase {
  // simular base de datos temporal
  private users: UserEntity[] = [
    new UserEntity('1', 'admin', 'admin123'), // luego será un hash
  ];

  execute(username: string, password: string): UserEntity | null {
    const user = this.users.find(
      (u) => u.username === username && u.passwordHash === password,
    );
    return user || null;
  }
}
