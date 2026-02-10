import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Layout from './components/Layout';
import Home from './pages/Home';
import Login from './pages/Login';
import Register from './pages/Register';
import Challenges from './pages/Challenges';
import ChallengeSolver from './pages/ChallengeSolver';
import Leaderboard from './pages/Leaderboard';
import Courses from './pages/Courses';
import CourseDetails from './pages/CourseDetails';
import CourseStudents from './pages/CourseStudents';
import ExamRunner from './pages/ExamRunner';
import Dashboard from './pages/Dashboard';
import CreateChallenge from './pages/CreateChallenge';
import CreateCourse from './pages/CreateCourse';
import EditCourse from './pages/EditCourse';
import JoinCourse from './pages/JoinCourse';
import BrowseCourses from './pages/BrowseCourses';
import Submissions from './pages/Submissions';
import ProtectedRoute from './components/ProtectedRoute';

// Main App component
function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Layout />}>
          <Route index element={<Home />} />
          <Route path="login" element={<Login />} />
          <Route path="register" element={<Register />} />
          <Route path="challenges" element={<ProtectedRoute><Challenges /></ProtectedRoute>} />
          <Route path="challenge/:id" element={<ProtectedRoute><ChallengeSolver /></ProtectedRoute>} />
          <Route path="leaderboard" element={<ProtectedRoute><Leaderboard /></ProtectedRoute>} />
          <Route path="courses" element={<ProtectedRoute><Courses /></ProtectedRoute>} />
          <Route path="courses/join" element={<ProtectedRoute><JoinCourse /></ProtectedRoute>} />
          <Route path="courses/browse" element={<ProtectedRoute><BrowseCourses /></ProtectedRoute>} />
          <Route path="courses/:id" element={<ProtectedRoute><CourseDetails /></ProtectedRoute>} />
          <Route path="courses/:id/students" element={<ProtectedRoute><CourseStudents /></ProtectedRoute>} />
          <Route path="exam/:id" element={<ProtectedRoute><ExamRunner /></ProtectedRoute>} />
          <Route path="dashboard" element={<ProtectedRoute><Dashboard /></ProtectedRoute>} />
          <Route path="submissions" element={<ProtectedRoute><Submissions /></ProtectedRoute>} />
          <Route path="challenges/create" element={<ProtectedRoute><CreateChallenge /></ProtectedRoute>} />
          <Route path="challenges/edit/:id" element={<ProtectedRoute><CreateChallenge /></ProtectedRoute>} />
          <Route path="courses/create" element={<ProtectedRoute><CreateCourse /></ProtectedRoute>} />
          <Route path="courses/edit/:id" element={<ProtectedRoute><EditCourse /></ProtectedRoute>} />
        </Route>
      </Routes>
    </Router>
  );
}

export default App;
