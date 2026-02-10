import { Submission } from '../entities/submission.entity';
import { ISubmissionRepo } from '../interfaces/submission.repo';
import { ISubmissionQueue } from '../interfaces/submission.queue';
import { IChallengeRepo } from '../../challenges/interfaces/challenge.repo';

type Input = { challengeId: string; userId: string; code: string; language: string; examId?: string };
type Output = Submission;

import * as fs from 'fs';
import * as path from 'path';

export class CreateSubmissionUseCase {
  constructor(
    private readonly repo: ISubmissionRepo,
    private readonly queue: ISubmissionQueue | undefined,
    private readonly challengeRepo: IChallengeRepo,
  ) { }

  async execute(input: Input): Promise<Output> {
    if (!input.challengeId) throw new Error('challengeId is required');
    if (!input.userId) throw new Error('userId is required');

    // ✅ Validar que exista el reto
    const challenge = await this.challengeRepo.findById(input.challengeId);
    if (!challenge) {
      throw new Error('challenge_not_found');
    }

    const sub = Submission.create({
      challengeId: input.challengeId,
      userId: input.userId,
      code: input.code,
      language: input.language,
      examId: input.examId,
    });

    await this.repo.save(sub);

    // 📝 Escribir archivos en disco para que el Worker los encuentre
    this.writeSubmissionFiles(sub);

    if (this.queue) await this.queue.enqueue(sub.id);

    return sub;
  }

  private writeSubmissionFiles(sub: Submission) {
    // __dirname está en src/core/Submission/use-cases
    // Queremos llegar a src/core/Submission/{id}
    const dir = path.join(__dirname, '..', sub.id);

    if (!fs.existsSync(dir)) {
      fs.mkdirSync(dir, { recursive: true });
    }

    // Create 'code' subdirectory
    const codeDir = path.join(dir, 'code');
    if (!fs.existsSync(codeDir)) {
      fs.mkdirSync(codeDir, { recursive: true });
    }

    let ext = 'txt';
    const lang = sub.language.toLowerCase();
    if (lang === 'python') ext = 'py';
    else if (lang === 'javascript' || lang === 'node') ext = 'js';
    else if (lang === 'java') ext = 'java';
    else if (lang === 'cpp' || lang === 'c++') ext = 'cpp';

    // Java requires capital M in Main.java for public class Main
    const filename = lang === 'java' ? `Main.${ext}` : `main.${ext}`;

    // Escribir el código en el subdirectorio 'code'
    fs.writeFileSync(path.join(codeDir, filename), sub.code);

    // Escribir meta.json en el directorio raíz de la submission
    const meta = {
      language: sub.language,
      codeFile: filename,
    };
    fs.writeFileSync(path.join(dir, 'meta.json'), JSON.stringify(meta, null, 2));
  }
}
