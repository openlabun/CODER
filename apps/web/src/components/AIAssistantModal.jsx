import { useState } from 'react';
import client from '../api/client';
import Swal from 'sweetalert2';
import './AIAssistantModal.css';

const AIAssistantModal = ({ onClose, onApplyIdea, onApplyTestCases, onApplyExam, initialTab = 'full' }) => {
    const [activeTab, setActiveTab] = useState(initialTab);
    const [loading, setLoading] = useState(false);

    // For generating ideas
    const [topic, setTopic] = useState('');
    const [difficulty, setDifficulty] = useState('medium');
    const [ideas, setIdeas] = useState([]);

    // For generating full challenge
    const [fullChallenge, setFullChallenge] = useState(null);

    // For generating test cases
    const [challengeDesc, setChallengeDesc] = useState('');
    const [inputFormat, setInputFormat] = useState('');
    const [outputFormat, setOutputFormat] = useState('');
    const [testCases, setTestCases] = useState(null);

    // For generating exam
    const [examResult, setExamResult] = useState(null);

    const handleGenerateIdeas = async () => {
        if (!topic.trim()) {
            Swal.fire({
                icon: 'warning',
                title: 'Campo requerido',
                text: 'Por favor, ingresa un tema o categoría',
                toast: true,
                position: 'top-end',
                showConfirmButton: false,
                timer: 4000
            });
            return;
        }

        setLoading(true);
        setIdeas([]);

        try {
            const response = await client.post('/ai/generate-challenge-ideas', {
                topic,
                difficulty: difficulty || undefined,
                count: 3
            });
            setIdeas(response.data.ideas);
        } catch (err) {
            Swal.fire({
                icon: 'error',
                title: 'Error de IA',
                text: err.response?.data?.error || 'No se pudieron generar ideas en este momento.',
                toast: true,
                position: 'top-end',
                showConfirmButton: false,
                timer: 4000,
                timerProgressBar: true
            });
            console.error(err);
        } finally {
            setLoading(false);
        }
    };

    const handleGenerateFullChallenge = async () => {
        if (!topic.trim()) {
            Swal.fire({
                icon: 'warning',
                title: 'Campo requerido',
                text: 'Por favor, ingresa un tema para el reto',
                toast: true,
                position: 'top-end',
                showConfirmButton: false,
                timer: 4000
            });
            return;
        }

        setLoading(true);
        setFullChallenge(null);

        try {
            const response = await client.post('/ai/generate-full-challenge', {
                topic,
                difficulty: difficulty || 'medium'
            });
            setFullChallenge(response.data.challenge);
        } catch (err) {
            Swal.fire({
                icon: 'error',
                title: 'Error creando reto',
                text: err.response?.data?.error || 'La IA tuvo un problema diseñando el reto completo.',
                toast: true,
                position: 'top-end',
                showConfirmButton: false,
                timer: 4000,
                timerProgressBar: true
            });
            console.error(err);
        } finally {
            setLoading(false);
        }
    };

    const handleGenerateTestCases = async () => {
        if (!challengeDesc.trim() || !inputFormat.trim() || !outputFormat.trim()) {
            Swal.fire({
                icon: 'warning',
                title: 'Campos incompletos',
                text: 'Por favor, completa todos los campos',
                toast: true,
                position: 'top-end',
                showConfirmButton: false,
                timer: 4000
            });
            return;
        }

        setLoading(true);
        setTestCases(null);

        try {
            const response = await client.post('/ai/generate-test-cases', {
                challengeDescription: challengeDesc,
                inputFormat,
                outputFormat,
                publicCount: 2,
                hiddenCount: 3
            });
            setTestCases(response.data);
        } catch (err) {
            Swal.fire({
                icon: 'error',
                title: 'Error en casos',
                text: err.response?.data?.error || 'No se pudo generar los casos de prueba.',
                toast: true,
                position: 'top-end',
                showConfirmButton: false,
                timer: 4000,
                timerProgressBar: true
            });
            console.error(err);
        } finally {
            setLoading(false);
        }
    };

    const handleGenerateExam = async () => {
        if (!topic.trim()) {
            Swal.fire({ icon: 'warning', title: 'Campo requerido', text: 'Ingresa un tema para el examen', toast: true, position: 'top-end', showConfirmButton: false, timer: 4000 });
            return;
        }

        setLoading(true);
        setExamResult(null);

        try {
            const response = await client.post('/ai/generate-exam', {
                topic,
                difficulty: difficulty || 'medium'
            });
            setExamResult(response.data.exam);
        } catch (err) {
            Swal.fire({
                icon: 'error',
                title: 'Error IA Examen',
                text: err.response?.data?.error || 'No se pudo generar el examen.',
                toast: true,
                position: 'top-end',
                showConfirmButton: false,
                timer: 4000
            });
        } finally {
            setLoading(false);
        }
    };

    const applyIdea = (idea) => {
        onApplyIdea(idea);
        onClose();
    };

    const applyTestCases = () => {
        if (testCases) {
            onApplyTestCases(testCases);
            onClose();
        }
    };

    const applyExam = () => {
        if (examResult) {
            onApplyExam(examResult);
            onClose();
        }
    };

    return (
        <div className="modal-overlay" onClick={onClose}>
            <div className="ai-modal-content" onClick={(e) => e.stopPropagation()}>
                <div className="ai-modal-header">
                    <h2>✨ Asistente Creativo IA</h2>
                    <button onClick={onClose} className="close-btn">×</button>
                </div>

                <div className="ai-tabs">
                    <button
                        className={activeTab === 'full' ? 'ai-tab active' : 'ai-tab'}
                        onClick={() => setActiveTab('full')}
                    >
                        ✨ Reto Completo
                    </button>
                    <button
                        className={activeTab === 'ideas' ? 'ai-tab active' : 'ai-tab'}
                        onClick={() => setActiveTab('ideas')}
                    >
                        💡 Ideas
                    </button>
                    <button
                        className={activeTab === 'testcases' ? 'ai-tab active' : 'ai-tab'}
                        onClick={() => setActiveTab('testcases')}
                    >
                        🧪 Casos de Prueba
                    </button>
                    {onApplyExam && (
                        <button
                            className={activeTab === 'exam' ? 'ai-tab active' : 'ai-tab'}
                            onClick={() => setActiveTab('exam')}
                        >
                            📝 Examen
                        </button>
                    )}
                </div>

                <div className="ai-modal-body">
                    {activeTab === 'full' && (
                        <div className="ai-section">
                            <div className="ai-input-group">
                                <label>Tema del Reto *</label>
                                <input
                                    type="text"
                                    value={topic}
                                    onChange={(e) => setTopic(e.target.value)}
                                    placeholder="Ej: Suma de dos números, Palíndromos, Grafos..."
                                    disabled={loading}
                                />
                            </div>
                            <div className="ai-input-group">
                                <label>Nivel Sugerido</label>
                                <select value={difficulty} onChange={(e) => setDifficulty(e.target.value)} disabled={loading}>
                                    <option value="easy">Fácil</option>
                                    <option value="medium">Medio</option>
                                    <option value="hard">Difícil</option>
                                </select>
                            </div>

                            <button
                                onClick={handleGenerateFullChallenge}
                                disabled={loading}
                                className="ai-generate-btn full-gen"
                            >
                                {loading ? '🧠 Generándolo todo...' : '✨ Crear Reto por IA'}
                            </button>

                            {fullChallenge && (
                                <div className="ai-results">
                                    <div className="ai-full-preview">
                                        <h3>{fullChallenge.title}</h3>
                                        <p>{fullChallenge.description}</p>
                                        <div className="preview-meta">
                                            <span>Entrada: {fullChallenge.inputFormat}</span>
                                            <span>Salida: {fullChallenge.outputFormat}</span>
                                        </div>
                                        <button
                                            onClick={() => applyIdea(fullChallenge)}
                                            className="ai-apply-btn primary"
                                        >
                                            ✅ Usar este reto completo
                                        </button>
                                    </div>
                                </div>
                            )}
                        </div>
                    )}
                    {activeTab === 'ideas' && (
                        <div className="ai-section">
                            <div className="ai-input-group">
                                <label>Tema o Categoría *</label>
                                <input
                                    type="text"
                                    value={topic}
                                    onChange={(e) => setTopic(e.target.value)}
                                    placeholder="Ej: Búsqueda Binaria, Programación Dinámica, Árboles"
                                    disabled={loading}
                                />
                            </div>

                            <div className="ai-input-group">
                                <label>Dificultad (Opcional)</label>
                                <select value={difficulty} onChange={(e) => setDifficulty(e.target.value)} disabled={loading}>
                                    <option value="">Cualquiera</option>
                                    <option value="easy">Fácil</option>
                                    <option value="medium">Medio</option>
                                    <option value="hard">Difícil</option>
                                </select>
                            </div>

                            <button
                                onClick={handleGenerateIdeas}
                                disabled={loading}
                                className="ai-generate-btn"
                            >
                                {loading ? '🔄 Generando...' : '✨ Generar Ideas'}
                            </button>

                            {ideas.length > 0 && (
                                <div className="ai-results">
                                    <h3>Ideas Generadas</h3>
                                    {ideas.map((idea, idx) => (
                                        <div key={idx} className="ai-idea-card">
                                            <div className="ai-idea-header">
                                                <h4>{idea.title}</h4>
                                                <span className={`difficulty-badge ${idea.difficulty}`}>
                                                    {idea.difficulty === 'easy' ? 'Fácil' : idea.difficulty === 'hard' ? 'Difícil' : 'Medio'}
                                                </span>
                                            </div>
                                            <p className="ai-idea-description">{idea.description}</p>
                                            <div className="ai-idea-tags">
                                                {idea.tags?.map((tag, i) => (
                                                    <span key={i} className="tag">{tag}</span>
                                                ))}
                                            </div>
                                            <button
                                                onClick={() => applyIdea(idea)}
                                                className="ai-apply-btn"
                                            >
                                                ✅ Usar esta idea
                                            </button>
                                        </div>
                                    ))}
                                </div>
                            )}
                        </div>
                    )}

                    {activeTab === 'testcases' && (
                        <div className="ai-section">
                            <div className="ai-input-group">
                                <label>Descripción del Reto *</label>
                                <textarea
                                    value={challengeDesc}
                                    onChange={(e) => setChallengeDesc(e.target.value)}
                                    placeholder="Describe el problema brevemente..."
                                    rows="4"
                                    disabled={loading}
                                />
                            </div>

                            <div className="ai-input-group">
                                <label>Formato de Entrada *</label>
                                <textarea
                                    value={inputFormat}
                                    onChange={(e) => setInputFormat(e.target.value)}
                                    placeholder="Describe el formato de entrada..."
                                    rows="3"
                                    disabled={loading}
                                />
                            </div>

                            <div className="ai-input-group">
                                <label>Formato de Salida *</label>
                                <textarea
                                    value={outputFormat}
                                    onChange={(e) => setOutputFormat(e.target.value)}
                                    placeholder="Describe la salida esperada..."
                                    rows="3"
                                    disabled={loading}
                                />
                            </div>

                            <button
                                onClick={handleGenerateTestCases}
                                disabled={loading}
                                className="ai-generate-btn"
                            >
                                {loading ? '🔄 Generando...' : '🧪 Generar Casos de Prueba'}
                            </button>

                            {testCases && (
                                <div className="ai-results">
                                    <h3>Casos de Prueba Generados</h3>

                                    <div className="ai-testcases-section">
                                        <h4>📖 Casos Públicos ({testCases.publicTestCases?.length || 0})</h4>
                                        {testCases.publicTestCases?.map((tc, idx) => (
                                            <div key={idx} className="ai-testcase-card">
                                                <strong>{tc.name}</strong>
                                                <div className="ai-testcase-io">
                                                    <div>
                                                        <label>Entrada:</label>
                                                        <pre>{tc.input}</pre>
                                                    </div>
                                                    <div>
                                                        <label>Salida:</label>
                                                        <pre>{tc.output}</pre>
                                                    </div>
                                                </div>
                                            </div>
                                        ))}
                                    </div>

                                    <div className="ai-testcases-section">
                                        <h4>🔒 Casos Ocultos ({testCases.hiddenTestCases?.length || 0})</h4>
                                        {testCases.hiddenTestCases?.map((tc, idx) => (
                                            <div key={idx} className="ai-testcase-card">
                                                <strong>{tc.name}</strong>
                                                <div className="ai-testcase-io">
                                                    <div>
                                                        <label>Entrada:</label>
                                                        <pre>{tc.input}</pre>
                                                    </div>
                                                    <div>
                                                        <label>Salida:</label>
                                                        <pre>{tc.output}</pre>
                                                    </div>
                                                </div>
                                            </div>
                                        ))}
                                    </div>

                                    <button
                                        onClick={applyTestCases}
                                        className="ai-apply-all-btn"
                                    >
                                        ✅ Aplicar todos los casos
                                    </button>
                                </div>
                            )}
                        </div>
                    )}

                    {activeTab === 'exam' && (
                        <div className="ai-section">
                            <div className="ai-input-group">
                                <label>Objetivo Académico / Tema *</label>
                                <input
                                    type="text"
                                    value={topic}
                                    onChange={(e) => setTopic(e.target.value)}
                                    placeholder="Ej: Conceptos básicos de Python, Programación Orientada a Objetos..."
                                    disabled={loading}
                                />
                            </div>
                            <div className="ai-input-group">
                                <label>Nivel de Examen</label>
                                <select value={difficulty} onChange={(e) => setDifficulty(e.target.value)} disabled={loading}>
                                    <option value="easy">Introductorio</option>
                                    <option value="medium">Intermedio</option>
                                    <option value="hard">Avanzado / Final</option>
                                </select>
                            </div>

                            <button
                                onClick={handleGenerateExam}
                                disabled={loading}
                                className="ai-generate-btn full-gen"
                            >
                                {loading ? '🧠 Diseñando evaluación...' : '✨ Generar Estructura de Examen'}
                            </button>

                            {examResult && (
                                <div className="ai-results">
                                    <div className="ai-full-preview">
                                        <h3>{examResult.title}</h3>
                                        <p>{examResult.description}</p>
                                        <div className="preview-meta">
                                            <span>⏱️ Sugerido: {examResult.time_limit} min</span>
                                            <span>🔄 Intentos: {examResult.try_limit}</span>
                                        </div>
                                        <button
                                            onClick={applyExam}
                                            className="ai-apply-btn primary"
                                        >
                                            ✅ Usar esta configuración
                                        </button>
                                    </div>
                                </div>
                            )}
                        </div>
                    )}
                </div>

                <div className="ai-modal-footer">
                    <p className="ai-disclaimer">
                        ⚠️ El contenido generado por IA debe ser revisado y validado antes de publicar.
                    </p>
                </div>
            </div>
        </div>
    );
};

export default AIAssistantModal;
