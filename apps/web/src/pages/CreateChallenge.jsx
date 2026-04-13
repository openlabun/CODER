import { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import client from '../api/client';
import AIAssistantModal from '../components/AIAssistantModal';
import Swal from 'sweetalert2';
import './CreateChallenge.css';
import './Challenges.css';

const CreateChallenge = () => {
    const { user } = useAuth();
    const navigate = useNavigate();
    const { id } = useParams();
    const isEditing = !!id;

    const [activeTab, setActiveTab] = useState('basic');
    const [showPreview, setShowPreview] = useState(false);
    const [showAIModal, setShowAIModal] = useState(false);

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
        status: 'draft'
    });

    const [publicTestCases, setPublicTestCases] = useState([]);
    const [hiddenTestCases, setHiddenTestCases] = useState([]);
    const [newTag, setNewTag] = useState('');
    const [loading, setLoading] = useState(false);
    const [loadingAction, setLoadingAction] = useState('published');
    const [fetching, setFetching] = useState(isEditing);

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
                    status: challenge.status || challenge.Status || 'draft'
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

        if (user) {
            if (isEditing) fetchChallengeForEdit();
        }
    }, [id, isEditing, user?.role]);

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

        const requestedStatus = status || formData.status;
        setLoadingAction(requestedStatus);
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
                status: requestedStatus,
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

            setLoading(false);

            const successTitle = isEditing ? 'Reto actualizado' : 'Reto creado';
            const successText = requestedStatus === 'draft'
                ? `El reto se guardo como borrador exitosamente.`
                : `El reto se creo y publico exitosamente.`;

            await Swal.fire({
                icon: 'success',
                title: successTitle,
                text: successText,
                confirmButtonText: 'Ir al repositorio de retos',
                allowOutsideClick: false,
                allowEscapeKey: false,
                customClass: { container: 'swal-ultra-high-z' }
            });

            navigate('/challenges');
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
            setLoadingAction('published');
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

    const loadingTitle = isEditing
        ? 'Actualizando reto...'
        : loadingAction === 'draft'
            ? 'Guardando borrador...'
            : 'Creando reto...';

    return (
        <div className="create-challenge-page">
            {loading && (
                <div className="challenge-save-overlay" role="status" aria-live="polite">
                    <div className="challenge-save-card">
                        <div className="challenge-save-spinner" />
                        <h2>{loadingTitle}</h2>
                        <p>Estamos procesando tu reto. Esto puede tardar unos segundos.</p>
                    </div>
                </div>
            )}

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
                                <div className="form-group">
                                    <label>Estado Inicial</label>
                                    <select name="status" value={formData.status === 'draft' ? 'private' : formData.status} onChange={handleChange}>
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
