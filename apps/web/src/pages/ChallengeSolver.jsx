import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import Editor from '@monaco-editor/react';
import client from '../api/client';
import { useAuth } from '../context/AuthContext';
import './ChallengeSolver.css';

const ChallengeSolver = () => {
    const { id } = useParams();
    const navigate = useNavigate();
    const { user } = useAuth();
    const [challenge, setChallenge] = useState(null);
    const [code, setCode] = useState('');
    const [language, setLanguage] = useState('javascript');
    const [output, setOutput] = useState('');
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchChallenge = async () => {
            try {
                const { data } = await client.get(`/challenges/${id}`);
                setChallenge(data);
                // Set default code template based on language
                setCode('// Write your solution here\n');
            } catch (error) {
                console.error('Error fetching challenge:', error);
            } finally {
                setLoading(false);
            }
        };
        fetchChallenge();
    }, [id]);

    const handleSubmit = async () => {
        try {
            const { data } = await client.post('/submissions', {
                challengeId: id,
                code,
                language
            });
            setOutput(data.output || 'Submission received');
        } catch (error) {
            setOutput('Error submitting solution');
        }
    };

    if (loading) return <div>Loading...</div>;
    if (!challenge) return <div>Challenge not found</div>;

    if (user?.role === 'professor' || user?.role === 'admin') {
        return (
            <div className="solver-container">
                <div className="problem-description full-width">
                    <h2>{challenge.title}</h2>
                    <p>{challenge.description}</p>
                    <div className="professor-actions" style={{ textAlign: 'center', marginTop: '2rem' }}>
                        <button
                            onClick={() => navigate(`/challenges/edit/${id}`)}
                            className="btn-primary"
                        >
                            ✏️ Editar Reto
                        </button>
                    </div>
                </div>
            </div>
        );
    }

    return (
        <div className="solver-container">
            <div className="problem-description">
                <h2>{challenge.title}</h2>
                <p>{challenge.description}</p>

                {challenge.inputFormat && (
                    <div className="challenge-section">
                        <h3>Input Format</h3>
                        <p>{challenge.inputFormat}</p>
                    </div>
                )}

                {challenge.outputFormat && (
                    <div className="challenge-section">
                        <h3>Output Format</h3>
                        <p>{challenge.outputFormat}</p>
                    </div>
                )}

                {challenge.constraints && (
                    <div className="challenge-section">
                        <h3>Constraints</h3>
                        <p>{challenge.constraints}</p>
                    </div>
                )}

                {challenge.publicTestCases && challenge.publicTestCases.length > 0 && (
                    <div className="challenge-section">
                        <h3>Public Test Cases</h3>
                        {challenge.publicTestCases.map((testCase, index) => (
                            <div key={index} className="test-case">
                                <h4>{testCase.name}</h4>
                                <div className="test-case-content">
                                    <div>
                                        <strong>Input:</strong>
                                        <pre>{testCase.input}</pre>
                                    </div>
                                    <div>
                                        <strong>Expected Output:</strong>
                                        <pre>{testCase.output}</pre>
                                    </div>
                                </div>
                            </div>
                        ))}
                    </div>
                )}
            </div>
            <div className="editor-container">
                <div className="editor-header">
                    <select value={language} onChange={(e) => setLanguage(e.target.value)}>
                        <option value="javascript">JavaScript</option>
                        <option value="python">Python</option>
                        <option value="cpp">C++</option>
                        <option value="java">Java</option>
                    </select>
                    <button onClick={handleSubmit} className="btn-primary">Submit</button>
                </div>
                <Editor
                    height="80vh"
                    theme="vs-dark"
                    language={language}
                    value={code}
                    onChange={(value) => setCode(value)}
                />
                <div className="output-panel">
                    <h3>Output</h3>
                    <pre>{output}</pre>
                </div>
            </div>
        </div>
    );
};

export default ChallengeSolver;
