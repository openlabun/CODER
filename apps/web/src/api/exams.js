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

export const updateExam = async (examId, examData) => {
  const res = await client.patch(`/exams/${examId}`, examData);
  return res.data;
};

export const deleteExam = async (examId) => {
  const res = await client.delete(`/exams/${examId}`);
  return res.data;
};

export const toggleExamVisibility = async (examId) => {
  const res = await client.post(`/exams/${examId}/visibility`);
  return res.data;
};

export const closeExam = async (examId) => {
  const res = await client.post(`/exams/${examId}/close`);
  return res.data;
};
