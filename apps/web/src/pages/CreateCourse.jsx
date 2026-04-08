import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import client from '../api/client';
import Swal from 'sweetalert2';
import './CreateCourse.css';

const CreateCourse = () => {
    const navigate = useNavigate();
    const currentYear = new Date().getFullYear().toString();
    const [formData, setFormData] = useState({
        name: '',
        code: '',
        period: `${currentYear}-1`, // Default period
        groupNumber: 1,
        description: '',
        color: '#00f0ff',
        enrollmentMethod: 'code',
        enrollmentCode: '',
        startDate: '',
        endDate: '',
        visibility: 'active'
    });
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');
    const [validationErrors, setValidationErrors] = useState({});

    const courseColors = [
        '#00f0ff', '#7000ff', '#ff0055', '#00ff9d',
        '#ffcc00', '#ff6b35', '#4ecdc4', '#95e1d3'
    ];

    const handleChange = (e) => {
        const { name, value } = e.target;
        setFormData(prev => ({ ...prev, [name]: value }));

        // Clear validation error for this field
        if (validationErrors[name]) {
            setValidationErrors(prev => ({ ...prev, [name]: '' }));
        }
    };

    const generateEnrollmentCode = () => {
        const code = formData.code.toUpperCase() + '-' +
            formData.period.split('-').join('') + 'G' + formData.groupNumber;
        setFormData(prev => ({ ...prev, enrollmentCode: code }));
    };

    const validateForm = () => {
        const errors = {};

        // Course code validation (letters + numbers only)
        if (!/^[A-Za-z0-9]+$/.test(formData.code)) {
            errors.code = 'El código del curso debe contener solo letras y números';
        }

        // Period validation (YYYY-X format)
        if (!/^\d{4}-[123]$/.test(formData.period)) {
            errors.period = 'El periodo debe tener el formato YYYY-X';
        }

        // Group number validation
        if (formData.groupNumber < 1) {
            errors.groupNumber = 'El número de grupo debe ser un entero positivo';
        }

        // Date validation
        if (formData.startDate && formData.endDate) {
            if (new Date(formData.startDate) >= new Date(formData.endDate)) {
                errors.endDate = 'La fecha de finalización debe ser posterior a la fecha de inicio';
            }
        }

        setValidationErrors(errors);
        return Object.keys(errors).length === 0;
    };

    const handleSubmit = async (status) => {
        setError('');

        if (!validateForm()) {
            setError('Por favor corrige los errores de validación antes de enviar');
            return;
        }

        // Generate enrollment code if using code method and not set
        if (formData.enrollmentMethod === 'code' && !formData.enrollmentCode) {
            generateEnrollmentCode();
        }

        setLoading(true);
        try {
            const [year, semesterCode] = formData.period.split('-');
            
            // Map semesters to backend constants (01: first, 02: intersemestral/summer, 03: second)
            // 1 -> 01, 2 -> 03, 3 -> 02
            const semesterMap = { '1': '01', '2': '03', '3': '02' };
            const semester = semesterMap[semesterCode] || '01';

            // Map frontend fields to backend DTO
            const payload = {
                name: formData.name,
                description: formData.description,
                visibility: formData.visibility === 'active' ? 'public' : 'private',
                visual_identity: formData.color,
                code: formData.code,
                year: parseInt(year),
                semester: semester,
                enrollment_code: formData.enrollmentCode
            };

            await client.post('/courses', payload);
            
            Swal.fire({
                icon: 'success',
                title: '¡Curso Creado!',
                text: 'El curso se ha registrado exitosamente',
                timer: 1000,
                timerProgressBar: true,
                showConfirmButton: false,
                position: 'top-end',
                toast: true
            });

            setTimeout(() => navigate('/courses'), 1000);
        } catch (err) {
            console.error('Error creating course:', err);
            Swal.fire({
                icon: 'error',
                title: 'Error de Creación',
                text: 'No se pudo crear el curso. Revisa el código y periodo.',
                timer: 1500,
                showConfirmButton: false,
                position: 'center'
            });
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="create-course-page">
            <div className="page-header">
                <h1>Crear Nuevo Curso</h1>
                <p className="subtitle">Configura un curso completo para tus estudiantes</p>
            </div>

            {error && <div className="error-message">{error}</div>}

            <form className="course-form">
                <div className="form-section">
                    <h2>Información Básica</h2>

                    <div className="form-group">
                        <label htmlFor="name">Nombre del Curso *</label>
                        <input
                            type="text"
                            id="name"
                            name="name"
                            value={formData.name}
                            onChange={handleChange}
                            placeholder="ej. Algoritmos y Estructuras de Datos"
                            required
                        />
                    </div>

                    <div className="form-row">
                        <div className="form-group">
                            <label htmlFor="code">Código del Curso *</label>
                            <input
                                type="text"
                                id="code"
                                name="code"
                                value={formData.code}
                                onChange={handleChange}
                                placeholder="ej. CS101"
                                required
                                className={validationErrors.code ? 'error' : ''}
                            />
                            {validationErrors.code ? (
                                <small className="error-text">{validationErrors.code}</small>
                            ) : (
                                <small>Formato: CS101, MAT202, INF251</small>
                            )}
                        </div>

                        <div className="form-group">
                            <label>Periodo Académico *</label>
                            <div className="period-selectors">
                                <select
                                    value={formData.period ? formData.period.split('-')[0] : new Date().getFullYear()}
                                    onChange={(e) => {
                                        const currentTerm = (formData.period && formData.period.includes('-')) ? formData.period.split('-')[1] : '1';
                                        handleChange({ target: { name: 'period', value: `${e.target.value}-${currentTerm}` }});
                                    }}
                                    className="year-select"
                                >
                                    {[...Array(6)].map((_, i) => {
                                        const year = new Date().getFullYear() - 1 + i;
                                        return <option key={year} value={year}>{year}</option>;
                                    })}
                                </select>
                                <select
                                    value={formData.period ? formData.period.split('-')[1] : '1'}
                                    onChange={(e) => {
                                        const currentYear = formData.period ? formData.period.split('-')[0] : new Date().getFullYear();
                                        handleChange({ target: { name: 'period', value: `${currentYear}-${e.target.value}` }});
                                    }}
                                    className="term-select"
                                >
                                    <option value="1">Primer Semestre (10)</option>
                                    <option value="2">Segundo Semestre (30)</option>
                                    <option value="3">Verano (20)</option>
                                </select>
                            </div>
                            {validationErrors.period ? (
                                <small className="error-text">{validationErrors.period}</small>
                            ) : (
                                <small>Seleccione el año y ciclo académico (Formato Interno: {formData.period})</small>
                            )}
                        </div>

                        <div className="form-group">
                            <label htmlFor="groupNumber">Número de Grupo *</label>
                            <input
                                type="number"
                                id="groupNumber"
                                name="groupNumber"
                                value={formData.groupNumber}
                                onChange={handleChange}
                                min="1"
                                required
                                className={validationErrors.groupNumber ? 'error' : ''}
                            />
                            {validationErrors.groupNumber ? (
                                <small className="error-text">{validationErrors.groupNumber}</small>
                            ) : (
                                <small>Debe ser un número entero positivo</small>
                            )}
                        </div>
                    </div>

                    <div className="form-group">
                        <label htmlFor="description">Descripción del Curso (Opcional)</label>
                        <textarea
                            id="description"
                            name="description"
                            value={formData.description}
                            onChange={handleChange}
                            placeholder="Proporciona una breve descripción de los temas del curso y los objetivos de aprendizaje..."
                            rows="5"
                        />
                        <small>Describe objetivos, temas principales y propósito académico</small>
                    </div>
                </div>

                <div className="form-section">
                    <h2>Identidad Visual</h2>

                    <div className="form-group">
                        <label>Color del Curso (Opcional)</label>
                        <div className="color-picker">
                            {courseColors.map(color => (
                                <button
                                    key={color}
                                    type="button"
                                    className={`color-option ${formData.color === color ? 'selected' : ''}`}
                                    style={{ backgroundColor: color }}
                                    onClick={() => setFormData(prev => ({ ...prev, color }))}
                                    title={color}
                                />
                            ))}
                        </div>
                        <small>Elige un color para identificar visualmente este curso</small>
                    </div>
                </div>

                <div className="form-section">
                    <h2>Configuración de Inscripción</h2>

                    <div className="form-group">
                        <label>Método de Inscripción *</label>
                        <div className="radio-group">
                            <label className="radio-option">
                                <input
                                    type="radio"
                                    name="enrollmentMethod"
                                    value="code"
                                    checked={formData.enrollmentMethod === 'code'}
                                    onChange={handleChange}
                                />
                                <div>
                                    <span>🔑 Código de Inscripción</span>
                                    <small>Los estudiantes se unen usando un código único</small>
                                </div>
                            </label>
                            <label className="radio-option">
                                <input
                                    type="radio"
                                    name="enrollmentMethod"
                                    value="link"
                                    checked={formData.enrollmentMethod === 'link'}
                                    onChange={handleChange}
                                />
                                <div>
                                    <span>🔗 Enlace Privado</span>
                                    <small>Los estudiantes se unen a través de un enlace de invitación</small>
                                </div>
                            </label>
                            <label className="radio-option">
                                <input
                                    type="radio"
                                    name="enrollmentMethod"
                                    value="automatic"
                                    checked={formData.enrollmentMethod === 'automatic'}
                                    onChange={handleChange}
                                />
                                <div>
                                    <span>⚙️ Automático</span>
                                    <small>Inscripción gestionada por la institución</small>
                                </div>
                            </label>
                        </div>
                    </div>

                    {formData.enrollmentMethod === 'code' && (
                        <div className="form-group">
                            <label htmlFor="enrollmentCode">Código de Inscripción</label>
                            <div className="code-input-group">
                                <input
                                    type="text"
                                    id="enrollmentCode"
                                    name="enrollmentCode"
                                    value={formData.enrollmentCode}
                                    onChange={handleChange}
                                    placeholder="Se generará automáticamente"
                                    readOnly
                                />
                                <button type="button" onClick={generateEnrollmentCode} className="btn-generate">
                                    Generar Código
                                </button>
                            </div>
                            <small>Los estudiantes usarán este código para unirse al curso</small>
                        </div>
                    )}
                </div>

                <div className="form-section">
                    <h2>Calendario</h2>

                    <div className="form-row">
                        <div className="form-group">
                            <label htmlFor="startDate">Fecha de Inicio</label>
                            <input
                                type="date"
                                id="startDate"
                                name="startDate"
                                value={formData.startDate}
                                onChange={handleChange}
                            />
                            <small>Cuándo se activa el curso</small>
                        </div>

                        <div className="form-group">
                            <label htmlFor="endDate">Fecha de Finalización</label>
                            <input
                                type="date"
                                id="endDate"
                                name="endDate"
                                value={formData.endDate}
                                onChange={handleChange}
                                className={validationErrors.endDate ? 'error' : ''}
                            />
                            {validationErrors.endDate ? (
                                <small className="error-text">{validationErrors.endDate}</small>
                            ) : (
                                <small>Cuándo termina el curso</small>
                            )}
                        </div>
                    </div>
                </div>

                <div className="form-section">
                    <h2>Visibilidad</h2>

                    <div className="form-group">
                        <label>Visibilidad del Curso *</label>
                        <div className="radio-group">
                            <label className="radio-option">
                                <input
                                    type="radio"
                                    name="visibility"
                                    value="active"
                                    checked={formData.visibility === 'active'}
                                    onChange={handleChange}
                                />
                                <div>
                                    <span>✅ Activo</span>
                                    <small>Visible para los estudiantes inscritos</small>
                                </div>
                            </label>
                            <label className="radio-option">
                                <input
                                    type="radio"
                                    name="visibility"
                                    value="hidden"
                                    checked={formData.visibility === 'hidden'}
                                    onChange={handleChange}
                                />
                                <div>
                                    <span>👁️ Oculto / Borrador</span>
                                    <small>Aún no es visible para los estudiantes</small>
                                </div>
                            </label>
                        </div>
                    </div>
                </div>

                <div className="form-actions">
                    <button type="button" onClick={() => navigate('/courses')} className="btn-secondary">
                        Cancelar
                    </button>
                    <button
                        type="button"
                        onClick={() => handleSubmit('draft')}
                        disabled={loading}
                        className="btn-draft"
                    >
                        {loading ? 'Guardando...' : '💾 Guardar como Borrador'}
                    </button>
                    <button
                        type="button"
                        onClick={() => handleSubmit('published')}
                        disabled={loading}
                        className="btn-primary"
                    >
                        {loading ? 'Creando...' : '🚀 Crear Curso'}
                    </button>
                </div>
            </form>

            <div className="info-box">
                <h3>📚 Próximos Pasos</h3>
                <p>Después de crear el curso, podrás:</p>
                <ul>
                    <li>Inscribir estudiantes manualmente o compartir el código de inscripción</li>
                    <li>Asignar retos al curso</li>
                    <li>Crear exámenes y evaluaciones</li>
                    <li>Ver tablas de clasificación y analíticas del curso</li>
                </ul>
            </div>
        </div>
    );
};

export default CreateCourse;
