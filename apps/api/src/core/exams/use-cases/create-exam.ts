import { Exam } from '../entities/exam.entity';
import { IExamRepo } from '../interfaces/exam.repo';

export class CreateExam {
    constructor(private readonly examRepo: IExamRepo) { }

    async execute(props: {
        title: string;
        description: string;
        courseId: string;
        startTime: Date;
        endTime: Date;
        durationMinutes: number;
        challenges: { challengeId: string; points: number; order: number }[];
    }): Promise<Exam> {
        const exam = Exam.create({
            title: props.title,
            description: props.description,
            courseId: props.courseId,
            startTime: props.startTime,
            endTime: props.endTime,
            durationMinutes: props.durationMinutes,
        });

        await this.examRepo.save(exam);

        for (const ch of props.challenges) {
            await this.examRepo.addChallengeToExam(exam.id, ch.challengeId, ch.points, ch.order);
        }

        return exam;
    }
}
