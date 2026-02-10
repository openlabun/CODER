import { Module } from '@nestjs/common';
import { AIController } from './ai.controller';

@Module({
    controllers: [AIController],
})
export class AIModule { }
