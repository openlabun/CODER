import { ApiProperty } from '@nestjs/swagger';

export class CreateCourseDto {
    @ApiProperty({ description: 'Course name', example: 'Data Structures' })
    name!: string;

    @ApiProperty({ description: 'Course code', example: 'CS201' })
    code!: string;

    @ApiProperty({ description: 'Academic period', example: '2026-1' })
    period!: string;

    @ApiProperty({ description: 'Group number', example: 1 })
    groupNumber!: number;
}

export class EnrollByCodeDto {
    @ApiProperty({ description: 'Enrollment code to join a course', example: 'CS201-20261G1' })
    enrollmentCode!: string;
}

export class EnrollStudentDto {
    @ApiProperty({ description: 'Student user ID to enroll', example: '550e8400-e29b-41d4-a716-446655440000' })
    studentId!: string;
}

export class AssignChallengeDto {
    @ApiProperty({ description: 'Challenge ID to assign to the course', example: '550e8400-e29b-41d4-a716-446655440000' })
    challengeId!: string;
}
