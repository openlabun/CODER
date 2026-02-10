import { IChallengeRepo } from '../interfaces/challenge.repo';

export class ListChallengesUseCase {
  constructor(private readonly repo: IChallengeRepo) {}
  async execute() {
    return this.repo.list();
  }
}
