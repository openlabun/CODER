import { Inject, Injectable, BadRequestException } from '@nestjs/common';
import { CreateSubmissionUseCase } from '../../core/Submission/use-cases/create-submission.use-case';
import { GetSubmissionUseCase } from '../../core/Submission/use-cases/get-submission.use-case';
import { ISubmissionRepo } from '../../core/Submission/interfaces/submission.repo';

@Injectable()
export class SubmissionsService {
  constructor(
    private readonly createUC: CreateSubmissionUseCase,
    private readonly getUC: GetSubmissionUseCase,
    @Inject('SubmissionRepo') private readonly repo: ISubmissionRepo, // 👈 inyectamos el repo ya registrado
  ) { }

  async create(input: { challengeId: string; userId: string; code: string; language: string; examId?: string }) {
    try {
      return await this.createUC.execute(input);
    } catch (e: any) {
      if (e?.message === 'challenge_not_found') {
        throw new BadRequestException('Invalid challengeId');
      }
      throw e;
    }
  }

  get(id: string) {
    return this.getUC.execute(id);
  }

  async list(params: { challengeId?: string; userId?: string; status?: string; limit?: number; offset?: number }) {
    const [items, total] = await Promise.all([
      this.repo.list(params),
      this.repo.count(params),
    ]);
    return {
      total,
      limit: params.limit ?? 20,
      offset: params.offset ?? 0,
      items: items.map(s => ({
        id: s.id,
        challengeId: s.challengeId,
        userId: s.userId,
        language: s.language,
        status: s.status,
        score: s.score,
        timeMsTotal: s.timeMsTotal,
        createdAt: s.createdAt,
        updatedAt: s.updatedAt,
      })),
    };
  }
}
