import {
  Injectable, CanActivate, ExecutionContext, UnauthorizedException, Inject,
} from '@nestjs/common';
import { TOKEN_SIGNER } from '../tokens';
import { TokenSigner } from '../token-signer';

@Injectable()
export class JwtAuthGuard implements CanActivate {
  constructor(@Inject(TOKEN_SIGNER) private readonly signer: TokenSigner) {} // 👈 clave

  canActivate(context: ExecutionContext): boolean {
    const req = context.switchToHttp().getRequest();
    const auth = req.headers['authorization'] || '';
    const header = Array.isArray(auth) ? auth[0] : auth;
    const m = header.match(/^Bearer\s+(.+)$/i);
    if (!m) throw new UnauthorizedException('Missing Bearer token');

    try {
      const payload = this.signer.verify(m[1]);
      req.user = payload;
      return true;
    } catch {
      throw new UnauthorizedException('Invalid token');
    }
  }
}
