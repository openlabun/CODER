import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import client from '../api/client';
import { useAuth } from '../context/AuthContext';
import { Trash2 } from 'lucide-react';
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
            await client.post(`/courses/${id}/students`, { studentID: email });
            
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
                <div className="admin-actions-card">
                    <h3>Añadir Estudiante</h3>
                    <p>Agrega un alumno directamente usando su correo electrónico.</p>
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
            )}

            {students.length === 0 ? (
                <div className="empty-state">
                    <div className="empty-state-icon">👥</div>
                    <h3 className="empty-state-title">No hay estudiantes inscritos</h3>
                    <p className="empty-state-description">
                        Nadie se ha unido a este curso todavía.
                    </p>
                    <p className="empty-state-hint">
                        Comparte el código de inscripción: <strong>{course?.enrollmentCode}</strong>
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
