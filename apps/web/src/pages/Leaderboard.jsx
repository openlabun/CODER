import { useState, useEffect } from 'react';
import client from '../api/client';
import './Leaderboard.css';

const Leaderboard = () => {
    const [rankings, setRankings] = useState([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchLeaderboard = async () => {
            try {
                // Assuming we have a general leaderboard endpoint or course specific
                // For now mocking or using a placeholder endpoint
                // const { data } = await client.get('/leaderboard/global'); 
                // setRankings(data);
                setRankings([
                    { rank: 1, username: 'neo', score: 1500, challenges: 15 },
                    { rank: 2, username: 'trinity', score: 1450, challenges: 14 },
                    { rank: 3, username: 'morpheus', score: 1300, challenges: 13 },
                ]);
            } catch (err) {
                console.error(err);
            } finally {
                setLoading(false);
            }
        };
        fetchLeaderboard();
    }, []);

    return (
        <div className="leaderboard-page">
            <div className="page-header">
                <h1>Global Leaderboard</h1>
            </div>

            <div className="leaderboard-container">
                {rankings.map((user, index) => (
                    <div key={user.username} className={`rank-card rank-${index + 1}`}>
                        <div className="rank-position">{user.rank}</div>
                        <div className="rank-info">
                            <div className="rank-user">{user.username}</div>
                            <div className="rank-stats">
                                <span>{user.challenges} Challenges Solved</span>
                            </div>
                        </div>
                        <div className="rank-score">
                            {user.score} <span className="pts">PTS</span>
                        </div>
                    </div>
                ))}
            </div>
        </div>
    );
};

export default Leaderboard;
