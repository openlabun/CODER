import { Controller, Post, Body, UseGuards } from '@nestjs/common';
import { ApiTags, ApiOperation, ApiResponse, ApiBearerAuth } from '@nestjs/swagger';
import { JwtAuthGuard } from '../auth/guards/jwt-auth.guard';
import { CurrentUser } from '../auth/decorators/current-user.decorator';
import { GeminiService } from '../../core/ai/gemini.service';
import { GenerateChallengeIdeasDto, GenerateTestCasesDto } from './dto/ai.dto';

@ApiTags('ai')
@Controller('ai')
export class AIController {
    private geminiService: GeminiService;

    constructor() {
        this.geminiService = new GeminiService();
    }

    @UseGuards(JwtAuthGuard)
    @ApiBearerAuth('JWT-auth')
    @Post('generate-challenge-ideas')
    @ApiOperation({ summary: 'Generate challenge ideas using AI (professor/admin only)' })
    @ApiResponse({ status: 200, description: 'AI-generated challenge ideas' })
    @ApiResponse({ status: 401, description: 'Unauthorized' })
    async generateChallengeIdeas(
        @Body() dto: GenerateChallengeIdeasDto,
        @CurrentUser() user: any
    ) {
        // Only professors and admins can use AI assistant
        if (user.role !== 'professor' && user.role !== 'admin') {
            throw new Error('Only professors and admins can use AI assistant');
        }

        const ideas = await this.geminiService.generateChallengeIdeas(
            dto.topic,
            dto.difficulty,
            dto.count || 3
        );

        return { ideas };
    }

    @UseGuards(JwtAuthGuard)
    @ApiBearerAuth('JWT-auth')
    @Post('generate-test-cases')
    @ApiOperation({ summary: 'Generate test cases using AI (professor/admin only)' })
    @ApiResponse({ status: 200, description: 'AI-generated test cases (public and hidden)' })
    @ApiResponse({ status: 401, description: 'Unauthorized' })
    async generateTestCases(
        @Body() dto: GenerateTestCasesDto,
        @CurrentUser() user: any
    ) {
        // Only professors and admins can use AI assistant
        if (user.role !== 'professor' && user.role !== 'admin') {
            throw new Error('Only professors and admins can use AI assistant');
        }

        const testCases = await this.geminiService.generateTestCases(
            dto.challengeDescription,
            dto.inputFormat,
            dto.outputFormat,
            dto.publicCount || 2,
            dto.hiddenCount || 3
        );

        return testCases;
    }
}
