import { Body, Controller, Get, NotFoundException, Param, Post, UseGuards, BadRequestException, Query } from '@nestjs/common';
import { SubmissionsService } from './submissions.service';
import { CreateSubmissionDto } from './dto/create-submission.dto';
import { JwtAuthGuard } from '../auth/guards/jwt-auth.guard';
import { CurrentUser } from '../auth/decorators/current-user.decorator';

@Controller('submissions')
export class SubmissionsController {
  constructor(private readonly svc: SubmissionsService) { }

  @UseGuards(JwtAuthGuard)
  @Post()
  async create(@Body() dto: CreateSubmissionDto, @CurrentUser() user: any) {
    try {
      if (!dto.challengeId) throw new BadRequestException('challengeId is required');
      const sub = await this.svc.create({
        challengeId: dto.challengeId,
        userId: String(user?.sub),
        code: dto.code,
        language: dto.language,
        examId: dto.examId,
      });
      return { id: sub.id, status: sub.status, createdAt: sub.createdAt };
    } catch (error) {
      console.error('Error creating submission:', error);
      throw error;
    }
  }

  @Get(':id')
  async getById(@Param('id') id: string) {
    const sub = await this.svc.get(id);
    if (!sub) throw new NotFoundException('Submission not found');
    return {
      id: sub.id,
      challengeId: sub.challengeId,
      userId: sub.userId,
      status: sub.status,
      createdAt: sub.createdAt,
      updatedAt: sub.updatedAt,
    };
  }

  // 👇 Nuevo: listado (por defecto “mis envíos”)
  @UseGuards(JwtAuthGuard)
  @Get()
  async list(
    @Query('challengeId') challengeId: string,
    @Query('status') status: string,
    @Query('limit') limit = '20',
    @Query('offset') offset = '0',
    @CurrentUser() user: any,
  ) {
    try {
      return await this.svc.list({
        challengeId: challengeId || undefined,
        status: status || undefined,
        userId: String(user?.sub), // “mis envíos”
        limit: Number(limit),
        offset: Number(offset),
      });
    } catch (error) {
      console.error('Error listing submissions:', error);
      throw error;
    }
  }
}
