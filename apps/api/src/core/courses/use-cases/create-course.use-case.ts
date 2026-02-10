import { Course } from '../entities/course.entity';
import { ICourseRepo } from '../interfaces/course.repo';

export class CreateCourseUseCase {
    constructor(private courseRepo: ICourseRepo) { }

    async execute(input: {
        name: string;
        code: string;
        period: string;
        groupNumber: number;
        professorId: string;
    }): Promise<Course> {
        // Generate enrollment code: COURSE-PERIODG# (e.g., CS101-20251G1)
        const enrollmentCode = `${input.code.toUpperCase()}-${input.period.split('-').join('')}G${input.groupNumber}`;

        const course = Course.create({
            name: input.name,
            code: input.code,
            period: input.period,
            groupNumber: input.groupNumber,
            professorId: input.professorId,
            enrollmentCode,
        });

        await this.courseRepo.save(course);
        return course;
    }
}
