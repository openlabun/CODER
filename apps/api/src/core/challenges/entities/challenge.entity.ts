export type ChallengeStatus = 'draft' | 'published' | 'archived';
export type ChallengeDifficulty = 'easy' | 'medium' | 'hard';

export class Challenge {
  private constructor(
    public readonly id: string,
    public title: string,
    public description: string,
    public status: ChallengeStatus,
    public timeLimit: number,
    public memoryLimit: number,
    public difficulty: ChallengeDifficulty,
    public tags: string[],
    public inputFormat: string,
    public outputFormat: string,
    public constraints: string,
    public readonly createdAt: Date,
    public updatedAt: Date,
  ) { }

  static fromPersistence(row: {
    id: string;
    title: string;
    description: string;
    status: ChallengeStatus;
    time_limit?: number;
    memory_limit?: number;
    difficulty?: ChallengeDifficulty;
    tags?: string[];
    input_format?: string;
    output_format?: string;
    constraints?: string;
    created_at: Date | string;
    updated_at: Date | string;
  }) {
    return new Challenge(
      row.id,
      row.title,
      row.description,
      row.status,
      row.time_limit || 1500,
      row.memory_limit || 256,
      row.difficulty || 'medium',
      row.tags || [],
      row.input_format || '',
      row.output_format || '',
      row.constraints || '',
      new Date(row.created_at),
      new Date(row.updated_at),
    );
  }

  static create(params: {
    id: string;
    title: string;
    description: string;
    timeLimit?: number;
    memoryLimit?: number;
    difficulty?: ChallengeDifficulty;
    tags?: string[];
    inputFormat?: string;
    outputFormat?: string;
    constraints?: string;
  }) {
    if (!params.title || params.title.trim().length < 3) {
      throw new Error('Title must be at least 3 characters');
    }
    const now = new Date();
    return new Challenge(
      params.id,
      params.title.trim(),
      params.description?.trim() ?? '',
      'draft',
      params.timeLimit || 1500,
      params.memoryLimit || 256,
      params.difficulty || 'medium',
      params.tags || [],
      params.inputFormat || '',
      params.outputFormat || '',
      params.constraints || '',
      now,
      now,
    );
  }

  publish() {
    if (this.status === 'archived') throw new Error('Cannot publish archived challenge');
    this.status = 'published';
    this.updatedAt = new Date();
  }

  archive() {
    this.status = 'archived';
    this.updatedAt = new Date();
  }

  rename(newTitle: string) {
    if (!newTitle || newTitle.trim().length < 3) throw new Error('Title too short');
    this.title = newTitle.trim();
    this.updatedAt = new Date();
  }

  updateDescription(newDesc: string) {
    this.description = (newDesc ?? '').trim();
    this.updatedAt = new Date();
  }

  updateDetails(inputFormat: string, outputFormat: string, constraints: string) {
    this.inputFormat = inputFormat;
    this.outputFormat = outputFormat;
    this.constraints = constraints;
    this.updatedAt = new Date();
  }

  updateLimits(timeLimit: number, memoryLimit: number) {
    if (timeLimit < 100 || timeLimit > 10000) {
      throw new Error('Time limit must be between 100ms and 10000ms');
    }
    if (memoryLimit < 64 || memoryLimit > 1024) {
      throw new Error('Memory limit must be between 64MB and 1024MB');
    }
    this.timeLimit = timeLimit;
    this.memoryLimit = memoryLimit;
    this.updatedAt = new Date();
  }

  updateDifficulty(difficulty: ChallengeDifficulty) {
    this.difficulty = difficulty;
    this.updatedAt = new Date();
  }

  updateTags(tags: string[]) {
    this.tags = tags.filter(t => t.trim().length > 0);
    this.updatedAt = new Date();
  }
}
