import React from 'react';
import './LeaderboardTable.css';

/**
 * LeaderboardTable renders a list of leaderboard entries.
 * Expected entry shape:
 *   { rank?: number, username: string, score: number, timeMs?: number }
 */
const LeaderboardTable = ({ entries }) => {
    // Ensure entries are sorted by score descending, then time ascending
    const sorted = [...entries].sort((a, b) => {
        if (b.score !== a.score) return b.score - a.score;
        if (a.timeMs && b.timeMs) return a.timeMs - b.timeMs;
        return 0;
    });

    return (
        <div className="leaderboard-table-wrapper">
            <table className="leaderboard-table">
                <thead>
                    <tr>
                        <th>#</th>
                        <th>User</th>
                        <th>Score</th>
                        <th>Time (ms)</th>
                    </tr>
                </thead>
                <tbody>
                    {sorted.map((entry, idx) => (
                        <tr key={entry.username + idx}>
                            <td>{idx + 1}</td>
                            <td>{entry.username}</td>
                            <td>{entry.score}</td>
                            <td>{entry.timeMs ? entry.timeMs : '—'}</td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    );
};

export default LeaderboardTable;
