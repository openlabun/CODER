import { ApiProperty } from '@nestjs/swagger';

export class RegisterDto {
    @ApiProperty({ example: 'john_doe', description: 'Username' })
    username!: string;

    @ApiProperty({ example: 'password123', description: 'Password' })
    password!: string;

    @ApiProperty({ example: 'student', enum: ['student', 'professor', 'admin'], description: 'User role' })
    role!: 'student' | 'professor' | 'admin';
}
