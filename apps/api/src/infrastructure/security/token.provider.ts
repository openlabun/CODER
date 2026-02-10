import jwt, { type SignOptions, type Secret, type JwtPayload as JwtPayloadType } from 'jsonwebtoken';

export const TOKEN_SIGNER = 'TOKEN_SIGNER';
export type JwtPayload = JwtPayloadType | Record<string, any>;

export function createTokenSigner() {
  const secret: Secret = (process.env.JWT_SECRET || 'dev-secret') as Secret;

  // Lee el env y tipa 'expiresIn' como lo que jsonwebtoken espera
  const raw = process.env.JWT_EXPIRES_IN ?? '1h';
  const expiresIn: SignOptions['expiresIn'] =
    /^\d+$/.test(raw) ? Number(raw) : (raw as any); // "3600" -> 3600; "1h" -> "1h"

  const options: SignOptions = { expiresIn };

  function sign(payload: JwtPayload) {
    return jwt.sign(payload as object, secret, options);
  }

  function verify<T = JwtPayload>(token: string) {
    return jwt.verify(token, secret) as T;
  }

  return { sign, verify };
}
