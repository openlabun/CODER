import { Controller, Post, Body, Get, Param, UseGuards, Req } from '@nestjs/common';
import { CreateExam } from '../../core/exams/use-cases/create-exam';
import { GetCourseExams } from '../../core/exams/use-cases/get-course-exams';
import { GetExamDetails } from '../../core/exams/use-cases/get-exam-details';
import { JwtAuthGuard } from '../auth/guards/jwt-auth.guard';
import { RolesGuard } from '../auth/guards/roles.guard';
import { Roles } from '../auth/decorators/roles.decorator';

@Controller('exams')
export class ExamsController {
    constructor(
        private readonly createExamUseCase: CreateExam,
        private readonly getCourseExamsUseCase: GetCourseExams,
        private readonly getExamDetailsUseCase: GetExamDetails,
    ) { }

    @Post()
    @UseGuards(JwtAuthGuard, RolesGuard)
    @Roles('ADMIN', 'PROFESSOR')
    async create(@Body() body: any) {
        return this.createExamUseCase.execute({
            title: body.title,
            description: body.description,
            courseId: body.courseId,
            startTime: new Date(body.startTime),
            endTime: new Date(body.endTime),
            durationMinutes: body.durationMinutes,
            challenges: body.challenges,
        });
    }

    @Get('course/:courseId')
    @UseGuards(JwtAuthGuard)
    async getByCourse(@Param('courseId') courseId: string) {
        return this.getCourseExamsUseCase.execute(courseId);
    }

    @Get(':id')
    @UseGuards(JwtAuthGuard)
    async getDetails(@Param('id') id: string) {
        return this.getExamDetailsUseCase.execute(id);
    }
}
