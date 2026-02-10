import { Course } from '../../../core/courses/entities/course.entity';
import { ICourseRepo } from '../../../core/courses/interfaces/course.repo';

export class ListCoursesUseCase {
    constructor(private courseRepo: ICourseRepo) { }

    async execute(filters: {
        userId?: string;
        role?: string;
    }): Promise<Course[]> {
        // If professor, show their courses
        if (filters.role === 'professor' && filters.userId) {
            return this.courseRepo.findByProfessor(filters.userId);
        }

        // If student, show enrolled courses
        if (filters.role === 'student' && filters.userId) {
            return this.courseRepo.findByStudent(filters.userId);
        }

        // Admin sees all
        return this.courseRepo.list();
    }
}
