import { IChallengeRepo } from '../interfaces/challenge.repo';
import { Challenge, ChallengeDifficulty } from '../entities/challenge.entity';

export class CreateChallengeUseCase {
  constructor(private readonly repo: IChallengeRepo) { }
  async execute(input: {
    id: string;
    title: string;
    description: string;
    difficulty?: ChallengeDifficulty;
    timeLimit?: number;
    memoryLimit?: number;
    tags?: string[];
  }) {
    const entity = Challenge.create(input);
    await this.repo.save(entity);
    return entity;
  }
}
