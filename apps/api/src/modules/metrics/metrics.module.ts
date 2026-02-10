import { Module } from '@nestjs/common';
import { MetricsController } from './metrics.controller';
import { DatabaseModule } from '../../infrastructure/database/database.module';

@Module({
    imports: [DatabaseModule],
    controllers: [MetricsController],
})
export class MetricsModule { }
