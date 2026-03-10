import { ApiProperty, ApiPropertyOptional } from '@nestjs/swagger';

export class TestCaseInputDto {
  @ApiProperty({ description: 'Test case name', example: 'basic_test_1' })
  name!: string;

  @ApiProperty({ description: 'Input data for the test case', example: '5\n3' })
  input!: string;

  @ApiProperty({ description: 'Expected output', example: '8' })
  output!: string;
}

export class CreateChallengeDto {
  @ApiProperty({ description: 'Challenge title', example: 'Sum of Two Numbers' })
  title!: string;

  @ApiProperty({ description: 'Problem description in markdown', example: 'Given two integers, return their sum.' })
  description!: string;

  @ApiPropertyOptional({ description: 'Difficulty level', enum: ['easy', 'medium', 'hard'], example: 'easy' })
  difficulty?: string;

  @ApiPropertyOptional({ description: 'Time limit in milliseconds', example: 1500 })
  timeLimit?: number;

  @ApiPropertyOptional({ description: 'Memory limit in MB', example: 256 })
  memoryLimit?: number;

  @ApiPropertyOptional({ description: 'Tags for categorization', example: ['math', 'beginner'], type: [String] })
  tags?: string[];

  @ApiPropertyOptional({ description: 'Expected input format description', example: 'Two integers separated by newline' })
  inputFormat?: string;

  @ApiPropertyOptional({ description: 'Expected output format description', example: 'A single integer' })
  outputFormat?: string;

  @ApiPropertyOptional({ description: 'Problem constraints', example: '1 <= N <= 10^6' })
  constraints?: string;

  @ApiPropertyOptional({ description: 'Challenge status', enum: ['draft', 'published', 'archived'] })
  status?: string;

  @ApiPropertyOptional({ description: 'Public (sample) test cases visible to students', type: [TestCaseInputDto] })
  publicTestCases?: Array<{ input: string; output: string; name: string }>;

  @ApiPropertyOptional({ description: 'Hidden test cases used for grading', type: [TestCaseInputDto] })
  hiddenTestCases?: Array<{ input: string; output: string; name: string }>;
}
