export interface LogEntry {
    timestamp: string;
    level: 'info' | 'warn' | 'error' | 'debug';
    message: string;
    context?: string;
    requestId?: string;
    submissionId?: string;
    userId?: string;
    challengeId?: string;
    duration?: number;
    error?: {
        message: string;
        stack?: string;
    };
    metadata?: Record<string, any>;
}

export class Logger {
    private context: string;

    constructor(context: string) {
        this.context = context;
    }

    private log(level: LogEntry['level'], message: string, metadata?: Record<string, any>) {
        const entry: LogEntry = {
            timestamp: new Date().toISOString(),
            level,
            message,
            context: this.context,
            ...metadata,
        };

        // Output as JSON
        console.log(JSON.stringify(entry));
    }

    info(message: string, metadata?: Record<string, any>) {
        this.log('info', message, metadata);
    }

    warn(message: string, metadata?: Record<string, any>) {
        this.log('warn', message, metadata);
    }

    error(message: string, error?: Error, metadata?: Record<string, any>) {
        this.log('error', message, {
            ...metadata,
            error: error ? {
                message: error.message,
                stack: error.stack,
            } : undefined,
        });
    }

    debug(message: string, metadata?: Record<string, any>) {
        this.log('debug', message, metadata);
    }
}

export function createLogger(context: string): Logger {
    return new Logger(context);
}
