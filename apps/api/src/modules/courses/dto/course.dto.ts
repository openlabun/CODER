export class CreateCourseDto {
    name!: string;
    code!: string;
    period!: string;
    groupNumber!: number;
}

export class EnrollStudentDto {
    studentId!: string;
}

export class AssignChallengeDto {
    challengeId!: string;
}
