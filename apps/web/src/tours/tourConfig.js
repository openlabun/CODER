// using driver.js from CDN to bypass npm permission issues
const driver = window.driver.js.driver;

export const studentTourSteps = [
  { element: '#tour-nav-dashboard', popover: { title: 'Panel Principal', description: 'Aquí puedes ver un resumen de tu actividad, envíos recientes y acceso rápido.', side: 'right', align: 'start' }},
  { element: '#tour-nav-challenges', popover: { title: 'Retos', description: 'Encuentra desafíos de programación para poner a prueba tus habilidades.', side: 'right', align: 'start' }},
  { element: '#tour-nav-courses', popover: { title: 'Tus Cursos', description: 'Accede a los cursos en los que estás inscrito y continúa aprendiendo.', side: 'right', align: 'start' }},
  { element: '#tour-nav-submissions', popover: { title: 'Tus Envíos', description: 'Revisa el historial de tus soluciones y su estado de evaluación.', side: 'right', align: 'start' }},
  { element: '#tour-nav-leaderboard', popover: { title: 'Clasificación', description: '¡Compite con otros estudiantes y sube en el ranking global!', side: 'right', align: 'start' }},
  { element: '.user-card', popover: { title: 'Tu Perfil', description: 'Aquí puedes ver tu rol actual y cerrar sesión cuando lo desees.', side: 'right', align: 'start' }}
];

export const professorTourSteps = [
  { element: '#tour-nav-dashboard', popover: { title: 'Panel Principal', description: 'Resumen de la plataforma y accesos rápidos como profesor.', side: 'right', align: 'start' }},
  { element: '#tour-nav-challenges', popover: { title: 'Retos', description: 'Explora y resuelve los retos disponibles en la plataforma.', side: 'right', align: 'start' }},
  { element: '#tour-nav-challenges-create', popover: { title: 'Crear Reto', description: 'Diseña nuevos retos de programación, define las entradas/salidas y los casos de prueba.', side: 'right', align: 'start' }},
  { element: '#tour-nav-courses', popover: { title: 'Gestión de Cursos', description: 'Crea nuevos cursos, administra los estudiantes inscritos y asigna retos.', side: 'right', align: 'start' }},
  { element: '#tour-nav-submissions', popover: { title: 'Envíos Globales', description: 'Monitorea los envíos de los estudiantes a través de la plataforma.', side: 'right', align: 'start' }},
  { element: '#tour-nav-leaderboard', popover: { title: 'Clasificación', description: 'Revisa quiénes son los mejores programadores de RobleCode.', side: 'right', align: 'start' }},
  { element: '.user-card', popover: { title: 'Tu Perfil', description: 'Aquí puedes ver tu rol actual y cerrar sesión cuando lo desees.', side: 'right', align: 'start' }}
];

export const startTour = (role) => {
  const steps = role === 'professor' || role === 'teacher' ? professorTourSteps : studentTourSteps;
  
  const tourDriver = driver({
    showProgress: true,
    steps: steps,
    nextBtnText: 'Siguiente',
    prevBtnText: 'Anterior',
    doneBtnText: 'Finalizar',
    animate: true,
  });

  tourDriver.drive();
};
