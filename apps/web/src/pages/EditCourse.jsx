import { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import client from '../api/client';
import './CreateCourse.css'; // Reuse CreateCourse styles

const EditCourse = () => {
    const { id } = useParams();
    const navigate = useNavigate();
    const [formData, setFormData] = useState({
        name: '',
        code: '',
        period: '',
        groupNumber: '',
    });
    const [loading, setLoading] = useState(true);
    const [saving, setSaving] = useState(false);
    const [error, setError] = useState('');

    useEffect(() => {
        const fetchCourse = async () => {
            try {
                const { data } = await client.get(`/courses/${id}`);
                setFormData({
                    name: data.name,
                    code: data.code,
                    period: data.period,
                    groupNumber: data.groupNumber,
                });
            } catch (err) {
                setError('Failed to load course details');
            } finally {
                setLoading(false);
            }
        };
        fetchCourse();
    }, [id]);

    const handleSubmit = async (e) => {
        e.preventDefault();
        setSaving(true);
        setError('');

        try {
            await client.post(`/courses/${id}`, {
                ...formData,
                groupNumber: parseInt(formData.groupNumber),
            });
            navigate('/courses');
        } catch (err) {
            setError(err.response?.data?.message || 'Failed to update course');
        } finally {
            setSaving(false);
        }
    };

    const handleChange = (e) => {
        setFormData({
            ...formData,
            [e.target.name]: e.target.value,
        });
    };

    if (loading) return <div className="loading">Loading...</div>;

    return (
        <div className="create-course-page">
            <div className="page-header">
                <h1>Edit Course</h1>
            </div>

            <form onSubmit={handleSubmit} className="create-course-form">
                {error && <div className="error-message">{error}</div>}

                <div className="form-group">
                    <label htmlFor="name">Course Name</label>
                    <input
                        type="text"
                        id="name"
                        name="name"
                        value={formData.name}
                        onChange={handleChange}
                        required
                        minLength={3}
                    />
                </div>

                <div className="form-row">
                    <div className="form-group">
                        <label htmlFor="code">Course Code</label>
                        <input
                            type="text"
                            id="code"
                            name="code"
                            value={formData.code}
                            onChange={handleChange}
                            required
                            minLength={2}
                        />
                    </div>

                    <div className="form-group">
                        <label htmlFor="period">Period</label>
                        <input
                            type="text"
                            id="period"
                            name="period"
                            value={formData.period}
                            onChange={handleChange}
                            placeholder="e.g. 2025-1"
                            required
                            minLength={4}
                        />
                    </div>

                    <div className="form-group">
                        <label htmlFor="groupNumber">Group Number</label>
                        <input
                            type="number"
                            id="groupNumber"
                            name="groupNumber"
                            value={formData.groupNumber}
                            onChange={handleChange}
                            required
                            min={1}
                        />
                    </div>
                </div>

                <div className="form-actions">
                    <button type="button" onClick={() => navigate('/courses')} className="btn-secondary">
                        Cancel
                    </button>
                    <button type="submit" className="btn-primary" disabled={saving}>
                        {saving ? 'Saving...' : 'Save Changes'}
                    </button>
                </div>
            </form>
        </div>
    );
};

export default EditCourse;
