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
        inputVariables: [{ name: 'stdin', type: 'string' }],
        outputVariable: { name: 'stdout', type: 'string' },
        constraints: '',
        status: 'draft',
        codeTemplates: {},
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

    const SUPPORTED_LANGUAGES = [
        { id: 'python', label: 'Python', defaultTemplate: 'def solve():\n    pass\n' },
        { id: 'javascript', label: 'JavaScript', defaultTemplate: 'function solve() {\n\n}\n' },
        { id: 'java', label: 'Java', defaultTemplate: 'class Solution {\n    public static void solve() {\n\n    }\n}\n' },
        { id: 'cpp', label: 'C++', defaultTemplate: '#include <iostream>\nusing namespace std;\n\nvoid solve() {\n\n}\n' },
        { id: 'go', label: 'Go', defaultTemplate: 'package main\n\nfunc solve() {\n\n}\n' }
    ];

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
                    inputVariables: challenge.inputVariables || challenge.InputVariables || [{ name: 'stdin', type: 'string' }],
                    outputVariable: challenge.outputVariable || challenge.OutputVariable || { name: 'stdout', type: 'string' },
                    constraints: challenge.constraints || challenge.Constraints || '',
                    status: challenge.status || challenge.Status || 'draft',
                    codeTemplates: challenge.code_templates || challenge.CodeTemplates || {},
                    courseId: challenge.courseId || challenge.CourseID || null,
                    examId: challenge.examId || challenge.ExamID || queryParams.get('examId') || null
                });

                try {
                    const { data: testCases } = await client.get(`/test-cases/challenge/${id}`);
                    const mappedCases = testCases.map(tc => {
                        const inputs = tc.input || tc.Input || [];
                        const inputValues = {};
                        inputs.forEach(i => inputValues[i.name] = i.value);
                        
                        const expectedOut = tc.expectedOutput || tc.ExpectedOutput || {};
                        const outputValue = expectedOut.value || expectedOut.Value || '';
                        
                        return {
                            ...tc,
                            inputValues,
                            outputValue,
                            type: (tc.type || tc.Type || (tc.is_sample || tc.isSample ? 'public' : 'hidden'))
                        };
                    });

                    setPublicTestCases(mappedCases.filter(tc => tc.type === 'public' || tc.is_sample || tc.isSample));
                    setHiddenTestCases(mappedCases.filter(tc => tc.type !== 'public' && !tc.is_sample && !tc.isSample));
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
        setPublicTestCases([...publicTestCases, { inputValues: {}, outputValue: '', name: `Example ${publicTestCases.length + 1}`, type: 'public' }]);
    };

    const addHiddenTestCase = () => {
        setHiddenTestCases([...hiddenTestCases, { inputValues: {}, outputValue: '', name: `Hidden ${hiddenTestCases.length + 1}`, type: 'hidden' }]);
    };

    const updateTestCase = (index, field, value, isPublic) => {
        const cases = isPublic ? [...publicTestCases] : [...hiddenTestCases];
        cases[index][field] = value;
        isPublic ? setPublicTestCases(cases) : setHiddenTestCases(cases);
    };

    const handleInputVarChange = (index, field, value) => {
        const vars = [...formData.inputVariables];
        vars[index][field] = value;
        setFormData(prev => ({ ...prev, inputVariables: vars }));
    };

    const handleOutputVarChange = (field, value) => {
        setFormData(prev => ({ ...prev, outputVariable: { ...prev.outputVariable, [field]: value } }));
    };

    const updateTestCaseInput = (index, varName, field, value, isPublic) => {
        const cases = isPublic ? [...publicTestCases] : [...hiddenTestCases];
        if (!cases[index].inputValues) cases[index].inputValues = {};
        
        const currentVal = cases[index].inputValues[varName] || {};
        cases[index].inputValues[varName] = { ...currentVal, name: varName, [field]: value };
        
        isPublic ? setPublicTestCases(cases) : setHiddenTestCases(cases);
    };

    const removeTestCase = (index, isPublic) => {
        const cases = isPublic ? [...publicTestCases] : [...hiddenTestCases];
        cases.splice(index, 1);
        isPublic ? setPublicTestCases(cases) : setHiddenTestCases(cases);
    };

    const addInputVar = () => setFormData(prev => ({ ...prev, inputVariables: [...prev.inputVariables, { name: '', type: 'string' }] }));
    const removeInputVar = (index) => {
        const vars = [...formData.inputVariables];
        vars.splice(index, 1);
        setFormData(prev => ({ ...prev, inputVariables: vars }));
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
        if (!formData.outputVariable.name?.trim()) {
            Swal.fire({ icon: 'warning', title: 'Variable de salida', text: 'La variable de salida debe tener un nombre.', timer: 1500, toast: true, position: 'top-end', showConfirmButton: false });
            return false;
        }

        for (let iv of formData.inputVariables) {
            if (!iv.name?.trim()) {
                Swal.fire({ icon: 'warning', title: 'Variables de entrada', text: 'Todas las variables de entrada deben tener un nombre.', timer: 1500, toast: true, position: 'top-end', showConfirmButton: false });
                return false;
            }
        }

        const allCases = [...publicTestCases, ...hiddenTestCases];
        for (let tc of allCases) {
            if (!tc.outputValue?.toString().trim()) {
                Swal.fire({ icon: 'warning', title: 'Casos incompletos', text: 'Todos los casos deben tener una salida esperada.', timer: 1500, toast: true, position: 'top-end', showConfirmButton: false });
                return false;
            }
        }
        return true;
    };

    const handleSubmit = async (status) => {
        if (!validateForm()) return;

        setLoading(true);
        try {
            const payload = {
                title: formData.title.trim(),
                description: formData.description.trim(),
                difficulty: formData.difficulty,
                worker_time_limit: parseInt(formData.timeLimit),
                worker_memory_limit: parseInt(formData.memoryLimit),
                tags: formData.tags,
                input_variables: formData.inputVariables.map(v => ({ 
                    name: v.name.trim(), 
                    type: v.type || 'string', 
                    value: '' 
                })),
                output_variable: { 
                    name: formData.outputVariable.name.trim(), 
                    type: formData.outputVariable.type || 'string', 
                    value: '' 
                },
                constraints: formData.constraints,
                status: status || formData.status,
                code_templates: formData.codeTemplates,
                user_id: user?.id || user?.ID || ''
            };

            let response;
            if (isEditing) {
                payload.challenge_id = id;
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

                const tcRequests = allCasesToSave.map(tc => {
                    const inputsDto = formData.inputVariables.map(v => {
                        const tcVal = tc.inputValues?.[v.name] || {};
                        return {
                            name: v.name.trim(),
                            type: tcVal.type || v.type || 'string',
                            value: tcVal.value?.toString() || ''
                        };
                    });
                    
                    return client.post('/test-cases', {
                        name: tc.name || `Case ${allCasesToSave.indexOf(tc) + 1}`,
                        input: inputsDto,
                        expected_output: { 
                            name: formData.outputVariable.name.trim(), 
                            type: formData.outputVariable.type || 'string', 
                            value: tc.outputValue?.toString() || ''
                        },
                        is_sample: tc.isSample || false,
                        points: tc.points || 0,
                        challenge_id: challengeId
                    });
                });

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
            inputVariables: idea.inputVariables || idea.input_variables || [{ name: 'stdin', type: 'string' }],
            outputVariable: idea.outputVariable || idea.output_variable || { name: 'stdout', type: 'string' },
            timeLimit: idea.workerTimeLimit || idea.worker_time_limit || 1000,
            memoryLimit: idea.workerMemoryLimit || idea.worker_memory_limit || 256,
            constraints: idea.constraints || '',
            status: 'draft'
        }));
        if (idea.publicTestCases || idea.public_test_cases) {
            const arr = idea.publicTestCases || idea.public_test_cases;
            setPublicTestCases(arr.map(tc => {
                const ivs = {};
                (tc.input || []).forEach(i => {
                    ivs[i.name] = { name: i.name, type: i.type, value: i.value };
                });
                return { name: tc.name, type: 'public', inputValues: ivs, outputValue: tc.output?.value || '' };
            }));
        }
        if (idea.hiddenTestCases || idea.hidden_test_cases) {
            const arr = idea.hiddenTestCases || idea.hidden_test_cases;
            setHiddenTestCases(arr.map(tc => {
                const ivs = {};
                (tc.input || []).forEach(i => {
                    ivs[i.name] = { name: i.name, type: i.type, value: i.value };
                });
                return { name: tc.name, type: 'hidden', inputValues: ivs, outputValue: tc.output?.value || '' };
            }));
        }
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
                        if (cases.publicTestCases) {
                            setPublicTestCases(cases.publicTestCases.map(tc => {
                                const ivs = {};
                                (tc.input || []).forEach(i => {
                                    ivs[i.name] = { name: i.name, type: i.type, value: i.value };
                                });
                                return { name: tc.name, type: 'public', inputValues: ivs, outputValue: tc.output?.value || '' };
                            }));
                        }
                        if (cases.hiddenTestCases) {
                            setHiddenTestCases(cases.hiddenTestCases.map(tc => {
                                const ivs = {};
                                (tc.input || []).forEach(i => {
                                    ivs[i.name] = { name: i.name, type: i.type, value: i.value };
                                });
                                return { name: tc.name, type: 'hidden', inputValues: ivs, outputValue: tc.output?.value || '' };
                            }));
                        }
                        setFormData(prev => ({ ...prev, status: 'draft' }));
                        setActiveTab('testcases');
                    }}
                />
            )}

            {!showPreview ? (
                <>
                    <div className="tabs">
                        <button className={activeTab === 'basic' ? 'tab active' : 'tab'} onClick={() => setActiveTab('basic')}>📝 Básicos</button>
                        <button className={activeTab === 'templates' ? 'tab active' : 'tab'} onClick={() => setActiveTab('templates')}>💻 Plantillas</button>
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

                                <div className="io-variables-section">
                                    <div className="section-header-mini">
                                        <h3>Variables Globales</h3>
                                        <button className="btn-add-mini" onClick={addInputVar}>+ Variable</button>
                                    </div>
                                    <div className="vars-grid">
                                        {formData.inputVariables.map((iv, idx) => (
                                            <div key={idx} className="var-item-card">
                                                <div className="var-main-info">
                                                    <input type="text" placeholder="Nombre" value={iv.name} onChange={(e) => handleInputVarChange(idx, 'name', e.target.value)} />
                                                    <select value={iv.type} onChange={(e) => handleInputVarChange(idx, 'type', e.target.value)}>
                                                        <option value="string">String</option>
                                                        <option value="int">Integer</option>
                                                        <option value="float">Float</option>
                                                        <option value="boolean">Boolean</option>
                                                        <option value="array">Array</option>
                                                    </select>
                                                </div>
                                                <button className="btn-remove-var" onClick={() => removeInputVar(idx)}>×</button>
                                            </div>
                                        ))}
                                    </div>
                                    
                                    <h3 style={{ marginTop: '1.5rem', fontSize: '1rem', color: 'var(--text-muted)' }}>Variable de Salida</h3>
                                    <div className="var-item-card output-var">
                                        <input type="text" placeholder="Nombre" value={formData.outputVariable.name} onChange={(e) => handleOutputVarChange('name', e.target.value)} />
                                        <select value={formData.outputVariable.type} onChange={(e) => handleOutputVarChange('type', e.target.value)}>
                                            <option value="string">String</option>
                                            <option value="int">Integer</option>
                                            <option value="float">Float</option>
                                            <option value="boolean">Boolean</option>
                                            <option value="array">Array</option>
                                        </select>
                                    </div>
                                    
                                    
                                    <div className="form-group" style={{ marginTop: '1.5rem' }}>
                                        <label>Restricciones (Constraints)</label>
                                        <textarea 
                                            name="constraints" 
                                            value={formData.constraints} 
                                            onChange={handleChange} 
                                            placeholder="Ej: 1 <= nums.length <= 10^4"
                                            rows="2"
                                            style={{ height: 'auto' }}
                                        />
                                    </div>

                                    <div className="form-group" style={{ marginTop: '1.5rem' }}>
                                        <label>Etiquetas (Tags)</label>
                                        <div style={{ display: 'flex', flexWrap: 'wrap', gap: '0.5rem', marginBottom: '0.5rem' }}>
                                            {formData.tags.map(tag => (
                                                <span key={tag} style={{ background: 'rgba(200,16,46,0.1)', color: '#c8102e', padding: '4px 8px', borderRadius: '4px', fontSize: '0.85rem', display: 'flex', alignItems: 'center', gap: '4px', fontWeight: 'bold' }}>
                                                    {tag} <button type="button" onClick={() => removeTag(tag)} style={{ background: 'none', border: 'none', color: '#c8102e', cursor: 'pointer', fontSize: '1rem', padding: '0' }}>&times;</button>
                                                </span>
                                            ))}
                                        </div>
                                        <div style={{ display: 'flex', gap: '0.5rem' }}>
                                            <input type="text" value={newTag} onChange={(e) => setNewTag(e.target.value)} placeholder="Ej: math, arrays" onKeyDown={(e) => { if(e.key === 'Enter') { e.preventDefault(); addTag(newTag); } }} style={{ flex: 1 }} />
                                            <button type="button" className="btn-secondary" onClick={() => addTag(newTag)}>Agregar</button>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        )}

                        {activeTab === 'templates' && (
                            <div className="form-section templates-view">
                                <div className="section-header-mini">
                                    <h3>Plantillas de Código</h3>
                                    <p style={{ margin: 0, fontSize: '0.85rem', color: '#666' }}>Activa los lenguajes permitidos para este reto y edita su plantilla inicial.</p>
                                </div>
                                <div className="templates-list" style={{ display: 'flex', flexDirection: 'column', gap: '1rem', marginTop: '1rem' }}>
                                    {SUPPORTED_LANGUAGES.map(lang => {
                                        const isEnabled = formData.codeTemplates.hasOwnProperty(lang.id);
                                        const templateCode = isEnabled ? formData.codeTemplates[lang.id] : lang.defaultTemplate;
                                        
                                        return (
                                            <div key={lang.id} className="template-card" style={{ border: '1px solid #e0e0e0', borderRadius: '8px', overflow: 'hidden' }}>
                                                <div className="template-header" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', padding: '0.75rem 1rem', background: '#f9f9fa', borderBottom: isEnabled ? '1px solid #e0e0e0' : 'none' }}>
                                                    <div style={{ fontWeight: 'bold', display: 'flex', alignItems: 'center', gap: '8px' }}>
                                                        <span style={{ fontSize: '1.2rem' }}>{lang.id === 'python' ? '🐍' : lang.id === 'javascript' ? '🟨' : lang.id === 'java' ? '☕' : lang.id === 'cpp' ? '⚙️' : '🐹'}</span>
                                                        {lang.label}
                                                    </div>
                                                    <label className="switch" style={{ display: 'flex', alignItems: 'center', cursor: 'pointer' }}>
                                                        <input 
                                                            type="checkbox" 
                                                            checked={isEnabled} 
                                                            onChange={(e) => {
                                                                const newTemplates = { ...formData.codeTemplates };
                                                                if (e.target.checked) {
                                                                    newTemplates[lang.id] = lang.defaultTemplate;
                                                                } else {
                                                                    delete newTemplates[lang.id];
                                                                }
                                                                setFormData(prev => ({ ...prev, codeTemplates: newTemplates }));
                                                            }} 
                                                            style={{ marginRight: '8px' }}
                                                        />
                                                        <span style={{ fontSize: '0.85rem', color: isEnabled ? '#10b981' : '#888', fontWeight: 'bold' }}>{isEnabled ? 'Habilitado' : 'Deshabilitado'}</span>
                                                    </label>
                                                </div>
                                                {isEnabled && (
                                                    <div className="template-editor" style={{ padding: '0' }}>
                                                        <textarea 
                                                            value={templateCode}
                                                            onChange={(e) => {
                                                                setFormData(prev => ({ ...prev, codeTemplates: { ...prev.codeTemplates, [lang.id]: e.target.value } }));
                                                            }}
                                                            style={{ width: '100%', height: '120px', border: 'none', padding: '1rem', fontFamily: 'monospace', fontSize: '0.9rem', resize: 'vertical', background: '#fff', color: '#333' }}
                                                            spellCheck="false"
                                                        />
                                                    </div>
                                                )}
                                            </div>
                                        );
                                    })}
                                </div>
                            </div>
                        )}

                        {activeTab === 'testcases' && (
                            <div className="form-section test-cases-view">
                                <div className="section-header-row">
                                    <div className="header-info">
                                        <h3>Casos Públicos</h3>
                                        <p>Ejemplos visibles en el enunciado</p>
                                    </div>
                                    <button onClick={addPublicTestCase} className="btn-add-rich">+ Caso Ejemplo</button>
                                </div>

                                <div className="tc-cards-grid">
                                    {publicTestCases.map((tc, idx) => (
                                        <div key={idx} className="tc-rich-card public">
                                            <div className="tc-card-header">
                                                <div className="tc-title">
                                                    <span className="idx-tag">#{idx + 1}</span>
                                                    <input 
                                                        value={tc.name} 
                                                        onChange={(e) => updateTestCase(idx, 'name', e.target.value, true)} 
                                                        placeholder="Nombre del caso"
                                                    />
                                                </div>
                                                <button className="btn-icon-trash" onClick={() => removeTestCase(idx, true)}>🗑️</button>
                                            </div>
                                            
                                            <div className="tc-variables-grid">
                                                {formData.inputVariables.map(iv => {
                                                    const tcVar = tc.inputValues?.[iv.name] || { type: iv.type, value: '' };
                                                    return (
                                                        <div key={iv.name} className="tc-var-row">
                                                            <div className="var-label">
                                                                <span className="var-name">{iv.name}</span>
                                                                <select 
                                                                    className="var-type-select"
                                                                    value={tcVar.type || iv.type} 
                                                                    onChange={(e) => updateTestCaseInput(idx, iv.name, 'type', e.target.value, true)}
                                                                >
                                                                    <option value="string">str</option>
                                                                    <option value="int">int</option>
                                                                    <option value="array">arr</option>
                                                                    <option value="boolean">bool</option>
                                                                </select>
                                                            </div>
                                                            <input 
                                                                value={tcVar.value || ''} 
                                                                onChange={(e) => updateTestCaseInput(idx, iv.name, 'value', e.target.value, true)} 
                                                                placeholder={`Valor para ${iv.name}`} 
                                                            />
                                                        </div>
                                                    );
                                                })}
                                                <div className="tc-var-row output-var">
                                                    <div className="var-label">
                                                        <span className="var-name">Salida ({formData.outputVariable.name})</span>
                                                    </div>
                                                    <input 
                                                        value={tc.outputValue} 
                                                        onChange={(e) => updateTestCase(idx, 'outputValue', e.target.value, true)} 
                                                        placeholder="Resultado esperado" 
                                                    />
                                                </div>
                                            </div>
                                        </div>
                                    ))}
                                </div>

                                <div className="section-header-row" style={{ marginTop: '3rem' }}>
                                    <div className="header-info">
                                        <h3>Casos Ocultos</h3>
                                        <p>Se usarán para la calificación final</p>
                                    </div>
                                    <button onClick={addHiddenTestCase} className="btn-add-rich yellow">+ Caso Oculto</button>
                                </div>

                                <div className="tc-cards-grid">
                                    {hiddenTestCases.map((tc, idx) => (
                                        <div key={idx} className="tc-rich-card hidden-tc">
                                            <div className="tc-card-header">
                                                <div className="tc-title">
                                                    <span className="idx-tag">#{idx + 1}</span>
                                                    <input 
                                                        value={tc.name} 
                                                        onChange={(e) => updateTestCase(idx, 'name', e.target.value, false)} 
                                                        placeholder="Nombre del caso"
                                                    />
                                                </div>
                                                <button className="btn-icon-trash" onClick={() => removeTestCase(idx, false)}>🗑️</button>
                                            </div>
                                            
                                            <div className="tc-variables-grid">
                                                {formData.inputVariables.map(iv => {
                                                    const tcVar = tc.inputValues?.[iv.name] || { type: iv.type, value: '' };
                                                    return (
                                                        <div key={iv.name} className="tc-var-row">
                                                            <div className="var-label">
                                                                <span className="var-name">{iv.name}</span>
                                                                <select 
                                                                    className="var-type-select"
                                                                    value={tcVar.type || iv.type} 
                                                                    onChange={(e) => updateTestCaseInput(idx, iv.name, 'type', e.target.value, false)}
                                                                >
                                                                    <option value="string">str</option>
                                                                    <option value="int">int</option>
                                                                    <option value="array">arr</option>
                                                                    <option value="boolean">bool</option>
                                                                </select>
                                                            </div>
                                                            <input 
                                                                value={tcVar.value || ''} 
                                                                onChange={(e) => updateTestCaseInput(idx, iv.name, 'value', e.target.value, false)} 
                                                                placeholder={`Valor para ${iv.name}`} 
                                                            />
                                                        </div>
                                                    );
                                                })}
                                                <div className="tc-var-row output-var">
                                                    <div className="var-label">
                                                        <span className="var-name">Salida ({formData.outputVariable.name})</span>
                                                    </div>
                                                    <input 
                                                        value={tc.outputValue} 
                                                        onChange={(e) => updateTestCase(idx, 'outputValue', e.target.value, false)} 
                                                        placeholder="Resultado esperado" 
                                                    />
                                                </div>
                                            </div>
                                        </div>
                                    ))}
                                </div>
                            </div>
                        )}

                        {activeTab === 'settings' && (
                            <div className="form-section">
                                <div style={{ display: 'none' }}>
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
                                    </div>
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
