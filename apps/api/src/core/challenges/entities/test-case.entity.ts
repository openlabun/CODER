import { randomUUID } from 'crypto';

export class TestCase {
    private constructor(
        public readonly id: string,
        public readonly challengeId: string,
        public name: string,
        public input: string,
        public expectedOutput: string,
        public isSample: boolean,
        public points: number,
        public readonly createdAt: Date,
    ) { }

    static create(params: {
        challengeId: string;
        name: string;
        input: string;
        expectedOutput: string;
        isSample?: boolean;
        points?: number;
    }) {
        if (!params.name || params.name.trim().length < 1) {
            throw new Error('Test case name is required');
        }
        if (!params.input) {
            throw new Error('Input is required');
        }
        if (!params.expectedOutput) {
            throw new Error('Expected output is required');
        }

        return new TestCase(
            randomUUID(),
            params.challengeId,
            params.name.trim(),
            params.input,
            params.expectedOutput,
            params.isSample ?? false,
            params.points ?? 10,
            new Date(),
        );
    }

    static fromPersistence(row: {
        id: string;
        challenge_id: string;
        name: string;
        input: string;
        expected_output: string;
        is_sample: boolean;
        points: number;
        created_at: Date | string;
    }) {
        return new TestCase(
            row.id,
            row.challenge_id,
            row.name,
            row.input,
            row.expected_output,
            row.is_sample,
            row.points,
            new Date(row.created_at),
        );
    }

    updateInput(newInput: string) {
        if (!newInput) throw new Error('Input cannot be empty');
        this.input = newInput;
    }

    updateExpectedOutput(newOutput: string) {
        if (!newOutput) throw new Error('Expected output cannot be empty');
        this.expectedOutput = newOutput;
    }

    updatePoints(newPoints: number) {
        if (newPoints < 0) throw new Error('Points must be non-negative');
        this.points = newPoints;
    }

    markAsSample() {
        this.isSample = true;
    }

    markAsHidden() {
        this.isSample = false;
    }
}
