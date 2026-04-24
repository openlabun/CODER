import { useState, useEffect, useContext } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import client from '../api/client';
import { createExam } from '../api/exams';
import { AuthContext } from '../context/AuthContext';
import Swal from 'sweetalert2';
import { Calendar, Clock, FileText, Layout, Save, X, Info, Sparkles, Globe, Users, BookOpen, Lock } from 'lucide-react';
import AIAssistantModal from '../components/AIAssistantModal';
import PageLoader from '../components/PageLoader';
import './CreateCourse.css'; // Reutilizamos estilos por consistencia

const CreateExam = () => {
    const navigate = useNavigate();
    const location = useLocation();
    const { user } = useContext(AuthContext);
    const queryParams = new URLSearchParams(location.search);
    const courseId = queryParams.get('courseId');

    const getTodayISO = (time) => {
        const d = new Date();
        const year = d.getFullYear();
        const month = String(d.getMonth() + 1).padStart(2, '0');
        const day = String(d.getDate()).padStart(2, '0');
        return `${year}-${month}-${day}T${time}`;
    };

    const [formData, setFormData] = useState({
        title: '',
        description: '',
        visibility: 'course',
        startTime: getTodayISO('00:00'),
        endTime: getTodayISO('23:59'),
        timeLimit: 60, // En minutos para el front, convertiremos a segundos
        tryLimit: 1,
        allowLateSubmissions: false
    });

    const [loading, setLoading] = useState(false);
    const [showAIModal, setShowAIModal] = useState(false);
    const [checkingAccess, setCheckingAccess] = useState(true);

    const isTeacher = user?.role === 'professor' || user?.role === 'teacher' || user?.role === 'admin';

    useEffect(() => {
        if (!user || !isTeacher) {
            navigate('/dashboard');
            return;
        }

        setCheckingAccess(false);
    }, [user, isTeacher, navigate]);

    if (checkingAccess) {
        return (
            <div className="create-course-page">
                <PageLoader message="Preparando formulario del examen..." />
            </div>
        );
    }

    const handleChange = (e) => {
        const { name, value, type, checked } = e.target;
        setFormData(prev => ({
            ...prev,
            [name]: type === 'checkbox' ? checked : value
        }));
    };

    const handleApplyAIExam = (examIdea) => {
        setFormData(prev => ({
            ...prev,
            title: examIdea.title,
            description: examIdea.description,
            timeLimit: examIdea.time_limit,
            tryLimit: examIdea.try_limit
        }));
        
        Swal.fire({
            icon: 'success',
            title: 'Propuesta de IA Aplicada',
            text: 'Revisa los campos y completa la programación.',
            timer: 2000,
            toast: true,
            position: 'top-end',
            showConfirmButton: false
        });
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setLoading(true);

        try {
            const payload = {
                course_id: courseId || null,
                title: formData.title,
                description: formData.description,
                visibility: formData.visibility,
                start_time: new Date(formData.startTime).toISOString(),
                end_time: formData.endTime ? new Date(formData.endTime).toISOString() : null,
                time_limit: parseInt(formData.timeLimit) * 60, // A segundos
                try_limit: parseInt(formData.tryLimit),
                allow_late_submissions: formData.allowLateSubmissions,
                professor_id: user.id || user.ID || ''
            };

            await createExam(payload);

            Swal.fire({
                icon: 'success',
                title: 'Examen Creado',
                text: 'El examen se ha configurado correctamente',
                timer: 2000,
                showConfirmButton: false,
                toast: true,
                position: 'top-end'
            });

            navigate(courseId ? `/courses/${courseId}` : '/dashboard');
        } catch (err) {
            console.error(err);
            Swal.fire({
                icon: 'error',
                title: 'Error',
                text: err.response?.data?.error || 'No se pudo crear el examen'
            });
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="create-course-page">
            <div className="page-header">
                <div className="header-content">
                    <h1>Crear Nuevo Examen</h1>
                    <p className="subtitle">Configura una evaluación para tus estudiantes</p>
                </div>
                <button 
                    className="ai-gen-header-btn" 
                    onClick={() => setShowAIModal(true)}
                >
                    <Sparkles size={16} /> Diseñar con IA
                </button>
            </div>

            <div className="form-container">
                <form onSubmit={handleSubmit} className="course-form">
                    <div className="form-section">
                        <div className="section-header">
                            <FileText size={20} />
                            <h2>Información General</h2>
                        </div>
                        
                        <div className="form-group">
                            <label>Título del Examen *</label>
                            <input
                                type="text"
                                name="title"
                                value={formData.title}
                                onChange={handleChange}
                                placeholder="ej. Parcial 1: Algoritmos"
                                required
                            />
                        </div>

                        <div className="form-group">
                            <label>Descripción / Instrucciones</label>
                            <textarea
                                name="description"
                                value={formData.description}
                                onChange={handleChange}
                                placeholder="Instrucciones para los estudiantes..."
                                rows="4"
                            />
                        </div>
                    </div>

                    <div className="form-section">
                        <div className="section-header">
                            <Calendar size={20} />
                            <h2>Programación</h2>
                        </div>

                        <div className="form-row">
                            <div className="form-group">
                                <label>Fecha y Hora de Inicio *</label>
                                <input
                                    type="datetime-local"
                                    name="startTime"
                                    value={formData.startTime}
                                    onChange={handleChange}
                                    required
                                />
                            </div>

                            <div className="form-group">
                                <label>Fecha y Hora de Cierre (Opcional)</label>
                                <input
                                    type="datetime-local"
                                    name="endTime"
                                    value={formData.endTime}
                                    onChange={handleChange}
                                />
                            </div>
                        </div>
                    </div>

                    <div className="form-section">
                        <div className="section-header">
                            <Clock size={20} />
                            <h2>Restricciones y Límites</h2>
                        </div>

                        <div className="form-row">
                            <div className="form-group">
                                <label>Duración (minutos) *</label>
                                <input
                                    type="number"
                                    name="timeLimit"
                                    value={formData.timeLimit}
                                    onChange={handleChange}
                                    min="1"
                                    required
                                />
                                <small>Tiempo máximo para completar el examen</small>
                            </div>

                            <div className="form-group">
                                <label>Límite de Intentos *</label>
                                <input
                                    type="number"
                                    name="tryLimit"
                                    value={formData.tryLimit}
                                    onChange={handleChange}
                                    min="1"
                                    required
                                />
                            </div>
                        </div>

                        <div className="checkbox-group">
                            <label className="checkbox-label">
                                <input
                                    type="checkbox"
                                    name="allowLateSubmissions"
                                    checked={formData.allowLateSubmissions}
                                    onChange={handleChange}
                                />
                                <span>Permitir entregas tardías</span>
                            </label>
                        </div>
                    </div>

                    <div className="form-section">
                        <div className="section-header">
                            <Layout size={20} />
                            <h2>Visibilidad</h2>
                        </div>

                        <div className="radio-group grid-2 visibility-radio-group">
                            <label className={`radio-card visibility-radio-card ${formData.visibility === 'course' ? 'active' : ''}`}>
                                <input
                                    type="radio"
                                    name="visibility"
                                    value="course"
                                    checked={formData.visibility === 'course'}
                                    onChange={handleChange}
                                />
                                <div className="radio-content visibility-radio-content">
                                    <div className="visibility-title-row">
                                        <BookOpen size={16} className="visibility-icon" />
                                        <span className="radio-title">Solo mi Curso</span>
                                    </div>
                                    <small>Visible solo para estudiantes inscritos</small>
                                </div>
                            </label>

                            <label className={`radio-card visibility-radio-card ${formData.visibility === 'public' ? 'active' : ''}`}>
                                <input
                                    type="radio"
                                    name="visibility"
                                    value="public"
                                    checked={formData.visibility === 'public'}
                                    onChange={handleChange}
                                />
                                <div className="radio-content visibility-radio-content">
                                    <div className="visibility-title-row">
                                        <Globe size={16} className="visibility-icon" />
                                        <span className="radio-title">Público Global</span>
                                    </div>
                                    <small>Visible para toda la comunidad RobleCode</small>
                                </div>
                            </label>

                            <label className={`radio-card visibility-radio-card ${formData.visibility === 'teachers' ? 'active' : ''}`}>
                                <input
                                    type="radio"
                                    name="visibility"
                                    value="teachers"
                                    checked={formData.visibility === 'teachers'}
                                    onChange={handleChange}
                                />
                                <div className="radio-content visibility-radio-content">
                                    <div className="visibility-title-row">
                                        <Users size={16} className="visibility-icon" />
                                        <span className="radio-title">Solo Profesores</span>
                                    </div>
                                    <small>Colabora con otros docentes</small>
                                </div>
                            </label>

                            <label className={`radio-card visibility-radio-card ${formData.visibility === 'private' ? 'active' : ''}`}>
                                <input
                                    type="radio"
                                    name="visibility"
                                    value="private"
                                    checked={formData.visibility === 'private'}
                                    onChange={handleChange}
                                />
                                <div className="radio-content visibility-radio-content">
                                    <div className="visibility-title-row">
                                        <Lock size={16} className="visibility-icon" />
                                        <span className="radio-title">Privado / Borrador</span>
                                    </div>
                                    <small>Solo tú puedes verlo y editarlo</small>
                                </div>
                            </label>
                        </div>
                    </div>

                    <div className="form-actions">
                        <button 
                            type="button" 
                            onClick={() => navigate(-1)} 
                            className="btn-secondary"
                        >
                            <X size={18} /> Cancelar
                        </button>
                        <button 
                            type="submit" 
                            className="btn-primary" 
                            disabled={loading}
                        >
                            {loading ? (
                                'Creando...'
                            ) : (
                                <><Save size={18} /> Crear Examen</>
                            )}
                        </button>
                    </div>
                </form>
            </div>

            <div className="info-box-alt">
                <div className="info-icon">
                    <Info size={20} />
                </div>
                <div className="info-text">
                    <h3>💡 ¿Cómo añadir retos?</h3>
                    <p>
                        Una vez creado el examen, puedes añadir retos desde la sección de "Retos" 
                        editando cada reto y asociándolo a este examen.
                    </p>
                </div>
            </div>

            {showAIModal && (
                <AIAssistantModal 
                    onClose={() => setShowAIModal(false)}
                    onApplyExam={handleApplyAIExam}
                    initialTab="exam"
                />
            )}
        </div>
    );
};

export default CreateExam;
