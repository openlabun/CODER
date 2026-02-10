import { useState, useEffect, useContext } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { getExamDetails } from '../api/exams';
import client from '../api/client';
import { AuthContext } from '../context/AuthContext';
import './ChallengeSolver.css'; // Reuse styles

const ExamRunner = () => {
    const { id } = useParams();
    const { token } = useContext(AuthContext);
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
                const data = await getExamDetails(id, token);
                setExam(data);
                if (data.challenges && data.challenges.length > 0) {
                    setCurrentChallenge(data.challenges[0]);
                }
            } catch (err) {
                console.error(err);
                alert('Failed to load exam');
                navigate('/courses');
            } finally {
                setLoading(false);
            }
        };
        fetchExam();
    }, [id, token, navigate]);

    const handleSubmit = async () => {
        if (!code.trim()) return;
        setSubmitting(true);
        setResult(null);
        try {
            const { data } = await client.post('/submissions', {
                challengeId: currentChallenge.id,
                code,
                language,
                examId: exam.id,
            });

            // Poll for result
            const pollInterval = setInterval(async () => {
                const res = await client.get(`/submissions/${data.id}`);
                if (res.data.status !== 'queued' && res.data.status !== 'running') {
                    clearInterval(pollInterval);
                    setResult(res.data);
                    setSubmitting(false);
                }
            }, 1000);

        } catch (err) {
            console.error(err);
            setSubmitting(false);
            alert('Error submitting code');
        }
    };

    if (loading) return <div>Loading Exam...</div>;
    if (!exam) return <div>Exam not found</div>;

    return (
        <div className="challenge-solver">
            <div className="solver-header">
                <h2>Exam: {exam.title}</h2>
                <div className="timer">Time Remaining: --:--</div>
            </div>

            <div className="solver-container">
                <div className="problem-description">
                    <h3>Challenges</h3>
                    <ul className="exam-challenge-list">
                        {exam.challenges.map(ch => (
                            <li
                                key={ch.id}
                                className={currentChallenge?.id === ch.id ? 'active' : ''}
                                onClick={() => setCurrentChallenge(ch)}
                            >
                                {ch.title} ({ch.points} pts)
                            </li>
                        ))}
                    </ul>

                    {currentChallenge && (
                        <>
                            <h3>{currentChallenge.title}</h3>
                            <p>{currentChallenge.description}</p>
                        </>
                    )}
                </div>

                <div className="code-editor-section">
                    <div className="editor-controls">
                        <select value={language} onChange={(e) => setLanguage(e.target.value)}>
                            <option value="python">Python</option>
                            <option value="javascript">Node.js</option>
                            <option value="cpp">C++</option>
                            <option value="java">Java</option>
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
