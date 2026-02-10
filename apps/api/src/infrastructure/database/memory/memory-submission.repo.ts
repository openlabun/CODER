import { ISubmissionRepo } from '../../../core/Submission/interfaces/submission.repo';
import { Submission } from '../../../core/Submission/entities/submission.entity';

export class InMemorySubmissionRepo implements ISubmissionRepo {
  private submissions = new Map<string, Submission>();

  async save(sub: Submission): Promise<void> {
    this.submissions.set(sub.id, sub);
  }

  async findById(id: string): Promise<Submission | null> {
    return this.submissions.get(id) ?? null;
  }

  async list(params: {
    challengeId?: string;
    userId?: string;
    status?: string;
    limit?: number;
    offset?: number;
  }): Promise<Submission[]> {
    const limit = Number.isFinite(params.limit) ? Number(params.limit) : 20;
    const offset = Number.isFinite(params.offset) ? Number(params.offset) : 0;

    // Filtrar en memoria
    let items = Array.from(this.submissions.values());

    if (params.challengeId) {
      items = items.filter(s => s.challengeId === params.challengeId);
    }
    if (params.userId) {
      items = items.filter(s => s.userId === params.userId);
    }
    if (params.status) {
      items = items.filter(s => s.status === params.status);
    }

    // Orden por fecha desc (como en PG)
    items.sort((a, b) => b.createdAt.getTime() - a.createdAt.getTime());

    // Paginación
    return items.slice(offset, offset + limit);
  }

  async count(params: { challengeId?: string; userId?: string; status?: string }): Promise<number> {
    let items = Array.from(this.submissions.values());

    if (params.challengeId) {
      items = items.filter(s => s.challengeId === params.challengeId);
    }
    if (params.userId) {
      items = items.filter(s => s.userId === params.userId);
    }
    if (params.status) {
      items = items.filter(s => s.status === params.status);
    }

    return items.length;
  }
}
