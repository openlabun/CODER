import { IChallengeRepo } from '../interfaces/challenge.repo';

export class PublishChallengeUseCase {
  constructor(private readonly repo: IChallengeRepo) {}

  async execute(id: string) {
    const c = await this.repo.findById(id);
    if (!c) throw new Error('Challenge not found');
    c.publish();                  // regla de dominio
    await this.repo.save(c);      // persistimos cambio
    return c;
  }
}
