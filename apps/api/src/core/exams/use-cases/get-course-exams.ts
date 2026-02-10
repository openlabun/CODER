import { IExamRepo } from '../interfaces/exam.repo';

export class GetCourseExams {
    constructor(private readonly examRepo: IExamRepo) { }

    async execute(courseId: string) {
        return this.examRepo.findByCourseId(courseId);
    }
}
