export class CreateTestCaseDto {
    challengeId!: string;
    name!: string;
    input!: string;
    expectedOutput!: string;
    isSample?: boolean;
    points?: number;
}
