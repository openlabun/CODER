import { Inject, Injectable, UnauthorizedException, ConflictException } from '@nestjs/common';
import { TOKEN_SIGNER } from './tokens';
import { TokenSigner } from './token-signer';
import { IUsersRepo } from './interfaces/users.repo';
import { User } from './entities/user.entity';
import * as bcrypt from 'bcrypt';

@Injectable()
export class AuthService {
  constructor(
    @Inject(TOKEN_SIGNER) private readonly signer: TokenSigner,
    @Inject('UsersRepo') private readonly usersRepo: IUsersRepo,
  ) { }

  async register(username: string, password: string, role: 'student' | 'professor') {
    // Check if user already exists
    const existing = await this.usersRepo.findByUsername(username);
    if (existing) {
      throw new ConflictException('Username already exists');
    }

    // Hash password
    const passwordHash = await bcrypt.hash(password, 10);

    // Create user
    const user = User.create({ username, passwordHash, role });
    await this.usersRepo.save(user);

    // Generate token
    const accessToken = this.signer.sign({ sub: user.id, username: user.username, role: user.role });
    return { accessToken };
  }

  async validateUser(username: string, password: string) {
    const user = await this.usersRepo.findByUsername(username);
    if (!user) return null;

    const isValid = await bcrypt.compare(password, user.passwordHash);
    if (!isValid) return null;

    return { sub: user.id, username: user.username, role: user.role };
  }

  async login(username: string, password: string) {
    const user = await this.validateUser(username, password);
    if (!user) throw new UnauthorizedException('Credenciales inválidas');
    const accessToken = this.signer.sign(user);
    return { accessToken };
  }

  verify(token: string) {
    return this.signer.verify(token);
  }
}
