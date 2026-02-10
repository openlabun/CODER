import { randomUUID } from 'crypto';

export type SubmissionStatus =
  | 'queued'
  | 'running'
  | 'accepted'
  | 'wrong_answer'
  | 'error';

export class Submission {
  constructor(
    public readonly id: string,
    public readonly challengeId: string,
    public readonly userId: string,
    public readonly code: string,
    public readonly language: string,
    public status: SubmissionStatus,
    public score: number,
    public timeMsTotal: number,
    public readonly createdAt: Date,
    public updatedAt: Date,
    public readonly examId?: string,
  ) { }

  static create(props: { challengeId: string; userId: string; code: string; language: string; examId?: string }) {
    const now = new Date();
    return new Submission(
      randomUUID(),
      props.challengeId,
      props.userId,
      props.code,
      props.language,
      'queued',
      0,
      0,
      now,
      now,
      props.examId,
    );
  }

  start() {
    this.status = 'running';
    this.updatedAt = new Date();
  }

  accept() {
    this.status = 'accepted';
    this.updatedAt = new Date();
  }

  reject() {
    this.status = 'wrong_answer';
    this.updatedAt = new Date();
  }

  fail() {
    this.status = 'error';
    this.updatedAt = new Date();
  }

  updateScore(score: number, timeMs: number) {
    this.score = score;
    this.timeMsTotal = timeMs;
    this.updatedAt = new Date();
  }

  static fromPersistence(row: any) {
    return new Submission(
      row.id,
      row.challenge_id,
      row.user_id,
      row.code,
      row.language,
      row.status,
      row.score || 0,
      row.time_ms_total || 0,
      new Date(row.created_at),
      new Date(row.updated_at),
      row.exam_id,
    );
  }
}
