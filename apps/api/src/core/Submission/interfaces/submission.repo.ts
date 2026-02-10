import { Submission } from '../entities/submission.entity';

export interface ISubmissionRepo {
  save(sub: Submission): Promise<void>;
  findById(id: string): Promise<Submission | null>;

  // 👇 Nuevos métodos para listados/paginación
  list(params: {
    challengeId?: string;
    userId?: string;
    status?: string;
    limit?: number;
    offset?: number;
  }): Promise<Submission[]>;

  count(params: {
    challengeId?: string;
    userId?: string;
    status?: string;
  }): Promise<number>;
}
