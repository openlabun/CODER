import { useState, useEffect, useContext } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import client from '../api/client';
import { getExamDetails } from '../api/exams';
import { AuthContext } from '../context/AuthContext';
import Swal from 'sweetalert2';
import {
    Save, X, Clock, Calendar, FileText, Layout, Trash2,
    PlusCircle, ChevronRight, Code, Target, Search, Info, Eye, Globe, Users, BookOpen, Lock
} from 'lucide-react';
import ExamPreview, { buildPreviewChallenges } from '../components/ExamPreview';
import PageLoader from '../components/PageLoader';
import './CreateCourse.css';
import './Challenges.css';

const ExamEditor = () => {
    const { id } = useParams();
    const navigate = useNavigate();
    const { user } = useContext(AuthContext);
    const userId = user?.id || user?.ID;

    const [exam, setExam] = useState(null);
    const [formData, setFormData] = useState({
        title: '', description: '', visibility: 'private',
        startTime: '', endTime: '', timeLimit: 60,
        tryLimit: 1, allowLateSubmissions: false,
        isTimeUnlimited: false, isTryUnlimited: false, resultsVisibility: 'after_mine'
    });
    const [examItems, setExamItems] = useState([]);
    const [challenges, setChallenges] = useState([]);
    const [searchChallenge, setSearchChallenge] = useState('');
    const [showAddPanel, setShowAddPanel] = useState(false);
    const [loading, setLoading] = useState(true);
    const [saving, setSaving] = useState(false);
    const [addingItem, setAddingItem] = useState(null);
    const [isPreviewMode, setIsPreviewMode] = useState(false);
    const [previewLoading, setPreviewLoading] = useState(false);
    const [previewCodeMap, setPreviewCodeMap] = useState({});
    const [previewCurrentIndex, setPreviewCurrentIndex] = useState(0);
    const [previewLanguage, setPreviewLanguage] = useState('python');
    const [publicTestCasesMap, setPublicTestCasesMap] = useState({});

    const isProfessor = user?.role === 'professor' || user?.role === 'teacher' || user?.role === 'admin';

    // --- Fetch Exam ---
    useEffect(() => {
        if (!isProfessor) { navigate('/public-exams'); return; }
        const fetchAll = async () => {
            try {
                const data = await getExamDetails(id);
                const e = data || {};
                const profId = e.professorID || e.ProfessorID || e.professor_id || '';

                if (profId && userId && profId !== userId) {
                    Swal.fire({ icon: 'error', title: 'Acceso denegado', text: 'Solo el creador de este examen puede editarlo.' });
                    navigate('/public-exams');
                    return;
                }

                setExam(e);
                const tl = e.timeLimit || e.TimeLimit || e.time_limit || 3600;
                const st = e.startTime || e.StartTime || e.start_time || '';
                const et = e.endTime || e.EndTime || e.end_time || '';
                setFormData({
                    title: e.title || e.Title || '',
                    description: e.description || e.Description || '',
                    visibility: e.visibility || e.Visibility || 'private',
                    startTime: st ? new Date(st).toISOString().slice(0, 16) : '',
                    endTime: et ? new Date(et).toISOString().slice(0, 16) : '',
                    timeLimit: tl === -1 ? 60 : Math.floor(tl / 60),
                    tryLimit: (e.tryLimit === -1 || e.TryLimit === -1 || e.try_limit === -1) ? 1 : (e.tryLimit || e.TryLimit || e.try_limit || 1),
                    allowLateSubmissions: e.allowLateSubmissions || e.AllowLateSubmissions || false,
                    isTimeUnlimited: tl === -1,
                    isTryUnlimited: e.tryLimit === -1 || e.TryLimit === -1 || e.try_limit === -1,
                    resultsVisibility: 'after_mine'
                });

                // Fetch exam items
                const itemsRes = await client.get(`/exams/${id}/items`);
                const items = Array.isArray(itemsRes.data) ? itemsRes.data : (itemsRes.data?.items || []);
                setExamItems(items);

                // Fetch available challenges for adding
                const chRes = await client.get('/challenges');
                const chList = Array.isArray(chRes.data) ? chRes.data : (chRes.data?.items || []);
                setChallenges(chList);
            } catch (err) {
                console.error(err);
                Swal.fire({ icon: 'error', title: 'Error', text: 'No se pudo cargar el examen.' });
                navigate('/public-exams');
            } finally {
                setLoading(false);
            }
        };
        fetchAll();
    }, [id, isProfessor, navigate, userId]);

    const handleChange = (e) => {
        const { name, value, type, checked } = e.target;
        setFormData(prev => ({ ...prev, [name]: type === 'checkbox' ? checked : value }));
    };

    const previewChallenges = buildPreviewChallenges(examItems);

    const loadPreviewTestCases = async () => {
        if (!previewChallenges.length) {
            setPublicTestCasesMap({});
            return true;
        }

        setPreviewLoading(true);
        try {
            const results = await Promise.all(
                previewChallenges.map(async (challenge) => {
                    const response = await client.get(`/test-cases/challenge/${challenge.id}?exam_id=${id}`);
                    const visibleCases = Array.isArray(response.data)
                        ? response.data.filter((testCase) => testCase.type === 'public' || testCase.is_sample || testCase.isSample)
                        : [];

                    return [challenge.id, visibleCases];
                })
            );

            setPublicTestCasesMap(Object.fromEntries(results));
            return true;
        } catch (error) {
            console.error('Error loading preview test cases:', error);
            Swal.fire({
                icon: 'error',
                title: 'No se pudo preparar el preview',
                text: 'Hubo un problema cargando los casos de prueba visibles del examen.'
            });
            return false;
        } finally {
            setPreviewLoading(false);
        }
    };

    const handleEnterPreview = async () => {
        if (!previewChallenges.length) {
            await Swal.fire({
                icon: 'info',
                title: 'Preview limitado',
                text: 'Este examen aun no tiene retos asignados. Primero agrega al menos un reto para ver la experiencia del estudiante.'
            });
            return;
        }

        const { isConfirmed } = await Swal.fire({
            icon: 'question',
            title: 'Entrar al modo preview',
            text: 'Vas a abrir una vista previa del examen. Podras navegar y escribir codigo, pero nada se guardará ni se enviará.',
            showCancelButton: true,
            confirmButtonText: 'Entrar al preview',
            cancelButtonText: 'Cancelar',
            confirmButtonColor: '#c8102e'
        });

        if (!isConfirmed) {
            return;
        }

        setPreviewCurrentIndex(0);
        setPreviewLanguage('python');
        setPreviewCodeMap({});
        setIsPreviewMode(true);
        const prepared = await loadPreviewTestCases();
        if (!prepared) {
            setIsPreviewMode(false);
            return;
        }
    };

    const handleExitPreview = async () => {
        const hasWrittenCode = Object.values(previewCodeMap).some((snippet) => String(snippet || '').trim().length > 0);

        if (hasWrittenCode) {
            const { isConfirmed } = await Swal.fire({
                icon: 'warning',
                title: 'Salir del modo preview',
                text: 'El codigo escrito en esta vista previa se descartara. La configuracion del examen seguira intacta.',
                showCancelButton: true,
                confirmButtonText: 'Salir del preview',
                cancelButtonText: 'Seguir revisando',
                confirmButtonColor: '#c8102e'
            });

            if (!isConfirmed) {
                return;
            }
        }

        setIsPreviewMode(false);
    };

    // --- Save Exam (PATCH) ---
    const handleSave = async () => {
        setSaving(true);
        try {
            const payload = {};
            if (formData.title.trim()) payload.title = formData.title.trim();
            if (formData.description.trim()) payload.description = formData.description.trim();
            if (formData.visibility) payload.visibility = formData.visibility;
            if (formData.startTime) payload.start_time = new Date(formData.startTime).toISOString();
            if (formData.endTime) payload.end_time = new Date(formData.endTime).toISOString();
            payload.time_limit = formData.isTimeUnlimited ? -1 : parseInt(formData.timeLimit) * 60;
            payload.try_limit = formData.isTryUnlimited ? -1 : parseInt(formData.tryLimit);
            payload.allow_late_submissions = formData.allowLateSubmissions;

            await client.patch(`/exams/${id}`, payload);
            Swal.fire({ icon: 'success', title: 'Examen Actualizado', timer: 1500, toast: true, position: 'top-end', showConfirmButton: false });
        } catch (err) {
            console.error(err);
            Swal.fire({ icon: 'error', title: 'Error', text: err.response?.data?.error || 'No se pudo actualizar el examen.' });
        } finally {
            setSaving(false);
        }
    };

    // --- Add Challenge to Exam ---
    const handleAddChallenge = async (challenge) => {
        const challengeId = challenge.id || challenge.ID;
        const currentTotalPoints = examItems.reduce((acc, item) => acc + (item.points || item.Points || 0), 0);
        const remainingPoints = Math.max(0, 100 - currentTotalPoints);

        const { value: formValues } = await Swal.fire({
            title: 'Configurar Reto',
            html: `
                <div class="exam-config-swal-body">
                    <p class="exam-config-swal-note">Define los valores iniciales del reto dentro del examen.</p>
                    <div class="exam-config-grid">
                        <label class="exam-config-field" for="swal-points">
                            <span>Puntos (máx 100)</span>
                            <input id="swal-points" type="number" min="0" max="100" value="${Math.min(100, Math.max(0, remainingPoints || 100))}" class="swal2-input exam-config-input">
                            <small>Disponibles para asignar: <strong>${remainingPoints}</strong></small>
                        </label>
                        <label class="exam-config-field" for="swal-order">
                            <span>Orden</span>
                            <input id="swal-order" type="number" min="1" value="${examItems.length + 1}" class="swal2-input exam-config-input">
                            <small>Posición sugerida según el listado actual.</small>
                        </label>
                    </div>
                </div>
            `,
            customClass: {
                popup: 'exam-config-swal'
            },
            focusConfirm: false,
            showCancelButton: true,
            confirmButtonText: 'Añadir',
            cancelButtonText: 'Cancelar',
            confirmButtonColor: '#c8102e',
            preConfirm: () => {
                const pts = parseInt(document.getElementById('swal-points').value) || 0;
                const ord = parseInt(document.getElementById('swal-order').value) || 1;
                if (pts < 0 || pts > 100) { Swal.showValidationMessage('Los puntos deben estar entre 0 y 100'); return false; }
                if (currentTotalPoints + pts > 100) { Swal.showValidationMessage(`Excedes los 100 puntos totales. Podrías dar hasta ${100 - currentTotalPoints} pts.`); return false; }
                return { points: pts, order: ord };
            }
        });

        if (!formValues) return;

        setAddingItem(challengeId);
        try {
            await client.post('/exam-items', {
                exam_id: id,
                challenge_id: challengeId,
                order: formValues.order,
                points: formValues.points
            });
            const itemsRes = await client.get(`/exams/${id}/items`);
            const items = Array.isArray(itemsRes.data) ? itemsRes.data : (itemsRes.data?.items || []);
            setExamItems(items);
            Swal.fire({ icon: 'success', title: 'Reto Añadido', timer: 1000, toast: true, position: 'top-end', showConfirmButton: false });
        } catch (err) {
            console.error(err);
            Swal.fire({ icon: 'error', title: 'Error', text: err.response?.data?.error || 'No se pudo añadir el reto.' });
        } finally {
            setAddingItem(null);
        }
    };

    // --- Update Exam Item (points / order) ---
    const handleUpdateItem = async (itemId, updates) => {
        try {
            await client.patch(`/exam-items/${itemId}`, updates);
            const itemsRes = await client.get(`/exams/${id}/items`);
            const items = Array.isArray(itemsRes.data) ? itemsRes.data : (itemsRes.data?.items || []);
            setExamItems(items);
            Swal.fire({ icon: 'success', title: 'Actualizado', timer: 800, toast: true, position: 'top-end', showConfirmButton: false });
        } catch (err) {
            console.error(err);
            Swal.fire({ icon: 'error', title: 'Error', text: err.response?.data?.error || 'No se pudo actualizar.' });
        }
    };

    // --- Remove Challenge from Exam ---
    const handleRemoveItem = async (itemId) => {
        const { isConfirmed } = await Swal.fire({
            title: '¿Quitar reto?', text: 'Se eliminará este reto del examen.',
            icon: 'warning', showCancelButton: true, confirmButtonText: 'Sí, quitar', cancelButtonText: 'Cancelar'
        });
        if (!isConfirmed) return;
        try {
            await client.delete(`/exam-items/${itemId}`);
            setExamItems(prev => prev.filter(i => (i.id || i.ID) !== itemId));
            Swal.fire({ icon: 'success', title: 'Reto Eliminado', timer: 1000, toast: true, position: 'top-end', showConfirmButton: false });
        } catch {
            Swal.fire({ icon: 'error', title: 'Error', text: 'No se pudo quitar el reto.' });
        }
    };

    // Challenges already linked
    const linkedChallengeIds = new Set(
        examItems.map(item => item.challenge?.id || item.challenge?.ID || item.challengeID || item.challenge_id || '')
    );

    const availableChallenges = challenges.filter(c => {
        const cid = c.id || c.ID;
        if (linkedChallengeIds.has(cid)) return false;
        if (!searchChallenge) return true;
        return (c.title || '').toLowerCase().includes(searchChallenge.toLowerCase());
    });

    if (loading) return (
        <div className="create-course-page">
            <PageLoader message="Cargando examen..." />
        </div>
    );

    if (isPreviewMode) {
        return (
            <div className="create-course-page" style={{ maxWidth: 'none' }}>
                <ExamPreview
                    examTitle={formData.title || exam?.title || exam?.Title}
                    examDescription={formData.description || exam?.description || exam?.Description}
                    timeLimitMinutes={Number(formData.timeLimit) || 0}
                    tryLimit={Number(formData.tryLimit) || 0}
                    challenges={previewChallenges}
                    publicTestCasesMap={publicTestCasesMap}
                    previewCodeMap={previewCodeMap}
                    setPreviewCodeMap={setPreviewCodeMap}
                    previewLanguage={previewLanguage}
                    setPreviewLanguage={setPreviewLanguage}
                    previewCurrentIndex={previewCurrentIndex}
                    setPreviewCurrentIndex={setPreviewCurrentIndex}
                    onExit={handleExitPreview}
                    isLoading={previewLoading}
                />
            </div>
        );
    }

    return (
        <div className="create-course-page">
            <div className="page-header">
                <div className="header-content">
                    <h1>Editar Examen</h1>
                    <p className="subtitle">Configura y añade retos a tu evaluación</p>
                </div>
                <div style={{ display: 'flex', gap: '0.75rem' }}>
                    <button className="btn-secondary" onClick={handleEnterPreview}>
                        <Eye size={18} /> Preview
                    </button>
                    <button className="btn-secondary" onClick={() => navigate(-1)}><X size={18} /> Cancelar</button>
                    <button className="btn-primary" onClick={handleSave} disabled={saving}>
                        {saving ? 'Guardando...' : <><Save size={18} /> Guardar Cambios</>}
                    </button>
                </div>
            </div>

            <div className="form-container">
                {/* --- Información General --- */}
                <div className="form-section">
                    <div className="section-header"><FileText size={20} /><h2>Información General</h2></div>
                    <div className="form-group">
                        <label>Título del Examen <span style={{ color: 'var(--primary)', marginLeft: '4px' }}>*</span></label>
                        <input type="text" name="title" value={formData.title} onChange={handleChange} placeholder="ej. Parcial 1: Algoritmos" required />
                    </div>
                    <div className="form-group">
                        <label>Descripción / Instrucciones *</label>
                        <textarea name="description" value={formData.description} onChange={handleChange} placeholder="Instrucciones para los estudiantes..." rows="4" />
                    </div>
                </div>

                {/* --- Programación --- */}
                <div className="form-section">
                    <div className="section-header"><Calendar size={20} /><h2>Programación</h2></div>
                    <div className="form-row">
                        <div className="form-group">
                            <label>Fecha y Hora de Inicio <span style={{ color: 'var(--primary)', marginLeft: '4px' }}>*</span></label>
                            <input type="datetime-local" name="startTime" value={formData.startTime} onChange={handleChange} required />
                        </div>
                        <div className="form-group">
                            <label>Fecha y Hora de Cierre (Opcional)</label>
                            <input type="datetime-local" name="endTime" value={formData.endTime} onChange={handleChange} />
                        </div>
                    </div>
                </div>

                {/* --- Restricciones --- */}
                <div className="form-section">
                    <div className="section-header"><Clock size={20} /><h2>Restricciones y Límites</h2></div>
                    <div className="form-row">
                        <div className="form-group">
                            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                                <label>Duración (minutos) <span style={{ color: 'var(--primary)', marginLeft: '4px' }}>*</span></label>
                                <label className="checkbox-label" style={{ margin: 0, fontSize: '0.85rem' }}>
                                    <input type="checkbox" name="isTimeUnlimited" checked={formData.isTimeUnlimited} onChange={handleChange} />
                                    <span>Sin límite</span>
                                </label>
                            </div>
                            <input type="number" name="timeLimit" value={formData.timeLimit} onChange={handleChange} min="1" disabled={formData.isTimeUnlimited} required={!formData.isTimeUnlimited} style={{ opacity: formData.isTimeUnlimited ? 0.5 : 1 }} />
                        </div>
                        <div className="form-group">
                            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                                <label>Límite de Intentos <span style={{ color: 'var(--primary)', marginLeft: '4px' }}>*</span></label>
                                <label className="checkbox-label" style={{ margin: 0, fontSize: '0.85rem' }}>
                                    <input type="checkbox" name="isTryUnlimited" checked={formData.isTryUnlimited} onChange={handleChange} />
                                    <span>Ilimitado</span>
                                </label>
                            </div>
                            <input type="number" name="tryLimit" value={formData.tryLimit} onChange={handleChange} min="1" disabled={formData.isTryUnlimited} required={!formData.isTryUnlimited} style={{ opacity: formData.isTryUnlimited ? 0.5 : 1 }} />
                        </div>
                    </div>

                    <div className="form-group" style={{ marginTop: '1rem' }}>
                        <label>Visibilidad de Resultados para el Estudiante</label>
                        <select
                            name="resultsVisibility"
                            value={formData.resultsVisibility}
                            onChange={handleChange}
                        >
                            <option value="none">No mostrar resultados nunca</option>
                            <option value="after_mine">Mostrar mis resultados al finalizar mi envío</option>
                            <option value="after_all">Mostrar resultados cuando finalice la actividad para todos</option>
                        </select>
                        <small>Nota: La funcionalidad de esta opción está sujeta a la conexión con el servidor.</small>
                    </div>
                    <div className="checkbox-group">
                        <label className="checkbox-label">
                            <input type="checkbox" name="allowLateSubmissions" checked={formData.allowLateSubmissions} onChange={handleChange} />
                            <span>Permitir entregas tardías</span>
                        </label>
                    </div>
                </div>

                {/* --- Visibilidad --- */}
                <div className="form-section">
                    <div className="section-header"><Layout size={20} /><h2>Visibilidad</h2></div>
                    <div className="radio-group grid-2 visibility-radio-group">
                        {[
                            { value: 'course', title: 'Solo mi Curso', desc: 'Visible solo para estudiantes inscritos', icon: BookOpen },
                            { value: 'public', title: 'Público Global', desc: 'Visible para toda la comunidad', icon: Globe },
                            { value: 'teachers', title: 'Solo Profesores', desc: 'Colabora con otros docentes', icon: Users },
                            { value: 'private', title: 'Privado / Borrador', desc: 'Solo tú puedes verlo', icon: Lock },
                        ].map(opt => (
                            <label key={opt.value} className={`radio-card visibility-radio-card ${formData.visibility === opt.value ? 'active' : ''}`}>
                                <input type="radio" name="visibility" value={opt.value} checked={formData.visibility === opt.value} onChange={handleChange} />
                                <div className="radio-content visibility-radio-content">
                                    <div className="visibility-title-row">
                                        <opt.icon size={16} className="visibility-icon" />
                                        <span className="radio-title">{opt.title}</span>
                                    </div>
                                    <small>{opt.desc}</small>
                                </div>
                            </label>
                        ))}
                    </div>
                </div>

                {/* =============  RETOS DEL EXAMEN  ============= */}
                <div className="form-section">
                    <div className="section-header exam-items-header">
                        <div className="exam-items-title-wrap">
                            <Target size={20} />
                            <h2>Retos del Examen ({examItems.length})</h2>
                            {examItems.length > 0 && (
                                <span className="exam-total-badge">
                                    Total: {examItems.reduce((sum, it) => sum + (it.points || it.Points || 0), 0)} pts
                                </span>
                            )}
                        </div>
                        <button
                            type="button"
                            className="btn-create-mini exam-add-toggle-btn"
                            onClick={() => setShowAddPanel(!showAddPanel)}
                        >
                            <PlusCircle size={16} />
                            <span>{showAddPanel ? 'Cerrar Panel' : 'Añadir Reto'}</span>
                        </button>
                    </div>

                    {/* LIST OF LINKED EXAM ITEMS */}
                    {examItems.length === 0 ? (
                        <div className="exam-items-empty">
                            <Target size={40} className="exam-items-empty-icon" />
                            <h3>Sin retos asignados</h3>
                            <p>Usa el botón "Añadir Reto" para vincular desafíos a este examen.</p>
                        </div>
                    ) : (
                        <div className="exam-items-grid">
                            {examItems.map((item, idx) => {
                                const itemId = item.id || item.ID;
                                const ch = item.challenge || {};
                                const title = ch.title || ch.Title || `Reto #${idx + 1}`;
                                const diff = (ch.difficulty || ch.Difficulty || 'medium').toLowerCase();
                                const points = item.points || item.Points || 0;
                                const order = item.order || item.Order || idx + 1;

                                return (
                                    <div key={itemId} className="exam-item-card">
                                        <div className={`exam-item-accent ${diff}`}></div>
                                        <div className="exam-item-card-main">
                                            <div className="exam-item-top-row">
                                                <div className="exam-item-title">
                                                    <Code size={16} />
                                                    <h3>{title}</h3>
                                                </div>
                                                <span className={`exam-item-diff-badge ${diff}`}>
                                                    {diff === 'easy' ? 'Fácil' : diff === 'hard' ? 'Difícil' : 'Medio'}
                                                </span>
                                            </div>
                                            <p className="exam-item-description">
                                                {(ch.description || ch.Description || 'Sin descripción.').slice(0, 100)}{(ch.description || '').length > 100 ? '…' : ''}
                                            </p>

                                            <div className="exam-item-controls">
                                                <div className="exam-item-fields">
                                                    <div className="exam-item-field">
                                                        <label className="exam-item-field-label">PTS:</label>
                                                        <input
                                                            className="exam-item-field-input"
                                                            type="number"
                                                            min="0"
                                                            max="100"
                                                            defaultValue={points}
                                                            onBlur={(e) => {
                                                                const val = Math.min(100, Math.max(0, parseInt(e.target.value) || 0));
                                                                const currentOthers = examItems.reduce((acc, it) => acc + ((it.id || it.ID) === itemId ? 0 : (it.points || it.Points || 0)), 0);
                                                                if (currentOthers + val > 100) {
                                                                    Swal.fire({ icon: 'error', title: 'Error', text: `La suma no puede exceder 100. Restan ${100 - currentOthers} pts.` });
                                                                    e.target.value = points;
                                                                    return;
                                                                }
                                                                if (val !== points) handleUpdateItem(itemId, { points: val });
                                                            }}
                                                        />
                                                    </div>
                                                    <div className="exam-item-field">
                                                        <label className="exam-item-field-label">ORDEN:</label>
                                                        <input
                                                            className="exam-item-field-input"
                                                            type="number"
                                                            min="1"
                                                            defaultValue={order}
                                                            onBlur={(e) => {
                                                                const val = Math.max(1, parseInt(e.target.value) || 1);
                                                                if (val !== order) handleUpdateItem(itemId, { order: val });
                                                            }}
                                                        />
                                                    </div>
                                                </div>

                                                <button
                                                    className="exam-item-delete-btn"
                                                    onClick={() => handleRemoveItem(itemId)}
                                                    title="Quitar del examen"
                                                    data-tooltip="Quitar reto"
                                                >
                                                    <Trash2 size={14} />
                                                </button>
                                            </div>
                                        </div>
                                    </div>
                                );
                            })}
                        </div>
                    )}

                    {/* --- ADD CHALLENGE PANEL --- */}
                    {showAddPanel && (
                        <div className="exam-add-panel">
                            <div className="exam-add-panel-search">
                                <Search size={18} />
                                <input
                                    type="text"
                                    placeholder="Buscar reto por título..."
                                    value={searchChallenge}
                                    onChange={(e) => setSearchChallenge(e.target.value)}
                                />
                            </div>
                            {availableChallenges.length === 0 ? (
                                <p className="exam-add-panel-empty">No hay retos disponibles para añadir.</p>
                            ) : (
                                <div className="exam-add-panel-list">
                                    {availableChallenges.map(ch => {
                                        const cid = ch.id || ch.ID;
                                        const diff = (ch.difficulty || 'medium').toLowerCase();
                                        return (
                                            <div key={cid} className="exam-add-panel-item">
                                                <div className="exam-add-panel-item-main">
                                                    <Code size={16} className="exam-add-panel-item-icon" />
                                                    <div>
                                                        <strong className="exam-add-panel-item-title">{ch.title}</strong>
                                                        <div className="exam-add-panel-meta">
                                                            <span className={`diff-pill ${diff} exam-add-panel-diff`}>
                                                                {diff === 'easy' ? 'Fácil' : diff === 'hard' ? 'Difícil' : 'Medio'}
                                                            </span>
                                                            <span className="exam-add-panel-status">{ch.status}</span>
                                                        </div>
                                                    </div>
                                                </div>
                                                <button
                                                    type="button"
                                                    className="exam-add-panel-add-btn"
                                                    onClick={() => handleAddChallenge(ch)}
                                                    disabled={addingItem === cid}
                                                >
                                                    {addingItem === cid ? 'Añadiendo...' : <><PlusCircle size={14} /> Añadir</>}
                                                </button>
                                            </div>
                                        );
                                    })}
                                </div>
                            )}
                        </div>
                    )}
                </div>
            </div>

            <div className="info-box-alt">
                <div className="info-icon"><Info size={20} /></div>
                <div className="info-text">
                    <h3>💡 Gestión de Retos</h3>
                    <p>
                        Añade retos existentes de tu repositorio a este examen. Cada reto se presenta como un ejercicio
                        dentro de la evaluación con su propio sistema de puntuación.
                    </p>
                </div>
            </div>
        </div>
    );
};

export default ExamEditor;
