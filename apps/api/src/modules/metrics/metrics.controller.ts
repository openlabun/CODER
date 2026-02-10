import { Controller, Get } from '@nestjs/common';
import { ApiTags, ApiOperation, ApiResponse } from '@nestjs/swagger';
import { metricsCollector } from '../../infrastructure/metrics/metrics';
import { Inject } from '@nestjs/common';
import { Pool } from 'pg';
import { PG_POOL } from '../../infrastructure/database/postgres.provider';

@ApiTags('metrics')
@Controller('metrics')
export class MetricsController {
    constructor(@Inject(PG_POOL) private pool: Pool) { }

    @Get()
    @ApiOperation({ summary: 'Get system metrics' })
    @ApiResponse({
        status: 200,
        description: 'Returns system metrics in JSON and Prometheus format',
        schema: {
            example: {
                metrics: {
                    submissions_total: 150,
                    submissions_accepted: 85,
                    submissions_rejected: 45,
                    submissions_failed: 20,
                    average_execution_time_ms: 1234,
                    challenges_total: 25,
                    courses_total: 5,
                    users_total: 120
                },
                prometheus: '# HELP submissions_total...'
            }
        }
    })
    async getMetrics() {
        // Update counts from database
        try {
            const challengesCount = await this.pool.query('SELECT COUNT(*) FROM challenges');
            metricsCollector.setChallengesTotal(parseInt(challengesCount.rows[0].count));

            const coursesCount = await this.pool.query('SELECT COUNT(*) FROM courses');
            metricsCollector.setCoursesTotal(parseInt(coursesCount.rows[0].count));

            const usersCount = await this.pool.query('SELECT COUNT(*) FROM users');
            metricsCollector.setUsersTotal(parseInt(usersCount.rows[0].count));
        } catch (error) {
            console.error('Error fetching metrics:', error);
        }

        const metrics = metricsCollector.getMetrics();

        // Return in Prometheus-like format
        return {
            metrics,
            prometheus: this.toPrometheusFormat(metrics),
        };
    }

    private toPrometheusFormat(metrics: any): string {
        return `
# HELP submissions_total Total number of submissions
# TYPE submissions_total counter
submissions_total ${metrics.submissions_total}

# HELP submissions_accepted Number of accepted submissions
# TYPE submissions_accepted counter
submissions_accepted ${metrics.submissions_accepted}

# HELP submissions_rejected Number of rejected submissions
# TYPE submissions_rejected counter
submissions_rejected ${metrics.submissions_rejected}

# HELP submissions_failed Number of failed submissions
# TYPE submissions_failed counter
submissions_failed ${metrics.submissions_failed}

# HELP average_execution_time_ms Average execution time in milliseconds
# TYPE average_execution_time_ms gauge
average_execution_time_ms ${metrics.average_execution_time_ms}

# HELP challenges_total Total number of challenges
# TYPE challenges_total gauge
challenges_total ${metrics.challenges_total}

# HELP courses_total Total number of courses
# TYPE courses_total gauge
courses_total ${metrics.courses_total}

# HELP users_total Total number of users
# TYPE users_total gauge
users_total ${metrics.users_total}
    `.trim();
    }
}
