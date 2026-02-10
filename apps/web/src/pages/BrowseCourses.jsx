import { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import client from '../api/client';
import './Courses.css';

const BrowseCourses = () => {
    const navigate = useNavigate();
    const [courses, setCourses] = useState([]);
    const [filteredCourses, setFilteredCourses] = useState([]);
    const [searchTerm, setSearchTerm] = useState('');
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');

    // Modal state
    const [showModal, setShowModal] = useState(false);
    const [selectedCourse, setSelectedCourse] = useState(null);
    const [enrollmentCode, setEnrollmentCode] = useState('');
    const [enrolling, setEnrolling] = useState(false);
    const [enrollError, setEnrollError] = useState('');

    useEffect(() => {
        const fetchCourses = async () => {
            try {
                const { data } = await client.get('/courses/browse');
                setCourses(data);
                setFilteredCourses(data);
            } catch (err) {
                console.error(err);
                setError('Failed to load courses');
            } finally {
                setLoading(false);
            }
        };
        fetchCourses();
    }, []);

    useEffect(() => {
        const results = courses.filter(course =>
            course.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
            course.code.toLowerCase().includes(searchTerm.toLowerCase())
        );
        setFilteredCourses(results);
    }, [searchTerm, courses]);

    const handleJoinClick = (course) => {
        setSelectedCourse(course);
        setEnrollmentCode('');
        setEnrollError('');
        setShowModal(true);
    };

    const handleJoinCourse = async (e) => {
        e.preventDefault();
        setEnrollError('');

        if (!enrollmentCode.trim()) {
            setEnrollError('Please enter an enrollment code');
            return;
        }

        setEnrolling(true);
        try {
            await client.post('/courses/enroll', { enrollmentCode });
            setShowModal(false);
            alert(`Successfully enrolled in ${selectedCourse.name}!`);
            navigate('/courses');
        } catch (err) {
            setEnrollError(err.response?.data?.message || 'Invalid enrollment code');
        } finally {
            setEnrolling(false);
        }
    };

    if (loading) return <div className="loading">Loading courses...</div>;
    if (error) return <div className="error">{error}</div>;

    return (
        <div className="courses-page">
            <div className="page-header">
                <h1>Browse Courses</h1>
                <input
                    type="text"
                    placeholder="Search courses..."
                    value={searchTerm}
                    onChange={(e) => setSearchTerm(e.target.value)}
                    className="search-input"
                />
            </div>

            <div className="courses-grid">
                {filteredCourses.map((course) => (
                    <div key={course.id} className="course-card">
                        <h3>{course.name}</h3>
                        <p>Code: {course.code}</p>
                        <p>Period: {course.period}</p>
                        <p>Professor: {course.professorName || 'Unknown'}</p>
                        <button
                            onClick={() => handleJoinClick(course)}
                            className="btn-primary"
                            style={{ marginTop: '10px', width: '100%' }}
                        >
                            Join Course
                        </button>
                    </div>
                ))}
            </div>

            {/* Enrollment Modal */}
            {showModal && (
                <div className="modal-overlay" onClick={() => setShowModal(false)}>
                    <div className="modal-content" onClick={(e) => e.stopPropagation()}>
                        <div className="modal-header">
                            <h2>🔑 Enrollment Code</h2>
                            <button onClick={() => setShowModal(false)} className="close-btn">×</button>
                        </div>
                        <p>Enter the unique code shared by your professor to join the course</p>

                        {enrollError && <div className="error-message">{enrollError}</div>}

                        <form onSubmit={handleJoinCourse}>
                            <div className="form-group">
                                <label htmlFor="enrollmentCode">Enrollment Code</label>
                                <input
                                    type="text"
                                    id="enrollmentCode"
                                    value={enrollmentCode}
                                    onChange={(e) => setEnrollmentCode(e.target.value.toUpperCase())}
                                    placeholder="e.g., CS101-20251G1"
                                    disabled={enrolling}
                                    autoFocus
                                />
                                <small>Format: COURSE-PERIODG# (e.g., CS101-20251G1)</small>
                            </div>

                            <div className="form-actions">
                                <button
                                    type="button"
                                    onClick={() => setShowModal(false)}
                                    className="btn-secondary"
                                    disabled={enrolling}
                                >
                                    Cancel
                                </button>
                                <button
                                    type="submit"
                                    className="btn-primary"
                                    disabled={enrolling}
                                >
                                    {enrolling ? 'Joining...' : 'Join Course'}
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
            )}
        </div>
    );
};

export default BrowseCourses;
