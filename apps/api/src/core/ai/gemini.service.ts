import { GoogleGenerativeAI } from '@google/generative-ai';

export class GeminiService {
  private genAI: GoogleGenerativeAI;
  private model: any;

  constructor() {
    const apiKey = process.env.GEMINI_API_KEY || '';
    this.genAI = new GoogleGenerativeAI(apiKey);
    // Using gemini-flash-latest as it is confirmed to be available in the API
    this.model = this.genAI.getGenerativeModel({ model: 'gemini-flash-latest' });
  }
  async generateChallengeIdeas(topic: string, difficulty?: string, count: number = 3) {
    const prompt = `You are an expert programming instructor creating coding challenges.

Generate ${count} programming challenge ideas about "${topic}"${difficulty ? ` with ${difficulty} difficulty` : ''}.

For each challenge, provide:
1. Title (concise and descriptive)
2. Description (clear problem statement, 2-3 paragraphs)
3. Input Format (how data is provided)
4. Output Format (expected output structure)
5. Constraints (limits on input size, ranges, etc.)
6. Difficulty (easy, medium, or hard)
7. Tags (relevant topics, max 5)
8. Test Cases (2 public examples, 2 hidden edge cases)

Format your response as a JSON array of objects with these exact keys:
[
  {
    "title": "Challenge Title",
    "description": "Problem description...",
    "inputFormat": "Input format description...",
    "outputFormat": "Output format description...",
    "constraints": "Constraints description...",
    "difficulty": "medium",
    "tags": ["arrays", "sorting"],
    "publicTestCases": [{"name": "Example 1", "input": "...", "output": "..."}],
    "hiddenTestCases": [{"name": "Hidden 1", "input": "...", "output": "..."}]
  }
]

Make sure the challenges are educational, well-defined, and suitable for a competitive programming platform.
Return ONLY the JSON array, no additional text.`;

    try {
      const result = await this.model.generateContent(prompt);
      const response = await result.response;
      const text = response.text();

      // Extract JSON from response (sometimes the model adds markdown code blocks)
      const jsonMatch = text.match(/\[[\s\S]*\]/);
      if (!jsonMatch) {
        throw new Error('Invalid response format from AI');
      }

      const ideas = JSON.parse(jsonMatch[0]);
      return ideas;
    } catch (error) {
      console.error('Error generating challenge ideas:', error);
      if (error instanceof Error) {
        console.error('Error message:', error.message);
        console.error('Error stack:', error.stack);
      }
      throw new Error('Failed to generate challenge ideas');
    }
  }

  async generateTestCases(
    challengeDescription: string,
    inputFormat: string,
    outputFormat: string,
    publicCount: number = 2,
    hiddenCount: number = 3
  ) {
    const prompt = `You are an expert at creating test cases for programming challenges.

Challenge Description:
${challengeDescription}

Input Format:
${inputFormat}

Output Format:
${outputFormat}

Generate ${publicCount} public test cases (simple examples) and ${hiddenCount} hidden test cases (edge cases, larger inputs).

Format your response as a JSON object with these exact keys:
{
  "publicTestCases": [
    {
      "name": "Example 1",
      "input": "input data here",
      "output": "expected output here"
    }
  ],
  "hiddenTestCases": [
    {
      "name": "Edge Case 1",
      "input": "input data here",
      "output": "expected output here"
    }
  ]
}

Guidelines:
-- Public test cases should be simple and help students understand the problem
-- Hidden test cases should cover edge cases, boundary conditions, and larger inputs
-- Ensure all outputs are correct and match the expected format
-- Make test cases diverse and comprehensive

Return ONLY the JSON object, no additional text.`;

    try {
      const result = await this.model.generateContent(prompt);
      const response = await result.response;
      const text = response.text();

      // Extract JSON from response
      const jsonMatch = text.match(/\{[\s\S]*\}/);
      if (!jsonMatch) {
        throw new Error('Invalid response format from AI');
      }

      const testCases = JSON.parse(jsonMatch[0]);
      return testCases;
    } catch (error) {
      console.error('Error generating test cases:', error);
      throw new Error('Failed to generate test cases');
    }
  }
}
