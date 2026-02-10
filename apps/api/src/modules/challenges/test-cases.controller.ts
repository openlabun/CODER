import { Body, Controller, Delete, Get, Param, Post, Query, UseGuards } from '@nestjs/common';
import { JwtAuthGuard } from '../auth/guards/jwt-auth.guard';
import { CurrentUser } from '../auth/decorators/current-user.decorator';
import { CreateTestCaseDto } from './dto/create-test-case.dto';
import { CreateTestCaseUseCase } from '../../core/challenges/use-cases/create-test-case.use-case';
import { ListTestCasesUseCase } from '../../core/challenges/use-cases/list-test-cases.use-case';
import { DeleteTestCaseUseCase } from '../../core/challenges/use-cases/delete-test-case.use-case';

@Controller('test-cases')
export class TestCasesController {
    constructor(
        private createTestCase: CreateTestCaseUseCase,
        private listTestCases: ListTestCasesUseCase,
        private deleteTestCase: DeleteTestCaseUseCase,
    ) { }

    @UseGuards(JwtAuthGuard)
    @Post()
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
    @Get('challenge/:challengeId')
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
    @Delete(':id')
    async delete(@Param('id') id: string, @CurrentUser() user: any) {
        // Only professors and admins can delete test cases
        if (user.role !== 'professor' && user.role !== 'admin') {
            throw new Error('Unauthorized');
        }

        await this.deleteTestCase.execute(id);
        return { message: 'Test case deleted successfully' };
    }
}
