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
    Award,
    Loader2
} from 'lucide-react';
import { useNavigate, useLocation } from 'react-router-dom';
import './Leaderboard.css';

const Leaderboard = ({ challengeId, courseId }) => {
    const [rankings, setRankings] = useState([]);
    const [loading, setLoading] = useState(true);
    const [searchTerm, setSearchTerm] = useState('');
    const [error, setError] = useState('');
    const navigate = useNavigate();

    useEffect(() => {
        const fetchLeaderboard = async () => {
            setLoading(true);
            setError('');
            try {
                let endpoint = '/leaderboard';
                if (challengeId) endpoint = `/leaderboard/challenge/${challengeId}`;
                else if (courseId) endpoint = `/leaderboard/course/${courseId}`;
                
                const { data } = await client.get(endpoint);
                
                // Transform backend data to frontend format
                // Backend might return { items: [...] } or just [...]
                const rawData = Array.isArray(data) ? data : (data.items || []);
                
                const formattedData = rawData.map((item, index) => ({
                    rank: item.rank || (index + 1),
                    username: item.username || item.Username || 'Anónimo',
                    score: item.score || item.TotalScore || item.Score || 0,
                    challenges: item.challengesSolved || item.ChallengesSolved || item.SolvedCount || 0,
                    avatar: (item.username || item.Username || 'A').charAt(0).toUpperCase(),
                    userId: item.userId || item.ID
                }));

                setRankings(formattedData);
            } catch (err) {
                console.error('Error loading leaderboard:', err);
                setError('No se pudo cargar el ranking en este momento.');
                // Fallback to empty if error
                setRankings([]);
            } finally {
                setLoading(false);
            }
        };
        fetchLeaderboard();
    }, [challengeId, courseId]);

    const filteredRankings = rankings.filter(user => 
        user.username.toLowerCase().includes(searchTerm.toLowerCase())
    );

    const podium = rankings.slice(0, 3);
    const restOfUsers = filteredRankings.filter(u => u.rank > 3);

    if (loading) return (
        <div className="leaderboard-page-new">
            <div className="page-loader">
                <Loader2 className="page-loader-spinner" size={48} />
                <p className="page-loader-text">Cargando clasificación...</p>
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
