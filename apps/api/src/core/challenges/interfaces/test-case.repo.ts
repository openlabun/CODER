import { TestCase } from '../entities/test-case.entity';

export interface ITestCaseRepo {
    save(testCase: TestCase): Promise<void>;
    findByChallengeId(challengeId: string): Promise<TestCase[]>;
    findSamplesByChallengeId(challengeId: string): Promise<TestCase[]>;
    findById(id: string): Promise<TestCase | null>;
    deleteById(id: string): Promise<void>;
}
