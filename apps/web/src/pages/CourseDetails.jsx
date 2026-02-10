import { useState, useEffect, useContext } from 'react';
import { useParams, Link, useNavigate } from 'react-router-dom';
import client from '../api/client';
import { getCourseExams } from '../api/exams';
import { AuthContext } from '../context/AuthContext';
import './Courses.css';
import './CourseActions.css';

const CourseDetails = () => {
    const { id } = useParams();
    const navigate = useNavigate();
    const { user } = useContext(AuthContext);
    const [course, setCourse] = useState(null);
    const [exams, setExams] = useState([]);
    const [challenges, setChallenges] = useState([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const [courseRes, examsData, challengesRes] = await Promise.all([
                    client.get(`/courses/${id}`),
                    getCourseExams(id),
                    client.get(`/courses/${id}/challenges`)
                ]);
                setCourse(courseRes.data);
                setExams(examsData);
                setChallenges(challengesRes.data.challenges || []);
            } catch (err) {
                console.error(err);
            } finally {
                setLoading(false);
            }
        };
        fetchData();
    }, [id]);

    if (loading) return <div className="loading">Loading...</div>;

    const isProfessor = user?.role === 'professor' || user?.role === 'admin';

    return (
        <div className="course-details-page">
            <div className="course-header">
                <div>
                    <h1>{course?.name}</h1>
                    <p className="course-meta">
                        {course?.code} - {course?.period} - Group {course?.groupNumber}
                    </p>
                </div>

                {isProfessor && (
                    <div className="course-actions">
                        <button
                            onClick={() => navigate(`/courses/${id}/students`)}
                            className="btn-action btn-students"
                        >
                            <span>👥</span> View Students
                        </button>
                        <button
                            onClick={() => navigate(`/challenges/create?courseId=${id}`)}
                            className="btn-action btn-create-challenge"
                        >
                            <span>➕</span> Create Challenge
                        </button>
                    </div>
                )}
            </div>

            <section className="challenges-section">
                <h2>Challenges</h2>
                {challenges.length === 0 ? (
                    <div className="empty-state">
                        <div className="empty-state-icon">🎯</div>
                        <h3 className="empty-state-title">No Challenges Yet</h3>
                        <p className="empty-state-description">
                            No challenges have been assigned to this course yet.
                        </p>
                        {isProfessor && (
                            <p className="empty-state-hint">
                                Click "Create Challenge" above to add a challenge to this course.
                            </p>
                        )}
                    </div>
                ) : (
                    <div className="challenges-grid">
                        {challenges.map(challenge => (
                            <Link key={challenge.id} to={`/challenges/${challenge.id}`} className="challenge-card">
                                <div className="challenge-header">
                                    <h3>{challenge.title}</h3>
                                    <span className={`difficulty-badge ${challenge.difficulty}`}>
                                        {challenge.difficulty}
                                    </span>
                                </div>
                                <p className="challenge-description">
                                    {challenge.description?.substring(0, 100)}...
                                </p>
                                <div className="challenge-meta">
                                    <span>⏱️ {challenge.timeLimit}ms</span>
                                    <span>💾 {challenge.memoryLimit}MB</span>
                                </div>
                            </Link>
                        ))}
                    </div>
                )}
            </section>

            <section className="exams-section">
                <h2>Exams</h2>
                {exams.length === 0 ? (
                    <div className="empty-state">
                        <div className="empty-state-icon">📝</div>
                        <h3 className="empty-state-title">No Exams Yet</h3>
                        <p className="empty-state-description">
                            There are no exams scheduled for this course at the moment.
                        </p>
                        <p className="empty-state-hint">
                            Your professor will create exams here when they're ready.
                        </p>
                    </div>
                ) : (
                    <ul className="exams-list">
                        {exams.map(exam => (
                            <li key={exam.id} className="exam-item">
                                <h3>{exam.title}</h3>
                                <p>Duration: {exam.durationMinutes} mins</p>
                                <p>Start: {new Date(exam.startTime).toLocaleString()}</p>
                                <Link to={`/exam/${exam.id}`} className="btn-start-exam">
                                    Start Exam
                                </Link>
                            </li>
                        ))}
                    </ul>
                )}
            </section>
        </div>
    );
};

export default CourseDetails;
