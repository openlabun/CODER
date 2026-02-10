import { ICourseRepo } from '../../../core/courses/interfaces/course.repo';

export class EnrollStudentUseCase {
    constructor(private courseRepo: ICourseRepo) { }

    async execute(courseId: string, studentId: string): Promise<void> {
        const course = await this.courseRepo.findById(courseId);
        if (!course) {
            throw new Error('Course not found');
        }

        await this.courseRepo.addStudent(courseId, studentId);
    }
}
