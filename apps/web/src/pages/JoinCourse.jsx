import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import client from '../api/client';
import './JoinCourse.css';

const JoinCourse = () => {
    const navigate = useNavigate();
    const [enrollmentCode, setEnrollmentCode] = useState('');
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');
    const [success, setSuccess] = useState('');

    const handleJoinWithCode = async (e) => {
        e.preventDefault();
        setError('');
        setSuccess('');

        if (!enrollmentCode.trim()) {
            setError('Please enter an enrollment code');
            return;
        }

        setLoading(true);
        try {
            await client.post('/courses/enroll', { enrollmentCode });
            setSuccess('Successfully enrolled in the course!');
            setTimeout(() => {
                navigate('/courses');
            }, 2000);
        } catch (err) {
            setError(err.response?.data?.message || 'Invalid enrollment code or course not found');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="join-course-page">
            <div className="page-header">
                <h1>Join a Course</h1>
                <p className="subtitle">Enter the enrollment code provided by your professor</p>
            </div>

            <div className="join-course-container">
                <div className="join-method-card">
                    <div className="method-icon">🔑</div>
                    <h2>Enrollment Code</h2>
                    <p>Enter the unique code shared by your professor to join the course</p>

                    {error && <div className="error-message">{error}</div>}
                    {success && <div className="success-message">{success}</div>}

                    <form onSubmit={handleJoinWithCode} className="join-form">
                        <div className="form-group">
                            <label htmlFor="enrollmentCode">Enrollment Code</label>
                            <input
                                type="text"
                                id="enrollmentCode"
                                value={enrollmentCode}
                                onChange={(e) => setEnrollmentCode(e.target.value.toUpperCase())}
                                placeholder="e.g., CS101-20251G1"
                                className="code-input"
                                disabled={loading}
                                autoFocus
                            />
                            <small>Format: COURSE-PERIODG# (e.g., CS101-20251G1)</small>
                        </div>

                        <div className="form-actions">
                            <button
                                type="button"
                                onClick={() => navigate('/courses')}
                                className="btn-secondary"
                                disabled={loading}
                            >
                                Cancel
                            </button>
                            <button
                                type="submit"
                                className="btn-primary"
                                disabled={loading}
                            >
                                {loading ? 'Joining...' : 'Join Course'}
                            </button>
                        </div>
                    </form>
                </div>

                <div className="info-section">
                    <h3>📚 How to Join</h3>
                    <ol>
                        <li>Get the enrollment code from your professor</li>
                        <li>Enter the code in the field above</li>
                        <li>Click "Join Course" to enroll</li>
                        <li>Access course materials and challenges</li>
                    </ol>

                    <div className="help-box">
                        <h4>Need Help?</h4>
                        <p>If you don't have an enrollment code, contact your professor or course administrator.</p>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default JoinCourse;
