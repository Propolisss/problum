import React from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';

import Login from './pages/Login';
import Register from './pages/Register';
import Dashboard from './pages/Dashboard';
import Lesson from './pages/Lesson';
import Problem from './pages/Problem';
import Profile from './pages/Profile';
import Attempts from './pages/Attempts';
import AttemptDetail from './pages/AttemptDetail';
import CourseIndexPage from './pages/CourseIndexPage';

import { ProtectedRoute } from './routes/ProtectedRoute';
import Layout from './components/layout/Layout';
import CourseLayout from './components/layout/CourseLayout';

export default function App() {
  return (
    <Routes>
      <Route path="/login" element={<Login />} />
      <Route path="/register" element={<Register />} />

      <Route element={<ProtectedRoute><Layout /></ProtectedRoute>}>

        <Route path="/" element={<Dashboard />} />
        <Route path="/profile" element={<Profile />} />
        <Route path="/attempts" element={<Attempts />} />
        <Route path="/attempts/:id" element={<AttemptDetail />} />

        <Route path="/courses/:courseId" element={<CourseLayout />}>
          <Route index element={<CourseIndexPage />} />
          <Route path="lessons/:lessonId" element={<Lesson />} />
          <Route path="problems/:problemId" element={<Problem />} />
        </Route>

      </Route>

      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  );
}
