import { Loader2 } from 'lucide-react';
import './PageLoader.css';

const PageLoader = ({
    message = 'Cargando...',
    minHeight = '60vh',
    className = '',
    size = 48,
    fullScreen = false,
    compact = false,
}) => {
    const rootClass = `rc-page-loader ${fullScreen ? 'rc-page-loader--full-screen' : ''} ${compact ? 'rc-page-loader--compact' : ''} ${className}`.trim();

    return (
        <div className={rootClass} style={{ minHeight }} role="status" aria-live="polite" aria-busy="true">
            <div className="rc-page-loader__ring" aria-hidden="true">
                <Loader2 className="rc-page-loader__spinner" size={size} />
            </div>
            <p className="rc-page-loader__text">{message}</p>
        </div>
    );
};

export default PageLoader;
