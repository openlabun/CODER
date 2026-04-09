import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import client from '../api/client';
import { 
    Key, 
    BookOpen, 
    ArrowRight, 
    ShieldCheck, 
    AlertCircle, 
    CheckCircle2,
    Loader2,
    ChevronLeft,
    HelpCircle
} from 'lucide-react';
import './JoinCourse.css';

const JoinCourse = () => {
    const navigate = useNavigate();
    const [enrollmentCode, setEnrollmentCode] = useState('');
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');
    const [success, setSuccess] = useState('');

    const handleJoinWithCode = async (e) => {
        e.preventDefault();
        setError('');
        setSuccess('');

        if (!enrollmentCode.trim()) {
            setError('Por favor, ingresa un código de inscripción');
            return;
        }

        setLoading(true);
        try {
            await client.post('/courses/enroll', { enrollment_code: enrollmentCode });
            setSuccess('¡Te has inscrito correctamente en el curso!');
            setTimeout(() => {
                navigate('/courses');
            }, 2000);
        } catch (err) {
            setError(err.response?.data?.message || 'Código de inscripción inválido o curso no encontrado');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="join-course-page-new">
            <div className="join-header-compact">
                <button onClick={() => navigate('/courses')} className="btn-back-mini">
                    <ChevronLeft size={18} />
                    <span>Volver a mis cursos</span>
                </button>
            </div>

            <div className="join-container-centered">
                <div className="join-card-premium">
                    <div className="join-card-glow"></div>
                    
                    <div className="card-top-icon">
                        <div className="icon-ring">
                            <Key size={32} className="main-icon" />
                        </div>
                    </div>

                    <div className="card-content-join">
                        <h2>Unirse a un Curso</h2>
                        <p>Ingresa el código único proporcionado por tu profesor para inscribirte automáticamente.</p>

                        <form onSubmit={handleJoinWithCode} className="modern-join-form">
                            <div className="form-field-group">
                                <label>Código de Inscripción</label>
                                <div className="input-with-icon-unique">
                                    <ShieldCheck size={20} className="input-prefix-icon" />
                                    <input
                                        type="text"
                                        value={enrollmentCode}
                                        onChange={(e) => setEnrollmentCode(e.target.value.toUpperCase())}
                                        placeholder="Ej: CS101-20251G1"
                                        disabled={loading || success}
                                        autoFocus
                                    />
                                </div>
                                <span className="input-hint-mini">Formato: CODIGO-PERIODO-GRUPO</span>
                            </div>

                            {error && (
                                <div className="feedback-alert error-unique">
                                    <AlertCircle size={18} />
                                    <span>{error}</span>
                                </div>
                            )}

                            {success && (
                                <div className="feedback-alert success-unique">
                                    <CheckCircle2 size={18} />
                                    <span>{success}</span>
                                </div>
                            )}

                            <button 
                                type="submit" 
                                className={`btn-join-primary ${loading ? 'loading' : ''}`}
                                disabled={loading || success}
                            >
                                {loading ? (
                                    <>
                                        <Loader2 size={18} className="spin-icon" />
                                        <span>Procesando...</span>
                                    </>
                                ) : (
                                    <>
                                        <span>Confirmar Inscripción</span>
                                        <ArrowRight size={18} />
                                    </>
                                )}
                            </button>
                        </form>
                    </div>

                    <div className="card-footer-help">
                        <HelpCircle size={16} />
                        <p>¿No tienes un código? Contacta a tu docente.</p>
                    </div>
                </div>

                <div className="join-instructions-mini">
                    <div className="instruction-step-mini">
                        <div className="step-num-mini">1</div>
                        <p>Obtén el código de tu profesor</p>
                    </div>
                    <div className="instruction-step-mini">
                        <div className="step-num-mini">2</div>
                        <p>Ingresa el código en el campo de arriba</p>
                    </div>
                    <div className="instruction-step-mini">
                        <div className="step-num-mini">3</div>
                        <p>Accede a tus nuevos desafíos</p>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default JoinCourse;
