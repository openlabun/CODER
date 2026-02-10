import { Course } from '../entities/course.entity';

export interface ICourseRepo {
    save(course: Course): Promise<void>;
    update(course: Course): Promise<void>;
    findById(id: string): Promise<Course | null>;
    findAll(): Promise<Course[]>;
    list(): Promise<Course[]>;
    findByProfessor(professorId: string): Promise<Course[]>;
    findByStudent(studentId: string): Promise<Course[]>;

    // Student management
    addStudent(courseId: string, studentId: string): Promise<void>;
    removeStudent(courseId: string, studentId: string): Promise<void>;
    getStudents(courseId: string): Promise<string[]>;

    // Challenge management
    addChallenge(courseId: string, challengeId: string): Promise<void>;
    removeChallenge(courseId: string, challengeId: string): Promise<void>;
    getChallenges(courseId: string): Promise<string[]>;
}
