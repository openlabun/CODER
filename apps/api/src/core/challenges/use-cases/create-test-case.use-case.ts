import { TestCase } from '../../../core/challenges/entities/test-case.entity';
import { ITestCaseRepo } from '../../../core/challenges/interfaces/test-case.repo';

export class CreateTestCaseUseCase {
    constructor(private testCaseRepo: ITestCaseRepo) { }

    async execute(input: {
        challengeId: string;
        name: string;
        input: string;
        expectedOutput: string;
        isSample?: boolean;
        points?: number;
    }): Promise<TestCase> {
        const testCase = TestCase.create({
            challengeId: input.challengeId,
            name: input.name,
            input: input.input,
            expectedOutput: input.expectedOutput,
            isSample: input.isSample,
            points: input.points,
        });

        await this.testCaseRepo.save(testCase);
        return testCase;
    }
}
