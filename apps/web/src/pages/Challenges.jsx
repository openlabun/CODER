import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import client from '../api/client';
import './Challenges.css';

const Challenges = () => {
    const [challenges, setChallenges] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');

    useEffect(() => {
        const fetchChallenges = async () => {
            try {
                const { data } = await client.get('/challenges');
                setChallenges(data);
            } catch (err) {
                setError('Failed to load challenges');
            } finally {
                setLoading(false);
            }
        };
        fetchChallenges();
    }, []);

    if (loading) return <div className="loading">Loading challenges...</div>;
    if (error) return <div className="error">{error}</div>;

    return (
        <div className="challenges-page">
            <div className="page-header">
                <h1>Challenges</h1>
            </div>

            {challenges.length === 0 ? (
                <div className="empty-state">
                    <h3>📝 No Challenges Available</h3>
                    <p>There are no challenges available at the moment. Check back later or contact your professor.</p>
                </div>
            ) : (
                <div className="challenges-grid">
                    {challenges.map((challenge) => (
                        <div key={challenge.id} className="challenge-card">
                            <div className="challenge-header">
                                <h3 className="challenge-title">{challenge.title}</h3>
                                <span className={`difficulty-badge ${(challenge.difficulty || 'medium').toLowerCase()}`}>
                                    {challenge.difficulty || 'Medium'}
                                </span>
                            </div>
                            <p className="challenge-desc">{challenge.description}</p>
                            <Link to={`/challenge/${challenge.id}`} className="btn-solve">
                                Solve Challenge
                            </Link>
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
};

export default Challenges;
