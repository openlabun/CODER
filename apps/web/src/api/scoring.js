import client from './client';

export const getUsersScoresByExam = async (examId, courseId) => {
  const params = {};
  if (courseId) params.course_id = courseId;
  const { data } = await client.get(`/submissions/scores/exam/${examId}/users`, { params });
  return data;
};

export const getExamScoringForUser = async (examId, userId) => {
  const { data } = await client.get(`/submissions/scores/exam/${examId}/user/${userId}`);
  return data;
};

export const getExamsScoresByUser = async (userId) => {
  const { data } = await client.get(`/submissions/scores/user/${userId}`);
  return data;
};

export const getExamScoreItemDetails = async (examScoreId) => {
  const { data } = await client.get(`/submissions/scores/exam-score/${examScoreId}`);
  return data;
};
