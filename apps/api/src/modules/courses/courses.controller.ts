import { Body, Controller, Delete, Get, Param, Post, UseGuards, NotFoundException, BadRequestException, UnauthorizedException } from '@nestjs/common';
import { ApiTags, ApiOperation, ApiResponse, ApiBearerAuth, ApiParam } from '@nestjs/swagger';
import { JwtAuthGuard } from '../auth/guards/jwt-auth.guard';
import { CurrentUser } from '../auth/decorators/current-user.decorator';
import { CreateCourseDto, EnrollByCodeDto, EnrollStudentDto, AssignChallengeDto } from './dto/course.dto';
import { CreateCourseUseCase } from '../../core/courses/use-cases/create-course.use-case';
import { ListCoursesUseCase } from '../../core/courses/use-cases/list-courses.use-case';
import { EnrollStudentUseCase } from '../../core/courses/use-cases/enroll-student.use-case';
import { AssignChallengeUseCase } from '../../core/courses/use-cases/assign-challenge.use-case';
import { ICourseRepo } from '../../core/courses/interfaces/course.repo';
import { Inject } from '@nestjs/common';
import { PG_POOL } from '../../infrastructure/database/postgres.provider';
import { Pool } from 'pg';

@ApiTags('courses')
@Controller('courses')
export class CoursesController {
    constructor(
        private createCourse: CreateCourseUseCase,
        private listCourses: ListCoursesUseCase,
        private enrollStudent: EnrollStudentUseCase,
        private assignChallenge: AssignChallengeUseCase,
        @Inject('CourseRepo') private courseRepo: ICourseRepo,
        @Inject(PG_POOL) private pool: Pool,
    ) { }

    @UseGuards(JwtAuthGuard)
    @ApiBearerAuth('JWT-auth')
    @Post('enroll')
    @ApiOperation({ summary: 'Enroll in a course using an enrollment code (students only)' })
    @ApiResponse({ status: 200, description: 'Successfully enrolled' })
    @ApiResponse({ status: 404, description: 'Invalid enrollment code' })
    @ApiResponse({ status: 401, description: 'Only students can enroll' })
    async enrollByCode(@Body() dto: EnrollByCodeDto, @CurrentUser() user: any) {
        // Only students can enroll via code
        if (user.role !== 'student') {
            throw new UnauthorizedException('Only students can enroll in courses');
        }

        // Find course by enrollment code
        const courses = await this.courseRepo.findAll();
        const course = courses.find(c => c.enrollmentCode === dto.enrollmentCode);

        if (!course) {
            throw new NotFoundException('Invalid enrollment code');
        }

        // Enroll the student
        await this.courseRepo.addStudent(course.id, user.sub);

        return {
            message: 'Successfully enrolled in course',
            courseId: course.id,
            courseName: course.name
        };
    }

    @UseGuards(JwtAuthGuard)
    @ApiBearerAuth('JWT-auth')
    @Post()
    @ApiOperation({ summary: 'Create a new course (professor/admin only)' })
    @ApiResponse({ status: 201, description: 'Course created with auto-generated enrollment code' })
    @ApiResponse({ status: 401, description: 'Unauthorized' })
    async create(@Body() dto: CreateCourseDto, @CurrentUser() user: any) {
        // Only professors and admins can create courses
        if (user.role !== 'professor' && user.role !== 'admin') {
            throw new UnauthorizedException('Only professors and admins can create courses');
        }

        const course = await this.createCourse.execute({
            name: dto.name,
            code: dto.code,
            period: dto.period,
            groupNumber: dto.groupNumber,
            professorId: user.sub,
        });

        return {
            id: course.id,
            name: course.name,
            code: course.code,
            period: course.period,
            groupNumber: course.groupNumber,
            enrollmentCode: course.enrollmentCode,
            professorId: course.professorId,
            createdAt: course.createdAt,
        };
    }

    @UseGuards(JwtAuthGuard)
    @ApiBearerAuth('JWT-auth')
    @Get('browse')
    @ApiOperation({ summary: 'Browse all available courses' })
    @ApiResponse({ status: 200, description: 'List of all courses' })
    async browse(@CurrentUser() user: any) {
        const courses = await this.courseRepo.findAll();
        return courses.map((c) => ({
            id: c.id,
            name: c.name,
            code: c.code,
            period: c.period,
            groupNumber: c.groupNumber,
            professorId: c.professorId,
            createdAt: c.createdAt,
        }));
    }

    @UseGuards(JwtAuthGuard)
    @ApiBearerAuth('JWT-auth')
    @Get()
    @ApiOperation({ summary: 'List courses for current user (by role: student enrolled, professor owned)' })
    @ApiResponse({ status: 200, description: 'List of user courses' })
    async list(@CurrentUser() user: any) {
        const courses = await this.listCourses.execute({
            userId: user.sub,
            role: user.role,
        });

        return courses.map((c) => ({
            id: c.id,
            name: c.name,
            code: c.code,
            period: c.period,
            groupNumber: c.groupNumber,
            enrollmentCode: c.enrollmentCode,
            professorId: c.professorId,
            createdAt: c.createdAt,
        }));
    }

    @UseGuards(JwtAuthGuard)
    @ApiBearerAuth('JWT-auth')
    @Get(':id')
    @ApiOperation({ summary: 'Get course details by ID' })
    @ApiParam({ name: 'id', description: 'Course UUID' })
    @ApiResponse({ status: 200, description: 'Course details' })
    @ApiResponse({ status: 404, description: 'Course not found' })
    async getOne(@Param('id') id: string, @CurrentUser() user: any) {
        const course = await this.courseRepo.findById(id);
        if (!course) {
            throw new NotFoundException('Course not found');
        }

        // Check access: professor owns it, student is enrolled, or admin
        if (user.role === 'student') {
            const students = await this.courseRepo.getStudents(id);
            if (!students.includes(user.sub)) {
                throw new UnauthorizedException('Not enrolled in this course');
            }
        } else if (user.role === 'professor' && course.professorId !== user.sub) {
            throw new UnauthorizedException('Not your course');
        }

        return {
            id: course.id,
            name: course.name,
            code: course.code,
            period: course.period,
            groupNumber: course.groupNumber,
            enrollmentCode: course.enrollmentCode,
            professorId: course.professorId,
            createdAt: course.createdAt,
        };
    }

    @UseGuards(JwtAuthGuard)
    @ApiBearerAuth('JWT-auth')
    @Post(':id')
    @ApiOperation({ summary: 'Update course information (professor owner or admin)' })
    @ApiParam({ name: 'id', description: 'Course UUID' })
    @ApiResponse({ status: 200, description: 'Course updated' })
    @ApiResponse({ status: 404, description: 'Course not found' })
    @ApiResponse({ status: 401, description: 'Unauthorized' })
    async update(@Param('id') id: string, @Body() dto: CreateCourseDto, @CurrentUser() user: any) {
        const course = await this.courseRepo.findById(id);
        if (!course) {
            throw new NotFoundException('Course not found');
        }

        if (user.role !== 'professor' && user.role !== 'admin') {
            throw new UnauthorizedException('Unauthorized');
        }

        if (user.role === 'professor' && course.professorId !== user.sub) {
            throw new UnauthorizedException('Not your course');
        }

        course.updateInfo(dto.name, dto.code, dto.period, dto.groupNumber);
        await this.courseRepo.update(course);

        return {
            id: course.id,
            name: course.name,
            code: course.code,
            period: course.period,
            groupNumber: course.groupNumber,
            enrollmentCode: course.enrollmentCode,
            professorId: course.professorId,
            updatedAt: course.updatedAt,
        };
    }

    @UseGuards(JwtAuthGuard)
    @ApiBearerAuth('JWT-auth')
    @Post(':id/students')
    @ApiOperation({ summary: 'Enroll a student in a course (professor/admin)' })
    @ApiParam({ name: 'id', description: 'Course UUID' })
    @ApiResponse({ status: 200, description: 'Student enrolled' })
    @ApiResponse({ status: 404, description: 'Course not found' })
    @ApiResponse({ status: 401, description: 'Unauthorized' })
    async enroll(@Param('id') id: string, @Body() dto: EnrollStudentDto, @CurrentUser() user: any) {
        // Only professor of the course or admin can enroll students
        const course = await this.courseRepo.findById(id);
        if (!course) {
            throw new NotFoundException('Course not found');
        }

        if (user.role === 'professor' && course.professorId !== user.sub) {
            throw new UnauthorizedException('Not your course');
        }

        if (user.role !== 'professor' && user.role !== 'admin') {
            throw new UnauthorizedException('Unauthorized');
        }

        await this.enrollStudent.execute(id, dto.studentId);
        return { message: 'Student enrolled successfully' };
    }

    @UseGuards(JwtAuthGuard)
    @ApiBearerAuth('JWT-auth')
    @Delete(':id/students/:studentId')
    @ApiOperation({ summary: 'Remove a student from a course (professor/admin)' })
    @ApiParam({ name: 'id', description: 'Course UUID' })
    @ApiParam({ name: 'studentId', description: 'Student UUID to remove' })
    @ApiResponse({ status: 200, description: 'Student removed' })
    @ApiResponse({ status: 404, description: 'Course not found' })
    @ApiResponse({ status: 401, description: 'Unauthorized' })
    async unenroll(@Param('id') id: string, @Param('studentId') studentId: string, @CurrentUser() user: any) {
        const course = await this.courseRepo.findById(id);
        if (!course) {
            throw new NotFoundException('Course not found');
        }

        if (user.role === 'professor' && course.professorId !== user.sub) {
            throw new UnauthorizedException('Not your course');
        }

        if (user.role !== 'professor' && user.role !== 'admin') {
            throw new UnauthorizedException('Unauthorized');
        }

        await this.courseRepo.removeStudent(id, studentId);
        return { message: 'Student removed successfully' };
    }

    @UseGuards(JwtAuthGuard)
    @ApiBearerAuth('JWT-auth')
    @Post(':id/challenges')
    @ApiOperation({ summary: 'Assign a challenge to a course (professor/admin)' })
    @ApiParam({ name: 'id', description: 'Course UUID' })
    @ApiResponse({ status: 200, description: 'Challenge assigned' })
    @ApiResponse({ status: 404, description: 'Course not found' })
    @ApiResponse({ status: 401, description: 'Unauthorized' })
    async assignChallengeToourse(@Param('id') id: string, @Body() dto: AssignChallengeDto, @CurrentUser() user: any) {
        const course = await this.courseRepo.findById(id);
        if (!course) {
            throw new NotFoundException('Course not found');
        }

        if (user.role === 'professor' && course.professorId !== user.sub) {
            throw new UnauthorizedException('Not your course');
        }

        if (user.role !== 'professor' && user.role !== 'admin') {
            throw new UnauthorizedException('Unauthorized');
        }

        await this.assignChallenge.execute(id, dto.challengeId);
        return { message: 'Challenge assigned successfully' };
    }

    @UseGuards(JwtAuthGuard)
    @ApiBearerAuth('JWT-auth')
    @Get(':id/students')
    @ApiOperation({ summary: 'List students enrolled in a course' })
    @ApiParam({ name: 'id', description: 'Course UUID' })
    @ApiResponse({ status: 200, description: 'List of enrolled students' })
    @ApiResponse({ status: 404, description: 'Course not found' })
    async getStudents(@Param('id') id: string, @CurrentUser() user: any) {
        const course = await this.courseRepo.findById(id);
        if (!course) {
            throw new NotFoundException('Course not found');
        }

        const studentIds = await this.courseRepo.getStudents(id);

        // Fetch user details for each student
        const students = await Promise.all(
            studentIds.map(async (studentId) => {
                try {
                    const userResult = await this.pool.query(
                        'SELECT id, username FROM users WHERE id = $1',
                        [studentId]
                    );
                    if (userResult.rows.length > 0) {
                        return {
                            id: userResult.rows[0].id,
                            username: userResult.rows[0].username,
                        };
                    }
                    return { id: studentId, username: 'Unknown' };
                } catch (error) {
                    console.error('Error fetching user:', error);
                    return { id: studentId, username: 'Unknown' };
                }
            })
        );

        return { students };
    }

    @UseGuards(JwtAuthGuard)
    @ApiBearerAuth('JWT-auth')
    @Get(':id/challenges')
    @ApiOperation({ summary: 'List challenges assigned to a course' })
    @ApiParam({ name: 'id', description: 'Course UUID' })
    @ApiResponse({ status: 200, description: 'List of course challenges' })
    @ApiResponse({ status: 404, description: 'Course not found' })
    async getChallenges(@Param('id') id: string, @CurrentUser() user: any) {
        const course = await this.courseRepo.findById(id);
        if (!course) {
            throw new NotFoundException('Course not found');
        }

        const challengeIds = await this.courseRepo.getChallenges(id);

        if (challengeIds.length === 0) {
            return { challenges: [] };
        }

        try {
            const query = `
                SELECT * FROM challenges 
                WHERE id = ANY($1)
            `;
            const result = await this.pool.query(query, [challengeIds]);

            const challenges = result.rows.map(row => ({
                id: row.id,
                title: row.title,
                description: row.description,
                difficulty: row.difficulty,
                timeLimit: row.time_limit,
                memoryLimit: row.memory_limit,
                status: row.status,
                createdAt: row.created_at
            }));

            return { challenges };
        } catch (error) {
            console.error('Error fetching course challenges:', error);
            return { challenges: [] };
        }
    }
}
