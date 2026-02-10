import type { Redis } from 'ioredis';
import { ISubmissionQueue } from '../../core/Submission/interfaces/submission.queue';

export const SUBMISSION_QUEUE = 'SUBMISSION_QUEUE';

export class RedisSubmissionQueue implements ISubmissionQueue {
  constructor(private readonly redis: Redis, private readonly queueKey = 'queue:submissions') {}

  async enqueue(submissionId: string): Promise<void> {
    // LPUSH para añadir a la cola (LPOP/RPOP la consumirá el worker)
    await this.redis.lpush(this.queueKey, submissionId);
  }
}
