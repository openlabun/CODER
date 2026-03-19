import { CreateChallengeUseCase } from './create-challenge.use-case';
import { IChallengeRepo } from '../interfaces/challenge.repo';
import { Challenge } from '../entities/challenge.entity';

describe('CreateChallengeUseCase', () => {
  let useCase: CreateChallengeUseCase;
  let repo: jest.Mocked<IChallengeRepo>;

  beforeEach(() => {
    repo = {
      save: jest.fn(),
      findById: jest.fn(),
      list: jest.fn(),
    };

    useCase = new CreateChallengeUseCase(repo);
  });

  it('crea challenge con datos base y persiste en repositorio', async () => {
    repo.save.mockResolvedValue(undefined);

    const result = await useCase.execute({
      id: 'ch-1',
      title: 'Two Sum',
      description: 'Find indices of two numbers that add up to target',
    });

    expect(repo.save).toHaveBeenCalledTimes(1);
    expect(repo.save).toHaveBeenCalledWith(result);
    expect(result).toBeInstanceOf(Challenge);
    expect(result.id).toBe('ch-1');
    expect(result.title).toBe('Two Sum');
    expect(result.description).toBe('Find indices of two numbers that add up to target');
    expect(result.status).toBe('draft');
  });

  it('aplica valores por defecto cuando no se envian dificultad ni limites', async () => {
    repo.save.mockResolvedValue(undefined);

    const result = await useCase.execute({
      id: 'ch-2',
      title: 'Palindrome',
      description: 'Check if string is palindrome',
    });

    expect(result.difficulty).toBe('medium');
    expect(result.timeLimit).toBe(1500);
    expect(result.memoryLimit).toBe(256);
    expect(result.tags).toEqual([]);
  });

  it('respeta dificultad, limites y tags personalizados', async () => {
    repo.save.mockResolvedValue(undefined);

    const result = await useCase.execute({
      id: 'ch-3',
      title: 'Graph Traversal',
      description: 'Run BFS over adjacency list',
      difficulty: 'hard',
      timeLimit: 3000,
      memoryLimit: 512,
      tags: ['graphs', 'bfs'],
    });

    expect(result.difficulty).toBe('hard');
    expect(result.timeLimit).toBe(3000);
    expect(result.memoryLimit).toBe(512);
    expect(result.tags).toEqual(['graphs', 'bfs']);
  });

  it('normaliza campos de texto con trim', async () => {
    repo.save.mockResolvedValue(undefined);

    const result = await useCase.execute({
      id: 'ch-4',
      title: '   Valid Title   ',
      description: '   Desc with spaces   ',
    });

    expect(result.title).toBe('Valid Title');
    expect(result.description).toBe('Desc with spaces');
  });

  it('lanza error cuando el titulo es demasiado corto y no persiste', async () => {
    await expect(
      useCase.execute({
        id: 'ch-5',
        title: 'ab',
        description: 'short title',
      }),
    ).rejects.toThrow('Title must be at least 3 characters');

    expect(repo.save).not.toHaveBeenCalled();
  });
});
