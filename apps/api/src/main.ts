// apps/api/src/main.ts
import 'reflect-metadata';
import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';
import { DocumentBuilder, SwaggerModule } from '@nestjs/swagger';

async function bootstrap() {
  const app = await NestFactory.create(AppModule);

  // Enable CORS for frontend
  app.enableCors({
    origin: process.env.CORS_ORIGIN || 'http://localhost:5173',
    credentials: true,
  });

  // Swagger/OpenAPI configuration
  const config = new DocumentBuilder()
    .setTitle('Juez Online API')
    .setDescription(
      'API documentation for the Online Judge platform. Current submission flow: client sends code to /submissions, API persists and enqueues submission ID in Redis (queue:submissions), worker consumes jobs, executes language runner containers, and persists final verdict and score.',
    )
    .setVersion('1.0')
    .addTag('auth', 'Authentication endpoints')
    .addTag('challenges', 'Challenge management')
    .addTag('test-cases', 'Test cases for challenges')
    .addTag('submissions', 'Code submissions')
    .addTag('courses', 'Course management')
    .addTag('exams', 'Exam management')
    .addTag('leaderboard', 'Rankings and leaderboards')
    .addTag('ai', 'AI-powered content generation')
    .addTag('metrics', 'System metrics and observability')
    .addTag('health', 'Health check endpoints')
    .addBearerAuth(
      {
        type: 'http',
        scheme: 'bearer',
        bearerFormat: 'JWT',
        name: 'JWT',
        description: 'Enter JWT token',
        in: 'header',
      },
      'JWT-auth',
    )
    .build();

  const document = SwaggerModule.createDocument(app, config);
  SwaggerModule.setup('docs', app, document);

  const port = process.env.PORT || 3000;
  await app.listen(port);
  // eslint-disable-next-line no-console
  console.log(`API listening on http://localhost:${port}`);
  console.log(`Swagger docs available at http://localhost:${port}/docs`);
}
bootstrap();

