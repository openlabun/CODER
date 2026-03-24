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
    const [language, setLanguage] = useState('python');
    const [output, setOutput] = useState('');
    const [loading, setLoading] = useState(true);
    const [submitting, setSubmitting] = useState(false);

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

    const sleep = (ms) => new Promise((resolve) => setTimeout(resolve, ms));

    const formatResultsOutput = (submission, results) => {
        if (!results.length) {
            return `Submission ${submission?.ID || submission?.id || ''} creada. Aun no hay resultados.`;
        }

        const lines = [];
        lines.push(`Submission: ${submission?.ID || submission?.id || 'N/A'}`);
        lines.push(`Lenguaje: ${submission?.Language || submission?.language || 'python'}`);
        lines.push('');
        lines.push('Resultados por caso de prueba:');

        results.forEach((result, index) => {
            const status = (result.Status || result.status || 'unknown').toLowerCase();
            const errorMessage = result.ErrorMessage || result.errorMessage || '';
            let line = `- Caso ${index + 1}: ${status}`;
            if (errorMessage) {
                line += ` | error: ${errorMessage}`;
            }
            lines.push(line);
        });

        return lines.join('\n');
    };

    const handleSubmit = async () => {
        if (!code.trim()) {
            setOutput('No puedes enviar codigo vacio.');
            return;
        }

        setSubmitting(true);
        setOutput('Enviando solucion...');

        try {
            const sessionId = localStorage.getItem('session_id') || undefined;
            const payload = {
                challengeID: id,
                code,
                language,
            };

            if (sessionId) {
                payload.sessionID = sessionId;
            }

            const { data } = await client.post('/submissions', payload);
            const submissionId = data?.id || data?.ID;

            if (!submissionId) {
                setOutput('La API no retorno un ID de submission.');
                return;
            }

            setOutput('Solucion enviada. Ejecutando pruebas...');

            const maxAttempts = 40;
            for (let attempt = 0; attempt < maxAttempts; attempt += 1) {
                const res = await client.get(`/submissions/${submissionId}`);
                const submission = res?.data?.Submission || res?.data?.submission;
                const results = res?.data?.Results || res?.data?.results || [];

                if (Array.isArray(results) && results.length > 0) {
                    const hasPendingResults = results.some((result) => {
                        const status = (result.Status || result.status || '').toLowerCase();
                        return status === 'queued' || status === 'running';
                    });

                    if (!hasPendingResults) {
                        setOutput(formatResultsOutput(submission, results));
                        return;
                    }
                }

                await sleep(1500);
            }

            setOutput('La solucion fue enviada, pero la evaluacion sigue en proceso. Revisa tu historial de envios en unos segundos.');
        } catch (error) {
            const apiMessage = error?.response?.data?.error || error?.message;
            setOutput(apiMessage ? `Error submitting solution: ${apiMessage}` : 'Error submitting solution');
        } finally {
            setSubmitting(false);
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
                        <option value="python">Python</option>
                    </select>
                    <button onClick={handleSubmit} className="btn-primary" disabled={submitting}>
                        {submitting ? 'Evaluando...' : 'Submit'}
                    </button>
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
