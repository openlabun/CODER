import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { getExamDetails } from '../api/exams';
import client from '../api/client';
import './ChallengeSolver.css'; // Reuse styles

const ExamRunner = () => {
    const { id } = useParams();
    const navigate = useNavigate();
    const [exam, setExam] = useState(null);
    const [currentChallenge, setCurrentChallenge] = useState(null);
    const [code, setCode] = useState('');
    const [language, setLanguage] = useState('python');
    const [loading, setLoading] = useState(true);
    const [submitting, setSubmitting] = useState(false);
    const [result, setResult] = useState(null);

    useEffect(() => {
        const fetchExam = async () => {
            try {
                const data = await getExamDetails(id);
                setExam(data);
                const challenges = data?.challenges || data?.Challenges || [];
                if (challenges.length > 0) {
                    setCurrentChallenge(challenges[0]);
                }
            } catch (err) {
                console.error(err);
                const apiMessage = err?.response?.data?.error || err?.message || 'Failed to load exam';
                alert(apiMessage);
                navigate('/courses');
            } finally {
                setLoading(false);
            }
        };

        fetchExam();

        // Setup session heartbeat
        const sessionId = localStorage.getItem('session_id');
        let heartbeatInterval;
        
        if (sessionId) {
            heartbeatInterval = setInterval(async () => {
                try {
                    await client.post(`/submissions/sessions/${sessionId}/heartbeat`);
                } catch (err) {
                    console.warn('Heartbeat failed:', err);
                }
            }, 120000); // Pulse every 2 minutes
        }

        return () => {
            if (heartbeatInterval) clearInterval(heartbeatInterval);
        };
    }, [id, navigate]);

    const handleSubmit = async () => {
        if (!code.trim()) return;

        const sessionId = localStorage.getItem('session_id');
        if (!sessionId) {
            alert('No hay una sesion de examen activa. Vuelve a entrar al examen desde el curso.');
            return;
        }

        const challengeId = currentChallenge?.id || currentChallenge?.ID;
        if (!challengeId) {
            alert('No se pudo identificar el reto actual del examen.');
            return;
        }

        setSubmitting(true);
        setResult(null);
        try {
            const { data } = await client.post('/submissions', {
                challengeID: challengeId,
                code,
                language,
                sessionID: sessionId,
            });

            // Poll for result
            const pollInterval = setInterval(async () => {
                const submissionId = data?.id || data?.ID;
                if (!submissionId) {
                    clearInterval(pollInterval);
                    setSubmitting(false);
                    alert('La API no devolvio ID de submission.');
                    return;
                }

                const res = await client.get(`/submissions/${submissionId}`);
                const submission = res?.data?.Submission || res?.data?.submission;
                const results = res?.data?.Results || res?.data?.results || [];

                if (!Array.isArray(results) || results.length === 0) return;

                const hasPending = results.some((r) => {
                    const status = String(r?.Status || r?.status || '').toLowerCase();
                    return status === 'queued' || status === 'running';
                });

                if (hasPending) return;

                const acceptedCount = results.filter((r) => String(r?.Status || r?.status || '').toLowerCase() === 'accepted').length;
                const score = Math.round((acceptedCount / results.length) * 100);

                clearInterval(pollInterval);
                setResult({
                    status: score === 100 ? 'accepted' : 'wrong_answer',
                    score,
                    submission,
                    results,
                });
                setSubmitting(false);
            }, 1000);

        } catch (err) {
            console.error(err);
            setSubmitting(false);
            const apiMessage = err?.response?.data?.error || err?.message || 'Error submitting code';
            alert(apiMessage);
        }
    };

    if (loading) return <div>Loading Exam...</div>;
    if (!exam) return <div>Exam not found</div>;

    return (
        <div className="challenge-solver">
            <div className="solver-header">
                <h2>Exam: {exam.title || exam.Title}</h2>
                <div className="timer">Time Remaining: --:--</div>
            </div>

            <div className="solver-container">
                <div className="problem-description">
                    <h3>Challenges</h3>
                    <ul className="exam-challenge-list">
                        {(exam.challenges || exam.Challenges || []).map(ch => (
                            <li
                                key={ch.id || ch.ID}
                                className={(currentChallenge?.id || currentChallenge?.ID) === (ch.id || ch.ID) ? 'active' : ''}
                                onClick={() => setCurrentChallenge(ch)}
                            >
                                {ch.title || ch.Title} ({ch.points || ch.Points || 0} pts)
                            </li>
                        ))}
                    </ul>

                    {currentChallenge && (
                        <>
                            <h3>{currentChallenge.title || currentChallenge.Title}</h3>
                            <p>{currentChallenge.description || currentChallenge.Description}</p>
                        </>
                    )}
                </div>

                <div className="code-editor-section">
                    <div className="editor-controls">
                        <select value={language} onChange={(e) => setLanguage(e.target.value)}>
                            <option value="python">Python</option>
                        </select>
                        <button onClick={handleSubmit} disabled={submitting}>
                            {submitting ? 'Running...' : 'Submit Solution'}
                        </button>
                    </div>
                    <textarea
                        className="code-editor"
                        value={code}
                        onChange={(e) => setCode(e.target.value)}
                        placeholder="Write your code here..."
                    />
                    {result && (
                        <div className={`result-box ${result.status.toLowerCase()}`}>
                            <h4>Result: {result.status}</h4>
                            <p>Score: {result.score}</p>
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
};

export default ExamRunner;
