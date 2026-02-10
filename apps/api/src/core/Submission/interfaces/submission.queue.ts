export interface ISubmissionQueue {
  enqueue(submissionId: string): Promise<void>;
}
