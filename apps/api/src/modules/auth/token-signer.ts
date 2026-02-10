export interface TokenSigner {
  sign(payload: any): string;
  verify(token: string): any;
}

export function createTokenSigner(): TokenSigner {
  const jwt = require('jsonwebtoken');
  const secret = process.env.JWT_SECRET ?? 'change-me';
  const expiresIn = process.env.JWT_EXPIRES_IN ?? '1h';
  return {
    sign(payload: any) {
      return jwt.sign(payload, secret, { expiresIn });
    },
    verify(token: string) {
      return jwt.verify(token, secret);
    },
  };
}
