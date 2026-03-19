import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import client from '../api/client';
import './CreateCourse.css';

const CreateCourse = () => {
    const navigate = useNavigate();
    const currentYear = new Date().getFullYear().toString();
    const [formData, setFormData] = useState({
        name: '',
        code: '',
        period: `${currentYear}-1`, // Default period
        groupNumber: 1,
        description: '',
        color: '#00f0ff',
        enrollmentMethod: 'code',
        enrollmentCode: '',
        startDate: '',
        endDate: '',
        visibility: 'active'
    });
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');
    const [validationErrors, setValidationErrors] = useState({});

    const courseColors = [
        '#00f0ff', '#7000ff', '#ff0055', '#00ff9d',
        '#ffcc00', '#ff6b35', '#4ecdc4', '#95e1d3'
    ];

    const handleChange = (e) => {
        const { name, value } = e.target;
        setFormData(prev => ({ ...prev, [name]: value }));

        // Clear validation error for this field
        if (validationErrors[name]) {
            setValidationErrors(prev => ({ ...prev, [name]: '' }));
        }
    };

    const generateEnrollmentCode = () => {
        const code = formData.code.toUpperCase() + '-' +
            formData.period.split('-').join('') + 'G' + formData.groupNumber;
        setFormData(prev => ({ ...prev, enrollmentCode: code }));
    };

    const validateForm = () => {
        const errors = {};

        // Course code validation (letters + numbers only)
        if (!/^[A-Za-z0-9]+$/.test(formData.code)) {
            errors.code = 'Course code must contain only letters and numbers';
        }

        // Period validation (YYYY-X format)
        if (!/^\d{4}-[123]$/.test(formData.period)) {
            errors.period = 'El periodo debe tener el formato YYYY-X';
        }

        // Group number validation
        if (formData.groupNumber < 1) {
            errors.groupNumber = 'Group number must be a positive integer';
        }

        // Date validation
        if (formData.startDate && formData.endDate) {
            if (new Date(formData.startDate) >= new Date(formData.endDate)) {
                errors.endDate = 'End date must be after start date';
            }
        }

        setValidationErrors(errors);
        return Object.keys(errors).length === 0;
    };

    const handleSubmit = async (status) => {
        setError('');

        if (!validateForm()) {
            setError('Please fix the validation errors before submitting');
            return;
        }

        // Generate enrollment code if using code method and not set
        if (formData.enrollmentMethod === 'code' && !formData.enrollmentCode) {
            generateEnrollmentCode();
        }

        setLoading(true);
        try {
            const payload = {
                ...formData,
                status
            };
            await client.post('/courses', payload);
            navigate('/courses');
        } catch (err) {
            setError('Failed to create course. Please try again.');
            console.error(err);
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="create-course-page">
            <div className="page-header">
                <h1>Create New Course</h1>
                <p className="subtitle">Set up a comprehensive course for your students</p>
            </div>

            {error && <div className="error-message">{error}</div>}

            <form className="course-form">
                <div className="form-section">
                    <h2>Basic Information</h2>

                    <div className="form-group">
                        <label htmlFor="name">Course Name *</label>
                        <input
                            type="text"
                            id="name"
                            name="name"
                            value={formData.name}
                            onChange={handleChange}
                            placeholder="e.g., Algorithms and Data Structures"
                            required
                        />
                    </div>

                    <div className="form-row">
                        <div className="form-group">
                            <label htmlFor="code">Course Code *</label>
                            <input
                                type="text"
                                id="code"
                                name="code"
                                value={formData.code}
                                onChange={handleChange}
                                placeholder="e.g., CS101"
                                required
                                className={validationErrors.code ? 'error' : ''}
                            />
                            {validationErrors.code ? (
                                <small className="error-text">{validationErrors.code}</small>
                            ) : (
                                <small>Format: CS101, MAT202, INF251</small>
                            )}
                        </div>

                        <div className="form-group">
                            <label>Periodo Académico *</label>
                            <div className="period-selectors">
                                <select
                                    value={formData.period ? formData.period.split('-')[0] : new Date().getFullYear()}
                                    onChange={(e) => {
                                        const currentTerm = (formData.period && formData.period.includes('-')) ? formData.period.split('-')[1] : '1';
                                        handleChange({ target: { name: 'period', value: `${e.target.value}-${currentTerm}` }});
                                    }}
                                    className="year-select"
                                >
                                    {[...Array(6)].map((_, i) => {
                                        const year = new Date().getFullYear() - 1 + i;
                                        return <option key={year} value={year}>{year}</option>;
                                    })}
                                </select>
                                <select
                                    value={formData.period ? formData.period.split('-')[1] : '1'}
                                    onChange={(e) => {
                                        const currentYear = formData.period ? formData.period.split('-')[0] : new Date().getFullYear();
                                        handleChange({ target: { name: 'period', value: `${currentYear}-${e.target.value}` }});
                                    }}
                                    className="term-select"
                                >
                                    <option value="1">Primer Semestre (10)</option>
                                    <option value="2">Segundo Semestre (30)</option>
                                    <option value="3">Verano (20)</option>
                                </select>
                            </div>
                            {validationErrors.period ? (
                                <small className="error-text">{validationErrors.period}</small>
                            ) : (
                                <small>Seleccione el año y ciclo académico (Formato Interno: {formData.period})</small>
                            )}
                        </div>

                        <div className="form-group">
                            <label htmlFor="groupNumber">Group Number *</label>
                            <input
                                type="number"
                                id="groupNumber"
                                name="groupNumber"
                                value={formData.groupNumber}
                                onChange={handleChange}
                                min="1"
                                required
                                className={validationErrors.groupNumber ? 'error' : ''}
                            />
                            {validationErrors.groupNumber ? (
                                <small className="error-text">{validationErrors.groupNumber}</small>
                            ) : (
                                <small>Must be a positive integer</small>
                            )}
                        </div>
                    </div>

                    <div className="form-group">
                        <label htmlFor="description">Course Description (Optional)</label>
                        <textarea
                            id="description"
                            name="description"
                            value={formData.description}
                            onChange={handleChange}
                            placeholder="Provide a short overview of the course topics and learning goals..."
                            rows="5"
                        />
                        <small>Describe objectives, main topics, and academic purpose</small>
                    </div>
                </div>

                <div className="form-section">
                    <h2>Visual Identity</h2>

                    <div className="form-group">
                        <label>Course Color (Optional)</label>
                        <div className="color-picker">
                            {courseColors.map(color => (
                                <button
                                    key={color}
                                    type="button"
                                    className={`color-option ${formData.color === color ? 'selected' : ''}`}
                                    style={{ backgroundColor: color }}
                                    onClick={() => setFormData(prev => ({ ...prev, color }))}
                                    title={color}
                                />
                            ))}
                        </div>
                        <small>Choose a color to visually identify this course</small>
                    </div>
                </div>

                <div className="form-section">
                    <h2>Enrollment Settings</h2>

                    <div className="form-group">
                        <label>Enrollment Method *</label>
                        <div className="radio-group">
                            <label className="radio-option">
                                <input
                                    type="radio"
                                    name="enrollmentMethod"
                                    value="code"
                                    checked={formData.enrollmentMethod === 'code'}
                                    onChange={handleChange}
                                />
                                <div>
                                    <span>🔑 Enrollment Code</span>
                                    <small>Students join using a unique code</small>
                                </div>
                            </label>
                            <label className="radio-option">
                                <input
                                    type="radio"
                                    name="enrollmentMethod"
                                    value="link"
                                    checked={formData.enrollmentMethod === 'link'}
                                    onChange={handleChange}
                                />
                                <div>
                                    <span>🔗 Private Link</span>
                                    <small>Students join via invitation link</small>
                                </div>
                            </label>
                            <label className="radio-option">
                                <input
                                    type="radio"
                                    name="enrollmentMethod"
                                    value="automatic"
                                    checked={formData.enrollmentMethod === 'automatic'}
                                    onChange={handleChange}
                                />
                                <div>
                                    <span>⚙️ Automatic</span>
                                    <small>Institution-managed enrollment</small>
                                </div>
                            </label>
                        </div>
                    </div>

                    {formData.enrollmentMethod === 'code' && (
                        <div className="form-group">
                            <label htmlFor="enrollmentCode">Enrollment Code</label>
                            <div className="code-input-group">
                                <input
                                    type="text"
                                    id="enrollmentCode"
                                    name="enrollmentCode"
                                    value={formData.enrollmentCode}
                                    onChange={handleChange}
                                    placeholder="Will be auto-generated"
                                    readOnly
                                />
                                <button type="button" onClick={generateEnrollmentCode} className="btn-generate">
                                    Generate Code
                                </button>
                            </div>
                            <small>Students will use this code to join the course</small>
                        </div>
                    )}
                </div>

                <div className="form-section">
                    <h2>Schedule</h2>

                    <div className="form-row">
                        <div className="form-group">
                            <label htmlFor="startDate">Start Date</label>
                            <input
                                type="date"
                                id="startDate"
                                name="startDate"
                                value={formData.startDate}
                                onChange={handleChange}
                            />
                            <small>When the course becomes active</small>
                        </div>

                        <div className="form-group">
                            <label htmlFor="endDate">End Date</label>
                            <input
                                type="date"
                                id="endDate"
                                name="endDate"
                                value={formData.endDate}
                                onChange={handleChange}
                                className={validationErrors.endDate ? 'error' : ''}
                            />
                            {validationErrors.endDate ? (
                                <small className="error-text">{validationErrors.endDate}</small>
                            ) : (
                                <small>When the course ends</small>
                            )}
                        </div>
                    </div>
                </div>

                <div className="form-section">
                    <h2>Visibility</h2>

                    <div className="form-group">
                        <label>Course Visibility *</label>
                        <div className="radio-group">
                            <label className="radio-option">
                                <input
                                    type="radio"
                                    name="visibility"
                                    value="active"
                                    checked={formData.visibility === 'active'}
                                    onChange={handleChange}
                                />
                                <div>
                                    <span>✅ Active</span>
                                    <small>Visible to enrolled students</small>
                                </div>
                            </label>
                            <label className="radio-option">
                                <input
                                    type="radio"
                                    name="visibility"
                                    value="hidden"
                                    checked={formData.visibility === 'hidden'}
                                    onChange={handleChange}
                                />
                                <div>
                                    <span>👁️ Hidden / Draft</span>
                                    <small>Not visible to students yet</small>
                                </div>
                            </label>
                        </div>
                    </div>
                </div>

                <div className="form-actions">
                    <button type="button" onClick={() => navigate('/courses')} className="btn-secondary">
                        Cancel
                    </button>
                    <button
                        type="button"
                        onClick={() => handleSubmit('draft')}
                        disabled={loading}
                        className="btn-draft"
                    >
                        {loading ? 'Saving...' : '💾 Save as Draft'}
                    </button>
                    <button
                        type="button"
                        onClick={() => handleSubmit('published')}
                        disabled={loading}
                        className="btn-primary"
                    >
                        {loading ? 'Creating...' : '🚀 Create Course'}
                    </button>
                </div>
            </form>

            <div className="info-box">
                <h3>📚 Next Steps</h3>
                <p>After creating the course, you'll be able to:</p>
                <ul>
                    <li>Enroll students manually or share the enrollment code</li>
                    <li>Assign challenges to the course</li>
                    <li>Create exams and assessments</li>
                    <li>View course leaderboards and analytics</li>
                </ul>
            </div>
        </div>
    );
};

export default CreateCourse;
