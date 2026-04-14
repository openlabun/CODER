import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import client from '../api/client';
import { 
    Search, 
    User, 
    Calendar, 
    Hash, 
    ArrowRight, 
    X, 
    Lock,
    AlertCircle,
    Loader2,
    Sparkles
} from 'lucide-react';
import './BrowseCourses.css';

const BrowseCourses = () => {
    const navigate = useNavigate();
    const [courses, setCourses] = useState([]);
    const [filteredCourses, setFilteredCourses] = useState([]);
    const [searchTerm, setSearchTerm] = useState('');
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');

    // Modal state
    const [showModal, setShowModal] = useState(false);
    const [selectedCourse, setSelectedCourse] = useState(null);
    const [enrollmentCode, setEnrollmentCode] = useState('');
    const [enrolling, setEnrolling] = useState(false);
    const [enrollError, setEnrollError] = useState('');

    useEffect(() => {
        const fetchCourses = async () => {
            try {
                // Backend doesn't have a browse endpoint, show enrolled courses
                const { data } = await client.get('/courses?scope=enrolled');
                const list = Array.isArray(data) ? data : (data?.items || []);
                setCourses(list);
                setFilteredCourses(list);
            } catch (err) {
                console.error(err);
                // If enrolled returns empty or errors, just show empty state
                setCourses([]);
                setFilteredCourses([]);
            } finally {
                setLoading(false);
            }
        };
        fetchCourses();
    }, []);

    useEffect(() => {
        const results = courses.filter(course =>
            course.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
            course.code.toLowerCase().includes(searchTerm.toLowerCase())
        );
        setFilteredCourses(results);
    }, [searchTerm, courses]);

    const handleJoinClick = (course) => {
        setSelectedCourse(course);
        setEnrollmentCode('');
        setEnrollError('');
        setShowModal(true);
    };

    const handleJoinCourse = async (e) => {
        if (e) e.preventDefault();
        setEnrollError('');

        if (!enrollmentCode.trim()) {
            setEnrollError('El código de inscripción es obligatorio');
            return;
        }

        setEnrolling(true);
        try {
            await client.post('/courses/enroll', { enrollment_code: enrollmentCode });
            setShowModal(false);
            navigate('/courses');
        } catch (err) {
            setEnrollError(err.response?.data?.message || 'Código de inscripción inválido');
        } finally {
            setEnrolling(false);
        }
    };

    if (loading) return (
        <div className="browse-courses-page">
            <header className="browse-hero-banner">
                <div className="hero-content-inner">
                    <div style={{width: '200px', height: '40px', background: '#334155', borderRadius: '100px', margin: '0 auto 20px'}}></div>
                    <div style={{width: '100%', height: '80px', background: '#334155', borderRadius: '20px', margin: '0 auto'}}></div>
                </div>
            </header>
            <div className="browse-content-area">
                <div className="browse-grid-premium">
                    {[...Array(6)].map((_, i) => (
                        <div key={i} className="skeleton-rc-card">
                            <div className="sk-line" style={{width: '60px', height: '24px'}}></div>
                            <div className="sk-line" style={{width: '100%', height: '32px'}}></div>
                            <div className="sk-line" style={{width: '80%', height: '20px'}}></div>
                            <div className="sk-line" style={{width: '100%', height: '48px', marginTop: 'auto'}}></div>
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );

    if (error) return (
        <div className="browse-courses-page flex-center-view">
            <div className="rc-modal-container">
                <AlertCircle size={64} color="#f87171" strokeWidth={1.5} />
                <h2>Oops, algo salió mal</h2>
                <p>{error}</p>
                <button onClick={() => window.location.reload()} className="btn-rc-submit" style={{width: '100%'}}>
                    Reintentar Carga
                </button>
            </div>
        </div>
    );

    return (
        <div className="browse-courses-page">
            <header className="browse-hero-banner">
                <div className="hero-content-inner">
                    <h1>Descubre tu camino</h1>
                    <p>Encuentra tus asignaturas y únete a la comunidad académica</p>
                    
                    <div className="search-wrapper-floating">
                        <Search className="search-icon-float" size={24} />
                        <input
                            type="text"
                            placeholder="Buscar por nombre, código o docente..."
                            value={searchTerm}
                            onChange={(e) => setSearchTerm(e.target.value)}
                            className="search-field"
                        />
                    </div>
                </div>
            </header>

            <div className="browse-content-area">
                <div className="browse-grid-premium">
                    {filteredCourses.length > 0 ? (
                        filteredCourses.map((course) => (
                            <div key={course.id} className="glass-course-card">
                                <div className="card-top-inner">
                                    <div className="card-label-row">
                                        <div className="tag-premium">
                                            Asignatura
                                        </div>
                                        <div className="course-dot-accent"></div>
                                    </div>
                                    <h3>{course.name}</h3>
                                    
                                    <div className="meta-grid-compact">
                                        <div className="meta-row-premium">
                                            <Hash size={18} />
                                            <span>{course.code}</span>
                                        </div>
                                        <div className="meta-row-premium">
                                            <Calendar size={18} />
                                            <span>Periodo {course.period ? `${course.period.year}-${course.period.semester}` : 'S/P'}</span>
                                        </div>
                                        <div className="meta-row-premium" style={{borderTop: '1px solid rgba(0,0,0,0.05)', paddingTop: '1rem', marginTop: '0.5rem'}}>
                                            <User size={18} />
                                            <span>{course.professor_name || 'Docente asignado'}</span>
                                        </div>
                                    </div>
                                </div>
                                
                                <div className="card-actions-area">
                                    <button
                                        onClick={() => handleJoinClick(course)}
                                        className="btn-premium-join"
                                    >
                                        Unirse ahora <ArrowRight size={18} />
                                    </button>
                                </div>
                            </div>
                        ))
                    ) : (
                        <div style={{gridColumn: '1/-1', textAlign: 'center', padding: '6rem 2rem'}}>
                            <div style={{background: 'white', display: 'inline-flex', padding: '2rem', borderRadius: '40px', marginBottom: '2rem'}}>
                                <Search size={48} color="#94a3b8" />
                            </div>
                            <h3>No se encontraron resultados</h3>
                            <p>Intenta con otros términos de búsqueda para encontrar tu curso.</p>
                        </div>
                    )}
                </div>
            </div>

            {/* Modal Ultra Premium */}
            {showModal && (
                <div className="rc-modal-overlay" onClick={() => setShowModal(false)}>
                    <div className="rc-modal-container" onClick={(e) => e.stopPropagation()}>
                        <button onClick={() => setShowModal(false)} className="rc-close-btn">
                            <X size={20} />
                        </button>
                        
                        <div className="rc-modal-icon">
                            <Lock size={36} strokeWidth={1.5} />
                        </div>
                        
                        <h2>Inscripción</h2>
                        <p>
                            Introduce el código único de <strong>{selectedCourse?.name || 'este curso'}</strong> para completar tu inscripción.
                        </p>

                        {enrollError && (
                            <div className="rc-error-msg">
                                <AlertCircle size={18} />
                                {enrollError}
                            </div>
                        )}

                        <form onSubmit={handleJoinCourse}>
                            <input
                                type="text"
                                value={enrollmentCode}
                                onChange={(e) => setEnrollmentCode(e.target.value.toUpperCase())}
                                placeholder="CÓDIGO-01"
                                disabled={enrolling}
                                className="rc-input-code"
                                autoFocus
                            />

                            <div className="rc-modal-footer">
                                <button
                                    type="button"
                                    onClick={() => setShowModal(false)}
                                    className="btn-rc-cancel"
                                    disabled={enrolling}
                                >
                                    Cerrar
                                </button>
                                <button
                                    type="submit"
                                    className="btn-rc-submit"
                                    disabled={enrolling || !enrollmentCode.trim()}
                                >
                                    {enrolling ? (
                                        <><Loader2 size={18} className="animate-spin" /> Uniendo...</>
                                    ) : (
                                        <>Confirmar <ArrowRight size={18} /></>
                                    )}
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
            )}
        </div>
    );
};

export default BrowseCourses;
