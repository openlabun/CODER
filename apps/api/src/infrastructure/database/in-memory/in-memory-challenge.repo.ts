import { IChallengeRepo } from '../../../core/challenges/interfaces/challenge.repo';
import { Challenge } from '../../../core/challenges/entities/challenge.entity';

export class InMemoryChallengeRepo implements IChallengeRepo {
  private store = new Map<string, Challenge>();

  async save(challenge: Challenge): Promise<void> {
    this.store.set(challenge.id, challenge);
  }

  async findById(id: string): Promise<Challenge | null> {
    return this.store.get(id) ?? null;
  }

  async list(): Promise<Challenge[]> {
    return Array.from(this.store.values())
      .sort((a, b) => b.createdAt.getTime() - a.createdAt.getTime());
  }
}
