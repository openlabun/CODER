// apps/api/src/main.ts
import 'reflect-metadata';
import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';
import { DocumentBuilder, SwaggerModule } from '@nestjs/swagger';

async function bootstrap() {
  const app = await NestFactory.create(AppModule);

  // Enable CORS for frontend
  app.enableCors({
    origin: 'http://localhost:5173',
    credentials: true,
  });

  // Swagger/OpenAPI configuration
  const config = new DocumentBuilder()
    .setTitle('Juez Online API')
    .setDescription('API documentation for the Online Judge platform')
    .setVersion('1.0')
    .addTag('auth', 'Authentication endpoints')
    .addTag('challenges', 'Challenge management')
    .addTag('test-cases', 'Test cases for challenges')
    .addTag('submissions', 'Code submissions')
    .addTag('courses', 'Course management')
    .addTag('leaderboard', 'Rankings and leaderboards')
    .addTag('metrics', 'System metrics and observability')
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
  SwaggerModule.setup('api/docs', app, document);

  const port = process.env.PORT || 3000;
  await app.listen(port);
  // eslint-disable-next-line no-console
  console.log(`API listening on http://localhost:${port}`);
  console.log(`Swagger docs available at http://localhost:${port}/api/docs`);
}
bootstrap();

