import { useState, useEffect, useContext } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import Editor from '@monaco-editor/react';
import client from '../api/client';
import { AuthContext } from '../context/AuthContext';
import { 
    Clock, 
    Target, 
    ChevronLeft, 
    Play, 
    Send, 
    Code, 
    FileText, 
    Info, 
    AlertTriangle,
    CheckCircle2
} from 'lucide-react';
import './ChallengeSolver.css';
import './Dashboard.css';

const ChallengeSolver = () => {
    const { id } = useParams();
    const navigate = useNavigate();
    const { user } = useContext(AuthContext);
    const [challenge, setChallenge] = useState(null);
    const [code, setCode] = useState('');
    const [language, setLanguage] = useState('javascript');
    const [output, setOutput] = useState('');
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchChallenge = async () => {
            setLoading(true);
            try {
                // Ensure ID is passed correctly and try to fetch
                const response = await client.get(`/challenges/${id}`);
                setChallenge(response.data);
                // Set default code template
                setCode('// Escribe tu solución aquí\n');
            } catch (error) {
                console.error('Error fetching challenge:', error);
                // If 404, maybe it's not published
                if (error.response?.status === 404) {
                    console.warn("Challenge not found or not accessible");
                }
            } finally {
                setLoading(false);
            }
        };
        if (id) fetchChallenge();
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

    if (loading) {
        return (
            <div className="dashboard-loading">
                <div className="loader-orbit">
                    <div className="orbit-dot"></div>
                </div>
                <p>Preparando desafío...</p>
            </div>
        );
    }
    
    if (!challenge) {
        return (
            <div className="dashboard-loading error">
                <div className="error-icon" style={{fontSize: '3rem', marginBottom: '1rem'}}>🎯</div>
                <h2>Desafío no encontrado</h2>
                <p>No se pudo cargar la información del reto o no tienes permisos.</p>
                <button 
                    onClick={() => navigate(-1)} 
                    className="btn-retry" 
                    style={{marginTop: '2rem'}}
                >
                    Volver Atrás
                </button>
            </div>
        );
    }

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
