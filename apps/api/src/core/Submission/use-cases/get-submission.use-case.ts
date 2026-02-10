import { Submission } from '../entities/submission.entity';
import { ISubmissionRepo } from '../interfaces/submission.repo';

type Output = Submission | null;

export class GetSubmissionUseCase {
  constructor(private readonly repo: ISubmissionRepo) {}

  async execute(id: string): Promise<Output> {
    if (!id) throw new Error('id is required');
    return this.repo.findById(id);
  }
}
