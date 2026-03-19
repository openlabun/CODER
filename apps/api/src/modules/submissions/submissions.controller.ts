import { Body, Controller, Get, NotFoundException, Param, Post, UseGuards, BadRequestException, Query } from '@nestjs/common';
import { ApiTags, ApiOperation, ApiResponse, ApiBearerAuth, ApiParam, ApiQuery, ApiBadRequestResponse, ApiUnauthorizedResponse, ApiNotFoundResponse } from '@nestjs/swagger';
import { SubmissionsService } from './submissions.service';
import { CreateSubmissionDto } from './dto/create-submission.dto';
import { JwtAuthGuard } from '../auth/guards/jwt-auth.guard';
import { CurrentUser } from '../auth/decorators/current-user.decorator';

@ApiTags('submissions')
@Controller('submissions')
export class SubmissionsController {
  constructor(private readonly svc: SubmissionsService) { }

  @UseGuards(JwtAuthGuard)
  @ApiBearerAuth('JWT-auth')
  @Post()
  @ApiOperation({ summary: 'Submit code for evaluation against a challenge' })
  @ApiResponse({
    status: 201,
    description: 'Submission queued for execution',
    schema: {
      example: {
        id: '550e8400-e29b-41d4-a716-446655440010',
        status: 'queued',
        createdAt: '2026-03-16T14:20:00.000Z',
      },
    },
  })
  @ApiBadRequestResponse({ description: 'Invalid input or challengeId does not exist' })
  @ApiUnauthorizedResponse({ description: 'Unauthorized' })
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
  @ApiOperation({ summary: 'Get submission details by ID (currently public endpoint)' })
  @ApiParam({ name: 'id', description: 'Submission UUID' })
  @ApiResponse({
    status: 200,
    description: 'Submission details',
    schema: {
      example: {
        id: '550e8400-e29b-41d4-a716-446655440010',
        challengeId: '550e8400-e29b-41d4-a716-446655440000',
        userId: '550e8400-e29b-41d4-a716-446655440001',
        status: 'accepted',
        createdAt: '2026-03-16T14:20:00.000Z',
        updatedAt: '2026-03-16T14:20:02.000Z',
      },
    },
  })
  @ApiNotFoundResponse({ description: 'Submission not found' })
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

  @UseGuards(JwtAuthGuard)
  @ApiBearerAuth('JWT-auth')
  @Get()
  @ApiOperation({ summary: 'List current user submissions with optional filters' })
  @ApiQuery({ name: 'challengeId', required: false, description: 'Filter by challenge ID' })
  @ApiQuery({ name: 'status', required: false, description: 'Filter by status', enum: ['queued', 'running', 'accepted', 'wrong_answer', 'error'] })
  @ApiQuery({ name: 'limit', required: false, description: 'Max results (default 20)' })
  @ApiQuery({ name: 'offset', required: false, description: 'Skip results (default 0)' })
  @ApiResponse({
    status: 200,
    description: 'List of submissions',
    schema: {
      example: {
        total: 2,
        limit: 20,
        offset: 0,
        items: [
          {
            id: '550e8400-e29b-41d4-a716-446655440010',
            challengeId: '550e8400-e29b-41d4-a716-446655440000',
            userId: '550e8400-e29b-41d4-a716-446655440001',
            language: 'python',
            status: 'accepted',
            score: 100,
            timeMsTotal: 25,
            createdAt: '2026-03-16T14:20:00.000Z',
            updatedAt: '2026-03-16T14:20:02.000Z',
          },
        ],
      },
    },
  })
  @ApiUnauthorizedResponse({ description: 'Unauthorized' })
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
