import { ApiProperty, ApiPropertyOptional } from '@nestjs/swagger';

export class GenerateChallengeIdeasDto {
    @ApiProperty({ description: 'Topic for challenge generation', example: 'binary search' })
    topic!: string;

    @ApiPropertyOptional({ description: 'Difficulty level', enum: ['easy', 'medium', 'hard'], example: 'medium' })
    difficulty?: string;

    @ApiPropertyOptional({ description: 'Number of ideas to generate', example: 3 })
    count?: number;
}

export class GenerateTestCasesDto {
    @ApiProperty({ description: 'Challenge description', example: 'Given an array of integers, find the two numbers that add up to a target' })
    challengeDescription!: string;

    @ApiProperty({ description: 'Expected input format', example: 'First line: N (array size), Second line: N integers, Third line: target' })
    inputFormat!: string;

    @ApiProperty({ description: 'Expected output format', example: 'Two integers (indices) separated by space' })
    outputFormat!: string;

    @ApiPropertyOptional({ description: 'Number of public test cases to generate', example: 2 })
    publicCount?: number;

    @ApiPropertyOptional({ description: 'Number of hidden test cases to generate', example: 3 })
    hiddenCount?: number;
}
