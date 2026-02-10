import { Global, Module } from '@nestjs/common';
import { AuthController } from './auth.controller';
import { JwtAuthGuard } from './guards/jwt-auth.guard';
import { TOKEN_SIGNER } from './tokens';
import { createTokenSigner } from './token-signer';
import { AuthService } from './auth.service';
import { DatabaseModule } from '../../infrastructure/database/database.module';

@Global() // puedes dejarlo global; si lo quitas, importa AuthModule donde lo uses.
@Module({
  imports: [DatabaseModule],
  controllers: [AuthController],
  providers: [
    { provide: TOKEN_SIGNER, useFactory: () => createTokenSigner() },
    AuthService,
    JwtAuthGuard,
  ],
  exports: [
    TOKEN_SIGNER,
    AuthService,
    JwtAuthGuard,
  ],
})
export class AuthModule { }
