import { Challenge } from '../entities/challenge.entity';

export interface IChallengeRepo {
  save(challenge: Challenge): Promise<void>;
  findById(id: string): Promise<Challenge | null>;
  list(): Promise<Challenge[]>;
}
