import { ApiProperty, ApiPropertyOptional } from '@nestjs/swagger';

export class CreateTestCaseDto {
    @ApiProperty({ description: 'Challenge ID this test case belongs to', example: '550e8400-e29b-41d4-a716-446655440000' })
    challengeId!: string;

    @ApiProperty({ description: 'Test case name', example: 'edge_case_1' })
    name!: string;

    @ApiProperty({ description: 'Input data', example: '10\n20' })
    input!: string;

    @ApiProperty({ description: 'Expected output', example: '30' })
    expectedOutput!: string;

    @ApiPropertyOptional({ description: 'Whether this test case is visible to students', default: false })
    isSample?: boolean;

    @ApiPropertyOptional({ description: 'Points awarded for passing this test case', example: 10 })
    points?: number;
}
