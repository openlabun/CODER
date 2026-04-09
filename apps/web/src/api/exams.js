import client from './client';

const unwrapPayload = (payload) => {
  if (!payload || typeof payload !== 'object') {
    return payload;
  }

  if (Object.prototype.hasOwnProperty.call(payload, 'body')) {
    return payload.body;
  }

  return payload;
};

export const getCourseExams = async (courseId) => {
  const res = await client.get(`/exams/course/${courseId}`);
  return unwrapPayload(res.data);
};

export const createExam = async (examData) => {
  const res = await client.post('/exams', examData);
  return unwrapPayload(res.data);
};

export const getExamDetails = async (examId) => {
  const res = await client.get(`/exams/${examId}`);
  return unwrapPayload(res.data);
};

export const getOwnedExams = async () => {
  const res = await client.get('/exams/owned');
  return unwrapPayload(res.data);
};

export const getPublicExams = async () => {
  const res = await client.get('/exams/public');
  return unwrapPayload(res.data);
};

export const changeExamVisibility = async (examId, visibility, professorIds = []) => {
  const payload = {
    visibility,
  };

  if (visibility === 'teachers') {
    payload.professor_ids = professorIds;
  }

  const res = await client.post(`/exams/${examId}/visibility`, payload);
  return unwrapPayload(res.data);
};

export const toggleExamVisibility = async (examId) => {
  const res = await client.post(`/exams/${examId}/visibility`);
  return unwrapPayload(res.data);
};

export const closeExam = async (examId) => {
  const res = await client.post(`/exams/${examId}/close`);
  return unwrapPayload(res.data);
};

export const deleteExam = async (examId) => {
  const res = await client.delete(`/exams/${examId}`);
  return unwrapPayload(res.data);
};

export const shareExamWithProfessor = async (examId, professorId) => {
  const res = await client.post(`/exams/${examId}/share`, {
    professor_id: professorId,
  });
  return unwrapPayload(res.data);
};

export const unshareExamWithProfessor = async (examId, professorId) => {
  const res = await client.delete(`/exams/${examId}/share/${professorId}`);
  return unwrapPayload(res.data);
};

export const createExamSession = async (examId, userId) => {
  const res = await client.post('/submissions/sessions', { user_id: userId, exam_id: examId });
  return unwrapPayload(res.data);
};

export const getChallengesByExam = async (examId) => {
  const res = await client.get('/challenges', { params: { examId } });
  return unwrapPayload(res.data);
};
