import { Controller, Post, Body, UseGuards } from '@nestjs/common';
import { JwtAuthGuard } from '../auth/guards/jwt-auth.guard';
import { CurrentUser } from '../auth/decorators/current-user.decorator';
import { GeminiService } from '../../core/ai/gemini.service';

@Controller('ai')
export class AIController {
    private geminiService: GeminiService;

    constructor() {
        this.geminiService = new GeminiService();
    }

    @UseGuards(JwtAuthGuard)
    @Post('generate-challenge-ideas')
    async generateChallengeIdeas(
        @Body() dto: { topic: string; difficulty?: string; count?: number },
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
    @Post('generate-test-cases')
    async generateTestCases(
        @Body() dto: {
            challengeDescription: string;
            inputFormat: string;
            outputFormat: string;
            publicCount?: number;
            hiddenCount?: number;
        },
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
