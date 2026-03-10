import { ApiProperty, ApiPropertyOptional } from '@nestjs/swagger';

export class CreateSubmissionDto {
  @ApiProperty({ description: 'ID of the challenge to submit to', example: '550e8400-e29b-41d4-a716-446655440000' })
  challengeId!: string;

  @ApiProperty({ description: 'Source code to evaluate', example: 'print(input())' })
  code!: string;

  @ApiProperty({ description: 'Programming language', enum: ['python', 'node', 'cpp', 'java'], example: 'python' })
  language!: string;

  @ApiPropertyOptional({ description: 'Exam ID if submission is part of an exam', example: '550e8400-e29b-41d4-a716-446655440001' })
  examId?: string;
}
