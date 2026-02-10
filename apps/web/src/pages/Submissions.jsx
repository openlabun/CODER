import { useState, useEffect, useContext } from 'react';
import client from '../api/client';
import { AuthContext } from '../context/AuthContext';
import './Submissions.css';

const Submissions = () => {
    const [submissions, setSubmissions] = useState([]);
    const [loading, setLoading] = useState(true);
    const { user } = useContext(AuthContext);

    useEffect(() => {
        const fetchSubmissions = async () => {
            try {
                const { data } = await client.get('/submissions'); // Assuming endpoint returns user's submissions
                setSubmissions(data.items || []);
            } catch (err) {
                console.error(err);
            } finally {
                setLoading(false);
            }
        };
        fetchSubmissions();
    }, []);

    if (loading) return <div className="loading">Loading submissions...</div>;

    return (
        <div className="submissions-page">
            <div className="page-header">
                <h1>My Submissions</h1>
            </div>

            <div className="submissions-table-container">
                <table className="submissions-table">
                    <thead>
                        <tr>
                            <th>ID</th>
                            <th>Challenge</th>
                            <th>Language</th>
                            <th>Status</th>
                            <th>Score</th>
                            <th>Time</th>
                            <th>Date</th>
                        </tr>
                    </thead>
                    <tbody>
                        {submissions.map((sub) => (
                            <tr key={sub.id}>
                                <td className="mono">{sub.id.slice(0, 8)}</td>
                                <td>{sub.challengeId}</td>
                                <td className="capitalize">{sub.language}</td>
                                <td>
                                    <span className={`status-badge ${sub.status.toLowerCase()}`}>
                                        {sub.status.replace('_', ' ')}
                                    </span>
                                </td>
                                <td>{sub.score}</td>
                                <td>{sub.timeMsTotal}ms</td>
                                <td>{new Date(sub.createdAt).toLocaleDateString()}</td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
        </div>
    );
};

export default Submissions;
