import { Injectable } from '@nestjs/common';
import { CreateChallengeUseCase } from '../../core/challenges/use-cases/create-challenge.use-case';
import { ListChallengesUseCase } from '../../core/challenges/use-cases/list-challenges.use-case';
import { GetChallengeUseCase } from '../../core/challenges/use-cases/get-challenge.use-case';
import { PublishChallengeUseCase } from '../../core/challenges/use-cases/publish-challenge.use-case';
import { ArchiveChallengeUseCase } from '../../core/challenges/use-cases/archive-challenge.use-case';

@Injectable()
export class ChallengesService {
  constructor(
    private readonly createUC: CreateChallengeUseCase,
    private readonly listUC: ListChallengesUseCase,
    private readonly getUC: GetChallengeUseCase,
    private readonly publishUC: PublishChallengeUseCase,
    private readonly archiveUC: ArchiveChallengeUseCase,
  ) { }

  create(input: {
    id: string;
    title: string;
    description: string;
    difficulty?: any;
    timeLimit?: number;
    memoryLimit?: number;
    tags?: string[];
    inputFormat?: string;
    outputFormat?: string;
    constraints?: string;
  }) {
    return this.createUC.execute(input);
  }
  list() {
    return this.listUC.execute();
  }
  get(id: string) {
    return this.getUC.execute(id);
  }
  publish(id: string) {
    return this.publishUC.execute(id);
  }
  archive(id: string) {
    return this.archiveUC.execute(id);
  }
}
