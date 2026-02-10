import { Body, Controller, Get, Post, UseGuards } from '@nestjs/common';
import { ApiTags, ApiOperation, ApiResponse, ApiBearerAuth } from '@nestjs/swagger';
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
  @ApiResponse({ status: 201, description: 'User registered successfully' })
  @ApiResponse({ status: 400, description: 'Invalid input or user already exists' })
  register(@Body() dto: RegisterDto) {
    return this.auth.register(dto.username, dto.password, dto.role);
  }

  @Post('login')
  @ApiOperation({ summary: 'Login with username and password' })
  @ApiResponse({ status: 200, description: 'Login successful, returns JWT token' })
  @ApiResponse({ status: 401, description: 'Invalid credentials' })
  login(@Body() dto: LoginDto) {
    return this.auth.login(dto.username, dto.password);
  }

  @UseGuards(JwtAuthGuard)
  @ApiBearerAuth('JWT-auth')
  @Get('me')
  @ApiOperation({ summary: 'Get current user information' })
  @ApiResponse({ status: 200, description: 'Returns current user data' })
  @ApiResponse({ status: 401, description: 'Unauthorized' })
  me(@CurrentUser() user: any) {
    return { sub: user?.sub, username: user?.username, role: user?.role };
  }
}

