import { useState, useEffect } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import client from '../api/client';
import AIAssistantModal from '../components/AIAssistantModal';
import './CreateChallenge.css';
import './Challenges.css'; // Imported to style the Card Preview correctly

const CreateChallenge = () => {
    const navigate = useNavigate();
    const location = useLocation();
    const [activeTab, setActiveTab] = useState('basic');
    const [showPreview, setShowPreview] = useState(false);
    const [showAIModal, setShowAIModal] = useState(false);

    // Get courseId from URL query params
    const queryParams = new URLSearchParams(location.search);
    const courseId = queryParams.get('courseId');

    const [formData, setFormData] = useState({
        title: '',
        description: '',
        difficulty: 'medium',
        timeLimit: 1000,
        memoryLimit: 256,
        tags: [],
        inputFormat: '',
        outputFormat: '',
        constraints: '',
        status: 'draft',
        assignedCourses: [],
        courseId: courseId || null
    });

    const [publicTestCases, setPublicTestCases] = useState([]);
    const [hiddenTestCases, setHiddenTestCases] = useState([]);
    const [newTag, setNewTag] = useState('');
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');
    const [successMessage, setSuccessMessage] = useState('');
    const [courses, setCourses] = useState([]);

    useEffect(() => {
        const fetchCourses = async () => {
            try {
                const { data } = await client.get('/courses');
                setCourses(data.items || data);
            } catch (err) {
                console.error('Error fetching courses:', err);
            }
        };
        fetchCourses();
    }, []);

    const predefinedTags = [
        'arrays', 'strings', 'math', 'hashing', 'greedy', 'dynamic-programming',
        'trees', 'graphs', 'sorting', 'searching', 'recursion', 'backtracking',
        'two-pointers', 'sliding-window', 'stack', 'queue', 'linked-list'
    ];

    const handleChange = (e) => {
        const { name, value } = e.target;
        setFormData(prev => ({ ...prev, [name]: value }));
    };

    const addPublicTestCase = () => {
        setPublicTestCases([...publicTestCases, { input: '', output: '', name: `Example ${publicTestCases.length + 1}` }]);
    };

    const addHiddenTestCase = () => {
        setHiddenTestCases([...hiddenTestCases, { input: '', output: '', name: `Hidden ${hiddenTestCases.length + 1}` }]);
    };

    const updateTestCase = (index, field, value, isPublic) => {
        const cases = isPublic ? [...publicTestCases] : [...hiddenTestCases];
        cases[index][field] = value;
        isPublic ? setPublicTestCases(cases) : setHiddenTestCases(cases);
    };

    const removeTestCase = (index, isPublic) => {
        const cases = isPublic ? [...publicTestCases] : [...hiddenTestCases];
        cases.splice(index, 1);
        isPublic ? setPublicTestCases(cases) : setHiddenTestCases(cases);
    };

    const addTag = (tag) => {
        if (tag && !formData.tags.includes(tag)) {
            setFormData(prev => ({ ...prev, tags: [...prev.tags, tag] }));
            setNewTag('');
        }
    };

    const removeTag = (tag) => {
        setFormData(prev => ({ ...prev, tags: prev.tags.filter(t => t !== tag) }));
    };

    const validateForm = () => {
        if (!formData.title || !formData.description) {
            setError('Title and description are required');
            return false;
        }
        if (hiddenTestCases.length < 3) {
            setError('At least 3 hidden test cases are required');
            return false;
        }
        // Validate all test cases have both input and output
        const allCases = [...publicTestCases, ...hiddenTestCases];
        for (let tc of allCases) {
            if (!tc.input.trim() || !tc.output.trim()) {
                setError('All test cases must have both input and output');
                return false;
            }
        }
        return true;
    };

    const handleSubmit = async (status) => {
        setError('');
        setSuccessMessage('');
        if (status === 'published' && !validateForm()) {
            return;
        }

        setLoading(true);
        try {
            const payload = {
                ...formData,
                status,
                publicTestCases,
                hiddenTestCases
            };

            const response = await client.post('/challenges', payload);

            // If challenge was created for a course, assign it
            if (formData.courseId && response.data.id) {
                await client.post(`/courses/${formData.courseId}/challenges`, {
                    challengeId: response.data.id
                });
            }
            
            setSuccessMessage(`¡Reto publicado exitosamente!`);
            window.scrollTo(0, 0);

            // Optional: Reset form here to allow creating another one freely
            setFormData({
                title: '',
                description: '',
                difficulty: 'medium',
                timeLimit: 1000,
                memoryLimit: 256,
                tags: [],
                inputFormat: '',
                outputFormat: '',
                constraints: '',
                status: 'draft',
                assignedCourses: [],
                courseId: courseId || null
            });
            setPublicTestCases([]);
            setHiddenTestCases([]);
            setShowPreview(false);
            
        } catch (err) {
            setError('Failed to create challenge. Please try again.');
            console.error(err);
        } finally {
            setLoading(false);
        }
    };

    const handleApplyIdea = (idea) => {
        setFormData(prev => ({
            ...prev,
            title: idea.title,
            description: idea.description,
            difficulty: idea.difficulty || 'medium',
            tags: idea.tags || [],
            inputFormat: idea.inputFormat || '',
            outputFormat: idea.outputFormat || '',
            constraints: idea.constraints || ''
        }));

        if (idea.publicTestCases) {
            setPublicTestCases(idea.publicTestCases);
        }
        if (idea.hiddenTestCases) {
            setHiddenTestCases(idea.hiddenTestCases);
        }

        // If we have test cases, switch to test cases tab to show them
        if (idea.publicTestCases?.length > 0 || idea.hiddenTestCases?.length > 0) {
            setActiveTab('testcases');
        } else {
            setActiveTab('details');
        }
    };

    const handleApplyTestCases = (cases) => {
        if (cases.publicTestCases) setPublicTestCases(cases.publicTestCases);
        if (cases.hiddenTestCases) setHiddenTestCases(cases.hiddenTestCases);
        setActiveTab('testcases');
    };

    return (
        <div className="create-challenge-page">
            <div className="page-header">
                <h1>Create New Challenge</h1>
                <p className="subtitle">Design a comprehensive coding challenge</p>
                <button
                    className="btn-ai-assist"
                    onClick={() => setShowAIModal(true)}
                >
                    ✨ AI Assistant
                </button>
            </div>

            {showAIModal && (
                <AIAssistantModal
                    onClose={() => setShowAIModal(false)}
                    onApplyIdea={handleApplyIdea}
                    onApplyTestCases={handleApplyTestCases}
                />
            )}

            {error && <div className="error-message">{error}</div>}
            {successMessage && <div className="success-message" style={{ background: 'rgba(40, 167, 69, 0.1)', border: '1px solid var(--success-color)', color: 'var(--success-color)', padding: '15px', borderRadius: '8px', marginBottom: '20px' }}>{successMessage}</div>}

            {!showPreview ? (
                <>
                    <div className="tabs">
                        <button className={activeTab === 'basic' ? 'tab active' : 'tab'} onClick={() => setActiveTab('basic')}>
                            📝 Basic Info
                        </button>
                        <button className={activeTab === 'details' ? 'tab active' : 'tab'} onClick={() => setActiveTab('details')}>
                            📋 Details & Format
                        </button>
                        <button className={activeTab === 'testcases' ? 'tab active' : 'tab'} onClick={() => setActiveTab('testcases')}>
                            🧪 Test Cases
                        </button>
                        <button className={activeTab === 'settings' ? 'tab active' : 'tab'} onClick={() => setActiveTab('settings')}>
                            ⚙️ Settings
                        </button>
                    </div>

                    <div className="tab-content">
                        {activeTab === 'basic' && (
                            <div className="form-section">
                                <h2>Basic Information</h2>

                                <div className="form-group">
                                    <label htmlFor="title">Challenge Title *</label>
                                    <input
                                        type="text"
                                        id="title"
                                        name="title"
                                        value={formData.title}
                                        onChange={handleChange}
                                        placeholder="e.g., Two Sum Problem"
                                        required
                                    />
                                </div>

                                <div className="form-group">
                                    <label htmlFor="description">Problem Statement *</label>
                                    <textarea
                                        id="description"
                                        name="description"
                                        value={formData.description}
                                        onChange={handleChange}
                                        placeholder="Describe the problem clearly..."
                                        rows="12"
                                        required
                                    />
                                    <small>Supports Markdown formatting</small>
                                </div>

                                <div className="form-row">
                                    <div className="form-group">
                                        <label htmlFor="difficulty">Difficulty *</label>
                                        <select id="difficulty" name="difficulty" value={formData.difficulty} onChange={handleChange}>
                                            <option value="easy">Easy</option>
                                            <option value="medium">Medium</option>
                                            <option value="hard">Hard</option>
                                        </select>
                                    </div>

                                    <div className="form-group">
                                        <label htmlFor="timeLimit">Time Limit (ms) *</label>
                                        <input type="number" id="timeLimit" name="timeLimit" value={formData.timeLimit} onChange={handleChange} min="100" max="10000" />
                                    </div>

                                    <div className="form-group">
                                        <label htmlFor="memoryLimit">Memory Limit (MB) *</label>
                                        <input type="number" id="memoryLimit" name="memoryLimit" value={formData.memoryLimit} onChange={handleChange} min="64" max="512" />
                                    </div>
                                </div>

                                <div className="form-group">
                                    <label>Tags</label>
                                    <div className="tags-container">
                                        {formData.tags.map(tag => (
                                            <span key={tag} className="tag">
                                                {tag}
                                                <button type="button" onClick={() => removeTag(tag)}>×</button>
                                            </span>
                                        ))}
                                    </div>
                                    <div className="tag-input-group">
                                        <input
                                            type="text"
                                            value={newTag}
                                            onChange={(e) => setNewTag(e.target.value)}
                                            onKeyPress={(e) => e.key === 'Enter' && (e.preventDefault(), addTag(newTag))}
                                            placeholder="Type and press Enter"
                                        />
                                        <button type="button" onClick={() => addTag(newTag)} className="btn-add-tag">Add</button>
                                    </div>
                                    <div className="predefined-tags">
                                        {predefinedTags.filter(t => !formData.tags.includes(t)).map(tag => (
                                            <button key={tag} type="button" onClick={() => addTag(tag)} className="predefined-tag">
                                                + {tag}
                                            </button>
                                        ))}
                                    </div>
                                </div>
                            </div>
                        )}

                        {activeTab === 'details' && (
                            <div className="form-section">
                                <h2>Input/Output Format & Constraints</h2>

                                <div className="form-group">
                                    <label htmlFor="inputFormat">Input Format</label>
                                    <textarea
                                        id="inputFormat"
                                        name="inputFormat"
                                        value={formData.inputFormat}
                                        onChange={handleChange}
                                        placeholder="Describe the input format..."
                                        rows="4"
                                    />
                                </div>

                                <div className="form-group">
                                    <label htmlFor="outputFormat">Output Format</label>
                                    <textarea
                                        id="outputFormat"
                                        name="outputFormat"
                                        value={formData.outputFormat}
                                        onChange={handleChange}
                                        placeholder="Describe the expected output format..."
                                        rows="4"
                                    />
                                </div>

                                <div className="form-group">
                                    <label htmlFor="constraints">Constraints</label>
                                    <textarea
                                        id="constraints"
                                        name="constraints"
                                        value={formData.constraints}
                                        onChange={handleChange}
                                        placeholder="e.g., 1 ≤ n ≤ 10^5, -10^9 ≤ arr[i] ≤ 10^9"
                                        rows="4"
                                    />
                                </div>
                            </div>
                        )}

                        {activeTab === 'testcases' && (
                            <div className="form-section">
                                <h2>Test Cases</h2>

                                <div className="testcases-section">
                                    <div className="testcase-header">
                                        <h3>📖 Public Test Cases</h3>
                                        <p>Students can see these examples</p>
                                        <button type="button" onClick={addPublicTestCase} className="btn-add">+ Add Public Case</button>
                                    </div>

                                    {publicTestCases.map((tc, idx) => (
                                        <div key={idx} className="testcase-item">
                                            <div className="testcase-title">
                                                <input
                                                    type="text"
                                                    value={tc.name}
                                                    onChange={(e) => updateTestCase(idx, 'name', e.target.value, true)}
                                                    placeholder="Case name"
                                                />
                                                <button type="button" onClick={() => removeTestCase(idx, true)} className="btn-remove">🗑️</button>
                                            </div>
                                            <div className="testcase-io">
                                                <div className="io-group">
                                                    <label>Input</label>
                                                    <textarea
                                                        value={tc.input}
                                                        onChange={(e) => updateTestCase(idx, 'input', e.target.value, true)}
                                                        placeholder="Input data..."
                                                        rows="3"
                                                    />
                                                </div>
                                                <div className="io-group">
                                                    <label>Expected Output</label>
                                                    <textarea
                                                        value={tc.output}
                                                        onChange={(e) => updateTestCase(idx, 'output', e.target.value, true)}
                                                        placeholder="Expected output..."
                                                        rows="3"
                                                    />
                                                </div>
                                            </div>
                                        </div>
                                    ))}
                                </div>

                                <div className="testcases-section">
                                    <div className="testcase-header">
                                        <h3>🔒 Hidden Test Cases</h3>
                                        <p>Used for evaluation (min. 3 required)</p>
                                        <button type="button" onClick={addHiddenTestCase} className="btn-add">+ Add Hidden Case</button>
                                    </div>

                                    {hiddenTestCases.map((tc, idx) => (
                                        <div key={idx} className="testcase-item hidden">
                                            <div className="testcase-title">
                                                <input
                                                    type="text"
                                                    value={tc.name}
                                                    onChange={(e) => updateTestCase(idx, 'name', e.target.value, false)}
                                                    placeholder="Case name"
                                                />
                                                <button type="button" onClick={() => removeTestCase(idx, false)} className="btn-remove">🗑️</button>
                                            </div>
                                            <div className="testcase-io">
                                                <div className="io-group">
                                                    <label>Input</label>
                                                    <textarea
                                                        value={tc.input}
                                                        onChange={(e) => updateTestCase(idx, 'input', e.target.value, false)}
                                                        placeholder="Input data..."
                                                        rows="3"
                                                    />
                                                </div>
                                                <div className="io-group">
                                                    <label>Expected Output</label>
                                                    <textarea
                                                        value={tc.output}
                                                        onChange={(e) => updateTestCase(idx, 'output', e.target.value, false)}
                                                        placeholder="Expected output..."
                                                        rows="3"
                                                    />
                                                </div>
                                            </div>
                                        </div>
                                    ))}

                                    {hiddenTestCases.length < 3 && (
                                        <div className="warning-box">
                                            ⚠️ You need at least {3 - hiddenTestCases.length} more hidden test case(s) to publish this challenge.
                                        </div>
                                    )}
                                </div>
                            </div>
                        )}

                        {activeTab === 'settings' && (
                            <div className="form-section">
                                <h2>Publication & Course Assignment</h2>

                                <div className="form-group">
                                    <label>Status</label>
                                    <div className="status-options">
                                        <label className="radio-option">
                                            <input
                                                type="radio"
                                                name="status"
                                                value="draft"
                                                checked={formData.status === 'draft'}
                                                onChange={handleChange}
                                            />
                                            <span>📝 Draft</span>
                                            <small>Save for later, not visible to students</small>
                                        </label>
                                        <label className="radio-option">
                                            <input
                                                type="radio"
                                                name="status"
                                                value="published"
                                                checked={formData.status === 'published'}
                                                onChange={handleChange}
                                            />
                                            <span>✅ Published</span>
                                            <small>Visible to assigned courses</small>
                                        </label>
                                        <label className="radio-option">
                                            <input
                                                type="radio"
                                                name="status"
                                                value="archived"
                                                checked={formData.status === 'archived'}
                                                onChange={handleChange}
                                            />
                                            <span>📦 Archived</span>
                                            <small>Hidden but preserved</small>
                                        </label>
                                    </div>
                                </div>

                                <div className="info-box" style={{ background: 'rgba(200, 16, 46, 0.05)', border: '1px solid var(--primary-color)' }}>
                                    <h3 style={{ color: 'var(--primary-color)' }}>📚 Asignar a Curso</h3>
                                    <p>Selecciona un curso al que quieres asignar este reto (Opcional).</p>
                                    <div className="form-group" style={{ marginTop: '15px' }}>
                                        <select
                                            name="courseId"
                                            value={formData.courseId || ''}
                                            onChange={handleChange}
                                        >
                                            <option value="">-- No asignar a ningún curso por ahora --</option>
                                            {courses.map(course => (
                                                <option key={course.id} value={course.id}>
                                                    {course.title}
                                                </option>
                                            ))}
                                        </select>
                                    </div>
                                </div>
                            </div>
                        )}
                    </div>

                    <div className="form-actions">
                        <button type="button" onClick={() => navigate('/challenges')} className="btn-secondary">
                            Cancel
                        </button>
                        <button type="button" onClick={() => setShowPreview(true)} className="btn-preview">
                            👁️ Preview
                        </button>
                        <button type="button" onClick={() => handleSubmit('draft')} disabled={loading} className="btn-draft">
                            {loading ? 'Saving...' : '💾 Save as Draft'}
                        </button>
                        <button type="button" onClick={() => handleSubmit('published')} disabled={loading} className="btn-publish">
                            {loading ? 'Publishing...' : '🚀 Publish Challenge'}
                        </button>
                    </div>
                </>
            ) : (
                <div className="preview-container">
                    <div className="preview-header">
                        <h2>Preview</h2>
                        <button onClick={() => setShowPreview(false)} className="btn-secondary">← Back to Edit</button>
                    </div>

                    <div className="preview-section-title" style={{ marginTop: '20px', marginBottom: '15px', color: 'var(--text-muted)' }}>
                        <h3>Card Preview (How it looks in the list)</h3>
                    </div>

                    <div className="challenge-card" style={{ maxWidth: '400px', marginBottom: '40px' }}>
                        <div className="challenge-header" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', marginBottom: '15px' }}>
                            <h3 className="challenge-title" style={{ fontSize: '1.2rem', fontWeight: 'bold', color: 'var(--text-color)', margin: 0 }}>
                                {formData.title || 'Untitled Challenge'}
                            </h3>
                            <span className={`difficulty-badge ${(formData.difficulty || 'medium').toLowerCase()}`}>
                                {formData.difficulty || 'Medium'}
                            </span>
                        </div>
                        <p className="challenge-desc" style={{ color: 'var(--text-muted)', fontSize: '0.9rem', marginBottom: '20px', display: '-webkit-box', WebkitLineClamp: 2, WebkitBoxOrient: 'vertical', overflow: 'hidden' }}>
                            {formData.description || 'No description provided'}
                        </p>
                        <div className="btn-solve" style={{ display: 'inline-block', width: '100%', padding: '12px', textAlign: 'center', background: 'transparent', border: '1px solid var(--primary-color)', color: 'var(--primary-color)', borderRadius: '4px', fontWeight: 600, textTransform: 'uppercase', letterSpacing: '0.5px' }}>
                            Solve Challenge
                        </div>
                    </div>

                    <div className="preview-section-title" style={{ marginBottom: '15px', color: 'var(--text-muted)' }}>
                        <h3>Detailed Preview (How it looks to the student)</h3>
                    </div>

                    <div className="challenge-preview">
                        <div className="preview-title">
                            <h1>{formData.title || 'Untitled Challenge'}</h1>
                            <span className={`difficulty-badge ${formData.difficulty}`}>{formData.difficulty}</span>
                        </div>

                        {formData.tags.length > 0 && (
                            <div className="preview-tags">
                                {formData.tags.map(tag => <span key={tag} className="tag">{tag}</span>)}
                            </div>
                        )}

                        <div className="preview-section">
                            <h3>Problem Statement</h3>
                            <div className="preview-content">{formData.description || 'No description provided'}</div>
                        </div>

                        {formData.inputFormat && (
                            <div className="preview-section">
                                <h3>Input Format</h3>
                                <div className="preview-content">{formData.inputFormat}</div>
                            </div>
                        )}

                        {formData.outputFormat && (
                            <div className="preview-section">
                                <h3>Output Format</h3>
                                <div className="preview-content">{formData.outputFormat}</div>
                            </div>
                        )}

                        {formData.constraints && (
                            <div className="preview-section">
                                <h3>Constraints</h3>
                                <div className="preview-content">{formData.constraints}</div>
                            </div>
                        )}

                        {publicTestCases.length > 0 && (
                            <div className="preview-section">
                                <h3>Examples</h3>
                                {publicTestCases.map((tc, idx) => (
                                    <div key={idx} className="example-case">
                                        <h4>{tc.name}</h4>
                                        <div className="example-io">
                                            <div>
                                                <strong>Input:</strong>
                                                <pre>{tc.input}</pre>
                                            </div>
                                            <div>
                                                <strong>Output:</strong>
                                                <pre>{tc.output}</pre>
                                            </div>
                                        </div>
                                    </div>
                                ))}
                            </div>
                        )}

                        <div className="preview-meta">
                            <div>⏱️ Time Limit: {formData.timeLimit}ms</div>
                            <div>💾 Memory Limit: {formData.memoryLimit}MB</div>
                            <div>🔒 Hidden Test Cases: {hiddenTestCases.length}</div>
                        </div>
                    </div>

                    <div className="preview-actions">
                        <button onClick={() => setShowPreview(false)} className="btn-secondary">
                            ✏️ Edit Challenge
                        </button>
                        <button onClick={() => handleSubmit(formData.status)} disabled={loading} className="btn-publish">
                            {loading ? 'Publishing...' : '✅ Confirm & Publish'}
                        </button>
                    </div>
                </div>
            )}
        </div>
    );
};

export default CreateChallenge;
