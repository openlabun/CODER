import { useState } from 'react';
import client from '../api/client';
import './AIAssistantModal.css';

const AIAssistantModal = ({ onClose, onApplyIdea, onApplyTestCases }) => {
    const [activeTab, setActiveTab] = useState('ideas');
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');

    // For generating ideas
    const [topic, setTopic] = useState('');
    const [difficulty, setDifficulty] = useState('');
    const [ideas, setIdeas] = useState([]);

    // For generating test cases
    const [challengeDesc, setChallengeDesc] = useState('');
    const [inputFormat, setInputFormat] = useState('');
    const [outputFormat, setOutputFormat] = useState('');
    const [testCases, setTestCases] = useState(null);

    const handleGenerateIdeas = async () => {
        if (!topic.trim()) {
            setError('Please enter a topic');
            return;
        }

        setLoading(true);
        setError('');
        setIdeas([]);

        try {
            const response = await client.post('/ai/generate-challenge-ideas', {
                topic,
                difficulty: difficulty || undefined,
                count: 3
            });
            setIdeas(response.data.ideas);
        } catch (err) {
            setError('Failed to generate ideas. Please try again.');
            console.error(err);
        } finally {
            setLoading(false);
        }
    };

    const handleGenerateTestCases = async () => {
        if (!challengeDesc.trim() || !inputFormat.trim() || !outputFormat.trim()) {
            setError('Please fill in all fields');
            return;
        }

        setLoading(true);
        setError('');
        setTestCases(null);

        try {
            const response = await client.post('/ai/generate-test-cases', {
                challengeDescription: challengeDesc,
                inputFormat,
                outputFormat,
                publicCount: 2,
                hiddenCount: 3
            });
            setTestCases(response.data);
        } catch (err) {
            setError('Failed to generate test cases. Please try again.');
            console.error(err);
        } finally {
            setLoading(false);
        }
    };

    const applyIdea = (idea) => {
        onApplyIdea(idea);
        onClose();
    };

    const applyTestCases = () => {
        if (testCases) {
            onApplyTestCases(testCases);
            onClose();
        }
    };

    return (
        <div className="modal-overlay" onClick={onClose}>
            <div className="ai-modal-content" onClick={(e) => e.stopPropagation()}>
                <div className="ai-modal-header">
                    <h2>✨ AI Creative Assistant</h2>
                    <button onClick={onClose} className="close-btn">×</button>
                </div>

                <div className="ai-tabs">
                    <button
                        className={activeTab === 'ideas' ? 'ai-tab active' : 'ai-tab'}
                        onClick={() => setActiveTab('ideas')}
                    >
                        💡 Generate Ideas
                    </button>
                    <button
                        className={activeTab === 'testcases' ? 'ai-tab active' : 'ai-tab'}
                        onClick={() => setActiveTab('testcases')}
                    >
                        🧪 Generate Test Cases
                    </button>
                </div>

                {error && <div className="ai-error-message">{error}</div>}

                <div className="ai-modal-body">
                    {activeTab === 'ideas' && (
                        <div className="ai-section">
                            <div className="ai-input-group">
                                <label>Topic or Category *</label>
                                <input
                                    type="text"
                                    value={topic}
                                    onChange={(e) => setTopic(e.target.value)}
                                    placeholder="e.g., Binary Search, Dynamic Programming, Trees"
                                    disabled={loading}
                                />
                            </div>

                            <div className="ai-input-group">
                                <label>Difficulty (Optional)</label>
                                <select value={difficulty} onChange={(e) => setDifficulty(e.target.value)} disabled={loading}>
                                    <option value="">Any</option>
                                    <option value="easy">Easy</option>
                                    <option value="medium">Medium</option>
                                    <option value="hard">Hard</option>
                                </select>
                            </div>

                            <button
                                onClick={handleGenerateIdeas}
                                disabled={loading}
                                className="ai-generate-btn"
                            >
                                {loading ? '🔄 Generating...' : '✨ Generate Ideas'}
                            </button>

                            {ideas.length > 0 && (
                                <div className="ai-results">
                                    <h3>Generated Ideas</h3>
                                    {ideas.map((idea, idx) => (
                                        <div key={idx} className="ai-idea-card">
                                            <div className="ai-idea-header">
                                                <h4>{idea.title}</h4>
                                                <span className={`difficulty-badge ${idea.difficulty}`}>
                                                    {idea.difficulty}
                                                </span>
                                            </div>
                                            <p className="ai-idea-description">{idea.description}</p>
                                            <div className="ai-idea-tags">
                                                {idea.tags?.map((tag, i) => (
                                                    <span key={i} className="tag">{tag}</span>
                                                ))}
                                            </div>
                                            <button
                                                onClick={() => applyIdea(idea)}
                                                className="ai-apply-btn"
                                            >
                                                ✅ Use This Idea
                                            </button>
                                        </div>
                                    ))}
                                </div>
                            )}
                        </div>
                    )}

                    {activeTab === 'testcases' && (
                        <div className="ai-section">
                            <div className="ai-input-group">
                                <label>Challenge Description *</label>
                                <textarea
                                    value={challengeDesc}
                                    onChange={(e) => setChallengeDesc(e.target.value)}
                                    placeholder="Describe the problem..."
                                    rows="4"
                                    disabled={loading}
                                />
                            </div>

                            <div className="ai-input-group">
                                <label>Input Format *</label>
                                <textarea
                                    value={inputFormat}
                                    onChange={(e) => setInputFormat(e.target.value)}
                                    placeholder="Describe the input format..."
                                    rows="3"
                                    disabled={loading}
                                />
                            </div>

                            <div className="ai-input-group">
                                <label>Output Format *</label>
                                <textarea
                                    value={outputFormat}
                                    onChange={(e) => setOutputFormat(e.target.value)}
                                    placeholder="Describe the expected output..."
                                    rows="3"
                                    disabled={loading}
                                />
                            </div>

                            <button
                                onClick={handleGenerateTestCases}
                                disabled={loading}
                                className="ai-generate-btn"
                            >
                                {loading ? '🔄 Generating...' : '🧪 Generate Test Cases'}
                            </button>

                            {testCases && (
                                <div className="ai-results">
                                    <h3>Generated Test Cases</h3>

                                    <div className="ai-testcases-section">
                                        <h4>📖 Public Test Cases ({testCases.publicTestCases?.length || 0})</h4>
                                        {testCases.publicTestCases?.map((tc, idx) => (
                                            <div key={idx} className="ai-testcase-card">
                                                <strong>{tc.name}</strong>
                                                <div className="ai-testcase-io">
                                                    <div>
                                                        <label>Input:</label>
                                                        <pre>{tc.input}</pre>
                                                    </div>
                                                    <div>
                                                        <label>Output:</label>
                                                        <pre>{tc.output}</pre>
                                                    </div>
                                                </div>
                                            </div>
                                        ))}
                                    </div>

                                    <div className="ai-testcases-section">
                                        <h4>🔒 Hidden Test Cases ({testCases.hiddenTestCases?.length || 0})</h4>
                                        {testCases.hiddenTestCases?.map((tc, idx) => (
                                            <div key={idx} className="ai-testcase-card">
                                                <strong>{tc.name}</strong>
                                                <div className="ai-testcase-io">
                                                    <div>
                                                        <label>Input:</label>
                                                        <pre>{tc.input}</pre>
                                                    </div>
                                                    <div>
                                                        <label>Output:</label>
                                                        <pre>{tc.output}</pre>
                                                    </div>
                                                </div>
                                            </div>
                                        ))}
                                    </div>

                                    <button
                                        onClick={applyTestCases}
                                        className="ai-apply-all-btn"
                                    >
                                        ✅ Apply All Test Cases
                                    </button>
                                </div>
                            )}
                        </div>
                    )}
                </div>

                <div className="ai-modal-footer">
                    <p className="ai-disclaimer">
                        ⚠️ AI-generated content should be reviewed and validated before publishing.
                    </p>
                </div>
            </div>
        </div>
    );
};

export default AIAssistantModal;
