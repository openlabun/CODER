export interface Metrics {
    submissions_total: number;
    submissions_accepted: number;
    submissions_rejected: number;
    submissions_failed: number;
    average_execution_time_ms: number;
    challenges_total: number;
    courses_total: number;
    users_total: number;
}

class MetricsCollector {
    private metrics: Metrics = {
        submissions_total: 0,
        submissions_accepted: 0,
        submissions_rejected: 0,
        submissions_failed: 0,
        average_execution_time_ms: 0,
        challenges_total: 0,
        courses_total: 0,
        users_total: 0,
    };

    private executionTimes: number[] = [];

    incrementSubmissionsTotal() {
        this.metrics.submissions_total++;
    }

    incrementSubmissionsAccepted() {
        this.metrics.submissions_accepted++;
    }

    incrementSubmissionsRejected() {
        this.metrics.submissions_rejected++;
    }

    incrementSubmissionsFailed() {
        this.metrics.submissions_failed++;
    }

    recordExecutionTime(timeMs: number) {
        this.executionTimes.push(timeMs);
        if (this.executionTimes.length > 1000) {
            this.executionTimes.shift(); // Keep only last 1000
        }
        const sum = this.executionTimes.reduce((a, b) => a + b, 0);
        this.metrics.average_execution_time_ms = Math.round(sum / this.executionTimes.length);
    }

    setChallengesTotal(count: number) {
        this.metrics.challenges_total = count;
    }

    setCoursesTotal(count: number) {
        this.metrics.courses_total = count;
    }

    setUsersTotal(count: number) {
        this.metrics.users_total = count;
    }

    getMetrics(): Metrics {
        return { ...this.metrics };
    }

    reset() {
        this.metrics = {
            submissions_total: 0,
            submissions_accepted: 0,
            submissions_rejected: 0,
            submissions_failed: 0,
            average_execution_time_ms: 0,
            challenges_total: 0,
            courses_total: 0,
            users_total: 0,
        };
        this.executionTimes = [];
    }
}

export const metricsCollector = new MetricsCollector();
