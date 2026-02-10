import { Link } from 'react-router-dom';
import './ChallengeCard.css';

const ChallengeCard = ({ challenge }) => {
    return (
        <div className="challenge-card">
            <h3>{challenge.title}</h3>
            <div className="challenge-meta">
                <span className={`status ${challenge.status.toLowerCase()}`}>
                    {challenge.status}
                </span>
            </div>
            <Link to={`/challenge/${challenge.id}`} className="btn-secondary">
                Solve Challenge
            </Link>
        </div>
    );
};

export default ChallengeCard;
