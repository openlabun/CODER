import { Module } from '@nestjs/common';
import { DatabaseModule } from '../../infrastructure/database/database.module';
import { PostgresExamRepo } from '../../infrastructure/database/postgres/postgres-exam.repo';
import { CreateExam } from '../../core/exams/use-cases/create-exam';
import { GetCourseExams } from '../../core/exams/use-cases/get-course-exams';
import { GetExamDetails } from '../../core/exams/use-cases/get-exam-details';
import { ExamsController } from './exams.controller';

@Module({
    imports: [DatabaseModule],
    controllers: [ExamsController],
    providers: [
        PostgresExamRepo,
        {
            provide: CreateExam,
            useFactory: (repo: PostgresExamRepo) => new CreateExam(repo),
            inject: [PostgresExamRepo],
        },
        {
            provide: GetCourseExams,
            useFactory: (repo: PostgresExamRepo) => new GetCourseExams(repo),
            inject: [PostgresExamRepo],
        },
        {
            provide: GetExamDetails,
            useFactory: (repo: PostgresExamRepo) => new GetExamDetails(repo),
            inject: [PostgresExamRepo],
        },
    ],
    exports: [PostgresExamRepo],
})
export class ExamsModule { }
