import { ICourseRepo } from '../../../core/courses/interfaces/course.repo';

export class AssignChallengeUseCase {
    constructor(private courseRepo: ICourseRepo) { }

    async execute(courseId: string, challengeId: string): Promise<void> {
        const course = await this.courseRepo.findById(courseId);
        if (!course) {
            throw new Error('Course not found');
        }

        await this.courseRepo.addChallenge(courseId, challengeId);
    }
}
