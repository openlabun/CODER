import client from './client';

export const getCourseExams = async (courseId) => {
  const res = await client.get(`/exams/course/${courseId}`);
  return res.data;
};

export const createExam = async (examData) => {
  const res = await client.post('/exams', examData);
  return res.data;
};

export const getExamDetails = async (examId) => {
  const res = await client.get(`/exams/${examId}`);
  return res.data;
};
