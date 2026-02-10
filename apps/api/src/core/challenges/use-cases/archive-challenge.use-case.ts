import { IChallengeRepo } from '../interfaces/challenge.repo';

export class ArchiveChallengeUseCase {
  constructor(private readonly repo: IChallengeRepo) {}

  async execute(id: string) {
    const c = await this.repo.findById(id);
    if (!c) throw new Error('Challenge not found');
    c.archive();                 // regla de dominio
    await this.repo.save(c);
    return c;
  }
}
