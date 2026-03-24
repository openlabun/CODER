import { useState, useEffect } from 'react';
import client from '../api/client';
import { 
    Trophy, 
    Medal, 
    User, 
    Star, 
    ChevronRight, 
    TrendingUp,
    Target,
    Zap,
    Search,
    Award
} from 'lucide-react';
import './Leaderboard.css';

const Leaderboard = () => {
    const [rankings, setRankings] = useState([]);
    const [loading, setLoading] = useState(true);
    const [searchTerm, setSearchTerm] = useState('');

    useEffect(() => {
        const fetchLeaderboard = async () => {
            try {
                // Mocking with enhanced data for a better visual representation
                const mockData = [
                    { rank: 1, username: 'neo_coder', score: 2850, challenges: 42, avatar: 'N' },
                    { rank: 2, username: 'trinity_dev', score: 2720, challenges: 38, avatar: 'T' },
                    { rank: 3, username: 'morpheus_root', score: 2600, challenges: 35, avatar: 'M' },
                    { rank: 4, username: 'cipher_hacker', score: 2100, challenges: 28, avatar: 'C' },
                    { rank: 5, username: 'agent_smith', score: 1950, challenges: 25, avatar: 'A' },
                    { rank: 6, username: 'the_oracle', score: 1800, challenges: 22, avatar: 'O' },
                    { rank: 7, username: 'keymaker', score: 1650, challenges: 20, avatar: 'K' },
                    { rank: 8, username: 'architect', score: 1500, challenges: 18, avatar: 'A' },
                    { rank: 9, username: 'niobe', score: 1420, challenges: 17, avatar: 'N' },
                    { rank: 10, username: 'seraph', score: 1380, challenges: 16, avatar: 'S' },
                ];
                setRankings(mockData);
            } catch (err) {
                console.error(err);
            } finally {
                setLoading(false);
            }
        };
        fetchLeaderboard();
    }, []);

    const filteredRankings = rankings.filter(user => 
        user.username.toLowerCase().includes(searchTerm.toLowerCase())
    );

    const podium = rankings.slice(0, 3);
    const restOfUsers = filteredRankings.filter(u => u.rank > 3);

    if (loading) return (
        <div className="leaderboard-page-new">
            <div className="skeleton title-skeleton-wide"></div>
            <div className="podium-skeleton"></div>
            <div className="list-skeleton-mini">
                {[...Array(5)].map((_, i) => <div key={i} className="skeleton-row-mini shimmer"></div>)}
            </div>
        </div>
    );

    return (
        <div className="leaderboard-page-new">
            <header className="leaderboard-header-new">
                <div className="header-info-new">
                    <h1>Ranking Mundial</h1>
                    <p>Los mejores programadores de la comunidad RobleCode</p>
                </div>
                
                <div className="search-bar-mini">
                    <Search size={18} />
                    <input 
                        type="text" 
                        placeholder="Buscar por usuario..." 
                        value={searchTerm}
                        onChange={(e) => setSearchTerm(e.target.value)}
                    />
                </div>
            </header>

            {/* Podium Section */}
            {!searchTerm && (
                <div className="podium-area-new">
                    {/* Rank 2 */}
                    <div className="podium-item second">
                        <div className="podium-rank">2</div>
                        <div className="podium-avatar">
                            <span>{podium[1]?.avatar}</span>
                            <Medal className="podium-medal silver" />
                        </div>
                        <div className="podium-name">{podium[1]?.username}</div>
                        <div className="podium-score">{podium[1]?.score} <span className="pts">PTS</span></div>
                        <div className="podium-bar"></div>
                    </div>

                    {/* Rank 1 */}
                    <div className="podium-item first">
                        <Trophy className="podium-crown" size={32} />
                        <div className="podium-rank">1</div>
                        <div className="podium-avatar">
                            <span>{podium[0]?.avatar}</span>
                            <Medal className="podium-medal gold" />
                        </div>
                        <div className="podium-name">{podium[0]?.username}</div>
                        <div className="podium-score">{podium[0]?.score} <span className="pts">PTS</span></div>
                        <div className="podium-bar"></div>
                    </div>

                    {/* Rank 3 */}
                    <div className="podium-item third">
                        <div className="podium-rank">3</div>
                        <div className="podium-avatar">
                            <span>{podium[2]?.avatar}</span>
                            <Medal className="podium-medal bronze" />
                        </div>
                        <div className="podium-name">{podium[2]?.username}</div>
                        <div className="podium-score">{podium[2]?.score} <span className="pts">PTS</span></div>
                        <div className="podium-bar"></div>
                    </div>
                </div>
            )}

            {/* List Section */}
            <div className="leaderboard-list-new">
                <div className="list-header-new">
                    <div className="col-rank">RANK</div>
                    <div className="col-user">USUARIO</div>
                    <div className="col-solved">RETOS</div>
                    <div className="col-score">PUNTOS</div>
                    <div className="col-trend">ESTADO</div>
                </div>

                {restOfUsers.length > 0 ? (
                    restOfUsers.map((user) => (
                        <div key={user.username} className="ranking-row-new">
                            <div className="col-rank">
                                <span className="rank-num">#{user.rank}</span>
                            </div>
                            <div className="col-user">
                                <div className="user-profile-mini">
                                    <div className="user-avatar-mini">{user.avatar}</div>
                                    <span className="username-text">{user.username}</span>
                                </div>
                            </div>
                            <div className="col-solved">
                                <div className="solved-badge-mini">
                                    <Target size={14} />
                                    <span>{user.challenges}</span>
                                </div>
                            </div>
                            <div className="col-score">
                                <span className="score-text-new">{user.score.toLocaleString()} <small>PTS</small></span>
                            </div>
                            <div className="col-trend">
                                <div className="trend-badge-mini">
                                    <TrendingUp size={14} />
                                    <span>UP</span>
                                </div>
                            </div>
                        </div>
                    ))
                ) : (
                    <div className="empty-search-mini">
                        <Zap size={40} className="icon-muted" />
                        <h3>No se encontraron resultados</h3>
                        <p>Intenta con otra búsqueda o usuario.</p>
                    </div>
                )}
            </div>
        </div>
    );
};

export default Leaderboard;
