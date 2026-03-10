import { Controller, Post, Body, Get, Param, UseGuards, Req } from '@nestjs/common';
import { ApiTags, ApiOperation, ApiResponse, ApiBearerAuth, ApiParam } from '@nestjs/swagger';
import { CreateExam } from '../../core/exams/use-cases/create-exam';
import { GetCourseExams } from '../../core/exams/use-cases/get-course-exams';
import { GetExamDetails } from '../../core/exams/use-cases/get-exam-details';
import { JwtAuthGuard } from '../auth/guards/jwt-auth.guard';
import { RolesGuard } from '../auth/guards/roles.guard';
import { Roles } from '../auth/decorators/roles.decorator';
import { CreateExamDto } from './dto/create-exam.dto';

@ApiTags('exams')
@Controller('exams')
export class ExamsController {
    constructor(
        private readonly createExamUseCase: CreateExam,
        private readonly getCourseExamsUseCase: GetCourseExams,
        private readonly getExamDetailsUseCase: GetExamDetails,
    ) { }

    @Post()
    @UseGuards(JwtAuthGuard, RolesGuard)
    @ApiBearerAuth('JWT-auth')
    @Roles('ADMIN', 'PROFESSOR')
    @ApiOperation({ summary: 'Create a new exam with challenges (professor/admin only)' })
    @ApiResponse({ status: 201, description: 'Exam created' })
    @ApiResponse({ status: 401, description: 'Unauthorized' })
    async create(@Body() body: CreateExamDto) {
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
    @ApiBearerAuth('JWT-auth')
    @ApiOperation({ summary: 'List all exams for a course' })
    @ApiParam({ name: 'courseId', description: 'Course UUID' })
    @ApiResponse({ status: 200, description: 'List of course exams' })
    async getByCourse(@Param('courseId') courseId: string) {
        return this.getCourseExamsUseCase.execute(courseId);
    }

    @Get(':id')
    @UseGuards(JwtAuthGuard)
    @ApiBearerAuth('JWT-auth')
    @ApiOperation({ summary: 'Get exam details with challenges, points and order' })
    @ApiParam({ name: 'id', description: 'Exam UUID' })
    @ApiResponse({ status: 200, description: 'Exam details with challenges' })
    async getDetails(@Param('id') id: string) {
        return this.getExamDetailsUseCase.execute(id);
    }
}
