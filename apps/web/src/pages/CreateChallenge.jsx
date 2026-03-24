import { useState, useEffect } from 'react';
import { useNavigate, useLocation, useParams } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import client from '../api/client';
import AIAssistantModal from '../components/AIAssistantModal';
import Swal from 'sweetalert2';
import './CreateChallenge.css';
import './Challenges.css';

const CreateChallenge = () => {
    const { user } = useAuth();
    const navigate = useNavigate();
    const location = useLocation();
    const { id } = useParams();
    const isEditing = !!id;

    const [activeTab, setActiveTab] = useState('basic');
    const [showPreview, setShowPreview] = useState(false);
    const [showAIModal, setShowAIModal] = useState(false);

    const queryParams = new URLSearchParams(location.search);
    const courseIdFromUrl = queryParams.get('courseId');

    const [formData, setFormData] = useState({
        title: '',
        description: '',
        difficulty: 'medium',
        timeLimit: 1000,
        memoryLimit: 256,
        tags: [],
        inputFormat: '',
        outputFormat: '',
        constraints: '',
        status: 'draft',
        courseId: courseIdFromUrl || null,
        examId: queryParams.get('examId') || null
    });

    const [publicTestCases, setPublicTestCases] = useState([]);
    const [hiddenTestCases, setHiddenTestCases] = useState([]);
    const [newTag, setNewTag] = useState('');
    const [loading, setLoading] = useState(false);
    const [fetching, setFetching] = useState(isEditing);
    const [courses, setCourses] = useState([]);
    const [exams, setExams] = useState([]);

    useEffect(() => {
        const fetchChallengeForEdit = async () => {
            if (!isEditing) return;
            setFetching(true);
            try {
                const { data: challenge } = await client.get(`/challenges/${id}`);

                setFormData({
                    title: challenge.title || challenge.Title || '',
                    description: challenge.description || challenge.Description || '',
                    difficulty: (challenge.difficulty || challenge.Difficulty || 'medium').toLowerCase(),
                    timeLimit: challenge.timeLimit || challenge.WorkerTimeLimit || 1000,
                    memoryLimit: challenge.memoryLimit || challenge.WorkerMemoryLimit || 256,
                    tags: challenge.tags || challenge.Tags || [],
                    inputFormat: challenge.inputFormat || challenge.InputFormat || '',
                    outputFormat: challenge.outputFormat || challenge.OutputFormat || '',
                    constraints: challenge.constraints || challenge.Constraints || '',
                    status: challenge.status || challenge.Status || 'draft',
                    courseId: challenge.courseId || challenge.CourseID || null,
                    examId: challenge.examId || challenge.ExamID || queryParams.get('examId') || null
                });

                try {
                    const { data: testCases } = await client.get(`/test-cases/challenge/${id}`);
                    const cases = Array.isArray(testCases) ? testCases : (testCases.items || []);

                    setPublicTestCases(cases.filter(tc => tc.type === 'public' || tc.Type === 'public' || tc.is_public || tc.isPublic));
                    setHiddenTestCases(cases.filter(tc => tc.type !== 'public' && tc.Type !== 'public' && !tc.is_public && !tc.isPublic));
                } catch (tcErr) {
                    console.warn('Test cases fetch failed:', tcErr);
                }
            } catch (err) {
                console.error('Error fetching challenge:', err);
                Swal.fire({
                    icon: 'error',
                    title: 'Error de carga',
                    text: 'No se pudo cargar el reto.',
                    toast: true,
                    position: 'top-end',
                    showConfirmButton: false,
                    timer: 4000
                });
            } finally {
                setFetching(false);
            }
        };

        const fetchCourses = async () => {
            try {
                const scope = (user?.role === 'professor' || user?.role === 'teacher' || user?.role === 'admin') ? '?scope=owned' : '';
                const { data } = await client.get(`/courses${scope}`);
                const coursesList = Array.isArray(data) ? data : (data.items || []);
                setCourses(coursesList);
            } catch (err) {
                console.error('Error fetching courses:', err);
            }
        };

        if (user) {
            fetchCourses();
            if (isEditing) fetchChallengeForEdit();
        }
    }, [id, isEditing, user?.role]);

    useEffect(() => {
        const fetchExams = async () => {
            if (!formData.courseId) {
                setExams([]);
                return;
            }
            try {
                const { data } = await client.get(`/exams/course/${formData.courseId}`);
                setExams(Array.isArray(data) ? data : (data.items || []));
            } catch (err) {
                console.error('Error fetching exams:', err);
                setExams([]);
            }
        };

        fetchExams();
    }, [formData.courseId]);

    const handleChange = (e) => {
        const { name, value } = e.target;
        setFormData(prev => ({ ...prev, [name]: value }));
    };

    const addPublicTestCase = () => {
        setPublicTestCases([...publicTestCases, { input: '', output: '', name: `Example ${publicTestCases.length + 1}`, type: 'public' }]);
    };

    const addHiddenTestCase = () => {
        setHiddenTestCases([...hiddenTestCases, { input: '', output: '', name: `Hidden ${hiddenTestCases.length + 1}`, type: 'hidden' }]);
    };

    const updateTestCase = (index, field, value, isPublic) => {
        const cases = isPublic ? [...publicTestCases] : [...hiddenTestCases];
        cases[index][field] = value;
        isPublic ? setPublicTestCases(cases) : setHiddenTestCases(cases);
    };

    const removeTestCase = (index, isPublic) => {
        const cases = isPublic ? [...publicTestCases] : [...hiddenTestCases];
        cases.splice(index, 1);
        isPublic ? setPublicTestCases(cases) : setHiddenTestCases(cases);
    };

    const addTag = (tag) => {
        if (tag && !formData.tags.includes(tag)) {
            setFormData(prev => ({ ...prev, tags: [...prev.tags, tag] }));
            setNewTag('');
        }
    };

    const removeTag = (tag) => {
        setFormData(prev => ({ ...prev, tags: prev.tags.filter(t => t !== tag) }));
    };

    const validateForm = () => {
        if (!formData.title || !formData.description) {
            Swal.fire({ icon: 'warning', title: 'Campos incompletos', text: 'El título y la descripción son requeridos.', timer: 1500, toast: true, position: 'top-end', showConfirmButton: false });
            return false;
        }
        if (hiddenTestCases.length < 3) {
            Swal.fire({ icon: 'warning', title: 'Faltan casos ocultos', text: 'Se requieren al menos 3 casos ocultos.', timer: 1500, toast: true, position: 'top-end', showConfirmButton: false });
            return false;
        }
        const allCases = [...publicTestCases, ...hiddenTestCases];
        for (let tc of allCases) {
            if (!tc.input.trim() || !tc.output.trim()) {
                Swal.fire({ icon: 'warning', title: 'Casos incompletos', text: 'Todos los casos deben tener entrada y salida.', timer: 1500, toast: true, position: 'top-end', showConfirmButton: false });
                return false;
            }
        }
        return true;
    };

    const handleSubmit = async (status) => {
        if (status === 'published' && !validateForm()) return;

        setLoading(true);
        try {
            const payload = {
                title: formData.title,
                description: formData.description,
                difficulty: formData.difficulty,
                workerTimeLimit: parseInt(formData.timeLimit),
                workerMemoryLimit: parseInt(formData.memoryLimit),
                tags: formData.tags,
                inputVariables: [{ name: 'stdin', type: 'string', value: '' }],
                outputVariable: { name: 'stdout', type: 'string', value: '' },
                constraints: formData.constraints,
                status: status || formData.status,
                courseId: formData.courseId,
                examId: formData.examId || null
            };

            let response;
            if (isEditing) {
                response = await client.patch(`/challenges/${id}`, payload);
            } else {
                response = await client.post('/challenges', payload);
            }

            const challengeId = response.data.id || response.data.ID || id;

            // Guardado automático y en cascada de Casos de Prueba
            if (!isEditing && challengeId) {
                const allCasesToSave = [
                    ...publicTestCases.map(tc => ({ ...tc, isSample: true, points: 0 })),
                    ...hiddenTestCases.map(tc => ({ ...tc, isSample: false, points: 10 }))
                ];

                const tcRequests = allCasesToSave.map(tc => 
                    client.post('/test-cases', {
                        name: tc.name,
                        input: [{ name: "stdin", type: "string", value: tc.input }],
                        expectedOutput: { name: "stdout", type: "string", value: tc.output },
                        isSample: tc.isSample,
                        points: tc.points,
                        challengeId: challengeId
                    })
                );

                await Promise.all(tcRequests);
            }

            Swal.fire({
                icon: 'success',
                title: isEditing ? '¡Actualizado!' : '¡Creado!',
                text: `Reto ${isEditing ? 'actualizado' : 'publicado'} exitosamente`,
                timer: 1000,
                showConfirmButton: false,
                toast: true,
                position: 'top-end'
            });

            if (formData.courseId && challengeId) {
                try {
                    await client.post(`/courses/${formData.courseId}/challenges`, { challengeId });
                } catch (assignErr) {
                    console.warn('Silent failure assigning challenge to course:', assignErr);
                }
            }
            
            setTimeout(() => navigate('/challenges'), 1000);
        } catch (err) {
            console.error('Error en handleSubmit:', err.response?.data || err);
            const serverMsg = err.response?.data?.error || err.response?.data?.message;
            Swal.fire({
                icon: 'error',
                title: 'No se pudo guardar',
                text: serverMsg ? `${serverMsg}` : 'Hubo un problema al guardar el reto.',
                toast: true,
                position: 'top-end',
                showConfirmButton: false,
                timer: 4000
            });
        } finally {
            setLoading(false);
        }
    };

    const handleApplyIdea = (idea) => {
        let diff = (idea.difficulty || 'medium').toLowerCase();
        if (diff === 'fácil' || diff === 'facil' || diff === 'easy') diff = 'easy';
        else if (diff === 'difícil' || diff === 'dificil' || diff === 'hard') diff = 'hard';
        else diff = 'medium';

        setFormData(prev => ({
            ...prev,
            title: idea.title,
            description: idea.description,
            difficulty: diff,
            tags: idea.tags || [],
            inputFormat: idea.inputFormat || '',
            outputFormat: idea.outputFormat || '',
            constraints: idea.constraints || '',
            status: 'draft'
        }));
        if (idea.publicTestCases) setPublicTestCases(idea.publicTestCases);
        if (idea.hiddenTestCases) setHiddenTestCases(idea.hiddenTestCases);
        setActiveTab('basic');
    };

    if (fetching) return <div className="loading">Cargando datos...</div>;

    return (
        <div className="create-challenge-page">
            <div className="page-header">
                <div>
                    <h1>{isEditing ? 'Editar Reto' : 'Crear Nuevo Reto'}</h1>
                    <p className="subtitle">{isEditing ? 'Modifica tu desafío' : 'Diseña un nuevo desafío de programación'}</p>
                </div>
                <button className="btn-ai-assist" onClick={() => setShowAIModal(true)}>✨ Asistente IA</button>
            </div>

            {showAIModal && (
                <AIAssistantModal
                    onClose={() => setShowAIModal(false)}
                    onApplyIdea={handleApplyIdea}
                    onApplyTestCases={(cases) => {
                        if (cases.publicTestCases) setPublicTestCases(cases.publicTestCases);
                        if (cases.hiddenTestCases) setHiddenTestCases(cases.hiddenTestCases);
                        setFormData(prev => ({ ...prev, status: 'draft' }));
                        setActiveTab('testcases');
                    }}
                />
            )}

            {!showPreview ? (
                <>
                    <div className="tabs">
                        <button className={activeTab === 'basic' ? 'tab active' : 'tab'} onClick={() => setActiveTab('basic')}>📝 Básicos</button>
                        <button className={activeTab === 'testcases' ? 'tab active' : 'tab'} onClick={() => setActiveTab('testcases')}>🧪 Pruebas</button>
                        <button className={activeTab === 'settings' ? 'tab active' : 'tab'} onClick={() => setActiveTab('settings')}>⚙️ Ajustes</button>
                    </div>

                    <div className="tab-content">
                        {activeTab === 'basic' && (
                            <div className="form-section">
                                <div className="form-group">
                                    <label>Título *</label>
                                    <input type="text" name="title" value={formData.title} onChange={handleChange} placeholder="Ej: Suma A+B" required />
                                </div>
                                <div className="form-group">
                                    <label>Descripción *</label>
                                    <textarea name="description" value={formData.description} onChange={handleChange} placeholder="Enunciado..." rows="8" required />
                                </div>
                                <div className="form-row">
                                    <div className="form-group">
                                        <label>Dificultad</label>
                                        <select name="difficulty" value={formData.difficulty} onChange={handleChange}>
                                            <option value="easy">Fácil</option>
                                            <option value="medium">Medio</option>
                                            <option value="hard">Difícil</option>
                                        </select>
                                    </div>
                                    <div className="form-group">
                                        <label>Tiempo (ms)</label>
                                        <input type="number" name="timeLimit" value={formData.timeLimit} onChange={handleChange} />
                                    </div>
                                </div>
                            </div>
                        )}

                        {activeTab === 'testcases' && (
                            <div className="form-section">
                                <h3>Casos Públicos (Ejemplos)</h3>
                                {publicTestCases.map((tc, idx) => (
                                    <div key={idx} className="testcase-item">
                                        <input value={tc.input} onChange={(e) => updateTestCase(idx, 'input', e.target.value, true)} placeholder="Entrada" />
                                        <input value={tc.output} onChange={(e) => updateTestCase(idx, 'output', e.target.value, true)} placeholder="Salida" />
                                        <button onClick={() => removeTestCase(idx, true)}>🗑️</button>
                                    </div>
                                ))}
                                <button onClick={addPublicTestCase} className="btn-add">+ Caso Público</button>

                                <h3 style={{ marginTop: '2rem' }}>Casos Ocultos (Evaluación)</h3>
                                {hiddenTestCases.map((tc, idx) => (
                                    <div key={idx} className="testcase-item">
                                        <input value={tc.input} onChange={(e) => updateTestCase(idx, 'input', e.target.value, false)} placeholder="Entrada" />
                                        <input value={tc.output} onChange={(e) => updateTestCase(idx, 'output', e.target.value, false)} placeholder="Salida" />
                                        <button onClick={() => removeTestCase(idx, false)}>🗑️</button>
                                    </div>
                                ))}
                                <button onClick={addHiddenTestCase} className="btn-add">+ Caso Oculto</button>
                            </div>
                        )}

                        {activeTab === 'settings' && (
                            <div className="form-section">
                                <div className="form-group">
                                    <label>Curso Destino</label>
                                    <select name="courseId" value={formData.courseId || ''} onChange={handleChange}>
                                        <option value="">Seleccionar curso...</option>
                                        {courses.map(c => <option key={c.id} value={c.id}>{c.name}</option>)}
                                    </select>
                                </div>
                                <div className="form-group">
                                    <label>Examen Asociado</label>
                                    <select name="examId" value={formData.examId || ''} onChange={handleChange}>
                                        <option value="">Ninguno / Autónomo</option>
                                        {exams.map(e => <option key={e.id} value={e.id}>{e.title}</option>)}
                                    </select>
                                    <small>Asocia este reto a un examen específico para que aparezca en la evaluación.</small>
                                </div>
                                <div className="form-group">
                                    <label>Estado Inicial</label>
                                    <select name="status" value={formData.status} onChange={handleChange}>
                                        <option value="draft">Borrador</option>
                                        <option value="published">Publicado</option>
                                        <option value="private">Privado</option>
                                    </select>
                                </div>
                            </div>
                        )}
                    </div>

                    <div className="form-actions">
                        <button onClick={() => navigate('/challenges')} className="btn-secondary">Cancelar</button>
                        <button onClick={() => handleSubmit('draft')} disabled={loading} className="btn-draft">Guardar Borrador</button>
                        <button onClick={() => handleSubmit('published')} disabled={loading} className="btn-publish">🚀 {isEditing ? 'Actualizar' : 'Publicar'}</button>
                    </div>
                </>
            ) : (
                <div className="preview-container">
                    {/* Simplified Preview */}
                    <button onClick={() => setShowPreview(false)}>Volver</button>
                    <h2>{formData.title}</h2>
                    <p>{formData.description}</p>
                </div>
            )}
        </div>
    );
};

export default CreateChallenge;
