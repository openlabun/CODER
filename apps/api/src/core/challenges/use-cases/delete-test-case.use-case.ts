import { ITestCaseRepo } from '../../../core/challenges/interfaces/test-case.repo';

export class DeleteTestCaseUseCase {
    constructor(private testCaseRepo: ITestCaseRepo) { }

    async execute(id: string): Promise<void> {
        const testCase = await this.testCaseRepo.findById(id);
        if (!testCase) {
            throw new Error('Test case not found');
        }
        await this.testCaseRepo.deleteById(id);
    }
}
