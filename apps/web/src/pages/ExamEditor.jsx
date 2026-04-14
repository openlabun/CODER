import { useState, useEffect, useContext } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import client from '../api/client';
import { getExamDetails } from '../api/exams';
import { AuthContext } from '../context/AuthContext';
import Swal from 'sweetalert2';
import {
    Save, X, Clock, Calendar, FileText, Layout, Trash2,
    PlusCircle, ChevronRight, Code, Target, Search, Info, Loader2
} from 'lucide-react';
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
        tryLimit: 1, allowLateSubmissions: false
    });
    const [examItems, setExamItems] = useState([]);
    const [challenges, setChallenges] = useState([]);
    const [searchChallenge, setSearchChallenge] = useState('');
    const [challengeVisibilityFilter, setChallengeVisibilityFilter] = useState('all');
    const [showAddPanel, setShowAddPanel] = useState(false);
    const [loading, setLoading] = useState(true);
    const [saving, setSaving] = useState(false);
    const [addingItem, setAddingItem] = useState(null);
    const [savingAllItems, setSavingAllItems] = useState(false);
    const [confirmRemoveItemId, setConfirmRemoveItemId] = useState(null);
    const [itemEdits, setItemEdits] = useState({});

    const isProfessor = user?.role === 'professor' || user?.role === 'teacher' || user?.role === 'admin';

    const getItemId = (item) => item.id || item.ID;
    const getItemOrder = (item, fallback = 1) => item.order || item.Order || fallback;
    const getItemPoints = (item) => item.points || item.Points || 0;

    const buildItemEdits = (items) => Object.fromEntries(
        items.map((item, index) => [getItemId(item), {
            order: String(getItemOrder(item, index + 1)),
            points: String(getItemPoints(item)),
        }])
    );

    const normalizeExamItems = (items) => [...items].sort((left, right) => {
        const leftOrder = getItemOrder(left, Number.MAX_SAFE_INTEGER);
        const rightOrder = getItemOrder(right, Number.MAX_SAFE_INTEGER);
        return leftOrder - rightOrder;
    });

    const refreshExamItems = async () => {
        const itemsRes = await client.get(`/exams/${id}/items`);
        const items = Array.isArray(itemsRes.data) ? itemsRes.data : (itemsRes.data?.items || []);
        const normalizedItems = normalizeExamItems(items);
        setExamItems(normalizedItems);
        setItemEdits(buildItemEdits(normalizedItems));
        return normalizedItems;
    };

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
                    timeLimit: Math.floor(tl / 60),
                    tryLimit: e.tryLimit || e.TryLimit || e.try_limit || 1,
                    allowLateSubmissions: e.allowLateSubmissions || e.AllowLateSubmissions || false,
                });

                // Fetch exam items
                await refreshExamItems();

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
    }, [id]);

    const handleChange = (e) => {
        const { name, value, type, checked } = e.target;
        setFormData(prev => ({ ...prev, [name]: type === 'checkbox' ? checked : value }));
    };

    const handleExamItemFieldChange = (itemId, field, value) => {
        setItemEdits(prev => ({
            ...prev,
            [itemId]: {
                ...prev[itemId],
                [field]: value,
            }
        }));
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
            payload.time_limit = parseInt(formData.timeLimit) * 60;
            payload.try_limit = parseInt(formData.tryLimit);
            payload.allow_late_submissions = formData.allowLateSubmissions;

            await client.patch(`/exams/${id}`, payload);
            Swal.fire({ icon: 'success', title: 'Examen Actualizado', timer: 1500, toast: true, position: 'top-end', showConfirmButton: false });
            navigate('/public-exams');
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

        const { value: formValues } = await Swal.fire({
            title: 'Configurar Reto',
            html:
                `<label style="display:block;text-align:left;font-weight:700;margin-bottom:4px">Puntos (máx 100)</label>` +
                `<input id="swal-points" type="number" min="0" max="100" value="100" class="swal2-input" style="margin:0 0 1rem 0">` +
                `<label style="display:block;text-align:left;font-weight:700;margin-bottom:4px">Orden</label>` +
                `<input id="swal-order" type="number" min="1" value="${examItems.length + 1}" class="swal2-input" style="margin:0">`,
            focusConfirm: false,
            showCancelButton: true,
            confirmButtonText: 'Añadir',
            cancelButtonText: 'Cancelar',
            confirmButtonColor: '#c8102e',
            preConfirm: () => {
                const pts = parseInt(document.getElementById('swal-points').value) || 0;
                const ord = parseInt(document.getElementById('swal-order').value) || 1;
                if (pts < 0 || pts > 100) { Swal.showValidationMessage('Los puntos deben estar entre 0 y 100'); return false; }
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
        try {
            await client.delete(`/exam-items/${itemId}`);
            setExamItems(prev => prev.filter(i => (i.id || i.ID) !== itemId));
            setItemEdits(prev => {
                const next = { ...prev };
                delete next[itemId];
                return next;
            });
            setConfirmRemoveItemId(null);
            Swal.fire({ icon: 'success', title: 'Reto Eliminado', timer: 1000, toast: true, position: 'top-end', showConfirmButton: false });
        } catch (err) {
            Swal.fire({ icon: 'error', title: 'Error', text: 'No se pudo quitar el reto.' });
        }
    };

    // Challenges already linked — use multiple fallbacks to handle DTO vs entity variations
    const linkedChallengeIds = new Set(
        examItems.flatMap(item => [
            item.challenge?.id,
            item.challenge?.ID,
            item.challenge_id,
            item.challengeID,
        ].filter(Boolean))
    );

    const normalizeStatus = (status) => String(status || 'draft').toLowerCase();

    const availableChallenges = challenges.filter(c => {
        const cid = c.id || c.ID;
        if (linkedChallengeIds.has(cid)) return false;
        
        // Filtro por título
        if (searchChallenge && !(c.title || '').toLowerCase().includes(searchChallenge.toLowerCase())) {
            return false;
        }
        
        // Filtro por visibilidad
        if (challengeVisibilityFilter !== 'all') {
            const status = normalizeStatus(c.status || c.Status);
            if (status !== challengeVisibilityFilter) return false;
        }
        
        return true;
    });

    if (loading) return (
        <div className="exam-loading-screen">
            <Loader2 className="exam-loading-spinner" size={48} />
            <p className="exam-loading-text">Cargando examen...</p>
        </div>
    );

    return (
        <div className="create-course-page">
            <div className="page-header">
                <div className="header-content">
                    <h1>Editar Examen</h1>
                    <p className="subtitle">Configura y añade retos a tu evaluación</p>
                </div>
                <div style={{ display: 'flex', gap: '0.75rem' }}>
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
                        <label>Título del Examen *</label>
                        <input type="text" name="title" value={formData.title} onChange={handleChange} placeholder="ej. Parcial 1: Algoritmos" required />
                    </div>
                    <div className="form-group">
                        <label>Descripción / Instrucciones</label>
                        <textarea name="description" value={formData.description} onChange={handleChange} placeholder="Instrucciones para los estudiantes..." rows="4" />
                    </div>
                </div>

                {/* --- Programación --- */}
                <div className="form-section">
                    <div className="section-header"><Calendar size={20} /><h2>Programación</h2></div>
                    <div className="form-row">
                        <div className="form-group">
                            <label>Fecha y Hora de Inicio *</label>
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
                            <label>Duración (minutos) *</label>
                            <input type="number" name="timeLimit" value={formData.timeLimit} onChange={handleChange} min="1" required />
                        </div>
                        <div className="form-group">
                            <label>Límite de Intentos *</label>
                            <input type="number" name="tryLimit" value={formData.tryLimit} onChange={handleChange} min="1" required />
                        </div>
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
                            { value: 'course', title: 'Solo mi Curso', desc: 'Visible solo para estudiantes inscritos' },
                            { value: 'public', title: 'Público Global', desc: 'Visible para toda la comunidad' },
                            { value: 'teachers', title: 'Solo Profesores', desc: 'Colabora con otros docentes' },
                            { value: 'private', title: 'Privado / Borrador', desc: 'Solo tú puedes verlo' },
                        ].map(opt => (
                            <label key={opt.value} className={`radio-card visibility-radio-card ${formData.visibility === opt.value ? 'active' : ''}`}>
                                <input type="radio" name="visibility" value={opt.value} checked={formData.visibility === opt.value} onChange={handleChange} />
                                <div className="radio-content visibility-radio-content">
                                    <span className="radio-title">{opt.title}</span>
                                    <small>{opt.desc}</small>
                                </div>
                            </label>
                        ))}
                    </div>
                </div>

                {/* =============  RETOS DEL EXAMEN  ============= */}
                <div className="form-section">
                    <div className="section-header" style={{ justifyContent: 'space-between' }}>
                        <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                            <Target size={20} />
                            <h2>Retos del Examen ({examItems.length})</h2>
                            {examItems.length > 0 && (
                                <span style={{ background: '#e0e7ff', color: '#4f46e5', padding: '3px 10px', borderRadius: '100px', fontSize: '0.75rem', fontWeight: 900, marginLeft: '0.5rem' }}>
                                    Total: {examItems.reduce((sum, it) => sum + (it.points || it.Points || 0), 0)} pts
                                </span>
                            )}
                        </div>
                        <div style={{ display: 'flex', gap: '0.5rem' }}>
                            {examItems.length > 0 && (
                                <button
                                    type="button"
                                    className="btn-primary"
                                    onClick={handleSaveAllItems}
                                    disabled={savingAllItems}
                                    style={{ height: '40px', fontSize: '0.85rem', padding: '0 1rem' }}
                                >
                                    {savingAllItems ? 'Guardando...' : <><Save size={15} /> Guardar orden y puntos</>}
                                </button>
                            )}
                            <button
                                type="button"
                                className="btn-create-mini"
                                onClick={() => setShowAddPanel(!showAddPanel)}
                                style={{ height: '40px', fontSize: '0.85rem' }}
                            >
                                <PlusCircle size={16} />
                                <span>{showAddPanel ? 'Cerrar Panel' : 'Añadir Reto'}</span>
                            </button>
                        </div>
                    </div>

                    {/* LIST OF LINKED EXAM ITEMS */}
                    {examItems.length === 0 ? (
                        <div style={{ textAlign: 'center', padding: '3rem 1rem', color: '#999' }}>
                            <Target size={40} style={{ marginBottom: '1rem', opacity: 0.3 }} />
                            <h3 style={{ fontWeight: 700, color: '#888' }}>Sin retos asignados</h3>
                            <p style={{ fontSize: '0.9rem' }}>Usa el botón "Añadir Reto" para vincular desafíos a este examen.</p>
                        </div>
                    ) : (
                        <div className="challenges-grid-compact" style={{ marginTop: '1rem' }}>
                            {examItems.map((item, idx) => {
                                const itemId = item.id || item.ID;
                                const ch = item.challenge || {};
                                const title = ch.title || ch.Title || `Reto #${idx + 1}`;
                                const diff = (ch.difficulty || ch.Difficulty || 'medium').toLowerCase();
                                const points = item.points || item.Points || 0;
                                const order = item.order || item.Order || idx + 1;

                                return (
                                    <div key={itemId} className="challenge-card-mini exam-item-card">
                                        <div className={`card-accent ${diff}`}></div>
                                        <div className="card-main exam-item-card-main">
                                            <div className="card-top">
                                                <div className="title-area exam-item-title">
                                                    <Code size={16} className="title-icon" />
                                                    <h3>{title}</h3>
                                                </div>
                                                <div className="badge-group">
                                                    <span className={`diff-pill ${diff}`}>
                                                        {diff === 'easy' ? 'Fácil' : diff === 'hard' ? 'Difícil' : 'Medio'}
                                                    </span>
                                                </div>
                                            </div>
                                            <p className="description-text" style={{ fontSize: '0.8rem', color: '#666', margin: '0.5rem 0' }}>
                                                {(ch.description || ch.Description || 'Sin descripción.').slice(0, 100)}{(ch.description || '').length > 100 ? '…' : ''}
                                            </p>
                                            <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem', marginTop: '0.75rem', flexWrap: 'wrap' }}>
                                                <div style={{ display: 'flex', alignItems: 'center', gap: '0.35rem' }}>
                                                    <label style={{ fontSize: '0.72rem', fontWeight: 800, color: '#64748b', textTransform: 'uppercase' }}>Pts:</label>
                                                    <input
                                                        type="number"
                                                        min="0"
                                                        max="100"
                                                        defaultValue={points}
                                                        onBlur={(e) => {
                                                            const val = Math.min(100, Math.max(0, parseInt(e.target.value) || 0));
                                                            if (val !== points) handleUpdateItem(itemId, { points: val });
                                                        }}
                                                        style={{ width: '55px', padding: '4px 6px', border: '1px solid #e2e8f0', borderRadius: '8px', fontSize: '0.8rem', fontWeight: 800, textAlign: 'center' }}
                                                    />
                                                </div>
                                                <div style={{ display: 'flex', alignItems: 'center', gap: '0.35rem' }}>
                                                    <label style={{ fontSize: '0.72rem', fontWeight: 800, color: '#64748b', textTransform: 'uppercase' }}>Orden:</label>
                                                    <input
                                                        type="number"
                                                        min="1"
                                                        defaultValue={order}
                                                        onBlur={(e) => {
                                                            const val = Math.max(1, parseInt(e.target.value) || 1);
                                                            if (val !== order) handleUpdateItem(itemId, { order: val });
                                                        }}
                                                        style={{ width: '50px', padding: '4px 6px', border: '1px solid #e2e8f0', borderRadius: '8px', fontSize: '0.8rem', fontWeight: 800, textAlign: 'center' }}
                                                    />
                                                </div>
                                                <div style={{ marginLeft: 'auto' }}>
                                                    <button
                                                        className="action-btn delete"
                                                        onClick={() => handleRemoveItem(itemId)}
                                                        title="Quitar del examen"
                                                        style={{ width: '32px', height: '32px', borderRadius: '50%', border: 'none', background: '#fee2e2', color: '#dc2626', cursor: 'pointer' }}
                                                    >
                                                        <Trash2 size={14} />
                                                    </button>
                                                </div>
                                            </div>

                                            <div className="meta-badges-row" style={{ marginTop: '0.5rem' }}>
                                                <span className={`diff-pill ${diff}`}>
                                                    {diff === 'easy' ? 'Fácil' : diff === 'hard' ? 'Difícil' : 'Medio'}
                                                </span>
                                            </div>

                                            <p className="description-text exam-item-description">
                                                {ch.description || ch.Description || 'Sin descripción.'}
                                            </p>

                                            <div className="card-footer-mini exam-item-fields">
                                                <div className="exam-item-field">
                                                    <span className="exam-item-field-label">Orden</span>
                                                    <input
                                                        type="number"
                                                        min="1"
                                                        value={order}
                                                        onChange={(event) => handleExamItemFieldChange(itemId, 'order', event.target.value)}
                                                        className="exam-item-field-input"
                                                    />
                                                </div>
                                                <div className="exam-item-field">
                                                    <span className="exam-item-field-label">Puntos</span>
                                                    <input
                                                        type="number"
                                                        min="0"
                                                        value={points}
                                                        onChange={(event) => handleExamItemFieldChange(itemId, 'points', event.target.value)}
                                                        className="exam-item-field-input"
                                                    />
                                                </div>
                                            </div>

                                            {confirmRemoveItemId === itemId && (
                                                <div className="exam-item-delete-confirm">
                                                    <span>¿Quitar este reto del examen?</span>
                                                    <div className="exam-item-delete-confirm-actions">
                                                        <button
                                                            type="button"
                                                            className="btn-secondary"
                                                            onClick={() => setConfirmRemoveItemId(null)}
                                                            style={{ padding: '6px 10px', fontSize: '0.78rem' }}
                                                        >
                                                            Cancelar
                                                        </button>
                                                        <button
                                                            type="button"
                                                            className="btn-primary"
                                                            onClick={() => handleRemoveItem(itemId)}
                                                            style={{ padding: '6px 10px', fontSize: '0.78rem' }}
                                                        >
                                                            Si, quitar
                                                        </button>
                                                    </div>
                                                </div>
                                            )}
                                        </div>
                                    </div>
                                );
                            })}
                        </div>
                    )}

                    {/* --- ADD CHALLENGE PANEL --- */}
                    {showAddPanel && (
                        <div style={{ marginTop: '1.5rem', background: '#f9fafb', borderRadius: '16px', padding: '1.5rem', border: '1px solid #e5e7eb' }}>
                            <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem', marginBottom: '1rem' }}>
                                <Search size={18} />
                                <input
                                    type="text"
                                    placeholder="Buscar reto por título..."
                                    value={searchChallenge}
                                    onChange={(e) => setSearchChallenge(e.target.value)}
                                    style={{ flex: 1, border: '1px solid #ddd', borderRadius: '10px', padding: '0.6rem 1rem', fontSize: '0.9rem' }}
                                />
                            </div>

                            {/* Visibility Filter Chips */}
                            <div style={{ display: 'flex', flexWrap: 'wrap', gap: '0.5rem', marginBottom: '1rem' }}>
                                {[
                                    { key: 'all', label: 'Todos' },
                                    { key: 'draft', label: 'Borrador' },
                                    { key: 'published', label: 'Publicado' },
                                    { key: 'private', label: 'Privado' },
                                ].map(filter => (
                                    <button
                                        key={filter.key}
                                        type="button"
                                        onClick={() => setChallengeVisibilityFilter(filter.key)}
                                        style={{
                                            padding: '0.4rem 0.8rem',
                                            borderRadius: '20px',
                                            border: challengeVisibilityFilter === filter.key ? '2px solid #c8102e' : '1px solid #ddd',
                                            background: challengeVisibilityFilter === filter.key ? '#fca5a5' : 'white',
                                            color: challengeVisibilityFilter === filter.key ? '#991b1b' : '#666',
                                            cursor: 'pointer',
                                            fontSize: '0.8rem',
                                            fontWeight: challengeVisibilityFilter === filter.key ? '600' : '500',
                                            transition: 'all 0.2s'
                                        }}
                                    >
                                        {filter.label}
                                    </button>
                                ))}
                            </div>
                            {availableChallenges.length === 0 ? (
                                <p style={{ textAlign: 'center', color: '#999', padding: '1rem' }}>
                                    {challenges.length === 0
                                        ? 'Aún no tienes retos en tu repositorio.'
                                        : 'Todos tus retos ya están añadidos a este examen.'}
                                </p>
                            ) : (
                                <div style={{ display: 'flex', flexDirection: 'column', gap: '0.5rem', maxHeight: '300px', overflowY: 'auto' }}>
                                    {availableChallenges.map(ch => {
                                        const cid = ch.id || ch.ID;
                                        const diff = (ch.difficulty || 'medium').toLowerCase();
                                        return (
                                            <div key={cid} style={{
                                                display: 'flex', justifyContent: 'space-between', alignItems: 'center',
                                                background: 'white', borderRadius: '12px', padding: '0.75rem 1rem',
                                                border: '1px solid #e5e7eb', transition: 'all 0.2s'
                                            }}>
                                                <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem' }}>
                                                    <Code size={16} style={{ color: '#c8102e' }} />
                                                    <div>
                                                        <strong style={{ fontSize: '0.9rem' }}>{ch.title}</strong>
                                                        <div style={{ display: 'flex', gap: '0.5rem', marginTop: '2px' }}>
                                                            <span className={`diff-pill ${diff}`} style={{ fontSize: '0.6rem', padding: '2px 8px' }}>
                                                                {diff === 'easy' ? 'Fácil' : diff === 'hard' ? 'Difícil' : 'Medio'}
                                                            </span>
                                                            <span style={{ fontSize: '0.7rem', color: '#999' }}>{ch.status}</span>
                                                        </div>
                                                    </div>
                                                </div>
                                                <button
                                                    onClick={() => handleAddChallenge(ch)}
                                                    disabled={addingItem === cid}
                                                    style={{
                                                        background: 'linear-gradient(135deg, #c8102e, #a00d25)', color: 'white',
                                                        border: 'none', borderRadius: '10px', padding: '0.5rem 1rem',
                                                        fontWeight: 700, fontSize: '0.8rem', cursor: 'pointer',
                                                        opacity: addingItem === cid ? 0.6 : 1
                                                    }}
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
