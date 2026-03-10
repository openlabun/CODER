import { ApiProperty, ApiPropertyOptional } from '@nestjs/swagger';

export class ExamChallengeDto {
    @ApiProperty({ description: 'Challenge ID', example: '550e8400-e29b-41d4-a716-446655440000' })
    challengeId!: string;

    @ApiProperty({ description: 'Points for this challenge in the exam', example: 100 })
    points!: number;

    @ApiProperty({ description: 'Display order in the exam', example: 1 })
    order!: number;
}

export class CreateExamDto {
    @ApiProperty({ description: 'Exam title', example: 'Midterm Exam' })
    title!: string;

    @ApiPropertyOptional({ description: 'Exam description', example: 'Covers arrays and linked lists' })
    description?: string;

    @ApiProperty({ description: 'Course ID this exam belongs to', example: '550e8400-e29b-41d4-a716-446655440000' })
    courseId!: string;

    @ApiProperty({ description: 'Exam start time', example: '2026-03-15T09:00:00Z' })
    startTime!: string;

    @ApiProperty({ description: 'Exam end time', example: '2026-03-15T11:00:00Z' })
    endTime!: string;

    @ApiProperty({ description: 'Duration in minutes', example: 120 })
    durationMinutes!: number;

    @ApiProperty({ description: 'Challenges included in the exam', type: [ExamChallengeDto] })
    challenges!: ExamChallengeDto[];
}
