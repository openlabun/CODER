import { ConflictException, UnauthorizedException } from '@nestjs/common';
import * as bcrypt from 'bcrypt';
import { AuthService } from './auth.service';
import { TokenSigner } from './token-signer';
import { IUsersRepo } from './interfaces/users.repo';
import { User } from './entities/user.entity';

jest.mock('bcrypt');
const bcryptMock = bcrypt as jest.Mocked<typeof bcrypt>;

describe('AuthService', () => {
  let service: AuthService;
  let usersRepo: jest.Mocked<IUsersRepo>;
  let signer: jest.Mocked<TokenSigner>;

  beforeEach(() => {
    usersRepo = {
      save: jest.fn(),
      findByUsername: jest.fn(),
      findById: jest.fn(),
    };

    signer = {
      sign: jest.fn().mockReturnValue('mocked-jwt-token'),
      verify: jest.fn(),
    };

    service = new AuthService(signer, usersRepo);
  });

  // ─── register ────────────────────────────────────────────────────────────────

  describe('register()', () => {
    it('registra un estudiante y devuelve accessToken', async () => {
      usersRepo.findByUsername.mockResolvedValue(null);
      usersRepo.save.mockResolvedValue(undefined);
      (bcryptMock.hash as jest.Mock).mockResolvedValue('hashed-password');

      const result = await service.register('alice', 'pass123', 'student');

      expect(usersRepo.findByUsername).toHaveBeenCalledWith('alice');
      expect(bcryptMock.hash).toHaveBeenCalledWith('pass123', 10);
      expect(usersRepo.save).toHaveBeenCalledTimes(1);
      expect(signer.sign).toHaveBeenCalledWith(
        expect.objectContaining({ username: 'alice', role: 'student' }),
      );
      expect(result).toEqual({ accessToken: 'mocked-jwt-token' });
    });

    it('registra un profesor y devuelve accessToken', async () => {
      usersRepo.findByUsername.mockResolvedValue(null);
      usersRepo.save.mockResolvedValue(undefined);
      (bcryptMock.hash as jest.Mock).mockResolvedValue('hashed-password');

      const result = await service.register('bob', 'pass456', 'professor');

      expect(signer.sign).toHaveBeenCalledWith(
        expect.objectContaining({ username: 'bob', role: 'professor' }),
      );
      expect(result).toEqual({ accessToken: 'mocked-jwt-token' });
    });

    it('lanza ConflictException si el username ya existe', async () => {
      const existingUser = User.create({
        username: 'alice',
        passwordHash: 'hashed-password',
        role: 'student',
      });
      usersRepo.findByUsername.mockResolvedValue(existingUser);

      await expect(service.register('alice', 'pass123', 'student')).rejects.toThrow(
        ConflictException,
      );
      expect(usersRepo.save).not.toHaveBeenCalled();
    });
  });

  // ─── validateUser ─────────────────────────────────────────────────────────────

  describe('validateUser()', () => {
    it('devuelve payload si las credenciales son correctas', async () => {
      const user = User.create({
        username: 'alice',
        passwordHash: 'hashed-password',
        role: 'student',
      });
      usersRepo.findByUsername.mockResolvedValue(user);
      (bcryptMock.compare as jest.Mock).mockResolvedValue(true);

      const result = await service.validateUser('alice', 'pass123');

      expect(result).toEqual({ sub: user.id, username: 'alice', role: 'student' });
    });

    it('devuelve null si el usuario no existe', async () => {
      usersRepo.findByUsername.mockResolvedValue(null);

      const result = await service.validateUser('nobody', 'pass');

      expect(result).toBeNull();
    });

    it('devuelve null si la contraseña es incorrecta', async () => {
      const user = User.create({
        username: 'alice',
        passwordHash: 'hashed-password',
        role: 'student',
      });
      usersRepo.findByUsername.mockResolvedValue(user);
      (bcryptMock.compare as jest.Mock).mockResolvedValue(false);

      const result = await service.validateUser('alice', 'wrong-pass');

      expect(result).toBeNull();
    });
  });

  // ─── login ────────────────────────────────────────────────────────────────────

  describe('login()', () => {
    it('devuelve accessToken con credenciales válidas', async () => {
      const user = User.create({
        username: 'alice',
        passwordHash: 'hashed-password',
        role: 'student',
      });
      usersRepo.findByUsername.mockResolvedValue(user);
      (bcryptMock.compare as jest.Mock).mockResolvedValue(true);

      const result = await service.login('alice', 'pass123');

      expect(signer.sign).toHaveBeenCalledWith(
        expect.objectContaining({ username: 'alice', role: 'student' }),
      );
      expect(result).toEqual({ accessToken: 'mocked-jwt-token' });
    });

    it('lanza UnauthorizedException si el usuario no existe', async () => {
      usersRepo.findByUsername.mockResolvedValue(null);

      await expect(service.login('nobody', 'pass')).rejects.toThrow(UnauthorizedException);
    });

    it('lanza UnauthorizedException si la contraseña es incorrecta', async () => {
      const user = User.create({
        username: 'alice',
        passwordHash: 'hashed-password',
        role: 'student',
      });
      usersRepo.findByUsername.mockResolvedValue(user);
      (bcryptMock.compare as jest.Mock).mockResolvedValue(false);

      await expect(service.login('alice', 'wrong-pass')).rejects.toThrow(UnauthorizedException);
    });
  });
});
