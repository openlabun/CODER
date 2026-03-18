import { Body, Controller, Get, Post, UseGuards } from '@nestjs/common';
import { ApiTags, ApiOperation, ApiResponse, ApiBearerAuth, ApiBadRequestResponse, ApiConflictResponse, ApiUnauthorizedResponse } from '@nestjs/swagger';
import { AuthService } from './auth.service';
import { LoginDto } from './dto/login.dto';
import { RegisterDto } from './dto/register.dto';
import { JwtAuthGuard } from './guards/jwt-auth.guard';
import { CurrentUser } from './decorators/current-user.decorator';

@ApiTags('auth')
@Controller('auth')
export class AuthController {
  constructor(private readonly auth: AuthService) { }

  @Post('register')
  @ApiOperation({ summary: 'Register a new user' })
  @ApiResponse({
    status: 201,
    description: 'User registered successfully',
    schema: {
      example: {
        accessToken: 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...',
      },
    },
  })
  @ApiBadRequestResponse({ description: 'Invalid input payload' })
  @ApiConflictResponse({ description: 'Username already exists' })
  register(@Body() dto: RegisterDto) {
    return this.auth.register(dto.username, dto.password, dto.role);
  }

  @Post('login')
  @ApiOperation({ summary: 'Login with username and password' })
  @ApiResponse({
    status: 200,
    description: 'Login successful, returns JWT token',
    schema: {
      example: {
        accessToken: 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...',
      },
    },
  })
  @ApiUnauthorizedResponse({ description: 'Invalid credentials' })
  login(@Body() dto: LoginDto) {
    return this.auth.login(dto.username, dto.password);
  }

  @UseGuards(JwtAuthGuard)
  @ApiBearerAuth('JWT-auth')
  @Get('me')
  @ApiOperation({ summary: 'Get current user information' })
  @ApiResponse({
    status: 200,
    description: 'Returns current user data',
    schema: {
      example: {
        sub: '550e8400-e29b-41d4-a716-446655440000',
        username: 'john_doe',
        role: 'student',
      },
    },
  })
  @ApiUnauthorizedResponse({ description: 'Unauthorized' })
  me(@CurrentUser() user: any) {
    return { sub: user?.sub, username: user?.username, role: user?.role };
  }
}

