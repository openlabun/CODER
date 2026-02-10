import { IExamRepo } from '../interfaces/exam.repo';

export class GetExamDetails {
    constructor(private readonly examRepo: IExamRepo) { }

    async execute(examId: string) {
        const exam = await this.examRepo.findById(examId);
        if (!exam) {
            throw new Error('Exam not found');
        }

        const challenges = await this.examRepo.getExamChallenges(examId);

        return {
            ...exam,
            challenges,
        };
    }
}
