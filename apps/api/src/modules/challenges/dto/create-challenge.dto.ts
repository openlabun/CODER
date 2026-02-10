export class CreateChallengeDto {
  title!: string;
  description!: string;
  difficulty?: string;
  timeLimit?: number;
  memoryLimit?: number;
  tags?: string[];
  inputFormat?: string;
  outputFormat?: string;
  constraints?: string;
  status?: string;
  publicTestCases?: Array<{ input: string; output: string; name: string }>;
  hiddenTestCases?: Array<{ input: string; output: string; name: string }>;
}
