import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import client from '../api/client';
import { useAuth } from '../context/AuthContext';
import { Trash2, Copy, Link as LinkIcon, UploadCloud, FileText, Key, ArrowLeft, Users } from 'lucide-react';
import PageLoader from '../components/PageLoader';
import Swal from 'sweetalert2';
import './Courses.css';
import './CourseActions.css';

const CourseStudents = () => {
    const { id } = useParams();
    const navigate = useNavigate();
    const { user } = useAuth();
    const [students, setStudents] = useState([]);
    const [course, setCourse] = useState(null);
    const [loading, setLoading] = useState(true);
    const [processingId, setProcessingId] = useState(null);

    const isTeacher = user?.role === 'professor' || user?.role === 'teacher' || user?.role === 'admin';
    const [searchEmail, setSearchEmail] = useState('');
    const [adding, setAdding] = useState(false);
    const [csvFile, setCsvFile] = useState(null);
    const [isUploadingCsv, setIsUploadingCsv] = useState(false);

    const fetchData = async () => {
        try {
            const [courseRes, studentsRes] = await Promise.all([
                client.get(`/courses/${id}`),
                client.get(`/courses/${id}/students`)
            ]);
            setCourse(courseRes.data);
            setStudents(Array.isArray(studentsRes.data) ? studentsRes.data : (studentsRes.data.students || []));
        } catch (error) {
            console.error('Error fetching students:', error);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchData();
    }, [id]);

    const handleAddStudent = async (e) => {
        if (e) e.preventDefault();
        let email = searchEmail.trim();
        if (email && !email.includes('@')) {
            email += '@uninorte.edu.co';
        }

        setAdding(true);
        try {
            await client.post(`/courses/${id}/students`, { studentEmail: email });
            
            Swal.fire({
                icon: 'success',
                title: 'Estudiante Agregado',
                text: `${email} ha sido unido al curso.`,
                timer: 1000,
                toast: true,
                position: 'top-end',
                showConfirmButton: false
            });

            setSearchEmail('');
            await fetchData();
        } catch (err) {
            console.error('Error adding student:', err);
            const errorMsg = err.response?.data?.message || err.response?.data?.error || 'No se pudo agregar al estudiante';
            
            Swal.fire({
                icon: errorMsg.toLowerCase().includes('inscrito') ? 'info' : 'error',
                title: errorMsg.toLowerCase().includes('inscrito') ? 'Ya inscrito' : 'Error',
                text: errorMsg,
                timer: 2000,
                toast: true,
                position: 'top-end',
                showConfirmButton: false
            });
        } finally {
            setAdding(false);
        }
    };

    const handleRemoveStudent = async (studentId) => {
        const { isConfirmed } = await Swal.fire({
            title: '¿Eliminar Estudiante?',
            text: 'El estudiante perderá el acceso a este curso inmediatamente.',
            icon: 'warning',
            showCancelButton: true,
            confirmButtonColor: '#d33',
            confirmButtonText: 'Sí, eliminar',
            cancelButtonText: 'Cancelar'
        });

        if (!isConfirmed) return;

        setProcessingId(studentId);
        try {
            await client.delete(`/courses/${id}/students/${studentId}`);
            Swal.fire({
                icon: 'success',
                title: 'Eliminado',
                timer: 1000,
                toast: true,
                position: 'top-end',
                showConfirmButton: false
            });
            await fetchData();
        } catch (err) {
            console.error('Error removing student:', err);
            Swal.fire({ icon: 'error', title: 'Error', text: 'No se pudo eliminar al estudiante.' });
        } finally {
            setProcessingId(null);
        }
    };

    const copyToClipboard = (text, type) => {
        navigator.clipboard.writeText(text);
        Swal.fire({
            icon: 'success',
            title: '¡Copiado!',
            text: `El ${type} ha sido copiado al portapapeles.`,
            timer: 1500,
            showConfirmButton: false,
            toast: true,
            position: 'top-end'
        });
    };

    const handleCsvUpload = async () => {
        if (!csvFile) return;
        setIsUploadingCsv(true);

        const reader = new FileReader();
        reader.onload = async (e) => {
            const text = e.target.result;
            // Basic parsing: split by newline, comma, semicolon and extract emails
            const words = text.split(/[\n\r,;]+/);
            const emails = words.map(w => w.trim()).filter(w => w.includes('@'));

            if (emails.length === 0) {
                Swal.fire({ icon: 'warning', title: 'Archivo sin correos', text: 'No se encontraron direcciones de correo en el archivo CSV.' });
                setIsUploadingCsv(false);
                return;
            }

            let successCount = 0;
            let failCount = 0;

            for (const email of emails) {
                try {
                    await client.post(`/courses/${id}/students`, { studentEmail: email });
                    successCount++;
                } catch (err) {
                    failCount++;
                }
            }

            Swal.fire({
                icon: 'info',
                title: 'Proceso CSV finalizado',
                html: `Se intentó inscribir a <b>${emails.length}</b> estudiantes.<br/><br/>
                       <span style="color: green">Éxito: ${successCount}</span><br/>
                       <span style="color: red">Fallidos / Ya inscritos: ${failCount}</span>`
            });
            
            setCsvFile(null);
            setIsUploadingCsv(false);
            fetchData();
        };
        reader.readAsText(csvFile);
    };

    if (loading) {
        return (
            <div className="course-students-page">
                <PageLoader message="Cargando estudiantes del curso..." minHeight="260px" />
            </div>
        );
    }

    return (
        <div className="course-students-page">
            <div className="page-header">
                <div>
                    <h1>Gestión de Estudiantes</h1>
                    <p className="subtitle">{course?.name} ({course?.code})</p>
                </div>
                <button onClick={() => navigate(`/courses/${id}`)} className="btn-back">
                    <ArrowLeft size={16} /> Volver al Curso
                </button>
            </div>

            {isTeacher && (
                <div className="course-students-admin-grid">
                    
                    {/* Add Single Student */}
                    <div className="admin-actions-card">
                        <h3>Añadir Estudiante</h3>
                        <p>Agrega un alumno usando su correo.</p>
                        <form className="add-student-form" onSubmit={handleAddStudent}>
                            <div className="input-with-helper">
                                <div className="input-group">
                                    <input
                                        type="text"
                                        placeholder="usuario o correo@uninorte.edu.co"
                                        value={searchEmail}
                                        onChange={(e) => setSearchEmail(e.target.value)}
                                        required
                                    />
                                    <button type="submit" className="btn-add-student" disabled={adding}>
                                        {adding ? 'Añadiendo...' : 'Añadir al Curso'}
                                    </button>
                                </div>
                                {!searchEmail.includes('@') && searchEmail.length > 2 && (
                                    <button 
                                        type="button" 
                                        className="helper-link"
                                        onClick={() => setSearchEmail(searchEmail.trim() + '@uninorte.edu.co')}
                                    >
                                        Completar con @uninorte.edu.co
                                    </button>
                                )}
                            </div>
                        </form>
                    </div>

                    {/* Auto CSV */}
                    <div className="admin-actions-card admin-actions-card--stacked">
                        <h3>Importar desde CSV</h3>
                        <p className="admin-card-copy">Sube un archivo .csv con los correos de los estudiantes. El archivo debe contener las direcciones de correo electrónico (separadas por comas o saltos de línea). Los estudiantes serán inscritos automáticamente al curso.</p>
                        <div className="csv-actions-row">
                            <input 
                                type="file" 
                                accept=".csv" 
                                id="csv-upload"
                                style={{ display: 'none' }}
                                onChange={(e) => setCsvFile(e.target.files[0])}
                            />
                            <label htmlFor="csv-upload" className={`csv-file-picker ${csvFile ? 'has-file' : ''}`}>
                                <FileText size={18} />
                                {csvFile ? csvFile.name : 'Seleccionar .csv'}
                            </label>
                            {csvFile && (
                                <button className="btn-action-filled" onClick={handleCsvUpload} disabled={isUploadingCsv}>
                                    <UploadCloud size={16} /> {isUploadingCsv ? 'Procesando...' : 'Cargar'}
                                </button>
                            )}
                        </div>
                    </div>

                    {/* Enrollment Methods Links */}
                    <div className="admin-actions-card">
                        <h3>Comparte el Curso</h3>
                        <p>Invita a estudiantes rápida y masivamente.</p>
                        
                        <div className="share-course-stack">
                            <div className="share-course-row">
                                <span className="share-course-icon"><LinkIcon size={16} /></span>
                                <div className="share-course-value" title={`${window.location.origin}/courses/join?code=${course?.enrollment_code}`}>
                                    {window.location.origin}/courses/join?code={course?.enrollment_code}
                                </div>
                                <button className="btn-add-student btn-share-copy" onClick={() => copyToClipboard(`${window.location.origin}/courses/join?code=${course?.enrollment_code}`, 'enlace')}>
                                    <Copy size={14} /> Link
                                </button>
                            </div>
                            <div className="share-course-row">
                                <span className="share-course-icon"><Key size={16} /></span>
                                <div className="share-course-value share-course-code">
                                    {course?.enrollment_code}
                                </div>
                                <button className="btn-add-student btn-share-copy btn-share-copy--code" onClick={() => copyToClipboard(course?.enrollment_code, 'código')}>
                                    <Copy size={14} /> Código
                                </button>
                            </div>
                        </div>

                    </div>

                </div>
            )}

            {students.length === 0 ? (
                <div className="empty-state">
                    <div className="empty-state-icon">👥</div>
                    <h3 className="empty-state-title">No hay estudiantes inscritos</h3>
                    <p className="empty-state-description">
                        Nadie se ha unido a este curso todavía.
                    </p>
                    <p className="empty-state-hint">
                        Comparte el código de inscripción: <strong>{course?.enrollment_code}</strong>
                    </p>
                </div>
            ) : (
                <div>
                    <div className="students-list-header">
                        <div className="students-list-title">
                            <Users size={18} />
                            <span>Estudiantes Inscritos</span>
                        </div>
                        <span className="students-count-pill">{students.length}</span>
                    </div>
                    <div className="students-grid">
                    {students.map((student, index) => (
                        <div key={student.id} className="student-card">
                            <div className="card-student-main">
                                <div className="student-avatar">
                                    {student.username?.charAt(0).toUpperCase() || (index + 1)}
                                </div>
                                <div className="student-info">
                                    <div className="student-name">{student.username || 'Unknown'}</div>
                                    <div className="student-email">{student.email}</div>
                                    <div className="student-meta">Inscrito en el curso</div>
                                </div>
                            </div>
                            {isTeacher && (
                                <button
                                    className="btn-remove-student"
                                    onClick={() => handleRemoveStudent(student.id)}
                                    disabled={processingId === student.id}
                                    title="Eliminar estudiante del curso"
                                >
                                    <Trash2 size={16} />
                                </button>
                            )}
                        </div>
                    ))}
                    </div>
                </div>
            )}

            <div className="students-stats">
                <div className="stat-card">
                    <div className="stat-value">{students.length}</div>
                    <div className="stat-label">Total de Estudiantes</div>
                </div>
            </div>
        </div>
    );
};

export default CourseStudents;
