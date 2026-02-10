import { ITestCaseRepo } from '../../../core/challenges/interfaces/test-case.repo';
import { TestCase } from '../../../core/challenges/entities/test-case.entity';

export class ListTestCasesUseCase {
    constructor(private testCaseRepo: ITestCaseRepo) { }

    async execute(challengeId: string, samplesOnly: boolean = false): Promise<TestCase[]> {
        if (samplesOnly) {
            return this.testCaseRepo.findSamplesByChallengeId(challengeId);
        }
        return this.testCaseRepo.findByChallengeId(challengeId);
    }
}
