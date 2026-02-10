export class UserEntity {
  constructor(
    public readonly id: string,
    public readonly username: string,
    public readonly passwordHash: string, // no guardamos contraseñas planas
  ) {}
}
