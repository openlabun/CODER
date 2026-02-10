import { useState, useEffect, useContext } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { AuthContext } from '../context/AuthContext';
import client from '../api/client';
import './Courses.css';

const Courses = () => {
    const [courses, setCourses] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    const { user } = useContext(AuthContext);
    const navigate = useNavigate();

    useEffect(() => {
        const fetchCourses = async () => {
            try {
                const { data } = await client.get('/courses');
                setCourses(data);
            } catch (err) {
                setError('Failed to load courses');
            } finally {
                setLoading(false);
            }
        };
        fetchCourses();
    }, []);

    if (loading) return <div className="loading">Loading courses...</div>;
    if (error) return <div className="error">{error}</div>;

    return (
        <div className="courses-page">
            <div className="page-header">
                <h1>My Courses</h1>
                <div className="header-actions">
                    {user?.role === 'student' && (
                        <>
                            <button onClick={() => navigate('/courses/browse')} className="btn-secondary" style={{ marginRight: '1rem' }}>
                                🔍 Browse Courses
                            </button>
                            <button onClick={() => navigate('/courses/join')} className="btn-primary">
                                🔑 Join Course
                            </button>
                        </>
                    )}
                    {user?.role === 'professor' && (
                        <button onClick={() => navigate('/courses/create')} className="btn-primary">
                            ➕ Create Course
                        </button>
                    )}
                </div>
            </div>
            {courses.length === 0 ? (
                <div className="empty-state">
                    <h3>📚 No Courses Available</h3>
                    {user?.role === 'student' ? (
                        <p>You haven't joined any courses yet. Click "Browse Courses" to find open courses or "Join Course" to enroll with a code.</p>
                    ) : (
                        <p>You haven't created any courses yet. Click "Create Course" to get started.</p>
                    )}
                </div>
            ) : (
                <div className="courses-grid">
                    {courses.map((course) => (
                        <Link key={course.id} to={`/courses/${course.id}`} className="course-card">
                            <h3>{course.name}</h3>
                            <p>Code: {course.code}</p>
                            <p>Period: {course.period}</p>
                            {user?.role === 'professor' && course.enrollmentCode && (
                                <p className="enrollment-code">🔑 {course.enrollmentCode}</p>
                            )}
                            {user?.role === 'professor' && (
                                <button
                                    onClick={(e) => {
                                        e.preventDefault();
                                        navigate(`/courses/edit/${course.id}`);
                                    }}
                                    className="btn-edit"
                                    style={{ marginTop: '1rem', width: '100%', padding: '0.5rem', background: 'rgba(255, 255, 255, 0.1)', border: 'none', borderRadius: '4px', color: 'white', cursor: 'pointer' }}
                                >
                                    ✏️ Edit Course
                                </button>
                            )}
                        </Link>
                    ))}
                </div>
            )}
        </div>
    );
};

export default Courses;
