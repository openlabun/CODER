import { IChallengeRepo } from '../interfaces/challenge.repo';
import { Challenge } from '../entities/challenge.entity';

export class GetChallengeUseCase {
  constructor(private readonly repo: IChallengeRepo) {}

  async execute(id: string): Promise<Challenge | null> {
    return this.repo.findById(id);
  }
}
