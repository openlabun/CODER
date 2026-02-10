import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import client from '../api/client';
import './Courses.css';
import './CourseActions.css';

const CourseStudents = () => {
    const { id } = useParams();
    const navigate = useNavigate();
    const [students, setStudents] = useState([]);
    const [course, setCourse] = useState(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const [courseRes, studentsRes] = await Promise.all([
                    client.get(`/courses/${id}`),
                    client.get(`/courses/${id}/students`)
                ]);
                setCourse(courseRes.data);
                setStudents(studentsRes.data.students || []);
                console.log('Students data:', studentsRes.data.students);
            } catch (error) {
                console.error('Error fetching students:', error);
            } finally {
                setLoading(false);
            }
        };
        fetchData();
    }, [id]);

    if (loading) return <div className="loading">Loading...</div>;

    return (
        <div className="course-students-page">
            <div className="page-header">
                <button onClick={() => navigate(`/courses/${id}`)} className="btn-back">
                    ← Back to Course
                </button>
                <h1>Students - {course?.name}</h1>
            </div>

            {students.length === 0 ? (
                <div className="empty-state">
                    <div className="empty-state-icon">👥</div>
                    <h3 className="empty-state-title">No Students Enrolled</h3>
                    <p className="empty-state-description">
                        No students have enrolled in this course yet.
                    </p>
                    <p className="empty-state-hint">
                        Share the enrollment code: <strong>{course?.enrollmentCode}</strong>
                    </p>
                </div>
            ) : (
                <div className="students-grid">
                    {students.map((student, index) => (
                        <div key={student.id} className="student-card">
                            <div className="student-avatar">
                                {student.username?.charAt(0).toUpperCase() || (index + 1)}
                            </div>
                            <div className="student-info">
                                <div className="student-name">{student.username || 'Unknown'}</div>
                            </div>
                        </div>
                    ))}
                </div>
            )}

            <div className="students-stats">
                <div className="stat-card">
                    <div className="stat-value">{students.length}</div>
                    <div className="stat-label">Total Students</div>
                </div>
            </div>
        </div>
    );
};

export default CourseStudents;
