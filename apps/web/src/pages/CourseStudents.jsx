import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import client from '../api/client';
import { useAuth } from '../context/AuthContext';
import { Trash2, Copy, Link as LinkIcon, UploadCloud, FileText, CheckCircle2, Key } from 'lucide-react';
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

    if (loading) return <div className="loading">Cargando...</div>;

    return (
        <div className="course-students-page">
            <div className="page-header">
                <div>
                    <h1>Gestión de Estudiantes</h1>
                    <p className="subtitle">{course?.name} ({course?.code})</p>
                </div>
                <button onClick={() => navigate(`/courses/${id}`)} className="btn-back">
                    ← Volver al Curso
                </button>
            </div>

            {isTeacher && (
                <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(320px, 1fr))', gap: '2rem', marginBottom: '2rem', alignItems: 'stretch' }}>
                    
                    {/* Add Single Student */}
                    <div className="admin-actions-card" style={{ margin: 0 }}>
                        <h3>Añadir Estudiante</h3>
                        <p>Agrega un alumno usando su correo.</p>
                        <form className="add-student-form" onSubmit={handleAddStudent} style={{ marginTop: '1rem' }}>
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
                    <div className="admin-actions-card" style={{ margin: 0, display: 'flex', flexDirection: 'column' }}>
                        <h3>Carga Automática CSV</h3>
                        <p>Inscribe múltiples estudiantes. Típicamente los correos separados por comas o por líneas.</p>
                        <div style={{ marginTop: 'auto', paddingTop: '1rem', display: 'flex', gap: '10px' }}>
                            <input 
                                type="file" 
                                accept=".csv" 
                                id="csv-upload"
                                style={{ display: 'none' }}
                                onChange={(e) => setCsvFile(e.target.files[0])}
                            />
                            <label htmlFor="csv-upload" style={{ 
                                flex: 1, border: '1px dashed #cbd5e1', borderRadius: '8px', padding: '10px', 
                                display: 'flex', alignItems: 'center', justifyContent: 'center', 
                                cursor: 'pointer', color: csvFile ? '#4f46e5' : '#64748b', fontWeight: 600, fontSize: '0.85rem' 
                            }}>
                                <FileText size={18} style={{ marginRight: '8px' }} />
                                {csvFile ? csvFile.name : 'Seleccionar .csv'}
                            </label>
                            {csvFile && (
                                <button className="btn-action-filled" onClick={handleCsvUpload} disabled={isUploadingCsv} style={{ whiteSpace: 'nowrap' }}>
                                    <UploadCloud size={16} /> {isUploadingCsv ? 'Procesando...' : 'Cargar'}
                                </button>
                            )}
                        </div>
                    </div>

                    {/* Enrollment Methods Links */}
                    <div className="admin-actions-card" style={{ margin: 0 }}>
                        <h3>Comparte el Curso</h3>
                        <p>Invita a estudiantes rápida y masivamente.</p>
                        
                        <div style={{ marginTop: '1rem', display: 'flex', flexDirection: 'column', gap: '0.8rem' }}>
                            <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', width: '100%' }}>
                                <span style={{ flexShrink: 0, background: '#f1f5f9', padding: '6px', borderRadius: '8px', color: '#64748b' }}><LinkIcon size={16} /></span>
                                <div style={{ flex: 1, minWidth: 0, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap', fontSize: '0.85rem', fontWeight: 600, color: '#334155' }}>
                                    {window.location.origin}/courses/join?code={course?.enrollment_code}
                                </div>
                                <button className="btn-add-student" onClick={() => copyToClipboard(`${window.location.origin}/courses/join?code=${course?.enrollment_code}`, 'enlace')} style={{ padding: '6px 12px', flexShrink: 0 }}>
                                    <Copy size={14} /> Link
                                </button>
                            </div>
                            <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', width: '100%' }}>
                                <span style={{ flexShrink: 0, background: '#f1f5f9', padding: '6px', borderRadius: '8px', color: '#64748b' }}><Key size={16} /></span>
                                <div style={{ flex: 1, minWidth: 0, fontSize: '0.85rem', fontWeight: 600, color: '#334155', letterSpacing: '0.05em' }}>
                                    {course?.enrollment_code}
                                </div>
                                <button className="btn-add-student" style={{ padding: '6px 12px', background: '#e2e8f0', color: '#475569', flexShrink: 0 }} onClick={() => copyToClipboard(course?.enrollment_code, 'código')}>
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
