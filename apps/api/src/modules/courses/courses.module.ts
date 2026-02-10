import { Module } from '@nestjs/common';
import { CoursesController } from './courses.controller';
import { CreateCourseUseCase } from '../../core/courses/use-cases/create-course.use-case';
import { ListCoursesUseCase } from '../../core/courses/use-cases/list-courses.use-case';
import { EnrollStudentUseCase } from '../../core/courses/use-cases/enroll-student.use-case';
import { AssignChallengeUseCase } from '../../core/courses/use-cases/assign-challenge.use-case';
import { DatabaseModule } from '../../infrastructure/database/database.module';

@Module({
    imports: [DatabaseModule],
    controllers: [CoursesController],
    providers: [
        {
            provide: CreateCourseUseCase,
            useFactory: (repo) => new CreateCourseUseCase(repo),
            inject: ['CourseRepo'],
        },
        {
            provide: ListCoursesUseCase,
            useFactory: (repo) => new ListCoursesUseCase(repo),
            inject: ['CourseRepo'],
        },
        {
            provide: EnrollStudentUseCase,
            useFactory: (repo) => new EnrollStudentUseCase(repo),
            inject: ['CourseRepo'],
        },
        {
            provide: AssignChallengeUseCase,
            useFactory: (repo) => new AssignChallengeUseCase(repo),
            inject: ['CourseRepo'],
        },
    ],
})
export class CoursesModule { }
