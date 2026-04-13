import { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { Loader2 } from 'lucide-react';
import client from '../api/client';
import './CreateCourse.css'; // Reuse CreateCourse styles

const EditCourse = () => {
    const { id } = useParams();
    const navigate = useNavigate();
    const [formData, setFormData] = useState({
        name: '',
        code: '',
        period: `${new Date().getFullYear()}-1`,
        groupNumber: '',
    });
    const [loading, setLoading] = useState(true);
    const [saving, setSaving] = useState(false);
    const [error, setError] = useState('');

    useEffect(() => {
        const fetchCourse = async () => {
            try {
                const { data } = await client.get(`/courses/${id}`);
                setFormData({
                    name: data.name,
                    code: data.code,
                    period: data.period ? `${data.period.year}-${data.period.semester}` : `${new Date().getFullYear()}-1`,
                    groupNumber: data.groupNumber,
                });
            } catch (err) {
                setError('Failed to load course details');
            } finally {
                setLoading(false);
            }
        };
        fetchCourse();
    }, [id]);

    const handleSubmit = async (e) => {
        e.preventDefault();
        setSaving(true);
        setError('');

        try {
            const [year, semester] = formData.period.split('-');
            await client.post(`/courses/${id}`, {
                ...formData,
                year: parseInt(year),
                semester: semester,
                groupNumber: parseInt(formData.groupNumber),
            });
            navigate('/courses');
        } catch (err) {
            setError(err.response?.data?.message || 'Failed to update course');
        } finally {
            setSaving(false);
        }
    };

    const handleChange = (e) => {
        setFormData({
            ...formData,
            [e.target.name]: e.target.value,
        });
    };

    if (loading) return (
        <div className="page-loader">
            <Loader2 className="page-loader-spinner" size={48} />
            <p className="page-loader-text">Cargando curso...</p>
        </div>
    );

    return (
        <div className="create-course-page">
            <div className="page-header">
                <h1>Editar Curso</h1>
            </div>

            <form onSubmit={handleSubmit} className="create-course-form">
                {error && <div className="error-message">{error}</div>}

                <div className="form-group">
                    <label htmlFor="name">Nombre del Curso</label>
                    <input
                        type="text"
                        id="name"
                        name="name"
                        value={formData.name}
                        onChange={handleChange}
                        required
                        minLength={3}
                    />
                </div>

                <div className="form-row">
                    <div className="form-group">
                        <label htmlFor="code">Código del Curso</label>
                        <input
                            type="text"
                            id="code"
                            name="code"
                            value={formData.code}
                            onChange={handleChange}
                            required
                            minLength={2}
                        />
                    </div>

                    <div className="form-group">
                        <label>Periodo Académico</label>
                        <div className="period-selectors" style={{ display: 'flex', gap: '10px' }}>
                            <select
                                value={formData.period ? formData.period.split('-')[0] : new Date().getFullYear()}
                                onChange={(e) => {
                                    const currentTerm = (formData.period && formData.period.includes('-')) ? formData.period.split('-')[1] : '1';
                                    handleChange({ target: { name: 'period', value: `${e.target.value}-${currentTerm}` } });
                                }}
                                className="year-select"
                                style={{ flex: 1, padding: '12px', borderRadius: '4px', border: '1px solid #ccc' }}
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
                                    handleChange({ target: { name: 'period', value: `${currentYear}-${e.target.value}` } });
                                }}
                                className="term-select"
                                style={{ flex: 2, padding: '12px', borderRadius: '4px', border: '1px solid #ccc' }}
                            >
                                <option value="1">Primer Semestre (10)</option>
                                <option value="2">Segundo Semestre (30)</option>
                                <option value="3">Verano (20)</option>
                            </select>
                        </div>
                    </div>

                    <div className="form-group">
                        <label htmlFor="groupNumber">Número de Grupo</label>
                        <input
                            type="number"
                            id="groupNumber"
                            name="groupNumber"
                            value={formData.groupNumber}
                            onChange={handleChange}
                            required
                            min={1}
                        />
                    </div>
                </div>

                <div className="form-actions">
                    <button type="button" onClick={() => navigate('/courses')} className="btn-secondary">
                        Cancelar
                    </button>
                    <button type="submit" className="btn-primary" disabled={saving}>
                        {saving ? 'Guardando...' : 'Guardar Cambios'}
                    </button>
                </div>
            </form>
        </div>
    );
};

export default EditCourse;
