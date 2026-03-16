import { Body, Controller, Delete, Get, Param, Post, Query, UseGuards } from '@nestjs/common';
import { ApiTags, ApiOperation, ApiResponse, ApiBearerAuth, ApiParam, ApiQuery, ApiUnauthorizedResponse } from '@nestjs/swagger';
import { JwtAuthGuard } from '../auth/guards/jwt-auth.guard';
import { CurrentUser } from '../auth/decorators/current-user.decorator';
import { CreateTestCaseDto } from './dto/create-test-case.dto';
import { CreateTestCaseUseCase } from '../../core/challenges/use-cases/create-test-case.use-case';
import { ListTestCasesUseCase } from '../../core/challenges/use-cases/list-test-cases.use-case';
import { DeleteTestCaseUseCase } from '../../core/challenges/use-cases/delete-test-case.use-case';

@ApiTags('test-cases')
@Controller('test-cases')
export class TestCasesController {
    constructor(
        private createTestCase: CreateTestCaseUseCase,
        private listTestCases: ListTestCasesUseCase,
        private deleteTestCase: DeleteTestCaseUseCase,
    ) { }

    @UseGuards(JwtAuthGuard)
    @ApiBearerAuth('JWT-auth')
    @Post()
    @ApiOperation({ summary: 'Create a test case for a challenge (professor/admin only)' })
    @ApiResponse({ status: 201, description: 'Test case created' })
    @ApiUnauthorizedResponse({ description: 'Unauthorized (missing/invalid JWT)' })
    @ApiResponse({ status: 500, description: 'Role validation currently throws generic Error in controller' })
    async create(@Body() dto: CreateTestCaseDto, @CurrentUser() user: any) {
        // Only professors and admins can create test cases
        if (user.role !== 'professor' && user.role !== 'admin') {
            throw new Error('Unauthorized');
        }

        const testCase = await this.createTestCase.execute({
            challengeId: dto.challengeId,
            name: dto.name,
            input: dto.input,
            expectedOutput: dto.expectedOutput,
            isSample: dto.isSample,
            points: dto.points,
        });

        return {
            id: testCase.id,
            challengeId: testCase.challengeId,
            name: testCase.name,
            isSample: testCase.isSample,
            points: testCase.points,
            createdAt: testCase.createdAt,
        };
    }

    @UseGuards(JwtAuthGuard)
    @ApiBearerAuth('JWT-auth')
    @Get('challenge/:challengeId')
    @ApiOperation({ summary: 'List test cases for a challenge (students see samples only)' })
    @ApiParam({ name: 'challengeId', description: 'Challenge UUID' })
    @ApiQuery({ name: 'samplesOnly', required: false, description: 'If true, returns only sample test cases' })
    @ApiResponse({ status: 200, description: 'List of test cases' })
    @ApiUnauthorizedResponse({ description: 'Unauthorized (missing/invalid JWT)' })
    async list(
        @Param('challengeId') challengeId: string,
        @Query('samplesOnly') samplesOnly: string,
        @CurrentUser() user: any,
    ) {
        // Students can only see sample test cases
        const onlySamples = user.role === 'student' || samplesOnly === 'true';

        const testCases = await this.listTestCases.execute(challengeId, onlySamples);

        return testCases.map((tc) => ({
            id: tc.id,
            challengeId: tc.challengeId,
            name: tc.name,
            input: onlySamples ? tc.input : undefined, // Hide input/output for hidden cases
            expectedOutput: onlySamples ? tc.expectedOutput : undefined,
            isSample: tc.isSample,
            points: tc.points,
            createdAt: tc.createdAt,
        }));
    }

    @UseGuards(JwtAuthGuard)
    @ApiBearerAuth('JWT-auth')
    @Delete(':id')
    @ApiOperation({ summary: 'Delete a test case (professor/admin only)' })
    @ApiParam({ name: 'id', description: 'Test case UUID' })
    @ApiResponse({ status: 200, description: 'Test case deleted' })
    @ApiUnauthorizedResponse({ description: 'Unauthorized (missing/invalid JWT)' })
    @ApiResponse({ status: 500, description: 'Role validation currently throws generic Error in controller' })
    async delete(@Param('id') id: string, @CurrentUser() user: any) {
        // Only professors and admins can delete test cases
        if (user.role !== 'professor' && user.role !== 'admin') {
            throw new Error('Unauthorized');
        }

        await this.deleteTestCase.execute(id);
        return { message: 'Test case deleted successfully' };
    }
}
