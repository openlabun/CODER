import { Exam } from '../entities/exam.entity';

export interface IExamRepo {
    save(exam: Exam): Promise<void>;
    findById(id: string): Promise<Exam | null>;
    findByCourseId(courseId: string): Promise<Exam[]>;
    addChallengeToExam(examId: string, challengeId: string, points: number, order: number): Promise<void>;
    getExamChallenges(examId: string): Promise<any[]>; // Returns challenge details with points
}
